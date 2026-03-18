package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string
type TaskStatus string

const (
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"

	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

type User struct {
	ID           int            `gorm:"primaryKey"`
	Username     string         `gorm:"unique;not null;size:50"`
	PasswordHash string         `gorm:"not null"`
	Role         Role           `gorm:"not null;type:varchar(20)"`
	Name         string         `gorm:"not null;size:100"`
	RefreshToken string         `gorm:"index"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Skills []Skill `gorm:"many2many:user_skills;"`
}

type Task struct {
	ID          int            `gorm:"primaryKey"`
	EmployeeID  int            `gorm:"not null"`
	CreatorID   int            `gorm:"not null"`
	Title       string         `gorm:"not null;size:200"`
	Description string         `gorm:"not null"`
	Deadline    time.Time      `gorm:"not null"`
	Status      TaskStatus     `gorm:"not null;type:varchar(20);default:pending"`
	Progress    int            `gorm:"not null;default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Employee       User                `gorm:"foreignKey:EmployeeID"`
	Creator        User                `gorm:"foreignKey:CreatorID"`
	Attachments    []FileAttachment    `gorm:"foreignKey:TaskID"`
	History        []TaskStatusHistory `gorm:"foreignKey:TaskID"`
	RequiredSkills []Skill             `gorm:"many2many:task_skills;"`
}

type TaskStatusHistory struct {
	ID        int        `gorm:"primaryKey"`
	TaskID    int        `gorm:"not null"`
	OldStatus TaskStatus `gorm:"not null;type:varchar(20)"`
	NewStatus TaskStatus `gorm:"not null;type:varchar(20)"`
	ChangedBy int        `gorm:"not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:ChangedBy"`
}

type FileAttachment struct {
	ID         int       `gorm:"primaryKey"`
	TaskID     int       `gorm:"not null"`
	FileName   string    `gorm:"not null"`
	FilePath   string    `gorm:"not null"`
	FileSize   int64     `gorm:"not null"`
	UploadedAt time.Time `gorm:"autoCreateTime"`
}

type Comment struct {
	ID        int            `gorm:"primaryKey"`
	TaskID    int            `gorm:"not null"`
	UserID    int            `gorm:"not null"`
	Text      string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Task Task `gorm:"foreignKey:TaskID"`
	User User `gorm:"foreignKey:UserID"`
}

type Skill struct {
	ID          int       `gorm:"primaryKey"`
	Name        string    `gorm:"unique;not null;size:100"`
	Description string    `gorm:"size:500"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	Users []User `gorm:"many2many:user_skills;"`
	Tasks []Task `gorm:"many2many:task_skills;"`
}

type UserSkill struct {
	UserID  int `gorm:"primaryKey"`
	SkillID int `gorm:"primaryKey"`
}

type TaskSkill struct {
	TaskID  int `gorm:"primaryKey"`
	SkillID int `gorm:"primaryKey"`
}
