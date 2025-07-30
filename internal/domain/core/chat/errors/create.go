package errors

import "errors"

var (
	ErrChatShortName         = errors.New("chat name is too short")
	ErrChatInvalidMembersLen = errors.New("chat requires at least one member")
)
