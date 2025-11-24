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

func (s *Storage) User() repository.UserRepositoryInterface              { return s }
func (s *Storage) Task() repository.TaskRepositoryInterface              { return s }
func (s *Storage) Comment() repository.CommentRepositoryInterface       { return s }
func (s *Storage) Project() repository.ProjectRepositoryInterface        { return s }
func (s *Storage) Skill() repository.SkillRepositoryInterface           { return s }
func (s *Storage) UserSkill() repository.UserSkillRepositoryInterface   { return s }
func (s *Storage) TaskSkill() repository.TaskSkillRepositoryInterface    { return s }
func (s *Storage) TaskAssignee() repository.TaskAssigneeRepositoryInterface { return s }
func (s *Storage) ProjectMember() repository.ProjectMemberRepositoryInterface { return s }

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash, role, name, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		user.Username, user.PasswordHash, user.Role, user.Name, user.Email).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT id, username, password_hash, role, name, email, created_at, updated_at FROM users`
	if err := s.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, role, name, email, created_at, updated_at FROM users WHERE id = $1`
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
	query := `SELECT id, username, password_hash, role, name, email, created_at, updated_at FROM users WHERE username = $1`
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
		    email = $5,
		    updated_at = NOW()
		WHERE id = $6`
	res, err := s.db.ExecContext(ctx, query,
		user.Username, user.PasswordHash, user.Role, user.Name, user.Email, user.ID)
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
		INSERT INTO tasks (project_id, title, description, deadline, status, progress, hours, priority, type, parent_task_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		task.ProjectID, task.Title, task.Description, task.Deadline, task.Status, task.Progress,
		task.Hours, task.Priority, task.Type, task.ParentTaskID).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

func (s *Storage) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	query := `SELECT id, project_id, title, description, deadline, status, progress, hours, priority, type, parent_task_id, created_at, updated_at FROM tasks WHERE id = $1`
	if err := s.db.GetContext(ctx, &task, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func (s *Storage) GetTasksByProjectID(ctx context.Context, projectID int) ([]models.Task, error) {
	var tasks []models.Task
	query := `
		SELECT id, project_id, title, description, deadline, status, progress, hours, priority, type, parent_task_id, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &tasks, query, projectID); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

func (s *Storage) GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error) {
	var tasks []models.Task
	query := `
		SELECT t.id, t.project_id, t.title, t.description, t.deadline, t.status, t.progress, t.hours, t.priority, t.type, t.parent_task_id, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN task_assignees ta ON t.id = ta.task_id
		WHERE ta.user_id = $1
		ORDER BY t.created_at DESC`
	if err := s.db.SelectContext(ctx, &tasks, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

func (s *Storage) GetChildTasks(ctx context.Context, parentTaskID int) ([]models.Task, error) {
	var tasks []models.Task
	query := `
		SELECT id, project_id, title, description, deadline, status, progress, hours, priority, type, parent_task_id, created_at, updated_at
		FROM tasks
		WHERE parent_task_id = $1
		ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &tasks, query, parentTaskID); err != nil {
		return nil, fmt.Errorf("failed to get child tasks: %w", err)
	}
	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, deadline = $3, status = $4, progress = $5, hours = $6, priority = $7, type = $8, parent_task_id = $9, updated_at = NOW()
		WHERE id = $10`
	result, err := s.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Deadline, task.Status, task.Progress,
		task.Hours, task.Priority, task.Type, task.ParentTaskID, task.ID)
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

// Project Repository Methods
func (s *Storage) CreateProject(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (name, description, manager_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		project.Name, project.Description, project.ManagerID, project.Status).
		Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

func (s *Storage) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	query := `SELECT id, name, description, manager_id, status, created_at, updated_at FROM projects WHERE id = $1`
	if err := s.db.GetContext(ctx, &project, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return &project, nil
}

func (s *Storage) GetProjectsByManagerID(ctx context.Context, managerID int) ([]models.Project, error) {
	var projects []models.Project
	query := `SELECT id, name, description, manager_id, status, created_at, updated_at FROM projects WHERE manager_id = $1 ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &projects, query, managerID); err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	return projects, nil
}

func (s *Storage) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project
	query := `SELECT id, name, description, manager_id, status, created_at, updated_at FROM projects ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &projects, query); err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	return projects, nil
}

func (s *Storage) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2, manager_id = $3, status = $4, updated_at = NOW()
		WHERE id = $5`
	result, err := s.db.ExecContext(ctx, query,
		project.Name, project.Description, project.ManagerID, project.Status, project.ID)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectNotFound
	}
	return nil
}

func (s *Storage) DeleteProject(ctx context.Context, id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectNotFound
	}
	return nil
}

// Skill Repository Methods
func (s *Storage) CreateSkill(ctx context.Context, skill *models.Skill) error {
	query := `
		INSERT INTO skills (name, description, category)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		skill.Name, skill.Description, skill.Category).
		Scan(&skill.ID, &skill.CreatedAt, &skill.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}
	return nil
}

