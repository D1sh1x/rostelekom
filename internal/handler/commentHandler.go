package handler

import (
	"SkillsTracker/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @Summary      Создать комментарий
// @Description  Создает новый комментарий к задаче. Автором комментария становится текущий авторизованный пользователь.
// @Tags         comments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.CommentRequest  true  "Данные комментария (task_id, text)"
// @Success      201  {object}  dto.CommentResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/comments [post]
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

// @Summary      Получить комментарий по id
// @Tags         comments
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  dto.CommentResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/comments/{id} [get]
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

// @Summary      Получить комментарии по task_id
// @Description  Возвращает все комментарии к указанной задаче, отсортированные по дате создания.
// @Tags         comments
// @Security     BearerAuth
// @Produce      json
// @Param        task_id  query     int  true  "ID задачи"
// @Success      200  {array}  dto.CommentResponse
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/comments [get]
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

// @Summary      Обновить комментарий
// @Tags         comments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      int                 true  "Comment ID"
// @Param        input  body      dto.CommentRequest  true  "Данные комментария"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /api/v1/comments/{id} [put]
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

// @Summary      Удалить комментарий
// @Tags         comments
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/comments/{id} [delete]
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
