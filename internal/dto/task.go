package dto

import (
	"time"
)

type TaskRequest struct {
	ProjectID    int     `json:"project_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Deadline     string  `json:"deadline"`
	Hours        int     `json:"hours"`
	Priority     string  `json:"priority"`
	Type         string  `json:"type"`
	ParentTaskID *int    `json:"parent_task_id,omitempty"`
	SkillIDs     []int   `json:"skill_ids,omitempty"`
	AssigneeIDs  []int   `json:"assignee_ids,omitempty"`
}

type TaskResponse struct {
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Deadline     time.Time `json:"deadline"`
	Status       string    `json:"status"`
	Progress     int       `json:"progress"`
	Hours        int       `json:"hours"`
	Priority     string    `json:"priority"`
	Type         string    `json:"type"`
	ParentTaskID *int      `json:"parent_task_id,omitempty"`
	SkillIDs     []int     `json:"skill_ids,omitempty"`
	AssigneeIDs  []int     `json:"assignee_ids,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
