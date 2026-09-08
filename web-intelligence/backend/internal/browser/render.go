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

	var html string
	err = chromedp.Run(renderCtx,
		chromedp.Navigate(url),

		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)
	if err != nil {
		if renderCtx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("browser render: timed out after %s", timeout)
		}
		return "", fmt.Errorf("browser render: %w", err)
	}

	return html, nil
}
