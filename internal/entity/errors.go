package entity

import "errors"

// User
var (
	ErrInvalidInput     = errors.New("user: invalid input")
	ErrNoChanges        = errors.New("user: no changes to update")
	ErrUserIDRequired   = errors.New("user: user id is required")
	ErrUserNotFound     = errors.New("user: user not found")
	ErrEmailExists      = errors.New("user: email already exists")
	ErrEmailNotFound    = errors.New("user: email not found")
	ErrPasswordRequired = errors.New("user: password is required")
	ErrPasswordHash     = errors.New("user: failed to hash password")
	ErrPasswordMismatch = errors.New("user: wrong password")
	ErrIDGeneration     = errors.New("user: failed to generate id")
	ErrRoleRequired     = errors.New("user: role is required")
	ErrNameRequired     = errors.New("user: name is required")
	ErrEmailRequired    = errors.New("user: email is required")
)
