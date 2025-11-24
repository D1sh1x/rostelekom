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
	logger := s.logger.With().Int("project_id", req.ProjectID).Int("user_id", userID).Logger()
	logger.Info().Msg("creating task")

	// Проверка существования проекта
	project, err := s.repo.Project().GetProjectByID(ctx, req.ProjectID)
	if err != nil {
		logger.Error().Err(err).Msg("project not found")
		return nil, errors.New("project not found")
	}

	// Проверка прав доступа (только менеджер проекта может создавать задачи)
	user, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Role != models.RoleManager && project.ManagerID != userID {
		return nil, errors.New("only project manager can create tasks")
	}

	// Проверка зависимости (parent task)
	if req.ParentTaskID != nil {
		parentTask, err := s.repo.Task().GetTaskByID(ctx, *req.ParentTaskID)
		if err != nil {
			return nil, errors.New("parent task not found")
		}
		if parentTask.Status != models.StatusCompleted {
			return nil, errors.New("parent task must be completed before creating dependent task")
		}
	}

	deadline, err := time.Parse("2006-01-02T15:04:05Z07:00", req.Deadline)
	if err != nil {
		logger.Error().Str("deadline", req.Deadline).Msg("invalid deadline format")
		return nil, errors.New("invalid deadline format")
	}

	// Валидация приоритета и типа
	priority := models.PriorityMedium
	if req.Priority != "" {
		priority = models.Priority(req.Priority)
		if priority != models.PriorityLow && priority != models.PriorityMedium &&
			priority != models.PriorityHigh && priority != models.PriorityUrgent {
			return nil, errors.New("invalid priority")
		}
	}

	taskType := models.TaskTypeTask
	if req.Type != "" {
		taskType = models.TaskType(req.Type)
		if taskType != models.TaskTypeFeature && taskType != models.TaskTypeBug &&
			taskType != models.TaskTypeTask && taskType != models.TaskTypeEpic {
			return nil, errors.New("invalid task type")
		}
	}

	task := &models.Task{
		ProjectID:    req.ProjectID,
		Title:        req.Title,
		Description:  req.Description,
		Deadline:     deadline,
		Status:       models.StatusPending,
		Progress:     0,
		Hours:        req.Hours,
		Priority:     priority,
		Type:         taskType,
		ParentTaskID: req.ParentTaskID,
	}

	if err := s.repo.Task().CreateTask(ctx, task); err != nil {
		logger.Error().Err(err).Msg("failed to create task")
		return nil, err
	}

	// Добавление навыков к задаче
	skillIDs := []int{}
	if len(req.SkillIDs) > 0 {
		for _, skillID := range req.SkillIDs {
			_, err := s.repo.Skill().GetSkillByID(ctx, skillID)
			if err != nil {
				logger.Warn().Int("skill_id", skillID).Msg("skill not found, skipping")
				continue
			}
			taskSkill := &models.TaskSkill{
				TaskID:  task.ID,
				SkillID: skillID,
			}
			if err := s.repo.TaskSkill().CreateTaskSkill(ctx, taskSkill); err != nil {
				logger.Warn().Err(err).Int("skill_id", skillID).Msg("failed to add skill to task")
			} else {
				skillIDs = append(skillIDs, skillID)
			}
		}
	}

	// Назначение сотрудников на задачу (с проверкой навыков)
	assigneeIDs := []int{}
	if len(req.AssigneeIDs) > 0 {
		requiredSkills, _ := s.repo.TaskSkill().GetTaskSkillsByTaskID(ctx, task.ID)
		requiredSkillMap := make(map[int]bool)
		for _, ts := range requiredSkills {
			requiredSkillMap[ts.SkillID] = true
		}

		for _, assigneeID := range req.AssigneeIDs {
			_, err := s.repo.User().GetUserByID(ctx, assigneeID)
			if err != nil {
				logger.Warn().Int("user_id", assigneeID).Msg("user not found, skipping")
				continue
			}

			// Проверка наличия необходимых навыков
			if len(requiredSkillMap) > 0 {
				userSkills, _ := s.repo.UserSkill().GetUserSkillsByUserID(ctx, assigneeID)
				userSkillMap := make(map[int]bool)
				for _, us := range userSkills {
					userSkillMap[us.SkillID] = true
				}

				hasAllSkills := true
				for skillID := range requiredSkillMap {
					if !userSkillMap[skillID] {
						hasAllSkills = false
						break
					}
				}

				if !hasAllSkills {
					logger.Warn().Int("user_id", assigneeID).Msg("user does not have required skills, skipping")
					continue
				}
			}

			taskAssignee := &models.TaskAssignee{
				TaskID: task.ID,
				UserID: assigneeID,
			}
			if err := s.repo.TaskAssignee().CreateTaskAssignee(ctx, taskAssignee); err != nil {
				logger.Warn().Err(err).Int("user_id", assigneeID).Msg("failed to assign user to task")
			} else {
				assigneeIDs = append(assigneeIDs, assigneeID)
			}
		}
	}

	logger.Info().Int("id", task.ID).Msg("task created")
	return &dto.TaskResponse{
		ID:           task.ID,
		ProjectID:    task.ProjectID,
		Title:        task.Title,
		Description:  task.Description,
		Deadline:     task.Deadline,
		Status:       string(task.Status),
		Progress:     task.Progress,
		Hours:        task.Hours,
		Priority:     string(task.Priority),
		Type:         string(task.Type),
		ParentTaskID: task.ParentTaskID,
		SkillIDs:     skillIDs,
		AssigneeIDs:  assigneeIDs,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
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

	// Получение навыков задачи
	taskSkills, _ := s.repo.TaskSkill().GetTaskSkillsByTaskID(ctx, taskID)
	skillIDs := make([]int, len(taskSkills))
	for i, ts := range taskSkills {
		skillIDs[i] = ts.SkillID
	}

	// Получение назначенных сотрудников
	taskAssignees, _ := s.repo.TaskAssignee().GetTaskAssigneesByTaskID(ctx, taskID)
	assigneeIDs := make([]int, len(taskAssignees))
	for i, ta := range taskAssignees {
		assigneeIDs[i] = ta.UserID
	}

	return &dto.TaskResponse{
		ID:           task.ID,
		ProjectID:    task.ProjectID,
		Title:        task.Title,
		Description:  task.Description,
		Deadline:     task.Deadline,
		Status:       string(task.Status),
		Progress:     task.Progress,
		Hours:        task.Hours,
		Priority:     string(task.Priority),
		Type:         string(task.Type),
		ParentTaskID: task.ParentTaskID,
		SkillIDs:     skillIDs,
		AssigneeIDs:  assigneeIDs,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}, nil
}

