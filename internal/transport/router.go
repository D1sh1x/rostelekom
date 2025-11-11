package transport

import (
	"SkillsTracker/docs"
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
				tasks.GET("", h.GetTasksByEmployeeID)
				tasks.PUT("/:id", h.UpdateTask)
				tasks.DELETE("/:id", h.DeleteTask)

				comments := protected.Group("/comments")
				comments.POST("", h.CreateComment)
				comments.GET("/:id", h.GetCommentByID)
				comments.GET("", h.GetCommentsByTaskID)
				comments.PUT("/:id", h.UpdateComment)
				comments.DELETE("/:id", h.DeleteComment)

				users := protected.Group("/users")
				users.GET("", h.GetUsers, customMiddleware.RequireRole("manager"))
				users.GET("/:id", h.GetUserByID, customMiddleware.RequireRole("manager"))
				users.PUT("/:id", h.UpdateUser, customMiddleware.RequireRole("manager"))
				users.DELETE("/:id", h.DeleteUser, customMiddleware.RequireRole("manager"))
			}
		}
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	return &http.Server{
		Addr:         config.HTTPServer.Port,
		Handler:      router,
		WriteTimeout: config.HTTPServer.WriteTimeout,
		ReadTimeout:  config.HTTPServer.ReadTimeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}
}
