package handler

import (
	"SkillsTracker/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateTask(c echo.Context) error {
	var req dto.TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID := c.Get("user_id").(int)

	resp, err := h.service.Task().CreateTask(c.Request().Context(), &req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	resp, err := h.service.Task().GetTaskByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTasksByEmployeeID(c echo.Context) error {
	employeeID, err := strconv.Atoi(c.QueryParam("employee_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid employee_id"})
	}

	tasks, err := h.service.Task().GetTasksByEmployeeID(c.Request().Context(), employeeID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *Handler) UpdateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req dto.TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID := c.Get("user_id").(int)

	if err := h.service.Task().UpdateTask(c.Request().Context(), id, &req, userID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task updated"})
}

func (h *Handler) DeleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	userID := c.Get("user_id").(int)

	if err := h.service.Task().DeleteTask(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task deleted"})
}
