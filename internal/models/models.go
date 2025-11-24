package models

import (
	"time"
)

type Role string
type TaskStatus string
type Priority string
type TaskType string

const (
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"

	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"

	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"

	TaskTypeFeature TaskType = "feature"
	TaskTypeBug     TaskType = "bug"
	TaskTypeTask    TaskType = "task"
	TaskTypeEpic    TaskType = "epic"
)

// User - пользователь системы
type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Project - проект, в котором находятся задачи
type Project struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	ManagerID   int       `db:"manager_id"` // менеджер, ответственный за проект
	Status      string    `db:"status"`     // например: active, completed, on_hold
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Skill - навык, который может быть у сотрудника
type Skill struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Category    string    `db:"category"` // например: technical, soft, domain
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// UserSkill - связь многие-ко-многим между пользователями и навыками
type UserSkill struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	SkillID   int       `db:"skill_id"`
	Level     int       `db:"level"` // уровень навыка (1-5, например)
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Task - задача в проекте
type Task struct {
	ID           int        `db:"id"`
	ProjectID    int        `db:"project_id"` // ссылка на проект
	Title        string     `db:"title"`
	Description  string     `db:"description"`
	Deadline     time.Time  `db:"deadline"`
	Status       TaskStatus `db:"status"`
	Progress     int        `db:"progress"`       // процент выполнения
	Hours        int        `db:"hours"`          // оценка в часах
	Priority     Priority   `db:"priority"`       // приоритет задачи
	Type         TaskType   `db:"type"`           // тип задачи
	ParentTaskID *int       `db:"parent_task_id"` // для поддержки зависимостей задач
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// TaskSkill - связь многие-ко-многим между задачами и навыками
type TaskSkill struct {
	ID        int       `db:"id"`
	TaskID    int       `db:"task_id"`
	SkillID   int       `db:"skill_id"`
	CreatedAt time.Time `db:"created_at"`
}

// TaskAssignee - назначение сотрудников на задачу (для поддержки нескольких исполнителей)
type TaskAssignee struct {
	ID        int       `db:"id"`
	TaskID    int       `db:"task_id"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

// Comment - комментарий к задаче
type Comment struct {
	ID        int       `db:"id"`
	TaskID    int       `db:"task_id"`
	UserID    int       `db:"user_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

// ProjectMember - участники проекта
type ProjectMember struct {
	ID        int       `db:"id"`
	ProjectID int       `db:"project_id"`
	UserID    int       `db:"user_id"`
	Role      string    `db:"role"` // project_manager, developer, tester и т.д.
	JoinedAt  time.Time `db:"joined_at"`
	CreatedAt time.Time `db:"created_at"`
}
