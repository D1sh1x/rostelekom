package dto

import (
	"time"
)

type TaskRequest struct {
	EmployeeID  int    `json:"employee_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Progress    int    `json:"progress"`
	Status      string `json:"status"`
}

type TaskResponse struct {
	ID          int       `json:"id"`
	EmployeeID  int       `json:"employee_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
