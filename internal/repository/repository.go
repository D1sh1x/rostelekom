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
	GetTasksByProjectID(ctx context.Context, projectID int) ([]models.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error)
	GetChildTasks(ctx context.Context, parentTaskID int) ([]models.Task, error)
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

type ProjectRepositoryInterface interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByManagerID(ctx context.Context, managerID int) ([]models.Project, error)
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id int) error
}

type SkillRepositoryInterface interface {
	CreateSkill(ctx context.Context, skill *models.Skill) error
	GetSkillByID(ctx context.Context, id int) (*models.Skill, error)
	GetAllSkills(ctx context.Context) ([]models.Skill, error)
	GetSkillsByCategory(ctx context.Context, category string) ([]models.Skill, error)
	UpdateSkill(ctx context.Context, skill *models.Skill) error
	DeleteSkill(ctx context.Context, id int) error
}

type UserSkillRepositoryInterface interface {
	CreateUserSkill(ctx context.Context, userSkill *models.UserSkill) error
	GetUserSkillByID(ctx context.Context, id int) (*models.UserSkill, error)
	GetUserSkillsByUserID(ctx context.Context, userID int) ([]models.UserSkill, error)
	GetUserSkillsBySkillID(ctx context.Context, skillID int) ([]models.UserSkill, error)
	GetUserSkill(ctx context.Context, userID, skillID int) (*models.UserSkill, error)
	UpdateUserSkill(ctx context.Context, userSkill *models.UserSkill) error
	DeleteUserSkill(ctx context.Context, id int) error
	DeleteUserSkillByUserAndSkill(ctx context.Context, userID, skillID int) error
}

type TaskSkillRepositoryInterface interface {
	CreateTaskSkill(ctx context.Context, taskSkill *models.TaskSkill) error
	GetTaskSkillByID(ctx context.Context, id int) (*models.TaskSkill, error)
	GetTaskSkillsByTaskID(ctx context.Context, taskID int) ([]models.TaskSkill, error)
	GetTaskSkillsBySkillID(ctx context.Context, skillID int) ([]models.TaskSkill, error)
	DeleteTaskSkill(ctx context.Context, id int) error
	DeleteTaskSkillByTaskAndSkill(ctx context.Context, taskID, skillID int) error
	DeleteTaskSkillsByTaskID(ctx context.Context, taskID int) error
}

type TaskAssigneeRepositoryInterface interface {
	CreateTaskAssignee(ctx context.Context, taskAssignee *models.TaskAssignee) error
	GetTaskAssigneeByID(ctx context.Context, id int) (*models.TaskAssignee, error)
	GetTaskAssigneesByTaskID(ctx context.Context, taskID int) ([]models.TaskAssignee, error)
	GetTaskAssigneesByUserID(ctx context.Context, userID int) ([]models.TaskAssignee, error)
	DeleteTaskAssignee(ctx context.Context, id int) error
	DeleteTaskAssigneeByTaskAndUser(ctx context.Context, taskID, userID int) error
	DeleteTaskAssigneesByTaskID(ctx context.Context, taskID int) error
}

type ProjectMemberRepositoryInterface interface {
	CreateProjectMember(ctx context.Context, projectMember *models.ProjectMember) error
	GetProjectMemberByID(ctx context.Context, id int) (*models.ProjectMember, error)
	GetProjectMembersByProjectID(ctx context.Context, projectID int) ([]models.ProjectMember, error)
	GetProjectMembersByUserID(ctx context.Context, userID int) ([]models.ProjectMember, error)
	GetProjectMember(ctx context.Context, projectID, userID int) (*models.ProjectMember, error)
	UpdateProjectMember(ctx context.Context, projectMember *models.ProjectMember) error
	DeleteProjectMember(ctx context.Context, id int) error
	DeleteProjectMemberByProjectAndUser(ctx context.Context, projectID, userID int) error
}

type RepositoryInterface interface {
	User() UserRepositoryInterface
	Task() TaskRepositoryInterface
	Comment() CommentRepositoryInterface
	Project() ProjectRepositoryInterface
	Skill() SkillRepositoryInterface
	UserSkill() UserSkillRepositoryInterface
	TaskSkill() TaskSkillRepositoryInterface
	TaskAssignee() TaskAssigneeRepositoryInterface
	ProjectMember() ProjectMemberRepositoryInterface
}
