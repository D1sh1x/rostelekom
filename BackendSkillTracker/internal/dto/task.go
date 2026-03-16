package dto

import "time"

type TaskRequest struct {
	EmployeeID  int    `json:"employee_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description"`
	Deadline    string `json:"deadline" validate:"required"`
	Progress    int    `json:"progress" validate:"min=0,max=100"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
}

type TaskResponse struct {
	ID          int       `json:"id"`
	EmployeeID  int       `json:"employee_id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AttachmentResponse struct {
	ID         int       `json:"id"`
	TaskID     int       `json:"task_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type TaskHistoryResponse struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	ChangedBy int       `json:"changed_by"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskFilter struct {
	Status     string `query:"status"`
	EmployeeID int    `query:"employee_id"`
	CreatorID  int    `query:"creator_id"`
	Search     string `query:"search"`
	FromDate   string `query:"from_date"`
	ToDate     string `query:"to_date"`
}
