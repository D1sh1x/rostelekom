package service

import (
	dto "SkillsTracker/internal/dto"
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"context"
	"errors"

	"github.com/rs/zerolog"
)

type skillService struct {
	repo   repository.RepositoryInterface
	logger zerolog.Logger
}

func newSkillService(repo repository.RepositoryInterface, logger zerolog.Logger) *skillService {
	return &skillService{
		repo:   repo,
		logger: logger,
	}
}

func (s *skillService) CreateSkill(ctx context.Context, req *dto.SkillRequest) (*dto.SkillResponse, error) {
	logger := s.logger.With().Str("name", req.Name).Logger()
	logger.Info().Msg("creating skill")

	skill := &models.Skill{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	}

	if err := s.repo.Skill().CreateSkill(ctx, skill); err != nil {
		logger.Error().Err(err).Msg("failed to create skill")
		return nil, err
	}

	logger.Info().Int("id", skill.ID).Msg("skill created")
	return &dto.SkillResponse{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		Category:    skill.Category,
		CreatedAt:   skill.CreatedAt,
		UpdatedAt:  skill.UpdatedAt,
	}, nil
}

func (s *skillService) GetSkillByID(ctx context.Context, skillID int) (*dto.SkillResponse, error) {
	logger := s.logger.With().Int("id", skillID).Logger()
	logger.Info().Msg("getting skill")

	skill, err := s.repo.Skill().GetSkillByID(ctx, skillID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get skill")
		return nil, err
	}

	return &dto.SkillResponse{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		Category:    skill.Category,
		CreatedAt:   skill.CreatedAt,
		UpdatedAt:  skill.UpdatedAt,
	}, nil
}

