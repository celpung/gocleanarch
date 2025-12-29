package errors

import stderrs "errors"

// User domain errors.
var (
	ErrInvalidInput     = stderrs.New("user: invalid input")
	ErrNoChanges        = stderrs.New("user: no changes to update")
	ErrUserIDRequired   = stderrs.New("user: user id is required")
	ErrUserNotFound     = stderrs.New("user: user not found")
	ErrEmailExists      = stderrs.New("user: email already exists")
	ErrEmailNotFound    = stderrs.New("user: email not found")
	ErrInvalidEmail     = stderrs.New("user: invalid email")
	ErrPasswordRequired = stderrs.New("user: password is required")
	ErrPasswordHash     = stderrs.New("user: failed to hash password")
	ErrPasswordMismatch = stderrs.New("user: wrong password")
	ErrIDGeneration     = stderrs.New("user: failed to generate id")
	ErrRoleRequired     = stderrs.New("user: role is required")
	ErrNameRequired     = stderrs.New("user: name is required")
	ErrEmailRequired    = stderrs.New("user: email is required")
)
