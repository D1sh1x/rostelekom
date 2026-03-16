package transport

import (
	"net/http"
	"skilltracker/internal/config"
	"skilltracker/internal/handler"
	m "skilltracker/internal/middleware"

	"context"
	"os"
	"os/signal"
	_ "skilltracker/docs"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewServer(jwtSecret []byte, h *handler.Handler, cfg *config.Config) *http.Server {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Public
	v1.POST("/login", h.Login)
	v1.POST("/refresh", h.RefreshToken)
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Protected
	auth := v1.Group("")
	auth.Use(m.AuthRequired([]byte(cfg.Auth.JWTSecret)))

	auth.POST("/logout", h.Logout)

	// Users (only manager)
	managerOnly := m.RoleRequired("manager")
	auth.POST("/users", h.CreateUser)
	auth.GET("/users", h.GetUsers, managerOnly)
	auth.GET("/users/:id", h.GetUserByID, managerOnly)
	auth.PUT("/users/:id", h.UpdateUser, managerOnly)
	auth.DELETE("/users/:id", h.DeleteUser, managerOnly)

	// Tasks
	auth.POST("/tasks", h.CreateTask, managerOnly)
	auth.GET("/tasks", h.ListTasks)
	auth.GET("/tasks/my", h.GetMyTasks)
	auth.GET("/tasks/:id", h.GetTaskByID)
	auth.PUT("/tasks/:id", h.UpdateTask)
	auth.DELETE("/tasks/:id", h.DeleteTask)
	auth.POST("/tasks/:id/attachments", h.UploadAttachment)
	auth.GET("/tasks/:id/history", h.GetTaskHistory)

	// Comments
	auth.POST("/comments", h.CreateComment)
	auth.GET("/tasks/:task_id/comments", h.GetCommentsByTaskID)
	auth.PUT("/comments/:id", h.UpdateComment)
	auth.DELETE("/comments/:id", h.DeleteComment)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return &http.Server{
		Addr:         cfg.HTTPServer.Port,
		Handler:      e,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
}

func Run(s *http.Server) error {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}
