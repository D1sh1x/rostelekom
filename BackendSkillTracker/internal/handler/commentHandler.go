package handler

import (
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
    "skilltracker/internal/dto"
)

// CreateComment godoc
// @Summary Create a new comment
// @Description Add a comment to a task
// @Tags comments
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param req body dto.CommentRequest true "Comment request"
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Router /comments [post]
func (h *Handler) CreateComment(c echo.Context) error {
	var req dto.CommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userID := c.Get("user_id").(int)
	res, err := h.service.Comment().CreateComment(c.Request().Context(), req.TaskID, userID, req.Text)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// GetCommentsByTaskID godoc
// @Summary Get comments by task ID
// @Description Retrieve all comments for a specific task
// @Tags comments
// @Security ApiKeyAuth
// @Produce json
// @Param task_id path int true "Task ID"
// @Success 200 {array} dto.CommentResponse
// @Router /tasks/{task_id}/comments [get]
func (h *Handler) GetCommentsByTaskID(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("task_id"))
	res, err := h.service.Comment().GetCommentsByTaskID(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateComment godoc
// @Summary Update comment
// @Description Update the text of a comment
// @Tags comments
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param req body dto.CommentRequest true "Update request"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /comments/{id} [put]
func (h *Handler) UpdateComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.CommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userID := c.Get("user_id").(int)
	if err := h.service.Comment().UpdateComment(c.Request().Context(), id, userID, req.Text); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "comment not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "updated"})
}

// DeleteComment godoc
// @Summary Delete comment
// @Description Remove a comment from a task
// @Tags comments
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /comments/{id} [delete]
func (h *Handler) DeleteComment(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    userID := c.Get("user_id").(int)
    if err := h.service.Comment().DeleteComment(c.Request().Context(), id, userID); err != nil {
        if err.Error() == "forbidden" { return c.JSON(http.StatusForbidden, map[string]string{"error":"forbidden"}) }
        return c.JSON(http.StatusNotFound, map[string]string{"error":"comment not found"})
    }
    return c.JSON(http.StatusOK, map[string]string{"message":"deleted"})
}
