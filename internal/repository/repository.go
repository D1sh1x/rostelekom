package repository

import (
	"SkillsTracker/internal/models"
	"context"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsers(ctx context.Context) ([]*models.User, error)
}

type TaskRepositoryInterface interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type CommentRepositoryInterface interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, id int) (*models.Comment, error)
	GetCommentsByTaskID(ctx context.Context, taskID int) ([]models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, id int) error
}

type RepositoryInterface interface {
	User() UserRepositoryInterface
	Task() TaskRepositoryInterface
	Comment() CommentRepositoryInterface
}
