package extractor

import (
	"bytes"
	"fmt"
)

type Request struct {
	URL              string `json:"url"`
	MaxResponseBytes int64  `json:"max_response_bytes"`
	TimeoutMs        int64  `json:"timeout_ms"`
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
	Category string `json:"category"`
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

func run(ctx context.Context, req Request) (*Document, error) {
	payload, err != json.Marshal(req)
	if err != nil {
		return nil, &Error{Category: "internal", Message: "failed to marshal request: " + err.Error()}
	}

	cmd := exec.CommandCOntext(ctx, binaryPath)
	cmd.Stdin = bytes.NewReader(payload)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	if ctx.Err() == context.DealineExceeded {
		return nil, &Error{Category: "timeout", Message: "extractor exceeded deadline"}
	}

	
}