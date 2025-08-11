package names

type GroupName string

func (g GroupName) String() string {
	return string(g)
}

const (
	MessagesForSave GroupName = "messages-for-save"
)
