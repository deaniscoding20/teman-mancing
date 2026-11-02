package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"teman-mancing/backend/internal/application/auth"
	appuser "teman-mancing/backend/internal/application/user"
	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/handlers"
	"teman-mancing/backend/internal/delivery/http/middleware"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
	"teman-mancing/backend/internal/infrastructure/storage"
)

// RouterDeps bundles router dependencies.
type RouterDeps struct {
	Config       *config.Config
	Logger       *zap.Logger
	DB           *database.PostgresDB
	MinIO        *storage.MinIOClient
	TokenService *security.TokenService
	UserRepo     appuser.Repository
	Version      string
}

// SetupRouter initializes the Gin engine and configures all routes.
func SetupRouter(deps RouterDeps) *gin.Engine {
	if deps.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Global Middlewares
	r.Use(middleware.ZapLogger(deps.Logger))
	r.Use(middleware.ZapRecovery(deps.Logger))
	r.Use(middleware.CORS(deps.Config.AllowedOrigins))

	// Root ping
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":     "Teman Mancing API",
			"status":  "running",
			"docs":    "/api/v1/version",
			"tagline": "Cari Spot. Cari Teman. Gas Mancing!",
		})
	})

	// Handlers
	healthHandler := handlers.NewHealthHandler(deps.DB, deps.MinIO)
	versionHandler := handlers.NewVersionHandler(deps.Version, deps.Config.AppEnv)

	userRepo := deps.UserRepo
	if userRepo == nil {
		if deps.DB != nil && deps.DB.DB != nil {
			userRepo = database.NewGormUserRepository(deps.DB.DB)
		} else {
			userRepo = database.NewInMemoryUserRepository()
		}
	}

	tokenService := deps.TokenService
	if tokenService == nil {
		tokenService = security.NewTokenService(deps.Config.JWTSecret, deps.Config.JWTExpiration)
	}

	authService := auth.NewAuthService(userRepo, tokenService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService, deps.MinIO)

	// API v1 Group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Status)
		v1.GET("/version", versionHandler.Info)

		// Auth Routes (README Section 13)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/me", middleware.RequireAuth(tokenService), authHandler.GetMe)
		}

		// User Routes (README Section 13)
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/:username", userHandler.GetByUsername)
			userGroup.PATCH("/me", middleware.RequireAuth(tokenService), userHandler.UpdateProfile)
			userGroup.POST("/me/avatar", middleware.RequireAuth(tokenService), userHandler.UploadAvatar)
		}

		// Feature route groups
		_ = v1.Group("/spots")
		_ = v1.Group("/events")
		_ = v1.Group("/conversations")
		_ = v1.Group("/notifications")
	}

	return r
}
