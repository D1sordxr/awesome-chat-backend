package validate

import (
	"errors"
	"strings"
)

var ErrInvalidInput = errors.New("invalid input")

type FieldViolation struct {
	Field   string
	Rule    string
	Message string
}

type Error struct {
	Violations []FieldViolation
}

func (e *Error) Error() string {
	if len(e.Violations) == 0 {
		return ErrInvalidInput.Error()
	}

	var builder strings.Builder
	builder.WriteString("validation failed: ")

	for i, violation := range e.Violations {
		if i > 0 {
			builder.WriteString("; ")
		}

		builder.WriteString(violation.Field)
		builder.WriteByte('(')
		builder.WriteString(violation.Rule)
		builder.WriteByte(')')
	}

	return builder.String()
}

func (e *Error) Unwrap() error {
	return ErrInvalidInput
}
