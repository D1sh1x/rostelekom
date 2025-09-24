package service

import (
	"SkillsTracker/internal/dto"
	"context"
)

type UserServiceInterface interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	CreateUser(ctx context.Context, req *dto.UserRequest) error
	GetUsers(ctx context.Context) ([]*dto.UserResponse, error)
	GetUserByID(ctx context.Context, userID int) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, userID int, req *dto.UserRequest) error
	DeleteUser(ctx context.Context, userID int) error
}

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, req *dto.TaskRequest, userID int) (*dto.TaskResponse, error)
	GetTaskByID(ctx context.Context, taskID int) (*dto.TaskResponse, error)
	GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]*dto.TaskResponse, error)
	UpdateTask(ctx context.Context, taskID int, req *dto.TaskRequest, userID int) error
	DeleteTask(ctx context.Context, taskID int, userID int) error
}

type CommentServiceInterface interface {
	CreateComment(ctx context.Context, req *dto.CommentRequest, userID int) (*dto.CommentResponse, error)
	GetCommentByID(ctx context.Context, commentID int) (*dto.CommentResponse, error)
	GetCommentsByTaskID(ctx context.Context, taskID int) ([]*dto.CommentResponse, error)
	UpdateComment(ctx context.Context, commentID int, req *dto.CommentRequest, userID int) error
	DeleteComment(ctx context.Context, commentID int, userID int) error
}

type ServiceInterface interface {
	User() UserServiceInterface
	Task() TaskServiceInterface
	Comment() CommentServiceInterface
}
