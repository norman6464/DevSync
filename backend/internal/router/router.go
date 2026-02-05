package router

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/middleware"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	origins := strings.Split(cfg.CORSOrigins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, userRepo)
	userHandler := handler.NewUserHandler(userRepo)

	// Public routes
	r.GET("/health", handler.HealthCheck)

	api := r.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(authService))
	{
		// Auth
		protected.GET("/auth/me", authHandler.Me)

		// Users
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
		}
	}

	return r
}
