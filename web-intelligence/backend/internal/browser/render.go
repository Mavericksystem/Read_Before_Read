package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func (p *Pool) Render(ctx context.Context, url string, timeout time.Duration) (string, error) {
	tabCtx, err := p.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer p.Release(tabCtx)

	renderCtx, cancel := context.WithTimeout(tabCtx, timeout)
	defer cancel()