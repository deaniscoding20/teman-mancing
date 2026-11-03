package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"teman-mancing/backend/internal/application/auth"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
)

type AuthHandler struct {
	authService        *auth.AuthService
	googleClientID     string
	googleClientSecret string
	googleRedirectURL  string
}

func NewAuthHandler(authService *auth.AuthService, googleClientID, googleClientSecret, googleRedirectURL string) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		googleRedirectURL:  googleRedirectURL,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid registration input", dto.FieldError{
			Field:   "validation",
			Message: err.Error(),
		})
		return
	}

	result, err := h.authService.Register(c.Request.Context(), auth.RegisterInput{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrEmailExists) {
			dto.Error(c, http.StatusConflict, "Email is already registered", dto.FieldError{Field: "email", Message: err.Error()})
			return
		}
		if errors.Is(err, auth.ErrUsernameExists) {
			dto.Error(c, http.StatusConflict, "Username is already taken", dto.FieldError{Field: "username", Message: err.Error()})
			return
		}
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusCreated, "User registered successfully", result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Format login tidak valid", dto.FieldError{
			Field:   "validation",
			Message: err.Error(),
		})
		return
	}

	ident := req.Identifier
	if ident == "" {
		ident = req.Email
	}
	if ident == "" {
		ident = req.Username
	}

	if ident == "" {
		dto.Error(c, http.StatusBadRequest, "Email atau username wajib diisi", dto.FieldError{
			Field:   "identifier",
			Message: "Email atau username tidak boleh kosong",
		})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), auth.LoginInput{
		Identifier: ident,
		Email:      req.Email,
		Username:   req.Username,
		Password:   req.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrEmailNotVerified) {
			dto.Error(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, auth.ErrInvalidAuth) {
			dto.Error(c, http.StatusUnauthorized, "Email/username atau password salah. Pastikan kredensial Anda benar.")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Terjadi kesalahan internal saat login")
		return
	}

	// Set HttpOnly session cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("teman_mancing_token", result.Token, 86400*7, "/", "", false, true)

	dto.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear session cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("teman_mancing_token", "", -1, "/", "", false, true)
	dto.Success(c, http.StatusOK, "Logged out successfully", nil)
}

type GoogleLoginRequest struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid request payload format", dto.FieldError{
			Field:   "body",
			Message: err.Error(),
		})
		return
	}

	result, err := h.authService.GoogleLogin(c.Request.Context(), auth.GoogleLoginInput{
		IDToken:     req.IDToken,
		AccessToken: req.AccessToken,
		Email:       req.Email,
		Name:        req.Name,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set HttpOnly session cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("teman_mancing_token", result.Token, 86400*7, "/", "", false, true)

	dto.Success(c, http.StatusOK, "Login Google berhasil", result)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	profile, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, http.StatusNotFound, "User not found")
		return
	}

	dto.Success(c, http.StatusOK, "Profile retrieved successfully", profile)
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	_ = c.ShouldBindJSON(&req)

	token := req.Token
	if token == "" {
		token = c.Query("token")
	}

	result, err := h.authService.VerifyEmail(c.Request.Context(), token, req.Email, req.OTP)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set session cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("teman_mancing_token", result.Token, 86400*7, "/", "", false, true)

	dto.Success(c, http.StatusOK, "Email berhasil diverifikasi", result)
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Format email tidak valid", dto.FieldError{
			Field:   "email",
			Message: err.Error(),
		})
		return
	}

	if err := h.authService.ResendVerificationEmail(c.Request.Context(), req.Email); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusOK, "Kode verifikasi baru telah dikirimkan ke email Anda", nil)
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Format email tidak valid", dto.FieldError{
			Field:   "email",
			Message: err.Error(),
		})
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		dto.Error(c, http.StatusInternalServerError, "Gagal memproses permintaan reset password")
		return
	}

	dto.Success(c, http.StatusOK, "Instruksi dan kode reset password telah dikirim ke email Anda", nil)
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	NewPassword string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Input reset password tidak valid", dto.FieldError{
			Field:   "validation",
			Message: err.Error(),
		})
		return
	}

	token := req.Token
	if token == "" {
		token = c.Query("token")
	}

	if err := h.authService.ResetPassword(c.Request.Context(), token, req.Email, req.OTP, req.NewPassword); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusOK, "Kata sandi Anda berhasil diperbarui. Silakan login kembali.", nil)
}

// GET /api/v1/auth/google/login
func (h *AuthHandler) GoogleOAuthRedirect(c *gin.Context) {
	if h.googleClientID == "" {
		dto.Error(c, http.StatusBadRequest, "Google Client ID belum dikonfigurasi di server")
		return
	}
	redirectURI := h.googleRedirectURL
	if redirectURI == "" {
		redirectURI = "https://temanmancing.my.id/api/v1/auth/google/callback"
	}
	scope := url.QueryEscape("openid email profile")
	googleAuthURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=offline&prompt=select_account",
		url.QueryEscape(h.googleClientID),
		url.QueryEscape(redirectURI),
		scope,
	)
	c.Redirect(http.StatusTemporaryRedirect, googleAuthURL)
}

// GET /api/v1/auth/google/callback
func (h *AuthHandler) GoogleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=Otorisasi+Google+dibatalkan")
		return
	}

	tokenURL := "https://oauth2.googleapis.com/token"
	redirectURI := h.googleRedirectURL
	if redirectURI == "" {
		redirectURI = "https://temanmancing.my.id/api/v1/auth/google/callback"
	}

	form := url.Values{
		"code":          {code},
		"client_id":     {h.googleClientID},
		"client_secret": {h.googleClientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm(tokenURL, form)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=Gagal+menukar+kode+Google")
		return
	}
	defer resp.Body.Close()

	var tokResult struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokResult); err != nil || tokResult.IDToken == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=Response+token+Google+tidak+valid")
		return
	}

	authRes, err := h.authService.GoogleLogin(c.Request.Context(), auth.GoogleLoginInput{
		IDToken: tokResult.IDToken,
	})
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error="+url.QueryEscape(err.Error()))
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("teman_mancing_token", authRes.Token, 86400*7, "/", "", false, true)

	userJSON, _ := json.Marshal(authRes.User)
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(
		"/auth/callback?token=%s&user=%s",
		url.QueryEscape(authRes.Token),
		url.QueryEscape(string(userJSON)),
	))
}
