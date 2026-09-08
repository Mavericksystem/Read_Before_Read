package pool

import (
	"fmt"
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
