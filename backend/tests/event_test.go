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
	domainevent "teman-mancing/backend/internal/domain/event"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
)

func setupEventTestRouter() (http.Handler, string, string, string) {
	cfg := &config.Config{
		AppEnv:         "test",
		AppPort:        8080,
		JWTSecret:      "event-test-secret-at-least-32-chars-long",
		JWTExpiration:  1 * time.Hour,
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	logger := zap.NewNop()
	userRepo := database.NewInMemoryUserRepository()
	spotRepo := database.NewInMemorySpotRepository()
	eventRepo := database.NewInMemoryEventRepository()
	tokenService := security.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)

	router := routes.SetupRouter(routes.RouterDeps{
		Config:       cfg,
		Logger:       logger,
		DB:           nil,
		MinIO:        nil,
		TokenService: tokenService,
		UserRepo:     userRepo,
		SpotRepo:     spotRepo,
		EventRepo:    eventRepo,
		Version:      "1.0.0-test",
	})

	authService := auth.NewAuthService(userRepo, tokenService, nil)

	u1, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Host Event",
		Username: "host_event",
		Email:    "host@example.com",
		Password: "password123",
	})
	u1.User.IsVerified = true
	_ = userRepo.Update(nil, u1.User)

	u2, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Participant 1",
		Username: "part1",
		Email:    "part1@example.com",
		Password: "password123",
	})
	u2.User.IsVerified = true
	_ = userRepo.Update(nil, u2.User)

	u3, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Participant 2",
		Username: "part2",
		Email:    "part2@example.com",
		Password: "password123",
	})
	u3.User.IsVerified = true
	_ = userRepo.Update(nil, u3.User)

	return router, u1.Token, u2.Token, u3.Token
}

