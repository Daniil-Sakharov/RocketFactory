package vo

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordEmpty             = errors.New("password cannot be empty")
	ErrPasswordTooShort          = errors.New("password must be at least 8 characters")
	ErrPasswordHashInvalidLength = errors.New("invalid bcrypt hash length")
	ErrPasswordInvalid           = errors.New("invalid password")
)

type Password struct {
	hash string
}

func NewPasswordFromPlaintext(plaintext string) (*Password, error) {
	if plaintext == "" {
		return nil, ErrPasswordEmpty
	}

	if len(plaintext) < 8 {
		return nil, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Password{hash: string(hash)}, nil
}

func NewPasswordFromHash(hash string) (*Password, error) {
	if hash == "" {
		return nil, ErrPasswordEmpty
	}

	if len(hash) != 60 {
		return nil, ErrPasswordHashInvalidLength
	}

	return &Password{hash: hash}, nil
}

func (p *Password) Hash() string {
	return p.hash
}

func (p *Password) CompareWith(plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext))
	if err != nil {
		return ErrPasswordInvalid
	}
	return nil
}

func (p *Password) Equals(other *Password) bool {
	if other == nil {
		return false
	}

	return p.hash == other.hash
}
