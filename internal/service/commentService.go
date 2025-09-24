package service

import (
	dto "SkillsTracker/internal/dto"
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"context"
	"errors"

	"github.com/rs/zerolog"
)

type commentService struct {
	repo   repository.RepositoryInterface
	logger zerolog.Logger
}

func newCommentService(repo repository.RepositoryInterface, logger zerolog.Logger) *commentService {
	return &commentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *commentService) CreateComment(ctx context.Context, req *dto.CommentRequest, userID int) (*dto.CommentResponse, error) {
	logger := s.logger.With().Int("task_id", req.TaskID).Int("user_id", userID).Logger()
	logger.Info().Msg("creating comment")

	_, err := s.repo.Task().GetTaskByID(ctx, req.TaskID)
	if err != nil {
		logger.Error().Msg("task not found")
		return nil, errors.New("task not found")
	}

	comment := &models.Comment{
		TaskID: req.TaskID,
		UserID: userID,
		Text:   req.Text,
	}

	if err := s.repo.Comment().CreateComment(ctx, comment); err != nil {
		logger.Error().Err(err).Msg("failed to create comment")
		return nil, err
	}

	logger.Info().Int("id", comment.ID).Msg("comment created")
	return &dto.CommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *commentService) GetCommentByID(ctx context.Context, commentID int) (*dto.CommentResponse, error) {
	logger := s.logger.With().Int("id", commentID).Logger()
	logger.Info().Msg("getting comment")

	comment, err := s.repo.Comment().GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment")
		return nil, err
	}

	return &dto.CommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *commentService) GetCommentsByTaskID(ctx context.Context, taskID int) ([]*dto.CommentResponse, error) {
	logger := s.logger.With().Int("task_id", taskID).Logger()
	logger.Info().Msg("getting comments")

	comments, err := s.repo.Comment().GetCommentsByTaskID(ctx, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments")
		return nil, err
	}

	responses := make([]*dto.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = &dto.CommentResponse{
			ID:        comment.ID,
			TaskID:    comment.TaskID,
			UserID:    comment.UserID,
			Text:      comment.Text,
			CreatedAt: comment.CreatedAt,
		}
	}

	return responses, nil
}

func (s *commentService) UpdateComment(ctx context.Context, commentID int, req *dto.CommentRequest, userID int) error {
	logger := s.logger.With().Int("id", commentID).Int("user_id", userID).Logger()
	logger.Info().Msg("updating comment")

	comment, err := s.repo.Comment().GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment")
		return err
	}

	if comment.UserID != userID {
		logger.Warn().Msg("user cannot edit comment")
		return errors.New("forbidden")
	}

	comment.Text = req.Text

	if err := s.repo.Comment().UpdateComment(ctx, comment); err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		return err
	}

	logger.Info().Int("id", comment.ID).Msg("comment updated")
	return nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentID int, userID int) error {
	logger := s.logger.With().Int("id", commentID).Int("user_id", userID).Logger()
	logger.Info().Msg("deleting comment")

	comment, err := s.repo.Comment().GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment")
		return err
	}

	if comment.UserID != userID {
		logger.Warn().Msg("user cannot delete comment")
		return errors.New("forbidden")
	}

	if err := s.repo.Comment().DeleteComment(ctx, commentID); err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}

	logger.Info().Msg("comment deleted")
	return nil
}