func (s *skillService) GetAllSkills(ctx context.Context) ([]*dto.SkillResponse, error) {
	logger := s.logger.With().Logger()
	logger.Info().Msg("getting all skills")

	skills, err := s.repo.Skill().GetAllSkills(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get skills")
		return nil, err
	}

	responses := make([]*dto.SkillResponse, len(skills))
	for i, skill := range skills {
		responses[i] = &dto.SkillResponse{
			ID:          skill.ID,
			Name:        skill.Name,
			Description: skill.Description,
			Category:    skill.Category,
			CreatedAt:   skill.CreatedAt,
			UpdatedAt:  skill.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *skillService) GetSkillsByCategory(ctx context.Context, category string) ([]*dto.SkillResponse, error) {
	logger := s.logger.With().Str("category", category).Logger()
	logger.Info().Msg("getting skills by category")

	skills, err := s.repo.Skill().GetSkillsByCategory(ctx, category)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get skills")
		return nil, err
	}

	responses := make([]*dto.SkillResponse, len(skills))
	for i, skill := range skills {
		responses[i] = &dto.SkillResponse{
			ID:          skill.ID,
			Name:        skill.Name,
			Description: skill.Description,
			Category:    skill.Category,
			CreatedAt:   skill.CreatedAt,
			UpdatedAt:  skill.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *skillService) UpdateSkill(ctx context.Context, skillID int, req *dto.SkillRequest) error {
	logger := s.logger.With().Int("id", skillID).Logger()
	logger.Info().Msg("updating skill")

	skill, err := s.repo.Skill().GetSkillByID(ctx, skillID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get skill")
		return err
	}

	if req.Name != "" {
		skill.Name = req.Name
	}
	if req.Description != "" {
		skill.Description = req.Description
	}
	if req.Category != "" {
		skill.Category = req.Category
	}

	if err := s.repo.Skill().UpdateSkill(ctx, skill); err != nil {
		logger.Error().Err(err).Msg("failed to update skill")
		return err
	}

	logger.Info().Int("id", skill.ID).Msg("skill updated")
	return nil
}

func (s *skillService) DeleteSkill(ctx context.Context, skillID int) error {
	logger := s.logger.With().Int("id", skillID).Logger()
	logger.Info().Msg("deleting skill")

	if err := s.repo.Skill().DeleteSkill(ctx, skillID); err != nil {
		logger.Error().Err(err).Msg("failed to delete skill")
		return err
	}

	logger.Info().Msg("skill deleted")
	return nil
}

func (s *skillService) AddUserSkill(ctx context.Context, userID int, req *dto.UserSkillRequest) (*dto.UserSkillResponse, error) {
	logger := s.logger.With().Int("user_id", userID).Int("skill_id", req.SkillID).Logger()
	logger.Info().Msg("adding user skill")

	// Проверка существования пользователя
	_, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Проверка существования навыка
	skill, err := s.repo.Skill().GetSkillByID(ctx, req.SkillID)
	if err != nil {
		return nil, err
	}

	// Проверка уровня (1-5)
	if req.Level < 1 || req.Level > 5 {
		return nil, errors.New("skill level must be between 1 and 5")
	}

	// Проверка, есть ли уже такой навык у пользователя
	existingUserSkill, err := s.repo.UserSkill().GetUserSkill(ctx, userID, req.SkillID)
	if err == nil {
		// Обновляем существующий навык
		existingUserSkill.Level = req.Level
		if err := s.repo.UserSkill().UpdateUserSkill(ctx, existingUserSkill); err != nil {
			logger.Error().Err(err).Msg("failed to update user skill")
			return nil, err
		}
		return &dto.UserSkillResponse{
			ID:        existingUserSkill.ID,
			UserID:    existingUserSkill.UserID,
			SkillID:   existingUserSkill.SkillID,
			SkillName: skill.Name,
			Level:     existingUserSkill.Level,
			CreatedAt: existingUserSkill.CreatedAt,
			UpdatedAt: existingUserSkill.UpdatedAt,
		}, nil
	}

	// Создаем новый навык пользователя
	userSkill := &models.UserSkill{
		UserID:  userID,
		SkillID: req.SkillID,
		Level:   req.Level,
	}

	if err := s.repo.UserSkill().CreateUserSkill(ctx, userSkill); err != nil {
		logger.Error().Err(err).Msg("failed to create user skill")
		return nil, err
	}

	logger.Info().Int("id", userSkill.ID).Msg("user skill added")
	return &dto.UserSkillResponse{
		ID:        userSkill.ID,
		UserID:    userSkill.UserID,
		SkillID:   userSkill.SkillID,
		SkillName: skill.Name,
		Level:     userSkill.Level,
		CreatedAt: userSkill.CreatedAt,
		UpdatedAt: userSkill.UpdatedAt,
	}, nil
}

func (s *skillService) GetUserSkills(ctx context.Context, userID int) ([]*dto.UserSkillResponse, error) {
	logger := s.logger.With().Int("user_id", userID).Logger()
	logger.Info().Msg("getting user skills")

	userSkills, err := s.repo.UserSkill().GetUserSkillsByUserID(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user skills")
		return nil, err
	}

	responses := make([]*dto.UserSkillResponse, len(userSkills))
	for i, us := range userSkills {
		skill, _ := s.repo.Skill().GetSkillByID(ctx, us.SkillID)
		skillName := ""
		if skill != nil {
			skillName = skill.Name
		}

		responses[i] = &dto.UserSkillResponse{
			ID:        us.ID,
			UserID:    us.UserID,
			SkillID:   us.SkillID,
			SkillName: skillName,
			Level:     us.Level,
			CreatedAt: us.CreatedAt,
			UpdatedAt: us.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *skillService) UpdateUserSkill(ctx context.Context, userID, skillID int, level int) error {
	logger := s.logger.With().Int("user_id", userID).Int("skill_id", skillID).Logger()
	logger.Info().Msg("updating user skill")

	if level < 1 || level > 5 {
		return errors.New("skill level must be between 1 and 5")
	}

	userSkill, err := s.repo.UserSkill().GetUserSkill(ctx, userID, skillID)
	if err != nil {
		return err
	}

	userSkill.Level = level
	if err := s.repo.UserSkill().UpdateUserSkill(ctx, userSkill); err != nil {
		logger.Error().Err(err).Msg("failed to update user skill")
		return err
	}

	logger.Info().Msg("user skill updated")
	return nil
}

func (s *skillService) RemoveUserSkill(ctx context.Context, userID, skillID int) error {
	logger := s.logger.With().Int("user_id", userID).Int("skill_id", skillID).Logger()
	logger.Info().Msg("removing user skill")

	if err := s.repo.UserSkill().DeleteUserSkillByUserAndSkill(ctx, userID, skillID); err != nil {
		logger.Error().Err(err).Msg("failed to remove user skill")
		return err
	}

	logger.Info().Msg("user skill removed")
	return nil
}

