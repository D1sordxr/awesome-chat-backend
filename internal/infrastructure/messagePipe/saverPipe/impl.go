package saverPipe

import (
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/infrastructure/messagePipe"
)

type Impl struct {
	Pipe chan vo.StreamMessage
}

func NewImpl() *Impl {
	return &Impl{
		Pipe: make(chan vo.StreamMessage, messagePipe.ChanBuff),
	}
}

func (i *Impl) GetReadChan() <-chan vo.StreamMessage {
	return i.Pipe
}

func (i *Impl) GetWriteChan() chan<- vo.StreamMessage {
	return i.Pipe
}
