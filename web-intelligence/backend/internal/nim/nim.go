package nim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	endpoint = "https://integrate.api.nvidia.com/v1/chat/completions",
	model	 = "nvidia/nemotron-3-ultra-550b-a55b"
)

type Client struct {
	apiKey	   string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	key := os.Getenv("NVIDIA_NIM_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("NVIDIA_NIM_API_KEY not set")
	}
	return &Client{
		apiKey:    key,
		httpClient: &http.Client{Timeout: 25 * time.Second},
	}, nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role   string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (c *Client) Analyze(ctx context.Context, title, content, question string) (string, error) {
	prompt := buildPrompt(title, content, question)

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
	}
	body, err :=json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal nim request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build nim request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("nim request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("nim auth failure: %s", truncate(string(respBody), 300))
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("nim rate limit: %s", truncate(string(respBody), 300))
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errrorf("nim returned %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}
}