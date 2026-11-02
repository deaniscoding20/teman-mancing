package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/infrastructure/database"
	"teman-mancing/backend/internal/infrastructure/storage"
)

// HealthHandler manages health and readiness probes.
type HealthHandler struct {
	db    *database.PostgresDB
	minio *storage.MinIOClient
}

// NewHealthHandler returns a new HealthHandler instance.
func NewHealthHandler(db *database.PostgresDB, minio *storage.MinIOClient) *HealthHandler {
	return &HealthHandler{
		db:    db,
		minio: minio,
	}
}

// Status checks the system, database, and MinIO readiness.
func (h *HealthHandler) Status(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "healthy"
	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			dbStatus = "unhealthy: " + err.Error()
		}
	} else {
		dbStatus = "not_configured"
	}

	minioStatus := "healthy"
	if h.minio != nil {
		if err := h.minio.Ping(ctx); err != nil {
			minioStatus = "unhealthy: " + err.Error()
		}
	} else {
		minioStatus = "not_configured"
	}

	isHealthy := (dbStatus == "healthy" || dbStatus == "not_configured") &&
		(minioStatus == "healthy" || minioStatus == "not_configured")

	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	dto.Success(c, statusCode, "System health status", gin.H{
		"status":   "running",
		"postgres": dbStatus,
		"minio":    minioStatus,
		"time":     time.Now().UTC().Format(time.RFC3339),
	})
}
