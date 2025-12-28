package entity

import "errors"

var (
	ErrEmailExists      = errors.New("email already exists")
	ErrEmailNotFound    = errors.New("email not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrPasswordMismatch = errors.New("password not match")
	ErrNoFieldToUpdate  = errors.New("no field to update")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordHash     = errors.New("failed to hash password")
	ErrIDGeneration     = errors.New("failed to generate ID")
)
