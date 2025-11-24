package transport

import (
	"SkillsTracker/internal/config"
	"SkillsTracker/internal/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	customMiddleware "SkillsTracker/internal/middleware"
)

func NewServer(jwtSecret []byte, h *handler.Handler, config *config.Config) *http.Server {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
	}))

	swaggerHandler := echoSwagger.WrapHandler
	router.GET("/swagger/*", swaggerHandler)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/login", h.Login)
			v1.POST("/register", h.CreateUser)

			protected := v1.Group("")
			protected.Use(customMiddleware.AuthRequired(jwtSecret))
			{

				tasks := protected.Group("/tasks")
				tasks.POST("", h.CreateTask)
				tasks.GET("/:id", h.GetTaskByID)
				tasks.GET("", h.GetTasks)
				tasks.PUT("/:id", h.UpdateTask)
				tasks.DELETE("/:id", h.DeleteTask)

				comments := protected.Group("/comments")
				comments.POST("", h.CreateComment)
				comments.GET("/:id", h.GetCommentByID)
				comments.GET("", h.GetCommentsByTaskID)
				comments.PUT("/:id", h.UpdateComment)
				comments.DELETE("/:id", h.DeleteComment)

				projects := protected.Group("/projects")
				projects.POST("", h.CreateProject)
				projects.GET("/:id", h.GetProjectByID)
				projects.GET("", h.GetProjects)
				projects.PUT("/:id", h.UpdateProject)
				projects.DELETE("/:id", h.DeleteProject)
				projects.POST("/:id/members", h.AddProjectMember)
				projects.DELETE("/:id/members", h.RemoveProjectMember)

				skills := protected.Group("/skills")
				skills.POST("", h.CreateSkill)
				skills.GET("/:id", h.GetSkillByID)
				skills.GET("", h.GetSkills)
				skills.PUT("/:id", h.UpdateSkill)
				skills.DELETE("/:id", h.DeleteSkill)

				users := protected.Group("/users")
				users.GET("", h.GetUsers, customMiddleware.RequireRole("manager"))
				users.GET("/:id", h.GetUserByID, customMiddleware.RequireRole("manager"))
				users.PUT("/:id", h.UpdateUser, customMiddleware.RequireRole("manager"))
				users.DELETE("/:id", h.DeleteUser, customMiddleware.RequireRole("manager"))
				users.POST("/:user_id/skills", h.AddUserSkill)
				users.GET("/:user_id/skills", h.GetUserSkills)
				users.PUT("/:user_id/skills/:skill_id", h.UpdateUserSkill)
				users.DELETE("/:user_id/skills/:skill_id", h.RemoveUserSkill)
			}
		}
	}

	return &http.Server{
		Addr:         config.HTTPServer.Port,
		Handler:      router,
		WriteTimeout: config.HTTPServer.WriteTimeout,
		ReadTimeout:  config.HTTPServer.ReadTimeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}
}
