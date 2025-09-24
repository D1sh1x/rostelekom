package service

import (
	"SkillsTracker/internal/repository"

	"github.com/rs/zerolog"
)

type service struct {
	user    *userService
	task    *taskService
	comment *commentService
}

func NewService(repo repository.RepositoryInterface, jwtSecret []byte, logger zerolog.Logger) ServiceInterface {
	return &service{
		user:    newUserService(repo, jwtSecret, logger),
		task:    newTaskService(repo, logger),
		comment: newCommentService(repo, logger),
	}
}

func (s *service) User() UserServiceInterface       { return s.user }
func (s *service) Task() TaskServiceInterface       { return s.task }
func (s *service) Comment() CommentServiceInterface { return s.comment }
