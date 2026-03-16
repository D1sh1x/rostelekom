package middleware

import (
    "net/http"
    "strings"
    "github.com/labstack/echo/v4"
    "skilltracker/internal/utils/jwt"
)

func AuthRequired(jwtSecret []byte) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if !strings.HasPrefix(authHeader, "Bearer ") {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing bearer"})
            }
            tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := jwt.ValidateToken(tokenStr, jwtSecret)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
            }
            c.Set("user_id", claims.UserID)
            c.Set("username", claims.Username)
            c.Set("role", claims.Role)
            return next(c)
        }
    }
}

func RoleRequired(roles ...string) echo.MiddlewareFunc {
    allowed := map[string]struct{}{}
    for _, r := range roles { allowed[r] = struct{}{} }
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            role, _ := c.Get("role").(string)
            if _, ok := allowed[role]; !ok {
                return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
            }
            return next(c)
        }
    }
}
