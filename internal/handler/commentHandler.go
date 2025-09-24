package handler

import (
	"SkillsTracker/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateComment(c echo.Context) error {
	var req dto.CommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID := c.Get("user_id").(int)

	resp, err := h.service.Comment().CreateComment(c.Request().Context(), &req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetCommentByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	resp, err := h.service.Comment().GetCommentByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "comment not found"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCommentsByTaskID(c echo.Context) error {
	taskID, err := strconv.Atoi(c.QueryParam("task_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task_id"})
	}

	comments, err := h.service.Comment().GetCommentsByTaskID(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, comments)
}

func (h *Handler) UpdateComment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req dto.CommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID := c.Get("user_id").(int)

	if err := h.service.Comment().UpdateComment(c.Request().Context(), id, &req, userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "comment not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "comment updated"})
}

func (h *Handler) DeleteComment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	userID := c.Get("user_id").(int)

	if err := h.service.Comment().DeleteComment(c.Request().Context(), id, userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "comment not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "comment deleted"})
}
