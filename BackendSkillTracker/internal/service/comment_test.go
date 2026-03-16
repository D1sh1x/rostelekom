package service

import (
	"context"
	"skilltracker/internal/models"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_CreateComment(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCommentRepo := new(MockCommentRepo)
	logger := zerolog.Nop()
	s := New(mockRepo, logger, []byte("secret"))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Comment").Return(mockCommentRepo)
		mockCommentRepo.On("CreateComment", ctx, mock.MatchedBy(func(c *models.Comment) bool {
			return c.Text == "test comment" && c.TaskID == 1
		})).Return(nil)

		res, err := s.Comment().CreateComment(ctx, 1, 2, "test comment")

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test comment", res.Text)
		mockCommentRepo.AssertExpectations(t)
	})
}

func TestCommentService_DeleteComment(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCommentRepo := new(MockCommentRepo)
	logger := zerolog.Nop()
	s := New(mockRepo, logger, []byte("secret"))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		comment := &models.Comment{ID: 1, UserID: 2}
		mockRepo.On("Comment").Return(mockCommentRepo)
		mockCommentRepo.On("GetCommentByID", ctx, 1).Return(comment, nil)
		mockCommentRepo.On("DeleteComment", ctx, 1).Return(nil)

		err := s.Comment().DeleteComment(ctx, 1, 2)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		comment := &models.Comment{ID: 1, UserID: 2}
		mockRepo.On("Comment").Return(mockCommentRepo)
		mockCommentRepo.On("GetCommentByID", ctx, 1).Return(comment, nil)

		err := s.Comment().DeleteComment(ctx, 1, 3)
		assert.Error(t, err)
		assert.Equal(t, "forbidden", err.Error())
	})
}
