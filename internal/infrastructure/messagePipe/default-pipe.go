package messagePipe

type MessagePipe[T any] struct {
	ch chan T
}

func NewMessagePipe[T any]() *MessagePipe[T] {
	return &MessagePipe[T]{
		ch: make(chan T, ChanBuff),
	}
}

func (p *MessagePipe[T]) GetReadChan() <-chan T {
	return p.ch
}

func (p *MessagePipe[T]) GetWriteChan() chan<- T {
	return p.ch
}

func (p *MessagePipe[T]) Close() {
	close(p.ch)
}
