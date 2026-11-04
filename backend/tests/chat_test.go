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
	u1.User.IsVerified = true
	_ = userRepo.Update(nil, u1.User)

	u2, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Angler Bob",
		Username: "bob_angler",
		Email:    "bob@example.com",
		Password: "password123",
	})
	u2.User.IsVerified = true
	_ = userRepo.Update(nil, u2.User)

	u3, _ := authService.Register(nil, auth.RegisterInput{
		Name:     "Intruder Chuck",
		Username: "chuck_angler",
		Email:    "chuck@example.com",
		Password: "password123",
	})
	u3.User.IsVerified = true
	_ = userRepo.Update(nil, u3.User)

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

func TestEventGroupChatInMemoryAndMembershipSync(t *testing.T) {
	router, tokenAlice, _, tokenBob, _, tokenChuck := setupChatTestRouter()

	// 1. Alice creates event
	eventPayload := map[string]any{
		"name":                "Mancing Mania Cirata",
		"description":         "Mancing bareng di waduk Cirata",
		"destination_name":    "Cirata",
		"destination_address": "Cianjur, Jawa Barat",
		"event_date":          "2026-12-20",
		"start_time":          "07:00:00",
		"meeting_point":       "Dermaga 1",
		"visibility":          "public",
	}
	eBody, _ := json.Marshal(eventPayload)

	reqCreateEv, _ := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(eBody))
	reqCreateEv.Header.Set("Authorization", "Bearer "+tokenAlice)
	reqCreateEv.Header.Set("Content-Type", "application/json")
	wCreateEv := httptest.NewRecorder()
	router.ServeHTTP(wCreateEv, reqCreateEv)

	if wCreateEv.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for event creation, got %d: %s", wCreateEv.Code, wCreateEv.Body.String())
	}

	var evRes struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wCreateEv.Body.Bytes(), &evRes)
	eventID := evRes.Data.ID

	// 2. Alice lists conversations -> group chat must appear in conversation list
	reqListAlice, _ := http.NewRequest(http.MethodGet, "/api/v1/conversations", nil)
	reqListAlice.Header.Set("Authorization", "Bearer "+tokenAlice)
	wListAlice := httptest.NewRecorder()
	router.ServeHTTP(wListAlice, reqListAlice)

	if wListAlice.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Alice conversation list, got %d", wListAlice.Code)
	}

	var listAliceRes struct {
		Data []*domainchat.Conversation `json:"data"`
	}
	_ = json.Unmarshal(wListAlice.Body.Bytes(), &listAliceRes)

	var groupConv *domainchat.Conversation
	for _, c := range listAliceRes.Data {
		if c.Type == domainchat.TypeGroup && c.EventID != nil && c.EventID.String() == eventID {
			groupConv = c
			break
		}
	}

	if groupConv == nil {
		t.Fatalf("expected event group conversation to be generated and listed for Alice (creator)")
	}

	groupConvID := groupConv.ID.String()

	// 3. Alice sends message in event group chat -> 201 Created
	msgPayload := map[string]string{"message": "Halo angler, selamat bergabung di grup event!"}
	mBody, _ := json.Marshal(msgPayload)
	reqSendAlice, _ := http.NewRequest(http.MethodPost, "/api/v1/conversations/"+groupConvID+"/messages", bytes.NewBuffer(mBody))
	reqSendAlice.Header.Set("Authorization", "Bearer "+tokenAlice)
	reqSendAlice.Header.Set("Content-Type", "application/json")
	wSendAlice := httptest.NewRecorder()
	router.ServeHTTP(wSendAlice, reqSendAlice)

	if wSendAlice.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created when Alice sends group chat message, got %d: %s", wSendAlice.Code, wSendAlice.Body.String())
	}

	// 4. Bob (not joined yet) tries to read event group chat -> 403 Forbidden
	reqReadBob, _ := http.NewRequest(http.MethodGet, "/api/v1/conversations/"+groupConvID+"/messages", nil)
	reqReadBob.Header.Set("Authorization", "Bearer "+tokenBob)
	wReadBob := httptest.NewRecorder()
	router.ServeHTTP(wReadBob, reqReadBob)

	if wReadBob.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when non-participant Bob reads group chat, got %d", wReadBob.Code)
	}

	// 5. Bob joins public event -> Approved automatically -> Now can read group chat
	reqJoinBob, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/join", nil)
	reqJoinBob.Header.Set("Authorization", "Bearer "+tokenBob)
	wJoinBob := httptest.NewRecorder()
	router.ServeHTTP(wJoinBob, reqJoinBob)

	if wJoinBob.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Bob joining event, got %d", wJoinBob.Code)
	}

	wReadBob2 := httptest.NewRecorder()
	router.ServeHTTP(wReadBob2, reqReadBob)
	if wReadBob2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for approved participant Bob reading group chat, got %d: %s", wReadBob2.Code, wReadBob2.Body.String())
	}

	// 6. Bob leaves event -> Now blocked from reading group chat -> 403 Forbidden
	reqLeaveBob, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+eventID+"/leave", nil)
	reqLeaveBob.Header.Set("Authorization", "Bearer "+tokenBob)
	wLeaveBob := httptest.NewRecorder()
	router.ServeHTTP(wLeaveBob, reqLeaveBob)

	if wLeaveBob.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Bob leaving event, got %d", wLeaveBob.Code)
	}

	wReadBob3 := httptest.NewRecorder()
	router.ServeHTTP(wReadBob3, reqReadBob)
	if wReadBob3.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for Bob reading group chat after leaving, got %d", wReadBob3.Code)
	}

	// 7. Private event approval synchronization (Issue #13 Acceptance Criteria)
	privPayload := map[string]any{
		"name":                "Trip Privat Kepulauan Seribu",
		"description":         "Mancing privat seru di Pulau Pari",
		"destination_name":    "Pulau Pari",
		"destination_address": "Kepulauan Seribu",
		"event_date":          "2026-12-25",
		"start_time":          "06:00:00",
		"meeting_point":       "Dermaga Marina Ancol",
		"visibility":          "private",
	}
	pBody, _ := json.Marshal(privPayload)
	reqCreatePriv, _ := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(pBody))
	reqCreatePriv.Header.Set("Authorization", "Bearer "+tokenAlice)
	reqCreatePriv.Header.Set("Content-Type", "application/json")
	wCreatePriv := httptest.NewRecorder()
	router.ServeHTTP(wCreatePriv, reqCreatePriv)

	var privRes struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wCreatePriv.Body.Bytes(), &privRes)
	privEventID := privRes.Data.ID

	// Chuck requests to join private event -> status pending
	reqJoinPriv, _ := http.NewRequest(http.MethodPost, "/api/v1/events/"+privEventID+"/join", nil)
	reqJoinPriv.Header.Set("Authorization", "Bearer "+tokenChuck)
	wJoinPriv := httptest.NewRecorder()
	router.ServeHTTP(wJoinPriv, reqJoinPriv)

	// Chuck lists conversations or tries to open group chat -> must be denied (403 Forbidden)
	reqOpenPrivChuck, _ := http.NewRequest(http.MethodPost, "/api/v1/conversations", bytes.NewBuffer([]byte(`{"type":"group","event_id":"`+privEventID+`"}`)))
	reqOpenPrivChuck.Header.Set("Authorization", "Bearer "+tokenChuck)
	reqOpenPrivChuck.Header.Set("Content-Type", "application/json")
	wOpenPrivChuck := httptest.NewRecorder()
	router.ServeHTTP(wOpenPrivChuck, reqOpenPrivChuck)

	if wOpenPrivChuck.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when pending participant Chuck tries to access private event group chat, got %d: %s", wOpenPrivChuck.Code, wOpenPrivChuck.Body.String())
	}
}
