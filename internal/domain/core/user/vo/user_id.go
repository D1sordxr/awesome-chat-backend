package vo

import "github.com/google/uuid"

type UserID uuid.UUID

func (u *UserID) String() string {
	return uuid.UUID(*u).String()
}

func (u *UserID) ToUUID() uuid.UUID {
	return uuid.UUID(*u)
}

type UserIDs []uuid.UUID

func (u *UserIDs) ToUUIDs() uuid.UUIDs {
	return uuid.UUIDs(*u)
}
