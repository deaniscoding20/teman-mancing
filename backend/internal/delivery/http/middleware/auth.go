package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/infrastructure/security"
)

const (
	CtxKeyUserID   = "auth_user_id"
	CtxKeyUsername = "auth_username"
)

// RequireAuth ensures that the incoming request contains a valid Bearer JWT.
func RequireAuth(tokenService *security.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenStr = strings.TrimSpace(parts[1])
			}
		}

		// Fallback to cookie
		if tokenStr == "" {
			if cookieVal, err := c.Cookie("teman_mancing_token"); err == nil && cookieVal != "" {
				tokenStr = strings.TrimSpace(cookieVal)
			}
		}

		// Fallback to query parameter (e.g. for WebSocket handshakes /api/v1/ws?token=...)
		if tokenStr == "" {
			if queryToken := strings.TrimSpace(c.Query("token")); queryToken != "" {
				tokenStr = queryToken
			}
		}

		if tokenStr == "" {
			dto.Error(c, http.StatusUnauthorized, "Sesi tidak ditemukan atau header otorisasi hilang")
			c.Abort()
			return
		}

		claims, err := tokenService.ValidateToken(tokenStr)
		if err != nil {
			dto.Error(c, http.StatusUnauthorized, "Token otorisasi tidak valid atau telah kedaluwarsa")
			c.Abort()
			return
		}

		c.Set(CtxKeyUserID, claims.UserID)
		c.Set(CtxKeyUsername, claims.Username)
		c.Next()
	}
}

// GetAuthUserID retrieves the authenticated user's UUID from Gin context.
func GetAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(CtxKeyUserID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// GetAuthUsername retrieves the authenticated username from Gin context.
func GetAuthUsername(c *gin.Context) (string, bool) {
	val, exists := c.Get(CtxKeyUsername)
	if !exists {
		return "", false
	}
	name, ok := val.(string)
	return name, ok
}
