package service

import (
	"context"
	"skilltracker/internal/models"
	"skilltracker/internal/repository"
	"skilltracker/internal/dto"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) User() repository.UserRepository {
	return m.Called().Get(0).(repository.UserRepository)
}

func (m *MockRepo) Task() repository.TaskRepository {
	return m.Called().Get(0).(repository.TaskRepository)
}

func (m *MockRepo) Comment() repository.CommentRepository {
	return m.Called().Get(0).(repository.CommentRepository)
}

func (m *MockRepo) File() repository.FileRepository {
	return m.Called().Get(0).(repository.FileRepository)
}

func (m *MockRepo) Skill() repository.SkillRepository {
	return m.Called().Get(0).(repository.SkillRepository)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, u *models.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, u *models.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepo) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByRefreshToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetEmployeesWithSkills(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) CreateTask(ctx context.Context, t *models.Task) error {
	return m.Called(ctx, t).Error(0)
}

func (m *MockTaskRepo) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepo) GetTasksByEmployeeID(ctx context.Context, empID int) ([]models.Task, error) {
	args := m.Called(ctx, empID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepo) UpdateTask(ctx context.Context, t *models.Task) error {
	return m.Called(ctx, t).Error(0)
}

func (m *MockTaskRepo) DeleteTask(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockTaskRepo) CreateHistory(ctx context.Context, h *models.TaskStatusHistory) error {
	return m.Called(ctx, h).Error(0)
}

func (m *MockTaskRepo) GetHistoryByTaskID(ctx context.Context, taskID int) ([]models.TaskStatusHistory, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]models.TaskStatusHistory), args.Error(1)
}

func (m *MockTaskRepo) ListTasks(ctx context.Context, filter dto.TaskFilter) ([]models.Task, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepo) AddSkillToTask(ctx context.Context, taskID int, skillID int) error {
	return m.Called(ctx, taskID, skillID).Error(0)
}

func (m *MockTaskRepo) RemoveSkillFromTask(ctx context.Context, taskID int, skillID int) error {
	return m.Called(ctx, taskID, skillID).Error(0)
}

func (m *MockTaskRepo) GetTaskSkills(ctx context.Context, taskID int) ([]models.Skill, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Skill), args.Error(1)
}

type MockSkillRepo struct {
	mock.Mock
}

func (m *MockSkillRepo) CreateSkill(ctx context.Context, skill *models.Skill) error {
	return m.Called(ctx, skill).Error(0)
}

func (m *MockSkillRepo) GetSkillByID(ctx context.Context, id int) (*models.Skill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Skill), args.Error(1)
}

func (m *MockSkillRepo) GetSkills(ctx context.Context) ([]models.Skill, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Skill), args.Error(1)
}

func (m *MockSkillRepo) DeleteSkill(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockSkillRepo) AssignSkillToUser(ctx context.Context, userID int, skillID int) error {
	return m.Called(ctx, userID, skillID).Error(0)
}

func (m *MockSkillRepo) RemoveSkillFromUser(ctx context.Context, userID int, skillID int) error {
	return m.Called(ctx, userID, skillID).Error(0)
}

func (m *MockSkillRepo) GetUserSkills(ctx context.Context, userID int) ([]models.Skill, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Skill), args.Error(1)
}

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) CreateComment(ctx context.Context, c *models.Comment) error {
	return m.Called(ctx, c).Error(0)
}

func (m *MockCommentRepo) GetCommentsByTaskID(ctx context.Context, taskID int) ([]models.Comment, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]models.Comment), args.Error(1)
}

func (m *MockCommentRepo) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentRepo) UpdateComment(ctx context.Context, c *models.Comment) error {
	return m.Called(ctx, c).Error(0)
}

func (m *MockCommentRepo) DeleteComment(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

type MockFileRepo struct {
	mock.Mock
}

func (m *MockFileRepo) CreateAttachment(ctx context.Context, f *models.FileAttachment) error {
	return m.Called(ctx, f).Error(0)
}
