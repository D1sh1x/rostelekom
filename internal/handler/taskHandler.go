package handler

import (
	"SkillsTracker/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @Summary      Создать задачу
// @Description  Создает новую задачу в проекте. Менеджер проекта может создавать задачи. При создании можно указать необходимые навыки и назначить сотрудников (с проверкой наличия навыков). Можно указать родительскую задачу для создания зависимостей.
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.TaskRequest  true  "Данные задачи (project_id, title, description, deadline, hours, priority, type, parent_task_id, skill_ids, assignee_ids)"
// @Success      201  {object}  dto.TaskResponse
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /api/v1/tasks [post]
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

// @Summary      Получить задачу по id
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  dto.TaskResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/tasks/{id} [get]
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

// @Summary      Получить задачи
// @Description  Получить список задач по project_id или user_id. Необходимо указать один из параметров.
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        project_id  query     int  false  "ID проекта (для получения всех задач проекта)"
// @Param        user_id     query     int  false  "ID пользователя (для получения всех задач пользователя)"
// @Success      200  {array}  dto.TaskResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/tasks [get]
func (h *Handler) GetTasks(c echo.Context) error {
	projectIDStr := c.QueryParam("project_id")
	userIDStr := c.QueryParam("user_id")

	if projectIDStr != "" {
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project_id"})
		}

		tasks, err := h.service.Task().GetTasksByProjectID(c.Request().Context(), projectID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, tasks)
	}

	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		}

		tasks, err := h.service.Task().GetTasksByUserID(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, tasks)
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"error": "project_id or user_id required"})
}

// @Summary      Обновить задачу
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Task ID"
// @Param        input  body      dto.TaskRequest true  "Данные задачи"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /api/v1/tasks/{id} [put]
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

// @Summary      Удалить задачу
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/tasks/{id} [delete]
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
