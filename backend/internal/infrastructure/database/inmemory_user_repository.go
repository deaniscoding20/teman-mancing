package database

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	appuser "teman-mancing/backend/internal/application/user"
	domainuser "teman-mancing/backend/internal/domain/user"
)

type InMemoryUserRepository struct {
	mu            sync.RWMutex
	users         map[uuid.UUID]*domainuser.User
	verifications map[uuid.UUID]*domainuser.EmailVerification
	resets        map[uuid.UUID]*domainuser.PasswordReset
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:         make(map[uuid.UUID]*domainuser.User),
		verifications: make(map[uuid.UUID]*domainuser.EmailVerification),
		resets:        make(map[uuid.UUID]*domainuser.PasswordReset),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, u *domainuser.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	email = strings.ToLower(strings.TrimSpace(email))
	for _, u := range r.users {
		if strings.EqualFold(u.Email, email) {
			return u, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	username = strings.ToLower(strings.TrimSpace(username))
	for _, u := range r.users {
		if strings.EqualFold(u.Username, username) {
			return u, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, u *domainuser.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) GetProfileView(ctx context.Context, id uuid.UUID) (*appuser.UserProfileView, error) {
	u, err := r.GetByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return &appuser.UserProfileView{
		User: u,
		Stats: appuser.UserStats{
			SpotsCount:         0,
			EventsCreatedCount: 0,
			EventsJoinedCount:  0,
		},
	}, nil
}

func (r *InMemoryUserRepository) CreateEmailVerification(ctx context.Context, v *domainuser.EmailVerification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.verifications[v.ID] = v
	return nil
}

func (r *InMemoryUserRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domainuser.EmailVerification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.verifications {
		if v.Token == token && v.ExpiresAt.After(time.Now().UTC()) {
			return v, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) GetEmailVerificationByOTP(ctx context.Context, userID uuid.UUID, otp string) (*domainuser.EmailVerification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.verifications {
		if v.UserID == userID && v.OTPCode == otp && v.ExpiresAt.After(time.Now().UTC()) {
			return v, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) DeleteEmailVerification(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.verifications, id)
	return nil
}

func (r *InMemoryUserRepository) CreatePasswordReset(ctx context.Context, reset *domainuser.PasswordReset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resets[reset.ID] = reset
	return nil
}

func (r *InMemoryUserRepository) GetPasswordResetByToken(ctx context.Context, token string) (*domainuser.PasswordReset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, reset := range r.resets {
		if reset.Token == token && reset.UsedAt == nil && reset.ExpiresAt.After(time.Now().UTC()) {
			return reset, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) GetPasswordResetByOTP(ctx context.Context, email, otp string) (*domainuser.PasswordReset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, reset := range r.resets {
		if strings.EqualFold(reset.Email, email) && reset.OTPCode == otp && reset.UsedAt == nil && reset.ExpiresAt.After(time.Now().UTC()) {
			return reset, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) MarkPasswordResetUsed(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if reset, ok := r.resets[id]; ok {
		now := time.Now().UTC()
		reset.UsedAt = &now
	}
	return nil
}
