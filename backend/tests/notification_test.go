package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"teman-mancing/backend/internal/application/auth"
	appnotif "teman-mancing/backend/internal/application/notification"
	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/routes"
	domainnotif "teman-mancing/backend/internal/domain/notification"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
)

func setupNotificationTestRouter() (http.Handler, string, uuid.UUID, appnotif.Repository) {
	cfg := &config.Config{
		AppEnv:         "test",
		AppPort:        8080,
		JWTSecret:      "notif-test-secret-at-least-32-chars-long",
		JWTExpiration:  1 * time.Hour,
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	logger := zap.NewNop()
	userRepo := database.NewInMemoryUserRepository()
	spotRepo := database.NewInMemorySpotRepository()
	eventRepo := database.NewInMemoryEventRepository()
	chatRepo := database.NewInMemoryChatRepository()
	notifRepo := database.NewInMemoryNotificationRepository()
	tokenService := security.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)

	router := routes.SetupRouter(routes.RouterDeps{
		Config:           cfg,
		Logger:           logger,
		DB:               nil,
		MinIO:            nil,
		TokenService:     tokenService,
		UserRepo:         userRepo,
		SpotRepo:         spotRepo,
		EventRepo:        eventRepo,
		ChatRepo:         chatRepo,
		NotificationRepo: notifRepo,
		Version:          "1.0.0-test",
	})

	authService := auth.NewAuthService(userRepo, tokenService)
	u, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Notif User",
		Username: "notif_user",
		Email:    "notif@example.com",
		Password: "password123",
	})

	return router, u.Token, u.User.ID, notifRepo
}

func TestNotificationsFlow(t *testing.T) {
	router, token, userID, notifRepo := setupNotificationTestRouter()

	// Seed 2 notifications
	n1 := &domainnotif.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      domainnotif.TypeSpotComment,
		Title:     "Komentar Baru",
		Message:   "Budi mengomentari spot mancing Anda",
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	}
	n2 := &domainnotif.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      domainnotif.TypeJoinRequest,
		Title:     "Permohonan Gabung Event",
		Message:   "Agus meminta izin bergabung ke trip Cirata",
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	}
	_ = notifRepo.Create(nil, n1)
	_ = notifRepo.Create(nil, n2)

	// 1. List Notifications -> 200 OK (2 unread)
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	reqList.Header.Set("Authorization", "Bearer "+token)
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for list notifications, got %d", wList.Code)
	}

	var listRes dto.Response
	_ = json.Unmarshal(wList.Body.Bytes(), &listRes)
	if !listRes.Success {
		t.Errorf("expected success true in notification response")
	}

	// 2. Mark single notification as read -> 200 OK
	reqRead1, _ := http.NewRequest(http.MethodPatch, "/api/v1/notifications/"+n1.ID.String()+"/read", nil)
	reqRead1.Header.Set("Authorization", "Bearer "+token)
	wRead1 := httptest.NewRecorder()
	router.ServeHTTP(wRead1, reqRead1)

	if wRead1.Code != http.StatusOK {
		t.Errorf("expected 200 OK for marking notification read, got %d", wRead1.Code)
	}

	// 3. Mark all as read -> 200 OK
	reqReadAll, _ := http.NewRequest(http.MethodPost, "/api/v1/notifications/read-all", nil)
	reqReadAll.Header.Set("Authorization", "Bearer "+token)
	wReadAll := httptest.NewRecorder()
	router.ServeHTTP(wReadAll, reqReadAll)

	if wReadAll.Code != http.StatusOK {
		t.Errorf("expected 200 OK for marking all read, got %d", wReadAll.Code)
	}
}
