package database

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"

	appuser "teman-mancing/backend/internal/application/user"
	domainuser "teman-mancing/backend/internal/domain/user"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*domainuser.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[uuid.UUID]*domainuser.User),
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
