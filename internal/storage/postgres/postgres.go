package postgres

import (
	"SkillsTracker/internal/models"
	"SkillsTracker/internal/repository"
	"SkillsTracker/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(dsn string) (repository.RepositoryInterface, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	st := &Storage{db: db}

	return st, nil
}

func (s *Storage) User() repository.UserRepositoryInterface       { return s }
func (s *Storage) Task() repository.TaskRepositoryInterface       { return s }
func (s *Storage) Comment() repository.CommentRepositoryInterface { return s }

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash, role, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		user.Username, user.PasswordHash, user.Role, user.Name).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT id, username, password_hash, role, name, created_at, updated_at FROM users`
	if err := s.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, role, name, created_at, updated_at FROM users WHERE id = $1`
	if err := s.db.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, role, name, created_at, updated_at FROM users WHERE username = $1`
	if err := s.db.GetContext(ctx, &user, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1,
		    password_hash = $2,
		    role = $3,
		    name = $4,
		    updated_at = NOW()
		WHERE id = $5`
	res, err := s.db.ExecContext(ctx, query,
		user.Username, user.PasswordHash, user.Role, user.Name, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if ra == 0 {
		return storage.ErrUserNotFound
	}
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrUserNotFound
	}
	return nil
}

func (s *Storage) CreateTask(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (employee_id, title, description, deadline, status, progress)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		task.EmployeeID, task.Title, task.Description, task.Deadline, task.Status, task.Progress).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

func (s *Storage) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	query := `SELECT id, employee_id, title, description, deadline, status, progress, created_at, updated_at FROM tasks WHERE id = $1`
	if err := s.db.GetContext(ctx, &task, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func (s *Storage) GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]models.Task, error) {
	var tasks []models.Task
	query := `
		SELECT id, employee_id, title, description, deadline, status, progress, created_at, updated_at
		FROM tasks
		WHERE employee_id = $1
		ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &tasks, query, employeeID); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, deadline = $3, status = $4, progress = $5, updated_at = NOW()
		WHERE id = $6`
	result, err := s.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Deadline, task.Status, task.Progress, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskNotUpdated
	}
	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}

func (s *Storage) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (task_id, user_id, text)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query,
		comment.TaskID, comment.UserID, comment.Text).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (s *Storage) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	var c models.Comment
	query := `SELECT id, task_id, user_id, text, created_at FROM comments WHERE id = $1`
	if err := s.db.GetContext(ctx, &c, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &c, nil
}

func (s *Storage) GetCommentsByTaskID(ctx context.Context, taskID int) ([]models.Comment, error) {
	var comments []models.Comment
	query := `SELECT id, task_id, user_id, text, created_at FROM comments WHERE task_id = $1 ORDER BY created_at ASC`
	if err := s.db.SelectContext(ctx, &comments, query, taskID); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	return comments, nil
}

func (s *Storage) UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET text = $1
		WHERE id = $2`
	res, err := s.db.ExecContext(ctx, query, comment.Text, comment.ID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if ra == 0 {
		return storage.ErrCommentNotFound
	}
	return nil
}

func (s *Storage) DeleteComment(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrCommentNotFound
	}
	return nil
}
