package vo

import "github.com/google/uuid"

type (
	ChatID  uuid.UUID
	ChatIDs uuid.UUIDs
)

func (c *ChatID) ToUUID() uuid.UUID {
	return uuid.UUID(*c)
}

func (c *ChatID) String() string {
	return uuid.UUID(*c).String()
}

func (c *ChatIDs) ToUUIDs() uuid.UUIDs {
	return uuid.UUIDs(*c)
}
