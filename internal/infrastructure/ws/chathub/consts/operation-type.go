package consts

type OperationType string

const (
	SendMessage OperationType = "send_message"
	Broadcast   OperationType = "broadcast"
	// GetMessages etc
)

func (t OperationType) String() string {
	return string(t)
}
