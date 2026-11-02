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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
		dto.Error(c, http.StatusBadRequest, "Invalid login credentials format", dto.FieldError{
			Field:   "validation",
			Message: err.Error(),
		})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidAuth) {
			dto.Error(c, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "An unexpected login error occurred")
		return
	}

	dto.Success(c, http.StatusOK, "Login successful", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
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
