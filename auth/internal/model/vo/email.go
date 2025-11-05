package vo

import (
	"errors"
	"regexp"
)

var (
	ErrEmailEmpty         = errors.New("email cannot be empty")
	ErrEmailInvalidFormat = errors.New("invalid email format")
	ErrEmailTooLong       = errors.New("email to long (max 320 characters)")
)

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	if value == "" {
		return nil, ErrEmailEmpty
	}

	if len(value) > 320 {
		return nil, ErrEmailTooLong
	}

	if !isValidEmailFormat(value) {
		return nil, ErrEmailInvalidFormat
	}

	return &Email{value: value}, nil
}

func (e *Email) Value() string {
	return e.value
}

func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

func isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
