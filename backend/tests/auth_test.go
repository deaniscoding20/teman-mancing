package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"teman-mancing/backend/internal/application/auth"
	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/routes"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
)

func setupAuthRouter() http.Handler {
	cfg := &config.Config{
		AppEnv:         "test",
		AppPort:        8080,
		JWTSecret:      "test-secret-min-32-characters-for-jwt-signing",
		JWTExpiration:  1 * time.Hour,
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	logger := zap.NewNop()
	userRepo := database.NewInMemoryUserRepository()
	tokenService := security.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)

	return routes.SetupRouter(routes.RouterDeps{
		Config:       cfg,
		Logger:       logger,
		DB:           nil,
		MinIO:        nil,
		TokenService: tokenService,
		UserRepo:     userRepo,
		Version:      "1.0.0-test",
	})
}

func TestAuthFlow(t *testing.T) {
	router := setupAuthRouter()

	// 1. Register User
	regPayload := auth.RegisterInput{
		Name:     "Budi Santoso",
		Username: "budimancing",
		Email:    "budi@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(regPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var regRes struct {
		Success bool            `json:"success"`
		Data    auth.AuthResult `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &regRes); err != nil {
		t.Fatalf("failed to decode register response: %v", err)
	}

	token := regRes.Data.Token
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	// 2. Duplicate Registration by Email
	reqDup, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	reqDup.Header.Set("Content-Type", "application/json")
	wDup := httptest.NewRecorder()
	router.ServeHTTP(wDup, reqDup)

	if wDup.Code != http.StatusConflict {
		t.Errorf("expected status 409 Conflict for duplicate registration, got %d", wDup.Code)
	}

	// 3. Login
	loginPayload := auth.LoginInput{
		Email:    "budi@example.com",
		Password: "password123",
	}
	lBody, _ := json.Marshal(loginPayload)
	reqLog, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(lBody))
	reqLog.Header.Set("Content-Type", "application/json")
	wLog := httptest.NewRecorder()
	router.ServeHTTP(wLog, reqLog)

	if wLog.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on login, got %d", wLog.Code)
	}

	// 4. Access Protected /api/v1/auth/me without token -> 401
	reqMeNoAuth, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	wMeNoAuth := httptest.NewRecorder()
	router.ServeHTTP(wMeNoAuth, reqMeNoAuth)

	if wMeNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", wMeNoAuth.Code)
	}

	// 5. Access Protected /api/v1/auth/me with Bearer token -> 200
	reqMe, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+token)
	wMe := httptest.NewRecorder()
	router.ServeHTTP(wMe, reqMe)

	if wMe.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /auth/me, got %d: %s", wMe.Code, wMe.Body.String())
	}

	// 6. Update Profile
	newBio := "Pemancing laut dan air tawar dari Surabaya."
	updatePayload := map[string]string{
		"name":     "Budi S. Angler",
		"bio":      newBio,
		"location": "Surabaya, Jawa Timur",
	}
	uBody, _ := json.Marshal(updatePayload)
	reqUpd, _ := http.NewRequest(http.MethodPatch, "/api/v1/users/me", bytes.NewBuffer(uBody))
	reqUpd.Header.Set("Authorization", "Bearer "+token)
	reqUpd.Header.Set("Content-Type", "application/json")
	wUpd := httptest.NewRecorder()
	router.ServeHTTP(wUpd, reqUpd)

	if wUpd.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for profile update, got %d", wUpd.Code)
	}

	// 7. Get Public Profile by Username
	reqPub, _ := http.NewRequest(http.MethodGet, "/api/v1/users/budimancing", nil)
	wPub := httptest.NewRecorder()
	router.ServeHTTP(wPub, reqPub)

	if wPub.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for public profile, got %d", wPub.Code)
	}

	var pubRes dto.Response
	_ = json.Unmarshal(wPub.Body.Bytes(), &pubRes)
	if !pubRes.Success {
		t.Errorf("expected success true on public profile")
	}
}