func (s *Storage) GetSkillByID(ctx context.Context, id int) (*models.Skill, error) {
	var skill models.Skill
	query := `SELECT id, name, description, category, created_at, updated_at FROM skills WHERE id = $1`
	if err := s.db.GetContext(ctx, &skill, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrSkillNotFound
		}
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return &skill, nil
}

func (s *Storage) GetAllSkills(ctx context.Context) ([]models.Skill, error) {
	var skills []models.Skill
	query := `SELECT id, name, description, category, created_at, updated_at FROM skills ORDER BY name ASC`
	if err := s.db.SelectContext(ctx, &skills, query); err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	return skills, nil
}

func (s *Storage) GetSkillsByCategory(ctx context.Context, category string) ([]models.Skill, error) {
	var skills []models.Skill
	query := `SELECT id, name, description, category, created_at, updated_at FROM skills WHERE category = $1 ORDER BY name ASC`
	if err := s.db.SelectContext(ctx, &skills, query, category); err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	return skills, nil
}

func (s *Storage) UpdateSkill(ctx context.Context, skill *models.Skill) error {
	query := `
		UPDATE skills
		SET name = $1, description = $2, category = $3, updated_at = NOW()
		WHERE id = $4`
	result, err := s.db.ExecContext(ctx, query,
		skill.Name, skill.Description, skill.Category, skill.ID)
	if err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrSkillNotFound
	}
	return nil
}

func (s *Storage) DeleteSkill(ctx context.Context, id int) error {
	query := `DELETE FROM skills WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrSkillNotFound
	}
	return nil
}

// UserSkill Repository Methods
func (s *Storage) CreateUserSkill(ctx context.Context, userSkill *models.UserSkill) error {
	query := `
		INSERT INTO user_skills (user_id, skill_id, level)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		userSkill.UserID, userSkill.SkillID, userSkill.Level).
		Scan(&userSkill.ID, &userSkill.CreatedAt, &userSkill.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user skill: %w", err)
	}
	return nil
}

func (s *Storage) GetUserSkillByID(ctx context.Context, id int) (*models.UserSkill, error) {
	var userSkill models.UserSkill
	query := `SELECT id, user_id, skill_id, level, created_at, updated_at FROM user_skills WHERE id = $1`
	if err := s.db.GetContext(ctx, &userSkill, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserSkillNotFound
		}
		return nil, fmt.Errorf("failed to get user skill: %w", err)
	}
	return &userSkill, nil
}

func (s *Storage) GetUserSkillsByUserID(ctx context.Context, userID int) ([]models.UserSkill, error) {
	var userSkills []models.UserSkill
	query := `SELECT id, user_id, skill_id, level, created_at, updated_at FROM user_skills WHERE user_id = $1 ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &userSkills, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get user skills: %w", err)
	}
	return userSkills, nil
}

func (s *Storage) GetUserSkillsBySkillID(ctx context.Context, skillID int) ([]models.UserSkill, error) {
	var userSkills []models.UserSkill
	query := `SELECT id, user_id, skill_id, level, created_at, updated_at FROM user_skills WHERE skill_id = $1 ORDER BY level DESC`
	if err := s.db.SelectContext(ctx, &userSkills, query, skillID); err != nil {
		return nil, fmt.Errorf("failed to get user skills: %w", err)
	}
	return userSkills, nil
}

func (s *Storage) GetUserSkill(ctx context.Context, userID, skillID int) (*models.UserSkill, error) {
	var userSkill models.UserSkill
	query := `SELECT id, user_id, skill_id, level, created_at, updated_at FROM user_skills WHERE user_id = $1 AND skill_id = $2`
	if err := s.db.GetContext(ctx, &userSkill, query, userID, skillID); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserSkillNotFound
		}
		return nil, fmt.Errorf("failed to get user skill: %w", err)
	}
	return &userSkill, nil
}

func (s *Storage) UpdateUserSkill(ctx context.Context, userSkill *models.UserSkill) error {
	query := `
		UPDATE user_skills
		SET level = $1, updated_at = NOW()
		WHERE id = $2`
	result, err := s.db.ExecContext(ctx, query, userSkill.Level, userSkill.ID)
	if err != nil {
		return fmt.Errorf("failed to update user skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrUserSkillNotFound
	}
	return nil
}

func (s *Storage) DeleteUserSkill(ctx context.Context, id int) error {
	query := `DELETE FROM user_skills WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrUserSkillNotFound
	}
	return nil
}

