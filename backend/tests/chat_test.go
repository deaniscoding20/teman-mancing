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
	domainchat "teman-mancing/backend/internal/domain/chat"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/security"
)

func setupChatTestRouter() (http.Handler, string, string, string, string, string) {
	cfg := &config.Config{
		AppEnv:         "test",
		AppPort:        8080,
		JWTSecret:      "chat-test-secret-at-least-32-chars-long",
		JWTExpiration:  1 * time.Hour,
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	logger := zap.NewNop()
	userRepo := database.NewInMemoryUserRepository()
	spotRepo := database.NewInMemorySpotRepository()
	eventRepo := database.NewInMemoryEventRepository()
	chatRepo := database.NewInMemoryChatRepository()
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
		ChatRepo:     chatRepo,
		Version:      "1.0.0-test",
	})

	authService := auth.NewAuthService(userRepo, tokenService, nil)

	u1, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Angler Alice",
		Username: "alice_angler",
		Email:    "alice@example.com",
		Password: "password123",
	})

	u2, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Angler Bob",
		Username: "bob_angler",
		Email:    "bob@example.com",
		Password: "password123",
	})

	u3, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Intruder Chuck",
		Username: "chuck_angler",
		Email:    "chuck@example.com",
		Password: "password123",
	})

	return router, u1.Token, u1.User.ID.String(), u2.Token, u2.User.ID.String(), u3.Token
}

func TestChatMembershipAuthorizationAndPersistence(t *testing.T) {
	router, tokenAlice, _, tokenBob, idBob, tokenChuck := setupChatTestRouter()

	// 1. Alice creates direct conversation with Bob
	convReqPayload := map[string]string{
		"type":           "direct",
		"participant_id": idBob,
	}
	body, _ := json.Marshal(convReqPayload)

	reqConv, _ := http.NewRequest(http.MethodPost, "/api/v1/conversations", bytes.NewBuffer(body))
	reqConv.Header.Set("Authorization", "Bearer "+tokenAlice)
	reqConv.Header.Set("Content-Type", "application/json")
	wConv := httptest.NewRecorder()
	router.ServeHTTP(wConv, reqConv)

	if wConv.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for conversation, got %d: %s", wConv.Code, wConv.Body.String())
	}

	var convRes struct {
		Data domainchat.Conversation `json:"data"`
	}
	_ = json.Unmarshal(wConv.Body.Bytes(), &convRes)
	convID := convRes.Data.ID.String()

	// 2. Alice sends message -> 201 Created (Persisted to storage)
	msgPayload := map[string]string{"message": "Halo Bob, akhir pekan ini mau mancing di mana?"}
	mBody, _ := json.Marshal(msgPayload)

	reqMsg, _ := http.NewRequest(http.MethodPost, "/api/v1/conversations/"+convID+"/messages", bytes.NewBuffer(mBody))
	reqMsg.Header.Set("Authorization", "Bearer "+tokenAlice)
	reqMsg.Header.Set("Content-Type", "application/json")
	wMsg := httptest.NewRecorder()
	router.ServeHTTP(wMsg, reqMsg)

	if wMsg.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for message, got %d: %s", wMsg.Code, wMsg.Body.String())
	}

	// 3. Bob reads message history -> 200 OK
	reqReadBob, _ := http.NewRequest(http.MethodGet, "/api/v1/conversations/"+convID+"/messages", nil)
	reqReadBob.Header.Set("Authorization", "Bearer "+tokenBob)
	wReadBob := httptest.NewRecorder()
	router.ServeHTTP(wReadBob, reqReadBob)

	if wReadBob.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Bob reading messages, got %d", wReadBob.Code)
	}

	var bobRes struct {
		Data []*domainchat.Message `json:"data"`
	}
	_ = json.Unmarshal(wReadBob.Body.Bytes(), &bobRes)
	if len(bobRes.Data) != 1 || bobRes.Data[0].Message != "Halo Bob, akhir pekan ini mau mancing di mana?" {
		t.Errorf("expected 1 message matching content, got %v", bobRes.Data)
	}

	// 4. Chuck (not a participant) tries to read messages -> 403 Forbidden (README Section 12.3)
	reqReadChuck, _ := http.NewRequest(http.MethodGet, "/api/v1/conversations/"+convID+"/messages", nil)
	reqReadChuck.Header.Set("Authorization", "Bearer "+tokenChuck)
	wReadChuck := httptest.NewRecorder()
	router.ServeHTTP(wReadChuck, reqReadChuck)

	if wReadChuck.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when unauthorized user reads chat, got %d", wReadChuck.Code)
	}

	// 5. Chuck tries to send message to conversation -> 403 Forbidden (README Section 12.3)
	chuckMsg := map[string]string{"message": "Menyusup masuk obrolan orang"}
	cBody, _ := json.Marshal(chuckMsg)
	reqSendChuck, _ := http.NewRequest(http.MethodPost, "/api/v1/conversations/"+convID+"/messages", bytes.NewBuffer(cBody))
	reqSendChuck.Header.Set("Authorization", "Bearer "+tokenChuck)
	reqSendChuck.Header.Set("Content-Type", "application/json")
	wSendChuck := httptest.NewRecorder()
	router.ServeHTTP(wSendChuck, reqSendChuck)

	if wSendChuck.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when unauthorized user posts to chat, got %d", wSendChuck.Code)
	}

	// 6. List conversations for Alice -> 200 OK
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/conversations", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenAlice)
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Errorf("expected 200 OK for list conversations, got %d", wList.Code)
	}

	var listRes dto.Response
	_ = json.Unmarshal(wList.Body.Bytes(), &listRes)
	if !listRes.Success {
		t.Errorf("expected success true in list conversations")
	}

	// 7. Test WebSocket Query Token Auth
	reqWSNoAuth, _ := http.NewRequest(http.MethodGet, "/api/v1/ws", nil)
	wWSNoAuth := httptest.NewRecorder()
	router.ServeHTTP(wWSNoAuth, reqWSNoAuth)
	if wWSNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for WebSocket without token, got %d", wWSNoAuth.Code)
	}

	reqWSQuery, _ := http.NewRequest(http.MethodGet, "/api/v1/ws?token="+tokenAlice, nil)
	wWSQuery := httptest.NewRecorder()
	router.ServeHTTP(wWSQuery, reqWSQuery)
	if wWSQuery.Code == http.StatusUnauthorized {
		t.Errorf("expected WebSocket with query token to pass authentication, but got 401")
	}
}
