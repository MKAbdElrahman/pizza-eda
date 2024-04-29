package models

import "errors"

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
