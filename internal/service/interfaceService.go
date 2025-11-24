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
	GetTasksByProjectID(ctx context.Context, projectID int) ([]*dto.TaskResponse, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]*dto.TaskResponse, error)
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

type ProjectServiceInterface interface {
	CreateProject(ctx context.Context, req *dto.ProjectRequest, managerID int) (*dto.ProjectResponse, error)
	GetProjectByID(ctx context.Context, projectID int) (*dto.ProjectResponse, error)
	GetProjectsByManagerID(ctx context.Context, managerID int) ([]*dto.ProjectResponse, error)
	GetAllProjects(ctx context.Context) ([]*dto.ProjectResponse, error)
	UpdateProject(ctx context.Context, projectID int, req *dto.ProjectRequest, managerID int) error
	DeleteProject(ctx context.Context, projectID int, managerID int) error
	AddProjectMember(ctx context.Context, projectID, userID int, role string, currentUserID int) error
	RemoveProjectMember(ctx context.Context, projectID, userID int) error
}

type SkillServiceInterface interface {
	CreateSkill(ctx context.Context, req *dto.SkillRequest) (*dto.SkillResponse, error)
	GetSkillByID(ctx context.Context, skillID int) (*dto.SkillResponse, error)
	GetAllSkills(ctx context.Context) ([]*dto.SkillResponse, error)
	GetSkillsByCategory(ctx context.Context, category string) ([]*dto.SkillResponse, error)
	UpdateSkill(ctx context.Context, skillID int, req *dto.SkillRequest) error
	DeleteSkill(ctx context.Context, skillID int) error
	AddUserSkill(ctx context.Context, userID int, req *dto.UserSkillRequest) (*dto.UserSkillResponse, error)
	GetUserSkills(ctx context.Context, userID int) ([]*dto.UserSkillResponse, error)
	UpdateUserSkill(ctx context.Context, userID, skillID int, level int) error
	RemoveUserSkill(ctx context.Context, userID, skillID int) error
}

type ServiceInterface interface {
	User() UserServiceInterface
	Task() TaskServiceInterface
	Comment() CommentServiceInterface
	Project() ProjectServiceInterface
	Skill() SkillServiceInterface
}
