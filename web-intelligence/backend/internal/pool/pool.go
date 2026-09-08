package pool

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	binaryPath string
	size       int
	free       chan *worker

	mu       sync.Mutex
	shutdown bool
}

func New(binaryPath string, size int) (*Pool, error) {
	if size < 1 {
		return nil, fmt.Errorf("pool: size must be >= 1, got %d", size)
	}

	p := &Pool{
		binaryPath: binaryPath,
		size:       size,
		free:       make(chan *worker, size),
	}

	started := make([]*worker, 0, size)
	for i := 0; i < size; i++ {
		w, err := newWorker(binaryPath)
		if err != nil {
			for _, sw := range started {
				sw.kill()
			}
			return nil, fmt.Errorf("pool: starting worker %d/%d: %w", i+1, size, err)
		}
		started = append(started, w)
		p.free <- w
	}

	return p, nil
}

func (p *Pool) Dispatch(ctx context.Context, reqLine []byte) ([]byte, error) {
	var w *worker
	select {
	case w = <-p.free:
		// got one
	case <-ctx.Done():
		return nil, fmt.Errorf("pool: %w waiting for a free worker", ctx.Err())
	}

	resp, err := w.send(ctx, reqLine)
	if err != nil || w.isDead() {

		go p.replace(w)
		if err != nil {
			return nil, err
		}
	} else {
		p.free <- w
	}

	return resp, nil
}

func (p *Pool) replace(w *worker) {
	w.kill()

	p.mu.Lock()
	down := p.shutdown
	p.mu.Unlock()
	if down {
		return
	}

	nw, err := newWorker(p.binaryPath)
	if err != nil {

		log.Printf("pool: failed to replace dead worker: %v (pool capacity reduced)", err)
		return
	}
	p.free <- nw
}

func (p *Pool) Shutdown() {
	p.mu.Lock()
	p.shutdown = true
	p.mu.Unlock()

	for i := 0; i < p.size; i++ {
		w := <-p.free
		w.kill()
	}
}
