package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"teman-mancing/backend/internal/config"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/routes"
	"teman-mancing/backend/internal/domain/event"
)

func setupTestRouter() http.Handler {
	cfg := &config.Config{
		AppEnv:         "test",
		AppPort:        8080,
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	logger := zap.NewNop()

	return routes.SetupRouter(routes.RouterDeps{
		Config:  cfg,
		Logger:  logger,
		DB:      nil,
		MinIO:   nil,
		Version: "1.0.0-test",
	})
}

func TestRootEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if res["app"] != "Teman Mancing API" {
		t.Errorf("expected app 'Teman Mancing API', got '%v'", res["app"])
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true, got false")
	}
}

func TestVersionEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/version", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if res.Message != "Teman Mancing Backend API" {
		t.Errorf("expected message 'Teman Mancing Backend API', got '%s'", res.Message)
	}
}

func TestEventCapacityRules(t *testing.T) {
	maxPart := 5
	ev := event.Event{
		Status:          event.StatusOpen,
		MaxParticipants: &maxPart,
	}

	// Under capacity
	if err := ev.CanAcceptParticipant(3); err != nil {
		t.Errorf("expected nil error when under capacity, got %v", err)
	}

	// At capacity
	if err := ev.CanAcceptParticipant(5); err != event.ErrEventCapacityFull {
		t.Errorf("expected ErrEventCapacityFull, got %v", err)
	}

	// Cancelled event
	ev.Status = event.StatusCancelled
	if err := ev.CanAcceptParticipant(2); err != event.ErrEventCancelled {
		t.Errorf("expected ErrEventCancelled, got %v", err)
	}
}
