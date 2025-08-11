package messagePipe

import (
	"context"
)

type closer interface {
	Close()
}
type PipeCloser struct {
	pipes []closer
}

func NewPipeCloser(pipes ...closer) *PipeCloser {
	return &PipeCloser{
		pipes: pipes,
	}
}

// Start is a stub for implementing app component interface
func (p *PipeCloser) Start(_ context.Context) error {
	return nil
}

func (p *PipeCloser) Shutdown(_ context.Context) error {
	for i := len(p.pipes) - 1; i >= 0; i-- {
		p.pipes[i].Close()
	}

	return nil
}