func (s *Storage) DeleteUserSkillByUserAndSkill(ctx context.Context, userID, skillID int) error {
	query := `DELETE FROM user_skills WHERE user_id = $1 AND skill_id = $2`
	result, err := s.db.ExecContext(ctx, query, userID, skillID)
	if err != nil {
		return fmt.Errorf("failed to delete user skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrUserSkillNotFound
	}
	return nil
}

// TaskSkill Repository Methods
func (s *Storage) CreateTaskSkill(ctx context.Context, taskSkill *models.TaskSkill) error {
	query := `
		INSERT INTO task_skills (task_id, skill_id)
		VALUES ($1, $2)
		RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query,
		taskSkill.TaskID, taskSkill.SkillID).
		Scan(&taskSkill.ID, &taskSkill.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task skill: %w", err)
	}
	return nil
}

func (s *Storage) GetTaskSkillByID(ctx context.Context, id int) (*models.TaskSkill, error) {
	var taskSkill models.TaskSkill
	query := `SELECT id, task_id, skill_id, created_at FROM task_skills WHERE id = $1`
	if err := s.db.GetContext(ctx, &taskSkill, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrTaskSkillNotFound
		}
		return nil, fmt.Errorf("failed to get task skill: %w", err)
	}
	return &taskSkill, nil
}

func (s *Storage) GetTaskSkillsByTaskID(ctx context.Context, taskID int) ([]models.TaskSkill, error) {
	var taskSkills []models.TaskSkill
	query := `SELECT id, task_id, skill_id, created_at FROM task_skills WHERE task_id = $1 ORDER BY created_at ASC`
	if err := s.db.SelectContext(ctx, &taskSkills, query, taskID); err != nil {
		return nil, fmt.Errorf("failed to get task skills: %w", err)
	}
	return taskSkills, nil
}

func (s *Storage) GetTaskSkillsBySkillID(ctx context.Context, skillID int) ([]models.TaskSkill, error) {
	var taskSkills []models.TaskSkill
	query := `SELECT id, task_id, skill_id, created_at FROM task_skills WHERE skill_id = $1 ORDER BY created_at ASC`
	if err := s.db.SelectContext(ctx, &taskSkills, query, skillID); err != nil {
		return nil, fmt.Errorf("failed to get task skills: %w", err)
	}
	return taskSkills, nil
}

func (s *Storage) DeleteTaskSkill(ctx context.Context, id int) error {
	query := `DELETE FROM task_skills WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskSkillNotFound
	}
	return nil
}

func (s *Storage) DeleteTaskSkillByTaskAndSkill(ctx context.Context, taskID, skillID int) error {
	query := `DELETE FROM task_skills WHERE task_id = $1 AND skill_id = $2`
	result, err := s.db.ExecContext(ctx, query, taskID, skillID)
	if err != nil {
		return fmt.Errorf("failed to delete task skill: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskSkillNotFound
	}
	return nil
}

func (s *Storage) DeleteTaskSkillsByTaskID(ctx context.Context, taskID int) error {
	query := `DELETE FROM task_skills WHERE task_id = $1`
	_, err := s.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task skills: %w", err)
	}
	return nil
}

// TaskAssignee Repository Methods
func (s *Storage) CreateTaskAssignee(ctx context.Context, taskAssignee *models.TaskAssignee) error {
	query := `
		INSERT INTO task_assignees (task_id, user_id)
		VALUES ($1, $2)
		RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query,
		taskAssignee.TaskID, taskAssignee.UserID).
		Scan(&taskAssignee.ID, &taskAssignee.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task assignee: %w", err)
	}
	return nil
}

func (s *Storage) GetTaskAssigneeByID(ctx context.Context, id int) (*models.TaskAssignee, error) {
	var taskAssignee models.TaskAssignee
	query := `SELECT id, task_id, user_id, created_at FROM task_assignees WHERE id = $1`
	if err := s.db.GetContext(ctx, &taskAssignee, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrTaskAssigneeNotFound
		}
		return nil, fmt.Errorf("failed to get task assignee: %w", err)
	}
	return &taskAssignee, nil
}

func (s *Storage) GetTaskAssigneesByTaskID(ctx context.Context, taskID int) ([]models.TaskAssignee, error) {
	var taskAssignees []models.TaskAssignee
	query := `SELECT id, task_id, user_id, created_at FROM task_assignees WHERE task_id = $1 ORDER BY created_at ASC`
	if err := s.db.SelectContext(ctx, &taskAssignees, query, taskID); err != nil {
		return nil, fmt.Errorf("failed to get task assignees: %w", err)
	}
	return taskAssignees, nil
}

func (s *Storage) GetTaskAssigneesByUserID(ctx context.Context, userID int) ([]models.TaskAssignee, error) {
	var taskAssignees []models.TaskAssignee
	query := `SELECT id, task_id, user_id, created_at FROM task_assignees WHERE user_id = $1 ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &taskAssignees, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get task assignees: %w", err)
	}
	return taskAssignees, nil
}

