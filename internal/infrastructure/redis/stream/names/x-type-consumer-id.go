package names

type ConsumerID string

func (c ConsumerID) String() string {
	return string(c)
}

const (
	MessageSaverID ConsumerID = "message-saver-0"
)
