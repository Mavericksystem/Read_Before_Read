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