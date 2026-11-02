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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			dto.Error(c, http.StatusUnauthorized, "Missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			dto.Error(c, http.StatusUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims, err := tokenService.ValidateToken(tokenStr)
		if err != nil {
			dto.Error(c, http.StatusUnauthorized, "Invalid or expired token")
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