func TestEventCreationAndCapacityValidation(t *testing.T) {
	router, tokenHost, tokenUser1, tokenUser2 := setupEventTestRouter()

	maxPart := 2 // Host is participant 1, so only 1 more slot available!
	createPayload := map[string]any{
		"name":                 "Trip Mancing Kakap Muara Baru",
		"description":          "Mancing bareng pagi hari menyusuri muara.",
		"destination_name":     "Muara Baru Jakarta Utara",
		"destination_address":  "Pelabuhan Muara Baru Dermaga 3",
		"event_date":           "2026-12-01",
		"start_time":           "06:00:00",
		"meeting_point":        "Pintu Masuk Gerbang Utama Dermaga",
		"difficulty":           "pemula",
		"visibility":           "public",
		"has_parking_fee":      true,
		"parking_fee":          10000,
		"requires_boat":        true,
		"boat_rental_fee":      150000,
		"max_participants":     maxPart,
	}
	body, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tokenHost)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for event, got %d: %s", w.Code, w.Body.String())
	}

	var createRes struct {
		Success bool              `json:"success"`
		Data    domainevent.Event `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &createRes)
	eventID := createRes.Data.ID.String()

	// 1. User 1 joins public event -> Should SUCCEED (slot 2 of 2 taken)
	reqJoin1, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/join", nil)
	reqJoin1.Header.Set("Authorization", "Bearer "+tokenUser1)
	wJoin1 := httptest.NewRecorder()
	router.ServeHTTP(wJoin1, reqJoin1)

	if wJoin1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for join, got %d: %s", wJoin1.Code, wJoin1.Body.String())
	}

	// 2. User 1 tries to join AGAIN -> 409 Conflict (Duplicate prevention, README Section 11.11, 20.3)
	reqJoinDup, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/join", nil)
	reqJoinDup.Header.Set("Authorization", "Bearer "+tokenUser1)
	wJoinDup := httptest.NewRecorder()
	router.ServeHTTP(wJoinDup, reqJoinDup)

	if wJoinDup.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict for duplicate join, got %d", wJoinDup.Code)
	}

	// 3. User 2 tries to join when capacity is full -> 422 Unprocessable (README Section 12.2, 20.3)
	reqJoin2, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/join", nil)
	reqJoin2.Header.Set("Authorization", "Bearer "+tokenUser2)
	wJoin2 := httptest.NewRecorder()
	router.ServeHTTP(wJoin2, reqJoin2)

	if wJoin2.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 Unprocessable Entity when event capacity is full, got %d: %s", wJoin2.Code, wJoin2.Body.String())
	}

	// 4. User 1 leaves event -> 200 OK
	reqLeave, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/leave", nil)
	reqLeave.Header.Set("Authorization", "Bearer "+tokenUser1)
	wLeave := httptest.NewRecorder()
	router.ServeHTTP(wLeave, reqLeave)

	if wLeave.Code != http.StatusOK {
		t.Errorf("expected 200 OK for leaving event, got %d", wLeave.Code)
	}

	// 5. User 2 can now join because a slot freed up!
	reqJoin2Retry, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/join", nil)
	reqJoin2Retry.Header.Set("Authorization", "Bearer "+tokenUser2)
	wJoin2Retry := httptest.NewRecorder()
	router.ServeHTTP(wJoin2Retry, reqJoin2Retry)

	if wJoin2Retry.Code != http.StatusOK {
		t.Errorf("expected 200 OK for join after slot freed, got %d", wJoin2Retry.Code)
	}

	// 6. Non-creator tries to update event -> 403 Forbidden (README Section 12.4)
	updPayload := map[string]any{"name": "Nama Diubah Ilegal"}
	updBody, _ := json.Marshal(updPayload)
	reqUpdIllegal, _ := http.NewRequest(http.MethodPatch, "/api/v1/events/"+eventID, bytes.NewBuffer(updBody))
	reqUpdIllegal.Header.Set("Authorization", "Bearer "+tokenUser2)
	reqUpdIllegal.Header.Set("Content-Type", "application/json")
	wUpdIllegal := httptest.NewRecorder()
	router.ServeHTTP(wUpdIllegal, reqUpdIllegal)

	if wUpdIllegal.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for non-creator updating event, got %d", wUpdIllegal.Code)
	}

	// 7. List Events
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/events", nil)
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Errorf("expected 200 OK for listing events, got %d", wList.Code)
	}

	var listRes dto.Response
	_ = json.Unmarshal(wList.Body.Bytes(), &listRes)
	if !listRes.Success {
		t.Errorf("expected success true in list events")
	}
}

func TestPrivateEventApprovalFlow(t *testing.T) {
	router, tokenHost, tokenUser1, _ := setupEventTestRouter()

	createPayload := map[string]any{
		"name":                "Private Angler Fishing Club",
		"description":         "Khusus anggota komunitas terverifikasi.",
		"destination_name":    "Danau Sunter Barat",
		"destination_address": "Jl. Danau Sunter Barat, Jakarta Utara",
		"event_date":          "2026-12-10",
		"start_time":          "07:00:00",
		"meeting_point":       "Indomaret Point Danau Sunter",
		"difficulty":          "menengah",
		"visibility":          "private",
	}
	body, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tokenHost)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createRes struct {
		Data domainevent.Event `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &createRes)
	eventID := createRes.Data.ID.String()

	// 1. User 1 sends request to join private event -> status must be 'pending'
	reqJoin, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/join", nil)
	reqJoin.Header.Set("Authorization", "Bearer "+tokenUser1)
	wJoin := httptest.NewRecorder()
	router.ServeHTTP(wJoin, reqJoin)

	if wJoin.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for request join, got %d", wJoin.Code)
	}

	var joinRes struct {
		Data domainevent.Participant `json:"data"`
	}
	_ = json.Unmarshal(wJoin.Body.Bytes(), &joinRes)
	if joinRes.Data.Status != domainevent.ParticipantStatusPending {
		t.Errorf("expected participant status 'pending' for private event, got '%s'", joinRes.Data.Status)
	}
	applicantID := joinRes.Data.UserID.String()

	// 2. Creator views requests -> should see 1 pending request
	reqListReq, _ := http.NewRequest(http.MethodGet, "/api/v1/events/"+eventID+"/requests", nil)
	reqListReq.Header.Set("Authorization", "Bearer "+tokenHost)
	wListReq := httptest.NewRecorder()
	router.ServeHTTP(wListReq, reqListReq)

	if wListReq.Code != http.StatusOK {
		t.Errorf("expected 200 OK for list requests, got %d", wListReq.Code)
	}

	// 3. Creator approves request -> 200 OK
	reqApprove, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/requests/"+applicantID+"/approve", nil)
	reqApprove.Header.Set("Authorization", "Bearer "+tokenHost)
	wApprove := httptest.NewRecorder()
	router.ServeHTTP(wApprove, reqApprove)

	if wApprove.Code != http.StatusOK {
		t.Errorf("expected 200 OK for approve request, got %d: %s", wApprove.Code, wApprove.Body.String())
	}
}
