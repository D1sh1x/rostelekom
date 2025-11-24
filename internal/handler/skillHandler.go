package handler

import (
	"SkillsTracker/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @Summary      Создать навык
// @Description  Создает новый навык в системе. Навыки используются для сопоставления с задачами и пользователями.
// @Tags         skills
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.SkillRequest  true  "Данные навыка (name, description, category)"
// @Success      201  {object}  dto.SkillResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/skills [post]
func (h *Handler) CreateSkill(c echo.Context) error {
	var req dto.SkillRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.Skill().CreateSkill(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

// @Summary      Получить навык по id
// @Tags         skills
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Skill ID"
// @Success      200  {object}  dto.SkillResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/skills/{id} [get]
func (h *Handler) GetSkillByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	resp, err := h.service.Skill().GetSkillByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "skill not found"})
	}

	return c.JSON(http.StatusOK, resp)
}

// @Summary      Получить навыки
// @Description  Получить список всех навыков или навыков по категории. Если указана категория, возвращаются только навыки этой категории.
// @Tags         skills
// @Security     BearerAuth
// @Produce      json
// @Param        category  query     string  false  "Категория навыка (например: technical, soft, domain)"
// @Success      200  {array}  dto.SkillResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/skills [get]
func (h *Handler) GetSkills(c echo.Context) error {
	category := c.QueryParam("category")

	if category != "" {
		skills, err := h.service.Skill().GetSkillsByCategory(c.Request().Context(), category)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, skills)
	}

	skills, err := h.service.Skill().GetAllSkills(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, skills)
}

// @Summary      Обновить навык
// @Tags         skills
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Skill ID"
// @Param        input  body      dto.SkillRequest true  "Данные навыка"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /api/v1/skills/{id} [put]
func (h *Handler) UpdateSkill(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req dto.SkillRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.service.Skill().UpdateSkill(c.Request().Context(), id, &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "skill updated"})
}

// @Summary      Удалить навык
// @Tags         skills
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Skill ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/skills/{id} [delete]
func (h *Handler) DeleteSkill(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.service.Skill().DeleteSkill(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "skill deleted"})
}

// @Summary      Добавить навык пользователю
// @Description  Добавляет навык пользователю с указанным уровнем (1-5). Если навык уже есть у пользователя, обновляется уровень.
// @Tags         skills
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Param        input    body      dto.UserSkillRequest true  "Данные навыка (skill_id, level от 1 до 5)"
// @Success      200  {object}  dto.UserSkillResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/users/{user_id}/skills [post]
func (h *Handler) AddUserSkill(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	var req dto.UserSkillRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.Skill().AddUserSkill(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// @Summary      Получить навыки пользователя
// @Tags         skills
// @Security     BearerAuth
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Success      200  {array}  dto.UserSkillResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/users/{user_id}/skills [get]
func (h *Handler) GetUserSkills(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	skills, err := h.service.Skill().GetUserSkills(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, skills)
}

// @Summary      Обновить уровень навыка пользователя
// @Description  Обновляет уровень навыка пользователя (1-5). Навык должен уже быть добавлен пользователю.
// @Tags         skills
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Param        skill_id path      int  true  "Skill ID"
// @Param        level    query     int  true  "Уровень навыка (1-5)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/users/{user_id}/skills/{skill_id} [put]
func (h *Handler) UpdateUserSkill(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	skillID, err := strconv.Atoi(c.Param("skill_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid skill_id"})
	}

	level, err := strconv.Atoi(c.QueryParam("level"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid level"})
	}

	if err := h.service.Skill().UpdateUserSkill(c.Request().Context(), userID, skillID, level); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user skill updated"})
}

// @Summary      Удалить навык у пользователя
// @Tags         skills
// @Security     BearerAuth
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Param        skill_id path      int  true  "Skill ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/users/{user_id}/skills/{skill_id} [delete]
func (h *Handler) RemoveUserSkill(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	skillID, err := strconv.Atoi(c.Param("skill_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid skill_id"})
	}

	if err := h.service.Skill().RemoveUserSkill(c.Request().Context(), userID, skillID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user skill removed"})
}
