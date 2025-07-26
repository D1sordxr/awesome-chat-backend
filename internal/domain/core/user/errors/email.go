package errors

import "errors"

var (
	ErrInvalidEmailLength = errors.New("invalid email length")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)
