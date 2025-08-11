package messagePipe

import (
	"context"
	"sync"
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
	var wg sync.WaitGroup

	for _, pipe := range p.pipes {
		wg.Add(1)
		go func(p closer) {
			defer wg.Done()
			pipe.Close()
		}(pipe)
	}

	return nil
}
