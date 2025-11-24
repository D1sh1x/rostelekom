package dto

import (
	"time"
)

type ProjectRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	MemberIDs   []int    `json:"member_ids,omitempty"`
}

type ProjectResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ManagerID   int       `json:"manager_id"`
	Status      string    `json:"status"`
	MemberIDs   []int     `json:"member_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

