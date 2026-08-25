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
}