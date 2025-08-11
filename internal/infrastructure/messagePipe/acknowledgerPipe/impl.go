package acknowledgerPipe

import (
	"awesome-chat/internal/infrastructure/messagePipe"
)

type Impl struct {
	Pipe chan string
}

func NewImpl() *Impl {
	return &Impl{
		Pipe: make(chan string, messagePipe.ChanBuff),
	}
}

func (i *Impl) GetReadChan() <-chan string {
	return i.Pipe
}

func (i *Impl) GetWriteChan() chan<- string {
	return i.Pipe
}

func (i *Impl) Close() {
	close(i.Pipe)
}
