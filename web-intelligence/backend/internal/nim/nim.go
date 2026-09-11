package nim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

const (
	endpoint = "https://integrate.api.nvidia.com/v1/chat/completions"
	model    = "nvidia/nemotron-3-ultra-550b-a55b"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

var nimLimiter = rate.NewLimiter(rate.Every(time.Minute/40), 1)

func NewClient() (*Client, error) {
	key := os.Getenv("NVIDIA_NIM_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("NVIDIA_NIM_API_KEY not set")
	}

	return &Client{
		apiKey: key,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func logDeadline(stage string, ctx context.Context) {
	deadline, ok := ctx.Deadline()
	if !ok {
		log.Printf("nim: stage=%s no deadline", stage)
		return
	}

	remaining := time.Until(deadline)

	log.Printf(
		"nim: stage=%s deadline=%s remaining=%s ctx_err=%v",
		stage,
		deadline.Format(time.RFC3339Nano),
		remaining.Round(time.Millisecond),
		ctx.Err(),
	)
}

func (c *Client) Analyze(
	ctx context.Context,
	title, content, question string,
) (string, error) {

	start := time.Now()

	logDeadline("enter-client-analyze", ctx)

	// ------------------------------------------------------------
	// Rate limiter
	// ------------------------------------------------------------

	limitStart := time.Now()

	logDeadline("before-rate-limit", ctx)

	if err := nimLimiter.Wait(ctx); err != nil {
		logDeadline("rate-limit-failed", ctx)

		return "", fmt.Errorf(
			"nim rate limiter wait failed after %s: %w",
			time.Since(limitStart).Round(time.Millisecond),
			err,
		)
	}

	log.Printf(
		"nim: rate-limit-acquired wait=%s",
		time.Since(limitStart).Round(time.Millisecond),
	)

	logDeadline("after-rate-limit", ctx)

	// ------------------------------------------------------------
	// Build request
	// ------------------------------------------------------------

	prompt := buildPrompt(title, content, question)

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal nim request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("build nim request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// ------------------------------------------------------------
	// Actual HTTP request
	// ------------------------------------------------------------

	logDeadline("before-http-do", ctx)

	httpStart := time.Now()

	resp, err := c.httpClient.Do(httpReq)

	httpDuration := time.Since(httpStart)

	if err != nil {
		logDeadline("http-do-failed", ctx)

		log.Printf(
			"nim: HTTP FAILED duration=%s total=%s err=%v",
			httpDuration.Round(time.Millisecond),
			time.Since(start).Round(time.Millisecond),
			err,
		)

		return "", fmt.Errorf(
			"nim request failed: %w",
			err,
		)
	}

	log.Printf(
		"nim: HTTP SUCCESS duration=%s status=%d total=%s",
		httpDuration.Round(time.Millisecond),
		resp.StatusCode,
		time.Since(start).Round(time.Millisecond),
	)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read nim response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf(
			"nim auth failure: %s",
			truncate(string(respBody), 300),
		)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf(
			"nim rate limit: %s",
			truncate(string(respBody), 300),
		)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"nim returned %d: %s",
			resp.StatusCode,
			truncate(string(respBody), 300),
		)
	}

	var parsed chatResponse

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf(
			"malformed nim response: %w, raw: %s",
			err,
			truncate(string(respBody), 300),
		)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("nim response had no choices")
	}

	return parsed.Choices[0].Message.Content, nil
}

func buildPrompt(title, content, question string) string {
	if question != "" {
		return fmt.Sprintf(
			"Page title: %s\n\nPage content:\n%s\n\nQuestion: %s\n\nAnswer the question based only on the content above.",
			title,
			content,
			question,
		)
	}

	return fmt.Sprintf(
		"Page title: %s\n\nPage content:\n%s\n\nSummarize this page in a few sentences.",
		title,
		content,
	)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n] + "...truncated"
}
