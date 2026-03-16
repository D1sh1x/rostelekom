package handler

import (
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
    "skilltracker/internal/dto"
    "io"
    "os"
    "path/filepath"
    "time"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Assign a new task to an employee
// @Tags tasks
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param req body dto.TaskRequest true "Task request"
// @Success 200 {object} dto.TaskResponse
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func (h *Handler) CreateTask(c echo.Context) error {
	var req dto.TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userID := c.Get("user_id").(int)
	res, err := h.service.Task().CreateTask(c.Request().Context(), &req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Retrieve details of a specific task
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.TaskResponse
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *Handler) GetTaskByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.service.Task().GetTaskByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}
	return c.JSON(http.StatusOK, res)
}

// GetMyTasks godoc
// @Summary Get my tasks
// @Description Retrieve tasks assigned to or created by the current user
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.TaskResponse
// @Router /tasks/my [get]
func (h *Handler) GetMyTasks(c echo.Context) error {
	userID := c.Get("user_id").(int)
	res, err := h.service.Task().GetTasksByEmployeeID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateTask godoc
// @Summary Update task
// @Description Update task details (status, progress, etc.)
// @Tags tasks
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param req body dto.TaskRequest true "Update request"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [put]
func (h *Handler) UpdateTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userID := c.Get("user_id").(int)
	if err := h.service.Task().UpdateTask(c.Request().Context(), id, &req, userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "updated"})
}

// DeleteTask godoc
// @Summary Delete task
// @Description Remove a task from the system
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteTask(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    userID := c.Get("user_id").(int)
    if err := h.service.Task().DeleteTask(c.Request().Context(), id, userID); err != nil {
        if err.Error() == "forbidden" { return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"}) }
        return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
    }
    return c.JSON(http.StatusOK, map[string]string{"message":"deleted"})
}

// UploadAttachment godoc
// @Summary Upload task attachment
// @Description Upload a file and attach it to a task
// @Tags tasks
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Task ID"
// @Param file formData file true "File to upload"
// @Success 200 {object} dto.AttachmentResponse
// @Failure 400 {object} map[string]string
// @Router /tasks/{id}/attachments [post]
func (h *Handler) UploadAttachment(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("id"))
	userID := c.Get("user_id").(int)

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid file"})
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Ensure upload dir exists
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// Create unique filename
	dstPath := filepath.Join(uploadDir, strconv.Itoa(int(time.Now().Unix()))+"_"+file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	res, err := h.service.Task().UploadAttachment(c.Request().Context(), taskID, userID, file.Filename, dstPath, file.Size)
	if err != nil {
		os.Remove(dstPath)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// GetTaskHistory godoc
// @Summary Get task status history
// @Description Retrieve a list of all status changes for a task
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {array} dto.TaskHistoryResponse
// @Router /tasks/{id}/history [get]
func (h *Handler) GetTaskHistory(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("id"))
	res, err := h.service.Task().GetTaskHistory(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// ListTasks godoc
// @Summary List tasks with filters
// @Description Retrieve a list of tasks with filtering options
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param status query string false "Status filter"
// @Param employee_id query int false "Employee ID filter"
// @Param creator_id query int false "Creator ID filter"
// @Param search query string false "Search term"
// @Param from_date query string false "From date (YYYY-MM-DD)"
// @Param to_date query string false "To date (YYYY-MM-DD)"
// @Success 200 {array} dto.TaskResponse
// @Router /tasks [get]
func (h *Handler) ListTasks(c echo.Context) error {
	var filter dto.TaskFilter
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid query params"})
	}

	res, err := h.service.Task().ListTasks(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
