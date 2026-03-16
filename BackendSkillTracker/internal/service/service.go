package service

import (
    "context"
    "errors"
    "time"
    "skilltracker/internal/dto"
    "skilltracker/internal/models"
    "skilltracker/internal/repository"
    jwtutil "skilltracker/internal/utils/jwt"
    "github.com/rs/zerolog"
    "golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
    User() UserService
    Task() TaskService
    Comment() CommentService
}

type UserService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID int) error
	CreateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error)
	GetUsers(ctx context.Context) ([]*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req *dto.UserRequest) error
	DeleteUser(ctx context.Context, id int) error
	GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error)
}

type TaskService interface {
    CreateTask(ctx context.Context, req *dto.TaskRequest, creatorID int) (*dto.TaskResponse, error)
    GetTaskByID(ctx context.Context, id int) (*dto.TaskResponse, error)
    GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]*dto.TaskResponse, error)
    UpdateTask(ctx context.Context, id int, req *dto.TaskRequest, userID int) error
    DeleteTask(ctx context.Context, id int, userID int) error
    UploadAttachment(ctx context.Context, taskID int, userID int, fileName string, filePath string, fileSize int64) (*dto.AttachmentResponse, error)
    GetTaskHistory(ctx context.Context, taskID int) ([]*dto.TaskHistoryResponse, error)
    ListTasks(ctx context.Context, filter dto.TaskFilter) ([]*dto.TaskResponse, error)
}

type CommentService interface {
    CreateComment(ctx context.Context, taskID int, userID int, text string) (*dto.CommentResponse, error)
    GetCommentsByTaskID(ctx context.Context, taskID int) ([]*dto.CommentResponse, error)
    UpdateComment(ctx context.Context, id int, userID int, text string) error
    DeleteComment(ctx context.Context, id int, userID int) error
}

type services struct {
    repo      repository.Repository
    logger    zerolog.Logger
    jwtSecret []byte
}

func New(repo repository.Repository, l zerolog.Logger, jwtSecret []byte) ServiceInterface {
    return &services{repo: repo, logger: l, jwtSecret: jwtSecret}
}

// USER

func (s *services) User() UserService { return s }

func (s *services) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.repo.User().GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := jwtutil.GenerateAccessToken(u.ID, u.Username, string(u.Role), s.jwtSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwtutil.GenerateRefreshToken(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	u.RefreshToken = refreshToken
	if err := s.repo.User().UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *services) RefreshToken(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error) {
	// Simple validation: find user by refresh token
	// In production, we should probably validate the token signature too, but since we store it as is,
	// if it matches it's valid enough for this basic impl.
	// However, let's at least check if it's expired via jwtutil.ValidateToken (even if it doesn't have custom claims)
	_, err := jwtutil.ValidateToken(req.RefreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find user with this token
	users, err := s.repo.User().GetUsers(ctx) // This is inefficient but let's see if we have GetUserByRefreshToken
	if err != nil {
		return nil, err
	}
	var targetUser *models.User
	for _, u := range users {
		if u.RefreshToken == req.RefreshToken {
			targetUser = u
			break
		}
	}

	if targetUser == nil {
		return nil, errors.New("invalid refresh token")
	}

	// Rotation
	newAccessToken, _ := jwtutil.GenerateAccessToken(targetUser.ID, targetUser.Username, string(targetUser.Role), s.jwtSecret)
	newRefreshToken, _ := jwtutil.GenerateRefreshToken(s.jwtSecret)

	targetUser.RefreshToken = newRefreshToken
	if err := s.repo.User().UpdateUser(ctx, targetUser); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *services) Logout(ctx context.Context, userID int) error {
	u, err := s.repo.User().GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	u.RefreshToken = ""
	return s.repo.User().UpdateUser(ctx, u)
}

func (s *services) CreateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error) {
    hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    u := &models.User{
        Username: req.Username,
        PasswordHash: string(hash),
        Role: models.Role(req.Role),
        Name: req.Name,
    }
    if err := s.repo.User().CreateUser(ctx, u); err != nil { return nil, err }
    return &dto.UserResponse{ID: u.ID, Username: u.Username, Role: string(u.Role), Name: u.Name}, nil
}

func (s *services) GetUsers(ctx context.Context) ([]*dto.UserResponse, error) {
    users, err := s.repo.User().GetUsers(ctx)
    if err != nil { return nil, err }
    out := make([]*dto.UserResponse, 0, len(users))
    for _, u := range users {
        out = append(out, &dto.UserResponse{ID: u.ID, Username: u.Username, Role: string(u.Role), Name: u.Name})
    }
    return out, nil
}

func (s *services) UpdateUser(ctx context.Context, id int, req *dto.UserRequest) error {
    u, err := s.repo.User().GetUserByID(ctx, id)
    if err != nil { return err }
    if req.Username != "" { u.Username = req.Username }
    if req.Password != "" {
        hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        u.PasswordHash = string(hash)
    }
    if req.Role != "" { u.Role = models.Role(req.Role) }
    if req.Name != "" { u.Name = req.Name }
    return s.repo.User().UpdateUser(ctx, u)
}

func (s *services) DeleteUser(ctx context.Context, id int) error {
    return s.repo.User().DeleteUser(ctx, id)
}

func (s *services) GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error) {
    u, err := s.repo.User().GetUserByID(ctx, id)
    if err != nil { return nil, err }
    return &dto.UserResponse{ID: u.ID, Username: u.Username, Role: string(u.Role), Name: u.Name}, nil
}

