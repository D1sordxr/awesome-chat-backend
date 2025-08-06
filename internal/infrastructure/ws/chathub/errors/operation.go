package errors

import "errors"

var (
	ErrUnsupportedOp   = errors.New("unsupported operation")
	ErrInvalidOpFormat = errors.New("invalid operation format")
)
