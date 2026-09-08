package validate

import (
	"errors"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func Proto(msg proto.Message) error {
	err := protovalidate.Validate(msg)
	if err == nil {
		return nil
	}

	var validationErr *protovalidate.ValidationError
	if !errors.As(err, &validationErr) {
		return fmt.Errorf("validate proto: %w", err)
	}

	violations := make([]FieldViolation, 0, len(validationErr.Violations))
	for _, violation := range validationErr.Violations {
		violations = append(violations, FieldViolation{
			Field:   protovalidate.FieldPathString(violation.Proto.GetField()),
			Rule:    violation.Proto.GetRuleId(),
			Message: violation.Proto.GetMessage(),
		})
	}

	return &Error{Violations: violations}
}
