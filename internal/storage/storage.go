package storage

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")

	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskNotUpdated = errors.New("task not updated")

	ErrCommentNotFound = errors.New("comment not found")

	ErrInvalidInput = errors.New("invalid input")
	ErrDatabase     = errors.New("database error")
	ErrConflict     = errors.New("conflict")
)
