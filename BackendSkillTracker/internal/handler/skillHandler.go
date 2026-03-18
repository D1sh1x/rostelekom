package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"skilltracker/internal/dto"
)

// CreateSkill godoc
// @Summary Create a skill
// @Tags skills
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param req body dto.SkillRequest true "Skill request"
// @Success 201 {object} dto.SkillResponse
// @Router /skills [post]
func (h *Handler) CreateSkill(c echo.Context) error {
	var req dto.SkillRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	res, err := h.service.Skill().CreateSkill(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, res)
}

// GetSkills godoc
// @Summary List all skills
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.SkillResponse
// @Router /skills [get]
func (h *Handler) GetSkills(c echo.Context) error {
	res, err := h.service.Skill().GetSkills(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// DeleteSkill godoc
// @Summary Delete a skill
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Skill ID"
// @Success 200 {object} map[string]string
// @Router /skills/{id} [delete]
func (h *Handler) DeleteSkill(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Skill().DeleteSkill(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "skill not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "deleted"})
}

// AssignSkillToUser godoc
// @Summary Assign skill to user
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Param skill_id path int true "Skill ID"
// @Success 200 {object} map[string]string
// @Router /users/{id}/skills/{skill_id} [post]
func (h *Handler) AssignSkillToUser(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	skillID, _ := strconv.Atoi(c.Param("skill_id"))
	if err := h.service.Skill().AssignSkillToUser(c.Request().Context(), userID, skillID); err != nil {
		if err.Error() == "user not found" || err.Error() == "skill not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "skill assigned"})
}

// RemoveSkillFromUser godoc
// @Summary Remove skill from user
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Param skill_id path int true "Skill ID"
// @Success 200 {object} map[string]string
// @Router /users/{id}/skills/{skill_id} [delete]
func (h *Handler) RemoveSkillFromUser(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	skillID, _ := strconv.Atoi(c.Param("skill_id"))
	if err := h.service.Skill().RemoveSkillFromUser(c.Request().Context(), userID, skillID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "skill removed"})
}

// GetUserSkills godoc
// @Summary Get user skills
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} dto.SkillResponse
// @Router /users/{id}/skills [get]
func (h *Handler) GetUserSkills(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	res, err := h.service.Skill().GetUserSkills(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// AddSkillToTask godoc
// @Summary Add required skill to task
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Param skill_id path int true "Skill ID"
// @Success 200 {object} map[string]string
// @Router /tasks/{id}/skills/{skill_id} [post]
func (h *Handler) AddSkillToTask(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("id"))
	skillID, _ := strconv.Atoi(c.Param("skill_id"))
	userID := c.Get("user_id").(int)
	if err := h.service.Task().AddSkillToTask(c.Request().Context(), taskID, skillID, userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		if err.Error() == "task not found" || err.Error() == "skill not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "skill added to task"})
}

// RemoveSkillFromTask godoc
// @Summary Remove required skill from task
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Param skill_id path int true "Skill ID"
// @Success 200 {object} map[string]string
// @Router /tasks/{id}/skills/{skill_id} [delete]
func (h *Handler) RemoveSkillFromTask(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("id"))
	skillID, _ := strconv.Atoi(c.Param("skill_id"))
	userID := c.Get("user_id").(int)
	if err := h.service.Task().RemoveSkillFromTask(c.Request().Context(), taskID, skillID, userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "skill removed from task"})
}

// GetTaskSkills godoc
// @Summary Get task required skills
// @Tags skills
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {array} dto.SkillResponse
// @Router /tasks/{id}/skills [get]
func (h *Handler) GetTaskSkills(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("id"))
	res, err := h.service.Task().GetTaskSkills(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// GetRecommendedEmployees godoc
// @Summary Get recommended employees for task
// @Description Returns employees sorted by skill match score
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {array} dto.RecommendedEmployeeResponse
// @Router /tasks/{id}/recommended-employees [get]
func (h *Handler) GetRecommendedEmployees(c echo.Context) error {
	taskID, _ := strconv.Atoi(c.Param("id"))
	res, err := h.service.Task().GetRecommendedEmployees(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
