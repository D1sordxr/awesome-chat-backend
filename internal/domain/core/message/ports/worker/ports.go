package worker

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
