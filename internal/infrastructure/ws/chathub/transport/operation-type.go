package transport

type OperationType string

const (
	SendMessage OperationType = "send_message"
	// GetMessages etc
)

func (t OperationType) String() string {
	return string(t)
}