// TASK

func (s *services) Task() TaskService { return s }

func (s *services) CreateTask(ctx context.Context, req *dto.TaskRequest, creatorID int) (*dto.TaskResponse, error) {
    if req.EmployeeID == 0 || req.Title == "" { return nil, errors.New("invalid input") }
    deadline, err := time.Parse(time.RFC3339, req.Deadline)
    if err != nil { return nil, errors.New("invalid deadline format") }
    t := &models.Task{
        EmployeeID: req.EmployeeID,
        CreatorID:  creatorID,
        Title: req.Title,
        Description: req.Description,
        Deadline: deadline,
        Status: models.TaskStatus(req.Status),
        Progress: req.Progress,
    }
    if t.Status == "" { t.Status = models.StatusPending }
    if err := s.repo.Task().CreateTask(ctx, t); err != nil { return nil, err }
    return &dto.TaskResponse{
        ID: t.ID, EmployeeID: t.EmployeeID, CreatorID: t.CreatorID,
        Title: t.Title, Description: t.Description, Deadline: t.Deadline,
        Status: string(t.Status), Progress: t.Progress, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
    }, nil
}

func (s *services) GetTaskByID(ctx context.Context, id int) (*dto.TaskResponse, error) {
    t, err := s.repo.Task().GetTaskByID(ctx, id)
    if err != nil { return nil, err }
    return &dto.TaskResponse{
        ID: t.ID, EmployeeID: t.EmployeeID, CreatorID: t.CreatorID,
        Title: t.Title, Description: t.Description, Deadline: t.Deadline,
        Status: string(t.Status), Progress: t.Progress, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
    }, nil
}

func (s *services) GetTasksByEmployeeID(ctx context.Context, employeeID int) ([]*dto.TaskResponse, error) {
    ts, err := s.repo.Task().GetTasksByEmployeeID(ctx, employeeID)
    if err != nil { return nil, err }
    out := make([]*dto.TaskResponse, 0, len(ts))
    for _, t := range ts {
        t2 := t
        out = append(out, &dto.TaskResponse{
            ID: t2.ID, EmployeeID: t2.EmployeeID, CreatorID: t2.CreatorID,
            Title: t2.Title, Description: t2.Description, Deadline: t2.Deadline,
            Status: string(t2.Status), Progress: t2.Progress, CreatedAt: t2.CreatedAt, UpdatedAt: t2.UpdatedAt,
        })
    }
    return out, nil
}

func (s *services) UpdateTask(ctx context.Context, id int, req *dto.TaskRequest, userID int) error {
    t, err := s.repo.Task().GetTaskByID(ctx, id)
    if err != nil { return err }
    if t.CreatorID != userID && t.EmployeeID != userID {
        return errors.New("forbidden")
    }

    oldStatus := t.Status

    if req.Title != "" { t.Title = req.Title }
    if req.Description != "" { t.Description = req.Description }
    if req.Status != "" { t.Status = models.TaskStatus(req.Status) }
    if req.Progress != 0 { t.Progress = req.Progress }
    if req.Deadline != "" {
        dl, err := time.Parse(time.RFC3339, req.Deadline)
        if err != nil { return errors.New("invalid deadline format") }
        t.Deadline = dl
    }

    if err := s.repo.Task().UpdateTask(ctx, t); err != nil {
        return err
    }

    if oldStatus != t.Status {
        h := &models.TaskStatusHistory{
            TaskID:    id,
            OldStatus: oldStatus,
            NewStatus: t.Status,
            ChangedBy: userID,
        }
        if err := s.repo.Task().CreateHistory(ctx, h); err != nil {
            s.logger.Error().Err(err).Msg("failed to record status history")
        }
    }

    return nil
}

