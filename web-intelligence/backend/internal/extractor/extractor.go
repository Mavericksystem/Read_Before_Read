package extractor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"web-intelligence/backend/internal/pool"
)

type Request struct {
	URL              string `json:"url"`
	MaxResponseBytes int64  `json:"max_response_bytes"`
	TimeoutMs        int64  `json:"timeout_ms"`
	PrerenderedHTML  string `json:"prerendered_html,omitempty"`
	FinalURL         string `json:"final_url,omitempty"`
}

type Document struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	ContentType        string `json:"content_type"`
	ContentLengthBytes int64  `json:"content_length_bytes"`
	FetchDurationMs    int64  `json:"fetch_duration_ms"`
}

type rustResponse struct {
	Status   string    `json:"status"`
	Document *Document `json:"document,omitempty"`
	Error    *RustErr  `json:"error,omitempty"`
}

type RustErr struct {
	Category string `json:"category"` // Aligning field name
	Message  string `json:"message"`
}

type Error struct {
	Category string
	Message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("extractor: %s: %s", e.Category, e.Message)
}

const binaryPath = "./bin/extractor"

type Pooled struct {
	pool *pool.Pool
}

// NewPooled starts `size` long-lived extractor worker processes.
func NewPooled(size int) (*Pooled, error) {
	p, err := pool.New(binaryPath, size)
	if err != nil {
		return nil, fmt.Errorf("extractor: %w", err)
	}
	return &Pooled{pool: p}, nil
}

// Run dispatches a normal fetch-and-extract job to the pool.
func (p *Pooled) Run(ctx context.Context, req Request) (*Document, error) {
	return p.dispatch(ctx, req)
}

// RunFromHTML dispatches an extract-only job using already-rendered HTML
// (from the headless browser fallback), skipping the worker's own fetch.
func (p *Pooled) RunFromHTML(ctx context.Context, html, finalURL string) (*Document, error) {
	return p.dispatch(ctx, Request{
		PrerenderedHTML: html,
		FinalURL:        finalURL,
	})
}

func (p *Pooled) dispatch(ctx context.Context, req Request) (*Document, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, &Error{Category: "internal", Message: "failed to marshal request: " + err.Error()}
	}

	respLine, err := p.pool.Dispatch(ctx, payload)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &Error{Category: "timeout", Message: "extractor exceeded deadline"}
		}
		return nil, &Error{Category: "internal", Message: fmt.Sprintf("pool dispatch failed: %v", err)}
	}

	var resp rustResponse
	if err := json.Unmarshal(respLine, &resp); err != nil {
		return nil, &Error{
			Category: "internal",
			Message:  fmt.Sprintf("malformed extractor output: %v, raw:%s", err, truncate(string(respLine), 500)),
		}
	}

	if resp.Status == "error" {
		if resp.Error == nil {
			return nil, &Error{Category: "internal", Message: "extractor reported error status with no error body"}
		}
		return nil, &Error{Category: resp.Error.Category, Message: resp.Error.Message}
	}

	if resp.Document == nil {
		return nil, &Error{Category: "internal", Message: "extractor reported ok status with no document"}
	}

	return resp.Document, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "...(truncated)"
}

const DefaultTimeout = 20 * time.Second
