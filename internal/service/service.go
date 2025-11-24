package service

import (
	"SkillsTracker/internal/repository"

	"github.com/rs/zerolog"
)

type service struct {
	user    *userService
	task    *taskService
	comment *commentService
	project *projectService
	skill   *skillService
}

func NewService(repo repository.RepositoryInterface, jwtSecret []byte, logger zerolog.Logger) ServiceInterface {
	return &service{
		user:    newUserService(repo, jwtSecret, logger),
		task:    newTaskService(repo, logger),
		comment: newCommentService(repo, logger),
		project: newProjectService(repo, logger),
		skill:   newSkillService(repo, logger),
	}
}

func (s *service) User() UserServiceInterface       { return s.user }
func (s *service) Task() TaskServiceInterface        { return s.task }
func (s *service) Comment() CommentServiceInterface { return s.comment }
func (s *service) Project() ProjectServiceInterface  { return s.project }
func (s *service) Skill() SkillServiceInterface      { return s.skill }
