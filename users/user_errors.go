package users

import "errors"

var (
	ErrInvalidPassword = errors.New("user: invalid password provided")
	ErrInvalidEmail    = errors.New("user: invalid email provided")
)
