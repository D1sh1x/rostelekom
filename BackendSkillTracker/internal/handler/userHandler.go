package handler

import (
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
    "skilltracker/internal/dto"
)

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access and refresh tokens using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param req body dto.RefreshRequest true "Refresh request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /refresh [post]
func (h *Handler) RefreshToken(c echo.Context) error {
	var req dto.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	res, err := h.service.User().RefreshToken(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate the refresh token
// @Tags auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /logout [post]
func (h *Handler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if err := h.service.User().Logout(c.Request().Context(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param req body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	res, err := h.service.User().Login(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}
	return c.JSON(http.StatusOK, res)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user (Manager only)
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param req body dto.UserRequest true "User request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var req dto.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	u, err := h.service.User().CreateUser(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, u)
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users (Manager only)
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Router /users [get]
func (h *Handler) GetUsers(c echo.Context) error {
    users, err := h.service.User().GetUsers(c.Request().Context())
    if err != nil { return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()}) }
    return c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve details of a specific user (Manager only)
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *Handler) GetUserByID(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    u, err := h.service.User().GetUserByID(c.Request().Context(), id)
    if err != nil { return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"}) }
    return c.JSON(http.StatusOK, u)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details (Manager only)
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param req body dto.UserRequest true "Update request"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	if err := h.validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.service.User().UpdateUser(c.Request().Context(), id, &req); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "updated"})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Remove a user from the system (Manager only)
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.service.User().DeleteUser(c.Request().Context(), id); err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
    }
    return c.JSON(http.StatusOK, map[string]string{"message": "deleted"})
}
