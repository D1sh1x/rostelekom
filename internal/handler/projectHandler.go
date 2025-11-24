package handler

import (
	"SkillsTracker/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @Summary      Создать проект
// @Description  Создает новый проект. Только менеджеры могут создавать проекты. Создатель автоматически становится менеджером проекта. Можно сразу добавить участников проекта.
// @Tags         projects
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.ProjectRequest  true  "Данные проекта (name, description, status, member_ids)"
// @Success      201  {object}  dto.ProjectResponse
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /api/v1/projects [post]
func (h *Handler) CreateProject(c echo.Context) error {
	var req dto.ProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	managerID := c.Get("user_id").(int)

	resp, err := h.service.Project().CreateProject(c.Request().Context(), &req, managerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

// @Summary      Получить проект по id
// @Tags         projects
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  dto.ProjectResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/projects/{id} [get]
func (h *Handler) GetProjectByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	resp, err := h.service.Project().GetProjectByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
	}

	return c.JSON(http.StatusOK, resp)
}

// @Summary      Получить проекты
// @Description  Получить список проектов. Если указан manager_id, возвращаются проекты конкретного менеджера. Иначе возвращаются все проекты.
// @Tags         projects
// @Security     BearerAuth
// @Produce      json
// @Param        manager_id  query     int  false  "ID менеджера (для фильтрации по менеджеру)"
// @Success      200  {array}  dto.ProjectResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/projects [get]
func (h *Handler) GetProjects(c echo.Context) error {
	managerIDStr := c.QueryParam("manager_id")

	if managerIDStr != "" {
		managerID, err := strconv.Atoi(managerIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid manager_id"})
		}

		projects, err := h.service.Project().GetProjectsByManagerID(c.Request().Context(), managerID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, projects)
	}

	projects, err := h.service.Project().GetAllProjects(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, projects)
}

// @Summary      Обновить проект
// @Tags         projects
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      int                true  "Project ID"
// @Param        input  body      dto.ProjectRequest true  "Данные проекта"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /api/v1/projects/{id} [put]
func (h *Handler) UpdateProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req dto.ProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	managerID := c.Get("user_id").(int)

	if err := h.service.Project().UpdateProject(c.Request().Context(), id, &req, managerID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "project updated"})
}

// @Summary      Удалить проект
// @Tags         projects
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/projects/{id} [delete]
func (h *Handler) DeleteProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	managerID := c.Get("user_id").(int)

	if err := h.service.Project().DeleteProject(c.Request().Context(), id, managerID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "project deleted"})
}

// @Summary      Добавить участника в проект
// @Description  Добавляет пользователя в проект. Только менеджер проекта может добавлять участников.
// @Tags         projects
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Param        user_id  query     int  true  "User ID (ID пользователя для добавления)"
// @Param        role     query     string  false  "Роль в проекте (по умолчанию: developer)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /api/v1/projects/{id}/members [post]
func (h *Handler) AddProjectMember(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	role := c.QueryParam("role")
	if role == "" {
		role = "developer"
	}

	currentUserID := c.Get("user_id").(int)

	if err := h.service.Project().AddProjectMember(c.Request().Context(), projectID, userID, role, currentUserID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "member added"})
}

// @Summary      Удалить участника из проекта
// @Description  Удаляет пользователя из проекта. Менеджера проекта удалить нельзя.
// @Tags         projects
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Param        user_id  query     int  true  "User ID (ID пользователя для удаления)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /api/v1/projects/{id}/members [delete]
func (h *Handler) RemoveProjectMember(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	if err := h.service.Project().RemoveProjectMember(c.Request().Context(), projectID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "member removed"})
}

