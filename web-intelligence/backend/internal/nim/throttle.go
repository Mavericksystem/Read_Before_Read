package nim

import (
	"context"
	"fmt"
	"log"
	"time"
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

func (t *ThrottledClient) Analyze(
	ctx context.Context,
	title, content, question string,
) (string, error) {

	start := time.Now()

	if deadline, ok := ctx.Deadline(); ok {
		log.Printf(
			"nim: stage=before-semaphore deadline=%s remaining=%s ctx_err=%v",
			deadline.Format(time.RFC3339Nano),
			time.Until(deadline).Round(time.Millisecond),
			ctx.Err(),
		)
	}

	select {
	case t.sem <- struct{}{}:

		log.Printf(
			"nim: semaphore-acquired wait=%s",
			time.Since(start).Round(time.Millisecond),
		)

		defer func() {
			<-t.sem
			log.Printf("nim: semaphore-released")
		}()

	case <-ctx.Done():

		log.Printf(
			"nim: semaphore-failed wait=%s err=%v",
			time.Since(start).Round(time.Millisecond),
			ctx.Err(),
		)

		return "", fmt.Errorf(
			"nim: %w waiting for a throttle slot",
			ctx.Err(),
		)
	}

	if deadline, ok := ctx.Deadline(); ok {
		log.Printf(
			"nim: stage=after-semaphore remaining=%s ctx_err=%v",
			time.Until(deadline).Round(time.Millisecond),
			ctx.Err(),
		)
	}

	return t.client.Analyze(ctx, title, content, question)
}
