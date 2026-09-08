package errors

import "errors"

var (
	ErrChatShortName         = errors.New("chat name is too short")
	ErrChatInvalidMembersLen = errors.New("chat requires at least one member")
)

var (
	ErrChatAccessDenied = errors.New("chat access denied")
	ErrChatMemberExists = errors.New("user is already a chat member")
)
