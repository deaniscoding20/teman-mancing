package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"teman-mancing/backend/internal/delivery/http/dto"
)

// VersionHandler returns build information.
type VersionHandler struct {
	version string
	env     string
}

// NewVersionHandler creates a new VersionHandler.
func NewVersionHandler(version, env string) *VersionHandler {
	return &VersionHandler{
		version: version,
		env:     env,
	}
}

// Info returns the service version and platform metadata.
func (h *VersionHandler) Info(c *gin.Context) {
	dto.Success(c, http.StatusOK, "Teman Mancing Backend API", gin.H{
		"name":        "Teman Mancing API",
		"version":     h.version,
		"environment": h.env,
		"tagline":     "Cari Spot. Cari Teman. Gas Mancing!",
	})
}
