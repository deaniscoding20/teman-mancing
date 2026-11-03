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
	domainspot "teman-mancing/backend/internal/domain/fishing_spot"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
)

func setupSpotTestRouter() (http.Handler, string, string) {
	cfg := &config.Config{
		AppEnv:         "test",
		AppPort:        8080,
		JWTSecret:      "spot-test-secret-at-least-32-chars-long-here",
		JWTExpiration:  1 * time.Hour,
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	logger := zap.NewNop()
	userRepo := database.NewInMemoryUserRepository()
	spotRepo := database.NewInMemorySpotRepository()
	tokenService := security.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)

	router := routes.SetupRouter(routes.RouterDeps{
		Config:       cfg,
		Logger:       logger,
		DB:           nil,
		MinIO:        nil,
		TokenService: tokenService,
		UserRepo:     userRepo,
		SpotRepo:     spotRepo,
		Version:      "1.0.0-test",
	})

	authService := auth.NewAuthService(userRepo, tokenService, nil)

	// User 1
	u1, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Pemancing A",
		Username: "angler_a",
		Email:    "a@example.com",
		Password: "password123",
	})
	u1.User.IsVerified = true
	_ = userRepo.Update(nil, u1.User)

	// User 2
	u2, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Pemancing B",
		Username: "angler_b",
		Email:    "b@example.com",
		Password: "password123",
	})
	u2.User.IsVerified = true
	_ = userRepo.Update(nil, u2.User)

	return router, u1.Token, u2.Token
}

func TestFishingSpotCRUDAndOwnership(t *testing.T) {
	router, tokenA, tokenB := setupSpotTestRouter()

	// 1. Create Spot by User A
	createPayload := map[string]any{
		"name":        "Dermaga Rawa Pening",
		"description": "Spot mancing rawa asri dengan populasi nila dan lele melimpah.",
		"address":     "Kabupaten Semarang, Jawa Tengah",
		"latitude":    -7.2801,
		"longitude":   110.4350,
		"fish_types":  []string{"Nila", "Lele", "Gabus"},
		"techniques":  []string{"Dasaran", "Pelampung"},
		"recommended_baits": []map[string]string{
			{"name": "Cacing Lur", "description": "Umpan jitu di pinggir eceng gondok"},
			{"name": "Lumut Segar", "description": "Cocok untuk target babon"},
		},
	}
	body, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for spot, got %d: %s", w.Code, w.Body.String())
	}

	var createRes struct {
		Success bool                   `json:"success"`
		Data    domainspot.FishingSpot `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &createRes)
	spotID := createRes.Data.ID.String()

	// 2. User B attempts to edit User A's spot -> 403 Forbidden (README Section 12.4)
	newName := "Spot Dibajak User B"
	updPayload := map[string]any{"name": newName}
	updBody, _ := json.Marshal(updPayload)

	reqUpdB, _ := http.NewRequest(http.MethodPatch, "/api/v1/spots/"+spotID, bytes.NewBuffer(updBody))
	reqUpdB.Header.Set("Authorization", "Bearer "+tokenB)
	reqUpdB.Header.Set("Content-Type", "application/json")
	wUpdB := httptest.NewRecorder()
	router.ServeHTTP(wUpdB, reqUpdB)

	if wUpdB.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when non-owner edits spot, got %d", wUpdB.Code)
	}

	// 3. User A successfully edits own spot -> 200 OK
	legitName := "Dermaga Rawa Pening - Update"
	legitPayload := map[string]any{"name": legitName}
	legitBody, _ := json.Marshal(legitPayload)

	reqUpdA, _ := http.NewRequest(http.MethodPatch, "/api/v1/spots/"+spotID, bytes.NewBuffer(legitBody))
	reqUpdA.Header.Set("Authorization", "Bearer "+tokenA)
	reqUpdA.Header.Set("Content-Type", "application/json")
	wUpdA := httptest.NewRecorder()
	router.ServeHTTP(wUpdA, reqUpdA)

	if wUpdA.Code != http.StatusOK {
		t.Errorf("expected 200 OK when owner edits spot, got %d: %s", wUpdA.Code, wUpdA.Body.String())
	}

	// 4. User B posts comment on Spot
	commentPayload := map[string]string{"comment": "Air sedang tenang, cocok untuk casting!"}
	cBody, _ := json.Marshal(commentPayload)

	reqComm, _ := http.NewRequest(http.MethodPost, "/api/v1/spots/"+spotID+"/comments", bytes.NewBuffer(cBody))
	reqComm.Header.Set("Authorization", "Bearer "+tokenB)
	reqComm.Header.Set("Content-Type", "application/json")
	wComm := httptest.NewRecorder()
	router.ServeHTTP(wComm, reqComm)

	if wComm.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for comment, got %d: %s", wComm.Code, wComm.Body.String())
	}

	var commentRes struct {
		Success bool               `json:"success"`
		Data    domainspot.Comment `json:"data"`
	}
	_ = json.Unmarshal(wComm.Body.Bytes(), &commentRes)
	commentID := commentRes.Data.ID.String()

	// 5. User A tries to delete User B's comment -> 403 Forbidden (README Section 12.4)
	reqDelCommA, _ := http.NewRequest(http.MethodDelete, "/api/v1/comments/"+commentID, nil)
	reqDelCommA.Header.Set("Authorization", "Bearer "+tokenA)
	wDelCommA := httptest.NewRecorder()
	router.ServeHTTP(wDelCommA, reqDelCommA)

	if wDelCommA.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when deleting someone else's comment, got %d", wDelCommA.Code)
	}

	// 6. User B deletes own comment -> 200 OK
	reqDelCommB, _ := http.NewRequest(http.MethodDelete, "/api/v1/comments/"+commentID, nil)
	reqDelCommB.Header.Set("Authorization", "Bearer "+tokenB)
	wDelCommB := httptest.NewRecorder()
	router.ServeHTTP(wDelCommB, reqDelCommB)

	if wDelCommB.Code != http.StatusOK {
		t.Errorf("expected 200 OK when author deletes comment, got %d", wDelCommB.Code)
	}

	// 7. Get Spot by ID
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/spots/"+spotID, nil)
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Errorf("expected 200 OK for spot detail, got %d", wGet.Code)
	}

	// 8. List Spots
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/spots?q=Rawa", nil)
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Errorf("expected 200 OK for spots list, got %d", wList.Code)
	}
	var listRes dto.Response
	_ = json.Unmarshal(wList.Body.Bytes(), &listRes)
	if !listRes.Success {
		t.Errorf("expected success true in spot list response")
	}

	// 9. User B tries to delete User A's spot -> 403 Forbidden
	reqDelB, _ := http.NewRequest(http.MethodDelete, "/api/v1/spots/"+spotID, nil)
	reqDelB.Header.Set("Authorization", "Bearer "+tokenB)
	wDelB := httptest.NewRecorder()
	router.ServeHTTP(wDelB, reqDelB)

	if wDelB.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when non-owner deletes spot, got %d", wDelB.Code)
	}

	// 10. User A deletes own spot -> 200 OK
	reqDelA, _ := http.NewRequest(http.MethodDelete, "/api/v1/spots/"+spotID, nil)
	reqDelA.Header.Set("Authorization", "Bearer "+tokenA)
	wDelA := httptest.NewRecorder()
	router.ServeHTTP(wDelA, reqDelA)

	if wDelA.Code != http.StatusOK {
		t.Errorf("expected 200 OK when owner deletes spot, got %d", wDelA.Code)
	}
}
