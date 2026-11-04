package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"teman-mancing/backend/internal/application/auth"
	appchat "teman-mancing/backend/internal/application/chat"
	appevent "teman-mancing/backend/internal/application/event"
	appspot "teman-mancing/backend/internal/application/fishing_spot"
	appnotif "teman-mancing/backend/internal/application/notification"
	appuser "teman-mancing/backend/internal/application/user"
	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/handlers"
	"teman-mancing/backend/internal/delivery/http/middleware"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/email"
	"teman-mancing/backend/internal/infrastructure/security"
	"teman-mancing/backend/internal/infrastructure/storage"
	infraWS "teman-mancing/backend/internal/infrastructure/websocket"
)

// RouterDeps bundles router dependencies.
type RouterDeps struct {
	Config           *config.Config
	Logger           *zap.Logger
	DB               *database.PostgresDB
	MinIO            *storage.MinIOClient
	Mailer           *email.Mailer
	TokenService     *security.TokenService
	UserRepo         appuser.Repository
	SpotRepo         appspot.Repository
	EventRepo        appevent.Repository
	ChatRepo         appchat.Repository
	NotificationRepo appnotif.Repository
	Version          string
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

	spotRepo := deps.SpotRepo
	if spotRepo == nil {
		if deps.DB != nil && deps.DB.DB != nil {
			spotRepo = database.NewGormSpotRepository(deps.DB.DB)
		} else {
			spotRepo = database.NewInMemorySpotRepository()
		}
	}

	chatRepo := deps.ChatRepo
	if chatRepo == nil {
		if deps.DB != nil && deps.DB.DB != nil {
			chatRepo = database.NewGormChatRepository(deps.DB.DB)
		} else {
			chatRepo = database.NewInMemoryChatRepository()
		}
	}

	eventRepo := deps.EventRepo
	if eventRepo == nil {
		if deps.DB != nil && deps.DB.DB != nil {
			eventRepo = database.NewGormEventRepository(deps.DB.DB)
		} else {
			eventRepo = database.NewInMemoryEventRepository(chatRepo)
		}
	} else if inMemEv, ok := eventRepo.(*database.InMemoryEventRepository); ok {
		inMemEv.SetChatRepository(chatRepo)
	}

	notifRepo := deps.NotificationRepo
	if notifRepo == nil {
		if deps.DB != nil && deps.DB.DB != nil {
			notifRepo = database.NewGormNotificationRepository(deps.DB.DB)
		} else {
			notifRepo = database.NewInMemoryNotificationRepository()
		}
	}

	tokenService := deps.TokenService
	if tokenService == nil {
		tokenService = security.NewTokenService(deps.Config.JWTSecret, deps.Config.JWTExpiration)
	}

	authService := auth.NewAuthService(userRepo, tokenService, deps.Mailer)
	authHandler := handlers.NewAuthHandler(
		authService,
		deps.Config.GoogleClientID,
		deps.Config.GoogleClientSecret,
		deps.Config.GoogleRedirectURL,
	)
	userHandler := handlers.NewUserHandler(authService, deps.MinIO)
	mediaHandler := handlers.NewMediaHandler(deps.MinIO)

	notifService := appnotif.NewService(notifRepo)
	notifHandler := handlers.NewNotificationHandler(notifService)

	spotService := appspot.NewService(spotRepo)
	spotHandler := handlers.NewSpotHandler(spotService, deps.MinIO, notifService)

	eventService := appevent.NewService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	chatService := appchat.NewService(chatRepo)
	wsHub := infraWS.NewHub(chatService, deps.Logger)
	chatHandler := handlers.NewChatHandler(chatService, eventRepo, wsHub)

	// API v1 Group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Status)
		v1.HEAD("/health", healthHandler.Status)
		v1.GET("/version", versionHandler.Info)
		v1.HEAD("/version", versionHandler.Info)

		// Public Media / Storage Proxy
		v1.GET("/media/:bucket/*objectKey", mediaHandler.ServeMedia)

		authMiddleware := middleware.RequireAuth(tokenService)
		verifiedMiddleware := middleware.RequireVerified(userRepo)

		// Auth Routes (README Section 13)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/verify-email", authHandler.VerifyEmail)
			authGroup.POST("/resend-verification", authHandler.ResendVerification)
			authGroup.POST("/forgot-password", authHandler.ForgotPassword)
			authGroup.POST("/reset-password", authHandler.ResetPassword)
			authGroup.POST("/google", authHandler.GoogleLogin)
			authGroup.GET("/google/login", authHandler.GoogleOAuthRedirect)
			authGroup.GET("/google/callback", authHandler.GoogleOAuthCallback)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/me", authMiddleware, authHandler.GetMe)
		}

		// User Routes (README Section 13)
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/:username", userHandler.GetByUsername)
			userGroup.PATCH("/me", authMiddleware, verifiedMiddleware, userHandler.UpdateProfile)
			userGroup.POST("/me/avatar", authMiddleware, verifiedMiddleware, userHandler.UploadAvatar)
		}

		// Fishing Spots Routes (README Section 13)
		spotsGroup := v1.Group("/spots")
		{
			spotsGroup.GET("", spotHandler.List)
			spotsGroup.POST("", authMiddleware, verifiedMiddleware, spotHandler.Create)
			spotsGroup.GET("/:id", spotHandler.GetByID)
			spotsGroup.PATCH("/:id", authMiddleware, verifiedMiddleware, spotHandler.Update)
			spotsGroup.DELETE("/:id", authMiddleware, verifiedMiddleware, spotHandler.Delete)
			spotsGroup.POST("/:id/photos", authMiddleware, verifiedMiddleware, spotHandler.UploadPhoto)
			spotsGroup.DELETE("/:id/photos/:photoId", authMiddleware, verifiedMiddleware, spotHandler.DeletePhoto)

			// Comments
			spotsGroup.GET("/:id/comments", spotHandler.ListComments)
			spotsGroup.POST("/:id/comments", authMiddleware, verifiedMiddleware, spotHandler.CreateComment)
			spotsGroup.POST("/:id/comments/photo", authMiddleware, verifiedMiddleware, spotHandler.UploadCommentPhoto)
		}

		// Comments direct endpoint (README Section 13)
		v1.DELETE("/comments/:id", authMiddleware, verifiedMiddleware, spotHandler.DeleteComment)

		// Events Routes (README Section 13)
		eventsGroup := v1.Group("/events")
		{
			eventsGroup.GET("", eventHandler.List)
			eventsGroup.POST("", authMiddleware, verifiedMiddleware, eventHandler.Create)
			eventsGroup.GET("/:id", eventHandler.GetByID)
			eventsGroup.PATCH("/:id", authMiddleware, verifiedMiddleware, eventHandler.Update)
			eventsGroup.DELETE("/:id", authMiddleware, verifiedMiddleware, eventHandler.Delete)

			// Join / Leave
			eventsGroup.POST("/:id/join", authMiddleware, verifiedMiddleware, eventHandler.Join)
			eventsGroup.POST("/:id/leave", authMiddleware, verifiedMiddleware, eventHandler.Leave)

			// Manage Requests
			eventsGroup.GET("/:id/requests", authMiddleware, eventHandler.ListRequests)
			eventsGroup.POST("/:id/requests/:userId/approve", authMiddleware, verifiedMiddleware, eventHandler.ApproveRequest)
			eventsGroup.POST("/:id/requests/:userId/reject", authMiddleware, verifiedMiddleware, eventHandler.RejectRequest)
		}

		// Conversations & Chat Routes (README Section 13)
		chatGroup := v1.Group("/conversations")
		{
			chatGroup.GET("", authMiddleware, chatHandler.ListConversations)
			chatGroup.POST("", authMiddleware, verifiedMiddleware, chatHandler.CreateConversation)
			chatGroup.GET("/:id/messages", authMiddleware, chatHandler.ListMessages)
			chatGroup.POST("/:id/messages", authMiddleware, verifiedMiddleware, chatHandler.SendMessage)
		}

		// Realtime WebSocket Endpoint (README Section 13 & 9.4)
		v1.GET("/ws", authMiddleware, chatHandler.WebSocket)

		// Notifications Routes (README Section 13)
		notifGroup := v1.Group("/notifications")
		{
			notifGroup.GET("", authMiddleware, notifHandler.List)
			notifGroup.PATCH("/:id/read", authMiddleware, notifHandler.MarkAsRead)
			notifGroup.POST("/read-all", authMiddleware, notifHandler.MarkAllAsRead)
		}
	}

	// Also support root /ws for WebSocket clients pointing to host:port/ws
	r.GET("/ws", middleware.RequireAuth(tokenService), chatHandler.WebSocket)

	return r
}
