package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, &Error{Category: "internal", Message: "failed to marshal request: " + err.Error()}
	}

	cmd := exec.CommandContext(ctx, binaryPath)
	cmd.Stdin = bytes.NewReader(payload)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return nil, &Error{Category: "timeout", Message: "extractor exceeded deadline"}
	}

	if runErr != nil {
		if stdout.Len() == 0 {
			return nil, &Error{
				Category: "internal",
				Message:  fmt.Sprintf("extractor failed: %v, stderr: %s", runErr, stderr.String()),
			}
		}
	}

	var resp rustResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, &Error{
			Category: "internal",
			Message:  fmt.Sprintf("malformed extractor output: %v, raw:%s", err, truncate(stdout.String(), 500)),
		}
	}

	if resp.Status == "error" {
		if resp.Error == nil {
			return nil, &Error{Category: "internal", Message: "extractor reported errorstatus with no error body"}
		}
		return nil, &Error{Category: resp.Error.Category, Message: resp.Error.Message}
	}

	if resp.Document == nil {
		return nil, &Error{Category: "internal", Message: "extractor reported ok status with no document"}
	}

	return resp.Document, nil
}
