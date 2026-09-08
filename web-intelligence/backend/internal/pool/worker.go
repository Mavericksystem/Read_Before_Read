package pool

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
)

type worker struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Reader
	dead   bool
}

func newWorker(binaryPath string) (*worker, error) {
	cmd := exec.Command(binaryPath)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("worker: stdin pipe: %w", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("worker: stdout pipe: %w", err)
	}

	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("worker: start: %w", err)
	}

	return &worker{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdinPipe),
		stdout: bufio.NewReader(stdoutPipe),
	}, nil
}

func (w *worker) send(ctx context.Context, reqLine []byte) ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.dead {
		return nil, fmt.Errorf("worker: process already dead")
	}

	if _, err := w.stdin.Write(reqLine); err != nil {
		w.dead = true
		return nil, fmt.Errorf("worker: write job: %w", err)
	}
	if err := w.stdin.WriteByte('\n'); err != nil {
		w.dead = true
		return nil, fmt.Errorf("worker: write newline: %w", err)
	}
	if err := w.stdin.Flush(); err != nil {
		w.dead = true
		return nil, fmt.Errorf("worker: flush: %w", err)
	}

	type result struct {
		line []byte
		err  error
	}
	resultCh := make(chan result, 1)

	go func() {
		line, err := w.stdout.ReadBytes('\n')
		resultCh <- result{line: line, err: err}
	}()

	select {
	case <-ctx.Done():
		w.dead = true // caller must kill() this worker; it may still be blocked reading
		return nil, ctx.Err()
	case res := <-resultCh:
		if res.err != nil {
			w.dead = true
			return nil, fmt.Errorf("worker: read response: %w", res.err)
		}
		return res.line, nil
	}
}
