package service

import (
	dto "SkillsTracker/internal/dto"
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"context"
	"errors"

	"github.com/rs/zerolog"
)

type projectService struct {
	repo   repository.RepositoryInterface
	logger zerolog.Logger
}

func newProjectService(repo repository.RepositoryInterface, logger zerolog.Logger) *projectService {
	return &projectService{
		repo:   repo,
		logger: logger,
	}
}

func (s *projectService) CreateProject(ctx context.Context, req *dto.ProjectRequest, managerID int) (*dto.ProjectResponse, error) {
	logger := s.logger.With().Int("manager_id", managerID).Logger()
	logger.Info().Msg("creating project")

	// Проверка, что пользователь - менеджер
	user, err := s.repo.User().GetUserByID(ctx, managerID)
	if err != nil {
		return nil, err
	}
	if user.Role != models.RoleManager {
		return nil, errors.New("only managers can create projects")
	}

	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		ManagerID:   managerID,
		Status:      status,
	}

	if err := s.repo.Project().CreateProject(ctx, project); err != nil {
		logger.Error().Err(err).Msg("failed to create project")
		return nil, err
	}

	// Добавление участников проекта
	memberIDs := []int{managerID} // Менеджер автоматически добавляется
	if len(req.MemberIDs) > 0 {
		for _, memberID := range req.MemberIDs {
			if memberID == managerID {
				continue // Менеджер уже добавлен
			}
			_, err := s.repo.User().GetUserByID(ctx, memberID)
			if err != nil {
				logger.Warn().Int("user_id", memberID).Msg("user not found, skipping")
				continue
			}

			projectMember := &models.ProjectMember{
				ProjectID: project.ID,
				UserID:    memberID,
				Role:      "developer", // По умолчанию
			}
			if err := s.repo.ProjectMember().CreateProjectMember(ctx, projectMember); err != nil {
				logger.Warn().Err(err).Int("user_id", memberID).Msg("failed to add member to project")
			} else {
				memberIDs = append(memberIDs, memberID)
			}
		}
	}

	// Добавление менеджера как участника проекта
	managerMember := &models.ProjectMember{
		ProjectID: project.ID,
		UserID:    managerID,
		Role:      "project_manager",
	}
	if err := s.repo.ProjectMember().CreateProjectMember(ctx, managerMember); err != nil {
		logger.Warn().Err(err).Msg("failed to add manager as project member")
	} else {
		// Убедимся, что менеджер в списке
		found := false
		for _, id := range memberIDs {
			if id == managerID {
				found = true
				break
			}
		}
		if !found {
			memberIDs = append(memberIDs, managerID)
		}
	}

	logger.Info().Int("id", project.ID).Msg("project created")
	return &dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		ManagerID:   project.ManagerID,
		Status:      project.Status,
		MemberIDs:   memberIDs,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (s *projectService) GetProjectByID(ctx context.Context, projectID int) (*dto.ProjectResponse, error) {
	logger := s.logger.With().Int("id", projectID).Logger()
	logger.Info().Msg("getting project")

	project, err := s.repo.Project().GetProjectByID(ctx, projectID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get project")
		return nil, err
	}

	// Получение участников проекта
	projectMembers, _ := s.repo.ProjectMember().GetProjectMembersByProjectID(ctx, projectID)
	memberIDs := make([]int, len(projectMembers))
	for i, pm := range projectMembers {
		memberIDs[i] = pm.UserID
	}

	return &dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		ManagerID:   project.ManagerID,
		Status:      project.Status,
		MemberIDs:   memberIDs,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (s *projectService) GetProjectsByManagerID(ctx context.Context, managerID int) ([]*dto.ProjectResponse, error) {
	logger := s.logger.With().Int("manager_id", managerID).Logger()
	logger.Info().Msg("getting projects by manager")

	projects, err := s.repo.Project().GetProjectsByManagerID(ctx, managerID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get projects")
		return nil, err
	}

	responses := make([]*dto.ProjectResponse, len(projects))
	for i, project := range projects {
		projectMembers, _ := s.repo.ProjectMember().GetProjectMembersByProjectID(ctx, project.ID)
		memberIDs := make([]int, len(projectMembers))
		for j, pm := range projectMembers {
			memberIDs[j] = pm.UserID
		}

		responses[i] = &dto.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			ManagerID:   project.ManagerID,
			Status:      project.Status,
			MemberIDs:   memberIDs,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *projectService) GetAllProjects(ctx context.Context) ([]*dto.ProjectResponse, error) {
	logger := s.logger.With().Logger()
	logger.Info().Msg("getting all projects")

	projects, err := s.repo.Project().GetAllProjects(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get projects")
		return nil, err
	}

	responses := make([]*dto.ProjectResponse, len(projects))
	for i, project := range projects {
		projectMembers, _ := s.repo.ProjectMember().GetProjectMembersByProjectID(ctx, project.ID)
		memberIDs := make([]int, len(projectMembers))
		for j, pm := range projectMembers {
			memberIDs[j] = pm.UserID
		}

		responses[i] = &dto.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			ManagerID:   project.ManagerID,
			Status:      project.Status,
			MemberIDs:   memberIDs,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *projectService) UpdateProject(ctx context.Context, projectID int, req *dto.ProjectRequest, managerID int) error {
	logger := s.logger.With().Int("id", projectID).Int("manager_id", managerID).Logger()
	logger.Info().Msg("updating project")

	project, err := s.repo.Project().GetProjectByID(ctx, projectID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get project")
		return err
	}

	// Проверка прав доступа
	if project.ManagerID != managerID {
		manager, err := s.repo.User().GetUserByID(ctx, managerID)
		if err != nil {
			return err
		}
		if manager.Role != models.RoleManager {
			return errors.New("only project manager can update project")
		}
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.repo.Project().UpdateProject(ctx, project); err != nil {
		logger.Error().Err(err).Msg("failed to update project")
		return err
	}

	// Обновление участников проекта
	if req.MemberIDs != nil {
		// Удаляем всех участников кроме менеджера
		existingMembers, _ := s.repo.ProjectMember().GetProjectMembersByProjectID(ctx, projectID)
		for _, member := range existingMembers {
			if member.UserID != project.ManagerID {
				s.repo.ProjectMember().DeleteProjectMemberByProjectAndUser(ctx, projectID, member.UserID)
			}
		}

		// Добавляем новых участников
		for _, memberID := range req.MemberIDs {
			if memberID == project.ManagerID {
				continue
			}
			_, err := s.repo.User().GetUserByID(ctx, memberID)
			if err != nil {
				continue
			}

			projectMember := &models.ProjectMember{
				ProjectID: projectID,
				UserID:    memberID,
				Role:      "developer",
			}
			s.repo.ProjectMember().CreateProjectMember(ctx, projectMember)
		}
	}

	logger.Info().Int("id", project.ID).Msg("project updated")
	return nil
}

func (s *projectService) DeleteProject(ctx context.Context, projectID int, managerID int) error {
	logger := s.logger.With().Int("id", projectID).Int("manager_id", managerID).Logger()
	logger.Info().Msg("deleting project")

	project, err := s.repo.Project().GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	// Проверка прав доступа
	if project.ManagerID != managerID {
		user, err := s.repo.User().GetUserByID(ctx, managerID)
		if err != nil {
			return err
		}
		if user.Role != models.RoleManager {
			return errors.New("only project manager can delete project")
		}
	}

	if err := s.repo.Project().DeleteProject(ctx, projectID); err != nil {
		logger.Error().Err(err).Msg("failed to delete project")
		return err
	}

	logger.Info().Msg("project deleted")
	return nil
}

func (s *projectService) AddProjectMember(ctx context.Context, projectID, userID int, role string, currentUserID int) error {
	logger := s.logger.With().Int("project_id", projectID).Int("user_id", userID).Logger()
	logger.Info().Msg("adding project member")

	project, err := s.repo.Project().GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	// Проверка прав доступа - только менеджер проекта может добавлять участников
	if project.ManagerID != currentUserID {
		currentUser, err := s.repo.User().GetUserByID(ctx, currentUserID)
		if err != nil {
			return err
		}
		if currentUser.Role != models.RoleManager {
			return errors.New("only project manager can add members")
		}
	}

	// Проверка, не является ли пользователь уже участником
	_, err = s.repo.ProjectMember().GetProjectMember(ctx, projectID, userID)
	if err == nil {
		return errors.New("user is already a project member")
	}

	if role == "" {
		role = "developer"
	}

	projectMember := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}

	if err := s.repo.ProjectMember().CreateProjectMember(ctx, projectMember); err != nil {
		logger.Error().Err(err).Msg("failed to add project member")
		return err
	}

	logger.Info().Msg("project member added")
	return nil
}

func (s *projectService) RemoveProjectMember(ctx context.Context, projectID, userID int) error {
	logger := s.logger.With().Int("project_id", projectID).Int("user_id", userID).Logger()
	logger.Info().Msg("removing project member")

	project, err := s.repo.Project().GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	// Нельзя удалить менеджера проекта
	if project.ManagerID == userID {
		return errors.New("cannot remove project manager")
	}

	if err := s.repo.ProjectMember().DeleteProjectMemberByProjectAndUser(ctx, projectID, userID); err != nil {
		logger.Error().Err(err).Msg("failed to remove project member")
		return err
	}

	logger.Info().Msg("project member removed")
	return nil
}

