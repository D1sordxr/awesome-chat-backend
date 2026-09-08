package errors

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	chatErrors "awesome-chat/internal/domain/core/chat/errors"
	userErrors "awesome-chat/internal/domain/core/user/errors"
	"awesome-chat/internal/pkg/validate"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if s, ok := status.FromError(err); ok {
		return s.Err()
	}

	var validationErr *validate.Error
	if errors.As(err, &validationErr) {
		return validationStatus(validationErr)
	}

	switch {
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timeout")
	case errors.Is(err, userErrors.ErrUserDoesNotExist):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, userErrors.ErrNotAllUsersExist):
		return status.Error(codes.NotFound, "some users do not exist")
	case errors.Is(err, userErrors.ErrInvalidEmailFormat),
		errors.Is(err, userErrors.ErrInvalidEmailLength):
		return status.Error(codes.InvalidArgument, "invalid email")
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, chatErrors.ErrChatAccessDenied):
		return status.Error(codes.PermissionDenied, "chat access denied")
	case errors.Is(err, chatErrors.ErrChatMemberExists):
		return status.Error(codes.AlreadyExists, "user is already a chat member")
	case errors.Is(err, chatErrors.ErrChatShortName),
		errors.Is(err, chatErrors.ErrChatInvalidMembersLen):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

func validationStatus(validationErr *validate.Error) error {
	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0, len(validationErr.Violations))
	for _, violation := range validationErr.Violations {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       violation.Field,
			Description: violation.Message,
			Reason:      violation.Rule,
		})
	}

	invalidArgument := status.New(codes.InvalidArgument, validationErr.Error())

	withDetails, err := invalidArgument.WithDetails(&errdetails.BadRequest{FieldViolations: fieldViolations})
	if err != nil {
		return invalidArgument.Err()
	}

	return withDetails.Err()
}
