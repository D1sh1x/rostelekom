package service

import (
	dto "SkillsTracker/internal/dto"
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository - мок репозитория для тестирования
type MockRepository struct {
	mock.Mock
}

type MockUserRepository struct {
	mock.Mock
}

type MockTaskRepository struct {
	mock.Mock
}

type MockCommentRepository struct {
	mock.Mock
}

// Реализация интерфейсов для моков
func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]models.Task, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentsByTaskID(ctx context.Context, taskID int) ([]models.Comment, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Comment), args.Error(1)
}

func (m *MockCommentRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) DeleteComment(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) User() repository.UserRepositoryInterface {
	args := m.Called()
	return args.Get(0).(repository.UserRepositoryInterface)
}

func (m *MockRepository) Task() repository.TaskRepositoryInterface {
	args := m.Called()
	return args.Get(0).(repository.TaskRepositoryInterface)
}

func (m *MockRepository) Comment() repository.CommentRepositoryInterface {
	args := m.Called()
	return args.Get(0).(repository.CommentRepositoryInterface)
}

// Вспомогательная функция для создания тестового сервиса
func setupTaskService() (*taskService, *MockRepository, *MockUserRepository, *MockTaskRepository) {
	mockRepo := new(MockRepository)
	mockUserRepo := new(MockUserRepository)
	mockTaskRepo := new(MockTaskRepository)

	mockRepo.On("User").Return(mockUserRepo)
	mockRepo.On("Task").Return(mockTaskRepo)

	logger := zerolog.Nop()
	service := newTaskService(mockRepo, logger)

	return service, mockRepo, mockUserRepo, mockTaskRepo
}

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.TaskRequest
		userID        int
		setupMocks    func(*MockUserRepository, *MockTaskRepository)
		expectedError string
		validate      func(*testing.T, *dto.TaskResponse, error)
	}{
		{
			name: "успешное создание задачи",
			req: &dto.TaskRequest{
				EmployeeID:  1,
				Title:       "Test Task",
				Description: "Test Description",
				Deadline:    "2024-12-31T23:59:59Z",
			},
			userID: 1,
			setupMocks: func(userRepo *MockUserRepository, taskRepo *MockTaskRepository) {
				userRepo.On("GetUserByID", mock.Anything, 1).Return(&models.User{
					ID:   1,
					Name: "Test Employee",
				}, nil)

				taskRepo.On("CreateTask", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task.EmployeeID == 1 &&
						task.Title == "Test Task" &&
						task.Description == "Test Description" &&
						task.Status == models.StatusPending &&
						task.Progress == 0
				})).Run(func(args mock.Arguments) {
					task := args.Get(1).(*models.Task)
					task.ID = 1
					task.CreatedAt = time.Now()
					task.UpdatedAt = time.Now()
				}).Return(nil)
			},
			validate: func(t *testing.T, resp *dto.TaskResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, 1, resp.ID)
				assert.Equal(t, 1, resp.EmployeeID)
				assert.Equal(t, "Test Task", resp.Title)
				assert.Equal(t, "Test Description", resp.Description)
				assert.Equal(t, "pending", resp.Status)
				assert.Equal(t, 0, resp.Progress)
			},
		},
		{
			name: "сотрудник не найден",
			req: &dto.TaskRequest{
				EmployeeID:  999,
				Title:       "Test Task",
				Description: "Test Description",
				Deadline:    "2024-12-31T23:59:59Z",
			},
			userID: 1,
			setupMocks: func(userRepo *MockUserRepository, taskRepo *MockTaskRepository) {
				userRepo.On("GetUserByID", mock.Anything, 999).Return(nil, errors.New("user not found"))
			},
			expectedError: "employee not found",
			validate: func(t *testing.T, resp *dto.TaskResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "employee not found")
			},
		},
		{
			name: "неверный формат deadline",
			req: &dto.TaskRequest{
				EmployeeID:  1,
				Title:       "Test Task",
				Description: "Test Description",
				Deadline:    "invalid-date",
			},
			userID: 1,
			setupMocks: func(userRepo *MockUserRepository, taskRepo *MockTaskRepository) {
				userRepo.On("GetUserByID", mock.Anything, 1).Return(&models.User{
					ID:   1,
					Name: "Test Employee",
				}, nil)
			},
			expectedError: "invalid deadline format",
			validate: func(t *testing.T, resp *dto.TaskResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "invalid deadline format")
			},
		},
		{
			name: "ошибка при создании задачи в репозитории",
			req: &dto.TaskRequest{
				EmployeeID:  1,
				Title:       "Test Task",
				Description: "Test Description",
				Deadline:    "2024-12-31T23:59:59Z",
			},
			userID: 1,
			setupMocks: func(userRepo *MockUserRepository, taskRepo *MockTaskRepository) {
				userRepo.On("GetUserByID", mock.Anything, 1).Return(&models.User{
					ID:   1,
					Name: "Test Employee",
				}, nil)

				taskRepo.On("CreateTask", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			validate: func(t *testing.T, resp *dto.TaskResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, mockUserRepo, mockTaskRepo := setupTaskService()
			tt.setupMocks(mockUserRepo, mockTaskRepo)

			ctx := context.Background()
			resp, err := service.CreateTask(ctx, tt.req, tt.userID)

			tt.validate(t, resp, err)
			mockUserRepo.AssertExpectations(t)
			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	tests := []struct {
		name       string
		taskID     int
		setupMocks func(*MockTaskRepository)
		validate   func(*testing.T, *dto.TaskResponse, error)
	}{
		{
			name:   "успешное получение задачи",
			taskID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				createdAt := time.Now()
				updatedAt := time.Now()

				taskRepo.On("GetTaskByID", mock.Anything, 1).Return(&models.Task{
					ID:          1,
					EmployeeID:  1,
					Title:       "Test Task",
					Description: "Test Description",
					Deadline:    deadline,
					Status:      models.StatusPending,
					Progress:    50,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}, nil)
			},
			validate: func(t *testing.T, resp *dto.TaskResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, 1, resp.ID)
				assert.Equal(t, 1, resp.EmployeeID)
				assert.Equal(t, "Test Task", resp.Title)
				assert.Equal(t, "Test Description", resp.Description)
				assert.Equal(t, "pending", resp.Status)
				assert.Equal(t, 50, resp.Progress)
			},
		},
		{
			name:   "задача не найдена",
			taskID: 999,
			setupMocks: func(taskRepo *MockTaskRepository) {
				taskRepo.On("GetTaskByID", mock.Anything, 999).Return(nil, errors.New("task not found"))
			},
			validate: func(t *testing.T, resp *dto.TaskResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "task not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, mockTaskRepo := setupTaskService()
			tt.setupMocks(mockTaskRepo)

			ctx := context.Background()
			resp, err := service.GetTaskByID(ctx, tt.taskID)

			tt.validate(t, resp, err)
			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetTasksByEmployeeID(t *testing.T) {
	tests := []struct {
		name       string
		employeeID int
		setupMocks func(*MockTaskRepository)
		validate   func(*testing.T, []*dto.TaskResponse, error)
	}{
		{
			name:       "успешное получение списка задач",
			employeeID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				createdAt := time.Now()
				updatedAt := time.Now()

				taskRepo.On("GetTasksByEmployeeID", mock.Anything, 1).Return([]models.Task{
					{
						ID:          1,
						EmployeeID:  1,
						Title:       "Task 1",
						Description: "Description 1",
						Deadline:    deadline,
						Status:      models.StatusPending,
						Progress:    0,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
					{
						ID:          2,
						EmployeeID:  1,
						Title:       "Task 2",
						Description: "Description 2",
						Deadline:    deadline,
						Status:      models.StatusInProgress,
						Progress:    50,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
				}, nil)
			},
			validate: func(t *testing.T, resp []*dto.TaskResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp, 2)
				assert.Equal(t, 1, resp[0].ID)
				assert.Equal(t, 2, resp[1].ID)
				assert.Equal(t, "pending", resp[0].Status)
				assert.Equal(t, "in_progress", resp[1].Status)
			},
		},
		{
			name:       "пустой список задач",
			employeeID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				taskRepo.On("GetTasksByEmployeeID", mock.Anything, 1).Return([]models.Task{}, nil)
			},
			validate: func(t *testing.T, resp []*dto.TaskResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp, 0)
			},
		},
		{
			name:       "ошибка при получении задач",
			employeeID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				taskRepo.On("GetTasksByEmployeeID", mock.Anything, 1).Return(nil, errors.New("database error"))
			},
			validate: func(t *testing.T, resp []*dto.TaskResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, mockTaskRepo := setupTaskService()
			tt.setupMocks(mockTaskRepo)

			ctx := context.Background()
			resp, err := service.GetTasksByEmployeeID(ctx, tt.employeeID)

			tt.validate(t, resp, err)
			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	tests := []struct {
		name       string
		taskID     int
		req        *dto.TaskRequest
		userID     int
		setupMocks func(*MockTaskRepository)
		validate   func(*testing.T, error)
	}{
		{
			name:   "успешное обновление задачи",
			taskID: 1,
			req: &dto.TaskRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
				Deadline:    "2025-01-31T23:59:59Z",
				Progress:    75,
				Status:      "in_progress",
			},
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				createdAt := time.Now()
				updatedAt := time.Now()

				taskRepo.On("GetTaskByID", mock.Anything, 1).Return(&models.Task{
					ID:          1,
					EmployeeID:  1,
					Title:       "Old Title",
					Description: "Old Description",
					Deadline:    deadline,
					Status:      models.StatusPending,
					Progress:    0,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}, nil)

				taskRepo.On("UpdateTask", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task.ID == 1 &&
						task.Title == "Updated Title" &&
						task.Description == "Updated Description" &&
						task.Status == models.StatusInProgress &&
						task.Progress == 75
				})).Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "частичное обновление (только title)",
			taskID: 1,
			req: &dto.TaskRequest{
				Title: "Updated Title",
			},
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				createdAt := time.Now()
				updatedAt := time.Now()

				taskRepo.On("GetTaskByID", mock.Anything, 1).Return(&models.Task{
					ID:          1,
					EmployeeID:  1,
					Title:       "Old Title",
					Description: "Old Description",
					Deadline:    deadline,
					Status:      models.StatusPending,
					Progress:    0,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}, nil)

				taskRepo.On("UpdateTask", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task.Title == "Updated Title" &&
						task.Description == "Old Description" // не изменилось
				})).Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "задача не найдена",
			taskID: 999,
			req: &dto.TaskRequest{
				Title: "Updated Title",
			},
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				taskRepo.On("GetTaskByID", mock.Anything, 999).Return(nil, errors.New("task not found"))
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "task not found")
			},
		},
		{
			name:   "неверный формат deadline",
			taskID: 1,
			req: &dto.TaskRequest{
				Deadline: "invalid-date",
			},
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				createdAt := time.Now()
				updatedAt := time.Now()

				taskRepo.On("GetTaskByID", mock.Anything, 1).Return(&models.Task{
					ID:          1,
					EmployeeID:  1,
					Title:       "Test Task",
					Description: "Test Description",
					Deadline:    deadline,
					Status:      models.StatusPending,
					Progress:    0,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}, nil)
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid deadline format")
			},
		},
		{
			name:   "ошибка при обновлении в репозитории",
			taskID: 1,
			req: &dto.TaskRequest{
				Title: "Updated Title",
			},
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				createdAt := time.Now()
				updatedAt := time.Now()

				taskRepo.On("GetTaskByID", mock.Anything, 1).Return(&models.Task{
					ID:          1,
					EmployeeID:  1,
					Title:       "Old Title",
					Description: "Old Description",
					Deadline:    deadline,
					Status:      models.StatusPending,
					Progress:    0,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				}, nil)

				taskRepo.On("UpdateTask", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, mockTaskRepo := setupTaskService()
			tt.setupMocks(mockTaskRepo)

			ctx := context.Background()
			err := service.UpdateTask(ctx, tt.taskID, tt.req, tt.userID)

			tt.validate(t, err)
			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		taskID     int
		userID     int
		setupMocks func(*MockTaskRepository)
		validate   func(*testing.T, error)
	}{
		{
			name:   "успешное удаление задачи",
			taskID: 1,
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				taskRepo.On("DeleteTask", mock.Anything, 1).Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "ошибка при удалении задачи",
			taskID: 999,
			userID: 1,
			setupMocks: func(taskRepo *MockTaskRepository) {
				taskRepo.On("DeleteTask", mock.Anything, 999).Return(errors.New("task not found"))
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "task not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, mockTaskRepo := setupTaskService()
			tt.setupMocks(mockTaskRepo)

			ctx := context.Background()
			err := service.DeleteTask(ctx, tt.taskID, tt.userID)

			tt.validate(t, err)
			mockTaskRepo.AssertExpectations(t)
		})
	}
}

