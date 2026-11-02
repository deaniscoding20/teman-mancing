package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"teman-mancing/backend/internal/delivery/http/dto"
)

// ZapRecovery catches any panics and returns a structured 500 error response.
func ZapRecovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errStr := fmt.Sprintf("%v", r)
				log.Error("Unhandled panic recovered in HTTP handler",
					zap.String("error", errStr),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				dto.Error(c, http.StatusInternalServerError, "An unexpected server error occurred", dto.FieldError{
					Field:   "server",
					Message: "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
