package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"teman-mancing/backend/internal/application/auth"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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
