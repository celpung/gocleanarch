package entity

import "errors"

// User
var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrNoChanges        = errors.New("no changes to update")
	ErrUserIDRequired   = errors.New("user id is required")
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("email already exists")
	ErrEmailNotFound    = errors.New("email not found")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordHash     = errors.New("failed to hash password")
	ErrPasswordMismatch = errors.New("wrong password")
	ErrIDGeneration     = errors.New("failed to generate id")
	ErrRoleRequired     = errors.New("role is required")
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
)
