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