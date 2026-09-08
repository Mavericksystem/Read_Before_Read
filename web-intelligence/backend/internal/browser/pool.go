package browser

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

type Pool struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc

	free chan context.Context
	// cancels holds the cancel func for each tab context, keyed by
	// position, so Shutdown can clean all of them up.
	cancels []context.CancelFunc
}

func New(size int) (*Pool, error) {
	if size < 1 {
		return nil, fmt.Errorf("browser pool: size must be >= 1, got %d", size)
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true), // required in many containerized/local setups
		)...,
	)