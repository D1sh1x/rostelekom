package postgres

import (
	"context"
	"skilltracker/internal/models"
	"skilltracker/internal/repository"
	"skilltracker/internal/dto"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func New(dsn string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.Comment{}, &models.FileAttachment{}, &models.TaskStatusHistory{}); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) User() repository.UserRepository       { return s }
func (s *Storage) Task() repository.TaskRepository       { return s }
func (s *Storage) Comment() repository.CommentRepository { return s }
func (s *Storage) File() repository.FileRepository       { return s }

// USERS

func (s *Storage) CreateUser(ctx context.Context, u *models.User) error {
	return s.db.WithContext(ctx).Create(u).Error
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetUserByRefreshToken(ctx context.Context, token string) (*models.User, error) {
	var u models.User
	if err := s.db.WithContext(ctx).Where("refresh_token = ?", token).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) UpdateUser(ctx context.Context, u *models.User) error {
	return s.db.WithContext(ctx).Save(u).Error
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (s *Storage) GetUsers(ctx context.Context) ([]*models.User, error) {
	var out []*models.User
	err := s.db.WithContext(ctx).Order("id").Find(&out).Error
	return out, err
}

// TASKS

func (s *Storage) CreateTask(ctx context.Context, t *models.Task) error {
	return s.db.WithContext(ctx).Create(t).Error
}

func (s *Storage) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	var t models.Task
	if err := s.db.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Storage) GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]models.Task, error) {
	var ts []models.Task
	err := s.db.WithContext(ctx).Where("employee_id = ?", employeeID).Order("created_at DESC").Find(&ts).Error
	return ts, err
}

func (s *Storage) UpdateTask(ctx context.Context, t *models.Task) error {
	return s.db.WithContext(ctx).Save(t).Error
}

func (s *Storage) DeleteTask(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

// COMMENTS

func (s *Storage) CreateComment(ctx context.Context, cmt *models.Comment) error {
	return s.db.WithContext(ctx).Create(cmt).Error
}

func (s *Storage) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	var c models.Comment
	if err := s.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Storage) GetCommentsByTaskID(ctx context.Context, taskID int) ([]models.Comment, error) {
	var cs []models.Comment
	err := s.db.WithContext(ctx).Where("task_id = ?", taskID).Order("id").Find(&cs).Error
	return cs, err
}

func (s *Storage) UpdateComment(ctx context.Context, c *models.Comment) error {
	return s.db.WithContext(ctx).Save(c).Error
}

func (s *Storage) DeleteComment(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}

// FILES

func (s *Storage) CreateAttachment(ctx context.Context, f *models.FileAttachment) error {
	return s.db.WithContext(ctx).Create(f).Error
}

func (s *Storage) CreateHistory(ctx context.Context, h *models.TaskStatusHistory) error {
	return s.db.WithContext(ctx).Create(h).Error
}

func (s *Storage) GetHistoryByTaskID(ctx context.Context, taskID int) ([]models.TaskStatusHistory, error) {
	var out []models.TaskStatusHistory
	err := s.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at DESC").Find(&out).Error
	return out, err
}

func (s *Storage) ListTasks(ctx context.Context, filter dto.TaskFilter) ([]models.Task, error) {
	query := s.db.WithContext(ctx)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.EmployeeID != 0 {
		query = query.Where("employee_id = ?", filter.EmployeeID)
	}
	if filter.CreatorID != 0 {
		query = query.Where("creator_id = ?", filter.CreatorID)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}
	if filter.FromDate != "" {
		if t, err := time.Parse("2006-01-02", filter.FromDate); err == nil {
			query = query.Where("deadline >= ?", t)
		}
	}
	if filter.ToDate != "" {
		if t, err := time.Parse("2006-01-02", filter.ToDate); err == nil {
			query = query.Where("deadline <= ?", t)
		}
	}

	var out []models.Task
	err := query.Order("created_at DESC").Find(&out).Error
	return out, err
}
