package storage

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")

	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskNotUpdated = errors.New("task not updated")

	ErrCommentNotFound = errors.New("comment not found")

	ErrProjectNotFound = errors.New("project not found")

	ErrSkillNotFound = errors.New("skill not found")

	ErrUserSkillNotFound = errors.New("user skill not found")

	ErrTaskSkillNotFound = errors.New("task skill not found")

	ErrTaskAssigneeNotFound = errors.New("task assignee not found")

	ErrProjectMemberNotFound = errors.New("project member not found")

	ErrInvalidInput = errors.New("invalid input")
	ErrDatabase     = errors.New("database error")
	ErrConflict     = errors.New("conflict")
)
