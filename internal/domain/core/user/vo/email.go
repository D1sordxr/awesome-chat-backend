package vo

import (
	"awesome-chat/internal/domain/core/user/errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type Email string

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)

	if len(email) == 0 || utf8.RuneCountInString(email) > 255 {
		return "", errors.ErrInvalidEmailLength
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", errors.ErrInvalidEmailFormat
	}
	return Email(email), nil
}

func (e *Email) String() string {
	return string(*e)
}
