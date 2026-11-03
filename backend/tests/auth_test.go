package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"teman-mancing/backend/internal/application/auth"
	appuser "teman-mancing/backend/internal/application/user"
	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/routes"
	domainuser "teman-mancing/backend/internal/domain/user"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
)

func setupAuthRouter() (http.Handler, appuser.Repository) {
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

	handler := routes.SetupRouter(routes.RouterDeps{
		Config:       cfg,
		Logger:       logger,
		DB:           nil,
		MinIO:        nil,
		Mailer:       nil,
		TokenService: tokenService,
		UserRepo:     userRepo,
		Version:      "1.0.0-test",
	})
	return handler, userRepo
}

func TestAuthFlow(t *testing.T) {
	router, userRepo := setupAuthRouter()

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

	// 3. Login before email verification -> 403 Forbidden
	loginPayload := auth.LoginInput{
		Email:    "budi@example.com",
		Password: "password123",
	}
	lBody, _ := json.Marshal(loginPayload)
	reqLog, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(lBody))
	reqLog.Header.Set("Content-Type", "application/json")
	wLog := httptest.NewRecorder()
	router.ServeHTTP(wLog, reqLog)

	if wLog.Code != http.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden before email verification, got %d", wLog.Code)
	}

	// 3a. Verify Email
	u, _ := userRepo.GetByEmail(nil, "budi@example.com")
	u.IsVerified = true
	_ = userRepo.Update(nil, u)

	// 3b. Login after email verification -> 200 OK
	reqLog2, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(lBody))
	reqLog2.Header.Set("Content-Type", "application/json")
	wLog2 := httptest.NewRecorder()
	router.ServeHTTP(wLog2, reqLog2)
	if wLog2.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on login after verification, got %d: %s", wLog2.Code, wLog2.Body.String())
	}

	// 3c. Verify Cookie was set on login
	cookies := wLog2.Result().Cookies()
	var authCookieFound bool
	for _, cookie := range cookies {
		if cookie.Name == "teman_mancing_token" && cookie.Value != "" && cookie.HttpOnly {
			authCookieFound = true
			break
		}
	}
	if !authCookieFound {
		t.Errorf("expected HttpOnly cookie teman_mancing_token to be set on login")
	}

	// 3c. Login using Username ("budimancing")
	uLoginPayload := map[string]string{
		"identifier": "budimancing",
		"password":   "password123",
	}
	uLBody, _ := json.Marshal(uLoginPayload)
	reqULog, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(uLBody))
	reqULog.Header.Set("Content-Type", "application/json")
	wULog := httptest.NewRecorder()
	router.ServeHTTP(wULog, reqULog)

	if wULog.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on login by username, got %d: %s", wULog.Code, wULog.Body.String())
	}

	// 3d. Login with wrong password -> 401
	wrongPwdPayload := map[string]string{
		"identifier": "budimancing",
		"password":   "wrong-password",
	}
	wpBody, _ := json.Marshal(wrongPwdPayload)
	reqWp, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(wpBody))
	reqWp.Header.Set("Content-Type", "application/json")
	wWp := httptest.NewRecorder()
	router.ServeHTTP(wWp, reqWp)

	if wWp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for wrong password, got %d", wWp.Code)
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

	var pubRes struct {
		Success bool `json:"success"`
		Data    struct {
			Username           string `json:"username"`
			SpotsCount         int64  `json:"spots_count"`
			EventsCreatedCount int64  `json:"events_created_count"`
			EventsJoinedCount  int64  `json:"events_joined_count"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wPub.Body.Bytes(), &pubRes)
	if !pubRes.Success || pubRes.Data.Username != "budimancing" {
		t.Errorf("expected success true and username budimancing, got %+v", pubRes)
	}

	// 8. Google OAuth Login endpoint
	googlePayload := map[string]string{
		"email": "angler.google@gmail.com",
		"name":  "Angler Google User",
	}
	gBody, _ := json.Marshal(googlePayload)
	reqGoogle, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/google", bytes.NewBuffer(gBody))
	reqGoogle.Header.Set("Content-Type", "application/json")
	wGoogle := httptest.NewRecorder()
	router.ServeHTTP(wGoogle, reqGoogle)

	if wGoogle.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Google login, got %d: %s", wGoogle.Code, wGoogle.Body.String())
	}

	var gRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
			User  struct {
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wGoogle.Body.Bytes(), &gRes)
	if !gRes.Success || gRes.Data.User.Email != "angler.google@gmail.com" || gRes.Data.Token == "" {
		t.Errorf("expected successful google login response with token, got %+v", gRes)
	}

	// 9. Forgot Password endpoint
	fpPayload := map[string]string{
		"email": "budi@example.com",
	}
	fpBody, _ := json.Marshal(fpPayload)
	reqFP, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewBuffer(fpBody))
	reqFP.Header.Set("Content-Type", "application/json")
	wFP := httptest.NewRecorder()
	router.ServeHTTP(wFP, reqFP)
	if wFP.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for forgot-password, got %d", wFP.Code)
	}

	// 10. Reset Password with token
	resetToken := "test-reset-token-123"
	_ = userRepo.CreatePasswordReset(nil, &domainuser.PasswordReset{
		ID:        uuid.New(),
		UserID:    u.ID,
		Email:     "budi@example.com",
		OTPCode:   "123456",
		Token:     resetToken,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	})

	rpPayload := map[string]string{
		"token":    resetToken,
		"password": "newpassword123",
	}
	rpBody, _ := json.Marshal(rpPayload)
	reqRP, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewBuffer(rpBody))
	reqRP.Header.Set("Content-Type", "application/json")
	wRP := httptest.NewRecorder()
	router.ServeHTTP(wRP, reqRP)
	if wRP.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for reset-password, got %d: %s", wRP.Code, wRP.Body.String())
	}
}
