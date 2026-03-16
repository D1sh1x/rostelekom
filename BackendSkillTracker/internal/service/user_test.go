package service

import (
	"context"
	"skilltracker/internal/dto"
	"skilltracker/internal/models"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Login(t *testing.T) {
	jwtSecret := []byte("secret")
	logger := zerolog.Nop()

	ctx := context.Background()
	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockUserRepo := new(MockUserRepo)
		s := New(mockRepo, logger, jwtSecret)

		user := &models.User{
			ID:           1,
			Username:     username,
			PasswordHash: string(hashedPassword),
			Role:         models.RoleEmployee,
		}

		mockRepo.On("User").Return(mockUserRepo)
		mockUserRepo.On("GetUserByUsername", ctx, username).Return(user, nil)
		mockUserRepo.On("UpdateUser", ctx, mock.Anything).Return(nil)

		res, err := s.User().Login(ctx, &dto.LoginRequest{
			Username: username,
			Password: password,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockUserRepo := new(MockUserRepo)
		s := New(mockRepo, logger, jwtSecret)

		mockRepo.On("User").Return(mockUserRepo)
		mockUserRepo.On("GetUserByUsername", ctx, username).Return(nil, assert.AnError)

		res, err := s.User().Login(ctx, &dto.LoginRequest{
			Username: username,
			Password: password,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

func TestUserService_Logout(t *testing.T) {
	jwtSecret := []byte("secret")
	logger := zerolog.Nop()
	ctx := context.Background()
	userID := 1

	mockRepo := new(MockRepo)
	mockUserRepo := new(MockUserRepo)
	s := New(mockRepo, logger, jwtSecret)

	mockRepo.On("User").Return(mockUserRepo)
	mockUserRepo.On("GetUserByID", ctx, userID).Return(&models.User{ID: userID}, nil)
	mockUserRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.RefreshToken == ""
	})).Return(nil)

	err := s.User().Logout(ctx, userID)
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockRepo)
	mockUserRepo := new(MockUserRepo)
	logger := zerolog.Nop()

	s := New(mockRepo, logger, []byte("secret"))

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		req := &dto.UserRequest{
			Username: "newuser",
			Password: "password123",
			Role:     "employee",
			Name:     "New User",
		}

		mockRepo.On("User").Return(mockUserRepo)
		mockUserRepo.On("CreateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.Username == req.Username && u.Name == req.Name
		})).Return(nil)

		res, err := s.User().CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Username, res.Username)
		mockUserRepo.AssertExpectations(t)
	})
}
