package worker

import "awesome-chat/internal/domain/core/message/vo"

type AckReadPipe interface {
	GetReadChan() <-chan string
}

type AckWritePipe interface {
	GetWriteChan() chan<- string
}

type SaverReadPipe interface {
	GetReadChan() <-chan vo.StreamMessage
}

type SaverWritePipe interface {
	GetWriteChan() chan<- vo.StreamMessage
}

type MessagePipe[T any] interface {
	GetReadChan() <-chan T
	GetWriteChan() chan<- T
	Close()
}

type AckPipeTx interface {
	Add()
	Confirm()
	Wait()
}
