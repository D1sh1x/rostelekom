package models

import (
	"time"
)

type Role string
type TaskStatus string

const (
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"

	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Task struct {
	ID          int        `db:"id"`
	EmployeeID  int        `db:"employee_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Deadline    time.Time  `db:"deadline"`
	Status      TaskStatus `db:"status"`
	Progress    int        `db:"progress"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type Comment struct {
	ID        int       `db:"id"`
	TaskID    int       `db:"task_id"`
	UserID    int       `db:"user_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
