package service

import (
	"context"
	"skilltracker/internal/dto"
	"skilltracker/internal/models"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := new(MockRepo)
	mockTaskRepo := new(MockTaskRepo)
	logger := zerolog.Nop()
	s := New(mockRepo, logger, []byte("secret"))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		req := &dto.TaskRequest{
			EmployeeID: 1,
			Title:      "Test Task",
			Deadline:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			Status:     "pending",
		}

		mockRepo.On("Task").Return(mockTaskRepo)
		mockTaskRepo.On("CreateTask", ctx, mock.MatchedBy(func(tk *models.Task) bool {
			return tk.Title == req.Title && tk.EmployeeID == req.EmployeeID
		})).Return(nil)

		res, err := s.Task().CreateTask(ctx, req, 2)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Title, res.Title)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := new(MockRepo)
	mockTaskRepo := new(MockTaskRepo)
	logger := zerolog.Nop()
	s := New(mockRepo, logger, []byte("secret"))
	ctx := context.Background()

	t.Run("success - creator updates", func(t *testing.T) {
		task := &models.Task{ID: 1, CreatorID: 2, EmployeeID: 3}
		req := &dto.TaskRequest{Title: "New Title"}

		mockRepo.On("Task").Return(mockTaskRepo)
		mockTaskRepo.On("GetTaskByID", ctx, 1).Return(task, nil)
		mockTaskRepo.On("UpdateTask", ctx, mock.Anything).Return(nil)

		err := s.Task().UpdateTask(ctx, 1, req, 2)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		task := &models.Task{ID: 1, CreatorID: 2, EmployeeID: 3}
		req := &dto.TaskRequest{Title: "New Title"}

		mockRepo.On("Task").Return(mockTaskRepo)
		mockTaskRepo.On("GetTaskByID", ctx, 1).Return(task, nil)

		err := s.Task().UpdateTask(ctx, 1, req, 4) // Other user
		assert.Error(t, err)
		assert.Equal(t, "forbidden", err.Error())
	})
}
