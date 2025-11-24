package dto

import (
	"time"
)

type SkillRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type SkillResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSkillRequest struct {
	SkillID int `json:"skill_id"`
	Level   int `json:"level"`
}

type UserSkillResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	SkillID   int       `json:"skill_id"`
	SkillName string    `json:"skill_name,omitempty"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

