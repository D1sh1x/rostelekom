package service

import (
	dto "SkillsTracker/internal/dto"
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

type taskService struct {
	repo   repository.RepositoryInterface
	logger zerolog.Logger
}

func newTaskService(repo repository.RepositoryInterface, logger zerolog.Logger) *taskService {
	return &taskService{
		repo:   repo,
		logger: logger,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req *dto.TaskRequest, userID int) (*dto.TaskResponse, error) {
	logger := s.logger.With().Int("employee_id", req.EmployeeID).Int("user_id", userID).Logger()
	logger.Info().Msg("creating task")

	_, err := s.repo.User().GetUserByID(ctx, req.EmployeeID)
	if err != nil {
		logger.Error().Msg("employee not found")
		return nil, errors.New("employee not found")
	}

	deadline, err := time.Parse("2006-01-02T15:04:05Z07:00", req.Deadline)
	if err != nil {
		logger.Error().Str("deadline", req.Deadline).Msg("invalid deadline format")
		return nil, errors.New("invalid deadline format")
	}

	task := &models.Task{
		EmployeeID:  req.EmployeeID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    deadline,
		Status:      models.StatusPending,
		Progress:    0,
	}

	if err := s.repo.Task().CreateTask(ctx, task); err != nil {
		logger.Error().Err(err).Msg("failed to create task")
		return nil, err
	}

	logger.Info().Int("id", task.ID).Msg("task created")
	return &dto.TaskResponse{
		ID:          task.ID,
		EmployeeID:  task.EmployeeID,
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      string(task.Status),
		Progress:    task.Progress,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, taskID int) (*dto.TaskResponse, error) {
	logger := s.logger.With().Int("id", taskID).Logger()
	logger.Info().Msg("getting task")

	task, err := s.repo.Task().GetTaskByID(ctx, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get task")
		return nil, err
	}

	return &dto.TaskResponse{
		ID:          task.ID,
		EmployeeID:  task.EmployeeID,
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      string(task.Status),
		Progress:    task.Progress,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *taskService) GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]*dto.TaskResponse, error) {
	logger := s.logger.With().Int("employee_id", employeeID).Logger()
	logger.Info().Msg("getting tasks")

	tasks, err := s.repo.Task().GetTasksByEmployeeID(ctx, employeeID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get tasks")
		return nil, err
	}

	responses := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = &dto.TaskResponse{
			ID:          task.ID,
			EmployeeID:  task.EmployeeID,
			Title:       task.Title,
			Description: task.Description,
			Deadline:    task.Deadline,
			Status:      string(task.Status),
			Progress:    task.Progress,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *taskService) UpdateTask(ctx context.Context, taskID int, req *dto.TaskRequest, userID int) error {
	logger := s.logger.With().Int("id", taskID).Int("user_id", userID).Logger()
	logger.Info().Msg("updating task")

	task, err := s.repo.Task().GetTaskByID(ctx, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get task")
		return err
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Deadline != "" {
		deadline, err := time.Parse("2006-01-02T15:04:05Z07:00", req.Deadline)
		if err != nil {
			logger.Error().Str("deadline", req.Deadline).Msg("invalid deadline format")
			return errors.New("invalid deadline format")
		}
		task.Deadline = deadline
	}

	if req.Progress != task.Progress {
		task.Progress = req.Progress
	}

	if req.Status != "" {
		task.Status = models.TaskStatus(req.Status)
	}

	if err := s.repo.Task().UpdateTask(ctx, task); err != nil {
		logger.Error().Err(err).Msg("failed to update task")
		return err
	}

	logger.Info().Int("id", task.ID).Msg("task updated")
	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID int, userID int) error {
	logger := s.logger.With().Int("id", taskID).Int("user_id", userID).Logger()
	logger.Info().Msg("deleting task")

	if err := s.repo.Task().DeleteTask(ctx, taskID); err != nil {
		logger.Error().Err(err).Msg("failed to delete task")
		return err
	}

	logger.Info().Msg("task deleted")
	return nil
}
