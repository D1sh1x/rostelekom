package dto

import "time"

type CommentRequest struct {
	TaskID int    `json:"task_id" validate:"required"`
	Text   string `json:"text" validate:"required"`
}

type CommentResponse struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
