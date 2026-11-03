package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"teman-mancing/backend/internal/application/common"
	appuser "teman-mancing/backend/internal/application/user"
	domainuser "teman-mancing/backend/internal/domain/user"
	"teman-mancing/backend/internal/infrastructure/security"
)

var (
	ErrEmailExists    = errors.New("email is already registered")
	ErrUsernameExists = errors.New("username is already taken")
	ErrInvalidAuth    = errors.New("email/username atau password salah")
	dummyBcryptHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-auth-salt-constant-time"), bcrypt.DefaultCost)
)

type RegisterInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type AuthResult struct {
	User  *domainuser.User `json:"user"`
	Token string           `json:"token"`
}

type UpdateProfileInput struct {
	Name     *string `json:"name"`
	Bio      *string `json:"bio"`
	Location *string `json:"location"`
}

type AuthService struct {
	userRepo     appuser.Repository
	tokenService *security.TokenService
}

func NewAuthService(userRepo appuser.Repository, tokenService *security.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(input.Email))
	cleanUsername := strings.ToLower(strings.TrimSpace(input.Username))

	existingByEmail, err := s.userRepo.GetByEmail(ctx, cleanEmail)
	if err != nil {
		return nil, err
	}
	if existingByEmail != nil {
		return nil, ErrEmailExists
	}

	existingByUsername, err := s.userRepo.GetByUsername(ctx, cleanUsername)
	if err != nil {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, ErrUsernameExists
	}

	newUser, err := domainuser.NewUser(input.Name, cleanUsername, cleanEmail, input.Password)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	token, err := s.tokenService.GenerateToken(newUser.ID, newUser.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  newUser,
		Token: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	identifier := strings.TrimSpace(input.Identifier)
	if identifier == "" {
		identifier = strings.TrimSpace(input.Email)
	}
	if identifier == "" {
		identifier = strings.TrimSpace(input.Username)
	}

	if identifier == "" {
		return nil, errors.New("email atau username wajib diisi")
	}
	if input.Password == "" {
		return nil, errors.New("password wajib diisi")
	}

	cleanIdentifier := strings.ToLower(identifier)
	var u *domainuser.User
	var err error

	// Support login by email OR username
	if strings.Contains(cleanIdentifier, "@") {
		u, err = s.userRepo.GetByEmail(ctx, cleanIdentifier)
		if err != nil {
			return nil, err
		}
		if u == nil {
			u, _ = s.userRepo.GetByUsername(ctx, cleanIdentifier)
		}
	} else {
		u, err = s.userRepo.GetByUsername(ctx, cleanIdentifier)
		if err != nil {
			return nil, err
		}
		if u == nil {
			u, _ = s.userRepo.GetByEmail(ctx, cleanIdentifier)
		}
	}

	// Constant-time execution against timing attacks (user enumeration)
	if u == nil {
		_ = bcrypt.CompareHashAndPassword(dummyBcryptHash, []byte(input.Password))
		return nil, ErrInvalidAuth
	}

	if !u.CheckPassword(input.Password) {
		return nil, ErrInvalidAuth
	}

	token, err := s.tokenService.GenerateToken(u.ID, u.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  u,
		Token: token,
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*appuser.UserProfileView, error) {
	view, err := s.userRepo.GetProfileView(ctx, userID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, common.ErrNotFound
	}
	return view, nil
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, common.ErrNotFound
	}
	return u, nil
}

func (s *AuthService) GetUserProfileByUsername(ctx context.Context, username string) (*appuser.UserProfileView, error) {
	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, common.ErrNotFound
	}
	view, err := s.userRepo.GetProfileView(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return &appuser.UserProfileView{
			User: u,
		}, nil
	}
	return view, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*domainuser.User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, common.ErrNotFound
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		u.Name = strings.TrimSpace(*input.Name)
	}
	if input.Bio != nil {
		u.Bio = input.Bio
	}
	if input.Location != nil {
		u.Location = input.Location
	}

	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthService) UpdateAvatar(ctx context.Context, userID uuid.UUID, objectKey string) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return common.ErrNotFound
	}

	u.AvatarKey = &objectKey
	return s.userRepo.Update(ctx, u)
}
