package errors

import stderrs "errors"

var (
	ErrEmailAlreadyExists = stderrs.New("email already exists")
	ErrInvalidCredentials = stderrs.New("invalid credentials")
	ErrUserNotFound       = stderrs.New("user not found")
	ErrWrongPassword      = stderrs.New("wrong password")
	ErrJwtFailure         = stderrs.New("Failed to generate JWT")
)
