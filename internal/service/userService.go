package service

import (
	dto "SkillsTracker/internal/dto"
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"SkillsTracker/internal/utils/jwt"
	"context"
	"errors"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      repository.RepositoryInterface
	jwtSecret []byte
	logger    zerolog.Logger
}

func newUserService(repo repository.RepositoryInterface, jwtSecret []byte, logger zerolog.Logger) *userService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	logger := s.logger.With().Str("username", req.Username).Logger()
	logger.Info().Msg("login attempt")

	user, err := s.repo.User().GetUserByUsername(ctx, req.Username)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user")
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Warn().Msg("invalid password")
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, string(user.Role), s.jwtSecret)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate token")
		return nil, err
	}

	logger.Info().Msg("login successful")
	return &dto.LoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req *dto.UserRequest) error {
	logger := s.logger.With().Str("username", req.Username).Logger()
	logger.Info().Msg("creating user")

	_, err := s.repo.User().GetUserByUsername(ctx, req.Username)
	if err == nil {
		logger.Warn().Msg("user already exists")
		return errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hash password")
		return err
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Role:         models.Role(req.Role),
		Name:         req.Name,
	}

	if err := s.repo.User().CreateUser(ctx, user); err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return err
	}

	logger.Info().Int("id", user.ID).Msg("user created")
	return nil
}

func (s *userService) GetUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	logger := s.logger.With().Logger()
	logger.Info().Msg("getting all users")

	users, err := s.repo.User().GetUsers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get users")
		return nil, err
	}

	resp := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		resp = append(resp, &dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return resp, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID int) (*dto.UserResponse, error) {
	logger := s.logger.With().Int("id", userID).Logger()
	logger.Info().Msg("getting user")

	user, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user")
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID int, req *dto.UserRequest) error {
	logger := s.logger.With().Int("id", userID).Logger()
	logger.Info().Msg("updating user")

	user, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user")
		return err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.repo.User().UpdateUser(ctx, user); err != nil {
		logger.Error().Err(err).Msg("failed to update user")
		return err
	}

	logger.Info().Int("id", user.ID).Msg("user updated")
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	logger := s.logger.With().Int("id", userID).Logger()
	logger.Info().Msg("deleting user")

	if err := s.repo.User().DeleteUser(ctx, userID); err != nil {
		logger.Error().Err(err).Msg("failed to delete user")
		return err
	}

	logger.Info().Msg("user deleted")
	return nil
}