func (s *services) DeleteTask(ctx context.Context, id int, userID int) error {
    t, err := s.repo.Task().GetTaskByID(ctx, id)
    if err != nil { return err }
    if t.CreatorID != userID && t.EmployeeID != userID {
        return errors.New("forbidden")
    }
    return s.repo.Task().DeleteTask(ctx, id)
}

// COMMENTS

func (s *services) Comment() CommentService { return s }

func (s *services) CreateComment(ctx context.Context, taskID int, userID int, text string) (*dto.CommentResponse, error) {
    c := &models.Comment{ TaskID: taskID, UserID: userID, Text: text }
    if err := s.repo.Comment().CreateComment(ctx, c); err != nil { return nil, err }
    return &dto.CommentResponse{ ID: c.ID, TaskID: c.TaskID, UserID: c.UserID, Text: c.Text, CreatedAt: c.CreatedAt }, nil
}

func (s *services) GetCommentsByTaskID(ctx context.Context, taskID int) ([]*dto.CommentResponse, error) {
    cs, err := s.repo.Comment().GetCommentsByTaskID(ctx, taskID)
    if err != nil { return nil, err }
    out := make([]*dto.CommentResponse, 0, len(cs))
    for _, c := range cs {
        c2 := c
        out = append(out, &dto.CommentResponse{ ID: c2.ID, TaskID: c2.TaskID, UserID: c2.UserID, Text: c2.Text, CreatedAt: c2.CreatedAt })
    }
    return out, nil
}

func (s *services) UpdateComment(ctx context.Context, id int, userID int, text string) error {
    c, err := s.repo.Comment().GetCommentByID(ctx, id)
    if err != nil { return err }
    if c.UserID != userID { return errors.New("forbidden") }
    c.Text = text
    return s.repo.Comment().UpdateComment(ctx, c)
}

func (s *services) DeleteComment(ctx context.Context, id int, userID int) error {
	c, err := s.repo.Comment().GetCommentByID(ctx, id)
	if err != nil {
		return err
	}
	if c.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Comment().DeleteComment(ctx, id)
}

func (s *services) UploadAttachment(ctx context.Context, taskID int, userID int, fileName string, filePath string, fileSize int64) (*dto.AttachmentResponse, error) {
	// Check if task exists and user has access
	t, err := s.repo.Task().GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	if t.CreatorID != userID && t.EmployeeID != userID {
		return nil, errors.New("forbidden")
	}

	f := &models.FileAttachment{
		TaskID:   taskID,
		FileName: fileName,
		FilePath: filePath,
		FileSize: fileSize,
	}

	if err := s.repo.File().CreateAttachment(ctx, f); err != nil {
		return nil, err
	}

	return &dto.AttachmentResponse{
		ID:         f.ID,
		TaskID:     f.TaskID,
		FileName:   f.FileName,
		FileSize:   f.FileSize,
		UploadedAt: f.UploadedAt,
	}, nil
}
func (s *services) GetTaskHistory(ctx context.Context, taskID int) ([]*dto.TaskHistoryResponse, error) {
	history, err := s.repo.Task().GetHistoryByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	var out []*dto.TaskHistoryResponse
	for _, h := range history {
		out = append(out, &dto.TaskHistoryResponse{
			ID:        h.ID,
			TaskID:    h.TaskID,
			OldStatus: string(h.OldStatus),
			NewStatus: string(h.NewStatus),
			ChangedBy: h.ChangedBy,
			CreatedAt: h.CreatedAt,
		})
	}
	return out, nil
}

func (s *services) ListTasks(ctx context.Context, filter dto.TaskFilter) ([]*dto.TaskResponse, error) {
	tasks, err := s.repo.Task().ListTasks(ctx, filter)
	if err != nil {
		return nil, err
	}

	var out []*dto.TaskResponse
	for _, t := range tasks {
		out = append(out, &dto.TaskResponse{
			ID:          t.ID,
			EmployeeID:  t.EmployeeID,
			CreatorID:   t.CreatorID,
			Title:       t.Title,
			Description: t.Description,
			Deadline:    t.Deadline,
			Status:      string(t.Status),
			Progress:    t.Progress,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	return out, nil
}
