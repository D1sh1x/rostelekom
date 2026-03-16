package repository

import (
    "context"
    "skilltracker/internal/models"
    "skilltracker/internal/dto"
)

type UserRepository interface {
    CreateUser(ctx context.Context, user *models.User) error
    GetUserByID(ctx context.Context, id int) (*models.User, error)
    GetUserByUsername(ctx context.Context, username string) (*models.User, error)
    UpdateUser(ctx context.Context, user *models.User) error
    DeleteUser(ctx context.Context, id int) error
    GetUsers(ctx context.Context) ([]*models.User, error)
}

type TaskRepository interface {
    CreateTask(ctx context.Context, task *models.Task) error
    GetTaskByID(ctx context.Context, id int) (*models.Task, error)
    GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]models.Task, error)
    UpdateTask(ctx context.Context, task *models.Task) error
    DeleteTask(ctx context.Context, id int) error
    ListTasks(ctx context.Context, filter dto.TaskFilter) ([]models.Task, error)
    CreateHistory(ctx context.Context, h *models.TaskStatusHistory) error
    GetHistoryByTaskID(ctx context.Context, taskID int) ([]models.TaskStatusHistory, error)
}

type CommentRepository interface {
    CreateComment(ctx context.Context, comment *models.Comment) error
    GetCommentByID(ctx context.Context, id int) (*models.Comment, error)
    GetCommentsByTaskID(ctx context.Context, taskID int) ([]models.Comment, error)
    UpdateComment(ctx context.Context, comment *models.Comment) error
    DeleteComment(ctx context.Context, id int) error
}

type FileRepository interface {
	CreateAttachment(ctx context.Context, f *models.FileAttachment) error
}

type Repository interface {
	User() UserRepository
	Task() TaskRepository
	Comment() CommentRepository
	File() FileRepository
}
