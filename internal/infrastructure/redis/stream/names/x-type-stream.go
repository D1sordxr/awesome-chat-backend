package names

type StreamName string

func (s StreamName) String() string {
	return string(s)
}

const (
	SentMessage StreamName = "sent-message"
)
