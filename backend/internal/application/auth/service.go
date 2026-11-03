package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"teman-mancing/backend/internal/application/common"
	appuser "teman-mancing/backend/internal/application/user"
	domainuser "teman-mancing/backend/internal/domain/user"
	"teman-mancing/backend/internal/infrastructure/email"
	"teman-mancing/backend/internal/infrastructure/security"
)

var (
	ErrEmailExists       = errors.New("email is already registered")
	ErrUsernameExists    = errors.New("username is already taken")
	ErrInvalidAuth       = errors.New("email/username atau password salah")
	ErrEmailNotVerified  = errors.New("email belum diverifikasi. Kami telah mengirimkan kode verifikasi baru ke email Anda")
	ErrInvalidOTPOrToken = errors.New("kode OTP atau token verifikasi tidak valid atau telah kedaluwarsa")
	dummyBcryptHash, _   = bcrypt.GenerateFromPassword([]byte("dummy-auth-salt-constant-time"), bcrypt.DefaultCost)
)

func generateOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return fmt.Sprintf("%06d", time.Now().UnixNano()%900000+100000)
	}
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

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

type GoogleLoginInput struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
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
	mailer       *email.Mailer
}

func NewAuthService(userRepo appuser.Repository, tokenService *security.TokenService, mailer *email.Mailer) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
		mailer:       mailer,
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

	// Email verification setup: default is_verified = false
	newUser.IsVerified = false

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Create verification token and OTP
	otpCode := generateOTP()
	verifyToken := uuid.New().String()
	ver := &domainuser.EmailVerification{
		ID:        uuid.New(),
		UserID:    newUser.ID,
		OTPCode:   otpCode,
		Token:     verifyToken,
		ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
		CreatedAt: time.Now().UTC(),
	}
	_ = s.userRepo.CreateEmailVerification(ctx, ver)

	// Send verification email in background
	if s.mailer != nil {
		go func() {
			_ = s.mailer.SendVerificationEmail(newUser.Email, newUser.Name, otpCode, verifyToken)
		}()
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

	// Check if email has been verified
	if !u.IsVerified {
		otpCode := generateOTP()
		verifyToken := uuid.New().String()
		ver := &domainuser.EmailVerification{
			ID:        uuid.New(),
			UserID:    u.ID,
			OTPCode:   otpCode,
			Token:     verifyToken,
			ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
			CreatedAt: time.Now().UTC(),
		}
		_ = s.userRepo.CreateEmailVerification(ctx, ver)

		if s.mailer != nil {
			go func() {
				_ = s.mailer.SendVerificationEmail(u.Email, u.Name, otpCode, verifyToken)
			}()
		}
		return nil, ErrEmailNotVerified
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

func (s *AuthService) VerifyEmail(ctx context.Context, token, emailStr, otp string) (*AuthResult, error) {
	var targetUserID uuid.UUID
	var verID uuid.UUID

	if strings.TrimSpace(token) != "" {
		v, err := s.userRepo.GetEmailVerificationByToken(ctx, strings.TrimSpace(token))
		if err != nil || v == nil {
			return nil, ErrInvalidOTPOrToken
		}
		targetUserID = v.UserID
		verID = v.ID
	} else if strings.TrimSpace(emailStr) != "" && strings.TrimSpace(otp) != "" {
		u, err := s.userRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(emailStr)))
		if err != nil || u == nil {
			return nil, ErrInvalidOTPOrToken
		}
		v, err := s.userRepo.GetEmailVerificationByOTP(ctx, u.ID, strings.TrimSpace(otp))
		if err != nil || v == nil {
			return nil, ErrInvalidOTPOrToken
		}
		targetUserID = u.ID
		verID = v.ID
	} else {
		return nil, errors.New("token atau email dan kode OTP wajib diisi")
	}

	u, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil || u == nil {
		return nil, common.ErrNotFound
	}

	u.IsVerified = true
	now := time.Now().UTC()
	u.EmailVerifiedAt = &now

	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	_ = s.userRepo.DeleteEmailVerification(ctx, verID)

	authToken, err := s.tokenService.GenerateToken(u.ID, u.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  u,
		Token: authToken,
	}, nil
}