func (s *Storage) DeleteTaskAssignee(ctx context.Context, id int) error {
	query := `DELETE FROM task_assignees WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task assignee: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskAssigneeNotFound
	}
	return nil
}

func (s *Storage) DeleteTaskAssigneeByTaskAndUser(ctx context.Context, taskID, userID int) error {
	query := `DELETE FROM task_assignees WHERE task_id = $1 AND user_id = $2`
	result, err := s.db.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task assignee: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskAssigneeNotFound
	}
	return nil
}

func (s *Storage) DeleteTaskAssigneesByTaskID(ctx context.Context, taskID int) error {
	query := `DELETE FROM task_assignees WHERE task_id = $1`
	_, err := s.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task assignees: %w", err)
	}
	return nil
}

// ProjectMember Repository Methods
func (s *Storage) CreateProjectMember(ctx context.Context, projectMember *models.ProjectMember) error {
	query := `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, joined_at, created_at`
	err := s.db.QueryRowContext(ctx, query,
		projectMember.ProjectID, projectMember.UserID, projectMember.Role).
		Scan(&projectMember.ID, &projectMember.JoinedAt, &projectMember.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project member: %w", err)
	}
	return nil
}

func (s *Storage) GetProjectMemberByID(ctx context.Context, id int) (*models.ProjectMember, error) {
	var projectMember models.ProjectMember
	query := `SELECT id, project_id, user_id, role, joined_at, created_at FROM project_members WHERE id = $1`
	if err := s.db.GetContext(ctx, &projectMember, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrProjectMemberNotFound
		}
		return nil, fmt.Errorf("failed to get project member: %w", err)
	}
	return &projectMember, nil
}

func (s *Storage) GetProjectMembersByProjectID(ctx context.Context, projectID int) ([]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	query := `SELECT id, project_id, user_id, role, joined_at, created_at FROM project_members WHERE project_id = $1 ORDER BY joined_at ASC`
	if err := s.db.SelectContext(ctx, &projectMembers, query, projectID); err != nil {
		return nil, fmt.Errorf("failed to get project members: %w", err)
	}
	return projectMembers, nil
}

func (s *Storage) GetProjectMembersByUserID(ctx context.Context, userID int) ([]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	query := `SELECT id, project_id, user_id, role, joined_at, created_at FROM project_members WHERE user_id = $1 ORDER BY joined_at DESC`
	if err := s.db.SelectContext(ctx, &projectMembers, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get project members: %w", err)
	}
	return projectMembers, nil
}

func (s *Storage) GetProjectMember(ctx context.Context, projectID, userID int) (*models.ProjectMember, error) {
	var projectMember models.ProjectMember
	query := `SELECT id, project_id, user_id, role, joined_at, created_at FROM project_members WHERE project_id = $1 AND user_id = $2`
	if err := s.db.GetContext(ctx, &projectMember, query, projectID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrProjectMemberNotFound
		}
		return nil, fmt.Errorf("failed to get project member: %w", err)
	}
	return &projectMember, nil
}

func (s *Storage) UpdateProjectMember(ctx context.Context, projectMember *models.ProjectMember) error {
	query := `
		UPDATE project_members
		SET role = $1
		WHERE id = $2`
	result, err := s.db.ExecContext(ctx, query, projectMember.Role, projectMember.ID)
	if err != nil {
		return fmt.Errorf("failed to update project member: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectMemberNotFound
	}
	return nil
}

func (s *Storage) DeleteProjectMember(ctx context.Context, id int) error {
	query := `DELETE FROM project_members WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project member: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectMemberNotFound
	}
	return nil
}

func (s *Storage) DeleteProjectMemberByProjectAndUser(ctx context.Context, projectID, userID int) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`
	result, err := s.db.ExecContext(ctx, query, projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete project member: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectMemberNotFound
	}
	return nil
}
