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

	p := &Pool{
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		free:        make(chan context.Context, size),
		cancels:     make([]context.CancelFunc, 0, size),
	}

	for i := 0; i < size; i++ {
		tabCtx, tabCancel := chromedp.NewContext(allocCtx)

		if err := chromedp.Run(tabCtx); err != nil {
			p.Shutdown()
			return nil, fmt.Errorf("browser pool: starting tab %d/%d: %w", i+1, size, err)
		}
		p.cancels = append(p.cancels, tabCancel)
		p.free <- tabCtx
	}

	return p, nil
}
