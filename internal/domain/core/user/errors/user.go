package errors

import "errors"

var (
	ErrUserDoesNotExist = errors.New("user does not exist")
	ErrNotAllUsersExist = errors.New("not all users exist")
)