func (s *AuthService) ResendVerificationEmail(ctx context.Context, emailStr string) error {
	cleanEmail := strings.ToLower(strings.TrimSpace(emailStr))
	if cleanEmail == "" {
		return errors.New("email tidak boleh kosong")
	}

	u, err := s.userRepo.GetByEmail(ctx, cleanEmail)
	if err != nil || u == nil {
		// Silent success against user enumeration
		return nil
	}

	if u.IsVerified {
		return errors.New("email akun ini sudah berstatus terverifikasi")
	}

	otpCode := generateOTP()
	verifyToken := uuid.New().String()
	ver := &domainuser.EmailVerification{
		ID:        uuid.New(),
		UserID:    u.ID,
		OTPCode:   otpCode,
		Token:     verifyToken,
		ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
		CreatedAt: time.Now().UTC(),
	}
	_ = s.userRepo.CreateEmailVerification(ctx, ver)

	if s.mailer != nil {
		return s.mailer.SendVerificationEmail(u.Email, u.Name, otpCode, verifyToken)
	}
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, emailStr string) error {
	cleanEmail := strings.ToLower(strings.TrimSpace(emailStr))
	if cleanEmail == "" {
		return errors.New("email wajib diisi")
	}

	u, err := s.userRepo.GetByEmail(ctx, cleanEmail)
	if err != nil || u == nil {
		// Constant time against timing attacks / email discovery
		return nil
	}

	otpCode := generateOTP()
	resetToken := uuid.New().String()
	reset := &domainuser.PasswordReset{
		ID:        uuid.New(),
		UserID:    u.ID,
		Email:     cleanEmail,
		OTPCode:   otpCode,
		Token:     resetToken,
		ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
		CreatedAt: time.Now().UTC(),
	}

	if err := s.userRepo.CreatePasswordReset(ctx, reset); err != nil {
		return err
	}

	if s.mailer != nil {
		go func() {
			_ = s.mailer.SendPasswordResetEmail(u.Email, u.Name, otpCode, resetToken)
		}()
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, emailStr, otp, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password baru harus minimal 8 karakter")
	}

	var reset *domainuser.PasswordReset
	var err error

	if strings.TrimSpace(token) != "" {
		reset, err = s.userRepo.GetPasswordResetByToken(ctx, strings.TrimSpace(token))
	} else if strings.TrimSpace(emailStr) != "" && strings.TrimSpace(otp) != "" {
		reset, err = s.userRepo.GetPasswordResetByOTP(ctx, strings.TrimSpace(emailStr), strings.TrimSpace(otp))
	} else {
		return errors.New("token reset atau email dan kode OTP wajib diisi")
	}

	if err != nil || reset == nil {
		return errors.New("token atau kode OTP reset password tidak valid atau telah kedaluwarsa")
	}

	u, err := s.userRepo.GetByID(ctx, reset.UserID)
	if err != nil || u == nil {
		return common.ErrNotFound
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashed)
	u.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	_ = s.userRepo.MarkPasswordResetUsed(ctx, reset.ID)
	return nil
}

func (s *AuthService) GoogleLogin(ctx context.Context, input GoogleLoginInput) (*AuthResult, error) {
	email := strings.TrimSpace(input.Email)
	name := strings.TrimSpace(input.Name)

	// If IDToken is provided, verify with Google TokenInfo API
	if strings.TrimSpace(input.IDToken) != "" {
		reqURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", strings.TrimSpace(input.IDToken))
		client := &http.Client{Timeout: 8 * time.Second}
		resp, err := client.Get(reqURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var tokenInfo struct {
				Email   string `json:"email"`
				Name    string `json:"name"`
				Picture string `json:"picture"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err == nil && tokenInfo.Email != "" {
				email = tokenInfo.Email
				if tokenInfo.Name != "" {
					name = tokenInfo.Name
				}
			}
		}
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	if cleanEmail == "" {
		return nil, errors.New("email Google tidak valid atau tidak dapat diverifikasi")
	}

	if name == "" {
		parts := strings.Split(cleanEmail, "@")
		name = parts[0]
	}

	// 1. Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, cleanEmail)
	if err == nil && existingUser != nil {
		if !existingUser.IsVerified {
			existingUser.IsVerified = true
			now := time.Now().UTC()
			existingUser.EmailVerifiedAt = &now
			_ = s.userRepo.Update(ctx, existingUser)
		}
		token, err := s.tokenService.GenerateToken(existingUser.ID, existingUser.Username)
		if err != nil {
			return nil, err
		}
		return &AuthResult{
			User:  existingUser,
			Token: token,
		}, nil
	}

	// 2. Generate unique username from email
	cleanBase := strings.Split(cleanEmail, "@")[0]
	reg := regexp.MustCompile("[^a-z0-9_]")
	cleanBase = reg.ReplaceAllString(strings.ToLower(cleanBase), "")
	if len(cleanBase) < 3 {
		cleanBase = "angler"
	}
	if len(cleanBase) > 20 {
		cleanBase = cleanBase[:20]
	}

	candidateUsername := cleanBase
	existingByUsername, _ := s.userRepo.GetByUsername(ctx, candidateUsername)
	if existingByUsername != nil {
		candidateUsername = fmt.Sprintf("%s%d", cleanBase, time.Now().Unix()%10000)
	}

	// 3. Create user with secure random password
	randomPassword := uuid.New().String() + "A1!"
	newUser, err := domainuser.NewUser(name, candidateUsername, cleanEmail, randomPassword)
	if err != nil {
		return nil, err
	}
	newUser.IsVerified = true
	now := time.Now().UTC()
	newUser.EmailVerifiedAt = &now

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