func (s *taskService) GetTasksByProjectID(ctx context.Context, projectID int) ([]*dto.TaskResponse, error) {
	logger := s.logger.With().Int("project_id", projectID).Logger()
	logger.Info().Msg("getting tasks by project")

	tasks, err := s.repo.Task().GetTasksByProjectID(ctx, projectID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get tasks")
		return nil, err
	}

	responses := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskSkills, _ := s.repo.TaskSkill().GetTaskSkillsByTaskID(ctx, task.ID)
		skillIDs := make([]int, len(taskSkills))
		for j, ts := range taskSkills {
			skillIDs[j] = ts.SkillID
		}

		taskAssignees, _ := s.repo.TaskAssignee().GetTaskAssigneesByTaskID(ctx, task.ID)
		assigneeIDs := make([]int, len(taskAssignees))
		for j, ta := range taskAssignees {
			assigneeIDs[j] = ta.UserID
		}

		responses[i] = &dto.TaskResponse{
			ID:           task.ID,
			ProjectID:    task.ProjectID,
			Title:        task.Title,
			Description:  task.Description,
			Deadline:     task.Deadline,
			Status:       string(task.Status),
			Progress:     task.Progress,
			Hours:        task.Hours,
			Priority:     string(task.Priority),
			Type:         string(task.Type),
			ParentTaskID: task.ParentTaskID,
			SkillIDs:     skillIDs,
			AssigneeIDs:  assigneeIDs,
			CreatedAt:    task.CreatedAt,
			UpdatedAt:    task.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *taskService) GetTasksByUserID(ctx context.Context, userID int) ([]*dto.TaskResponse, error) {
	logger := s.logger.With().Int("user_id", userID).Logger()
	logger.Info().Msg("getting tasks by user")

	tasks, err := s.repo.Task().GetTasksByUserID(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get tasks")
		return nil, err
	}

	responses := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskSkills, _ := s.repo.TaskSkill().GetTaskSkillsByTaskID(ctx, task.ID)
		skillIDs := make([]int, len(taskSkills))
		for j, ts := range taskSkills {
			skillIDs[j] = ts.SkillID
		}

		taskAssignees, _ := s.repo.TaskAssignee().GetTaskAssigneesByTaskID(ctx, task.ID)
		assigneeIDs := make([]int, len(taskAssignees))
		for j, ta := range taskAssignees {
			assigneeIDs[j] = ta.UserID
		}

		responses[i] = &dto.TaskResponse{
			ID:           task.ID,
			ProjectID:    task.ProjectID,
			Title:        task.Title,
			Description:  task.Description,
			Deadline:     task.Deadline,
			Status:       string(task.Status),
			Progress:     task.Progress,
			Hours:        task.Hours,
			Priority:     string(task.Priority),
			Type:         string(task.Type),
			ParentTaskID: task.ParentTaskID,
			SkillIDs:     skillIDs,
			AssigneeIDs:  assigneeIDs,
			CreatedAt:    task.CreatedAt,
			UpdatedAt:    task.UpdatedAt,
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

	// Проверка прав доступа
	user, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	project, err := s.repo.Project().GetProjectByID(ctx, task.ProjectID)
	if err != nil {
		return err
	}
	if user.Role != models.RoleManager && project.ManagerID != userID {
		return errors.New("only project manager can update tasks")
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
	if req.Hours > 0 {
		task.Hours = req.Hours
	}
	if req.Priority != "" {
		task.Priority = models.Priority(req.Priority)
	}
	if req.Type != "" {
		task.Type = models.TaskType(req.Type)
	}
	if req.ParentTaskID != nil {
		// Проверка зависимости
		if *req.ParentTaskID != 0 {
			parentTask, err := s.repo.Task().GetTaskByID(ctx, *req.ParentTaskID)
			if err != nil {
				return errors.New("parent task not found")
			}
			if parentTask.Status != models.StatusCompleted {
				return errors.New("parent task must be completed")
			}
		}
		task.ParentTaskID = req.ParentTaskID
	}

	// Обновление навыков
	if req.SkillIDs != nil {
		// Удаляем старые навыки
		s.repo.TaskSkill().DeleteTaskSkillsByTaskID(ctx, taskID)
		// Добавляем новые
		for _, skillID := range req.SkillIDs {
			_, err := s.repo.Skill().GetSkillByID(ctx, skillID)
			if err != nil {
				continue
			}
			taskSkill := &models.TaskSkill{
				TaskID:  taskID,
				SkillID: skillID,
			}
			s.repo.TaskSkill().CreateTaskSkill(ctx, taskSkill)
		}
	}

	// Обновление назначенных сотрудников
	if req.AssigneeIDs != nil {
		// Удаляем старых назначенных
		s.repo.TaskAssignee().DeleteTaskAssigneesByTaskID(ctx, taskID)
		// Добавляем новых (с проверкой навыков)
		requiredSkills, _ := s.repo.TaskSkill().GetTaskSkillsByTaskID(ctx, taskID)
		requiredSkillMap := make(map[int]bool)
		for _, ts := range requiredSkills {
			requiredSkillMap[ts.SkillID] = true
		}

		for _, assigneeID := range req.AssigneeIDs {
			if len(requiredSkillMap) > 0 {
				userSkills, _ := s.repo.UserSkill().GetUserSkillsByUserID(ctx, assigneeID)
				userSkillMap := make(map[int]bool)
				for _, us := range userSkills {
					userSkillMap[us.SkillID] = true
				}

				hasAllSkills := true
				for skillID := range requiredSkillMap {
					if !userSkillMap[skillID] {
						hasAllSkills = false
						break
					}
				}

				if !hasAllSkills {
					continue
				}
			}

			taskAssignee := &models.TaskAssignee{
				TaskID: taskID,
				UserID: assigneeID,
			}
			s.repo.TaskAssignee().CreateTaskAssignee(ctx, taskAssignee)
		}
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

	task, err := s.repo.Task().GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	// Проверка прав доступа
	user, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	project, err := s.repo.Project().GetProjectByID(ctx, task.ProjectID)
	if err != nil {
		return err
	}
	if user.Role != models.RoleManager && project.ManagerID != userID {
		return errors.New("only project manager can delete tasks")
	}

	// Проверка зависимостей
	childTasks, _ := s.repo.Task().GetChildTasks(ctx, taskID)
	if len(childTasks) > 0 {
		return errors.New("cannot delete task with dependent tasks")
	}

	if err := s.repo.Task().DeleteTask(ctx, taskID); err != nil {
		logger.Error().Err(err).Msg("failed to delete task")
		return err
	}

	logger.Info().Msg("task deleted")
	return nil
}
