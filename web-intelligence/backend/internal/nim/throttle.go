package nim

import (
	"context"
	"fmt"
)

type ThrottledClient struct {
	client *Client
	sem    chan struct{}
}

func NewThrottled(client *Client, maxInFlight int) *ThrottledClient {
	if maxInFlight < 1 {
		maxInFlight = 1
	}
	return &ThrottledClient{
		client: client,
		sem:    make(chan struct{}, maxInFlight),
	}
}

func (t *ThrottledClient) Analyze(ctx context.Context, title, content, question string) (string, error) {
	select {
	case t.sem <- struct{}{}:
		defer func() { <-t.sem }()
	case <-ctx.Done():
		return "", fmt.Errorf("nim: %w waiting for a throttle slot", ctx.Err())
	}

	return t.client.Analyze(ctx, title, content, question)
}
