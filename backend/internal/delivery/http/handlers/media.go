package handlers

import (
	"context"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/infrastructure/storage"
)

type MediaHandler struct {
	minioClient *storage.MinIOClient
}

func NewMediaHandler(minioClient *storage.MinIOClient) *MediaHandler {
	return &MediaHandler{
		minioClient: minioClient,
	}
}

// GET /media/:bucket/*objectKey
func (h *MediaHandler) ServeMedia(c *gin.Context) {
	bucket := c.Param("bucket")
	objectKey := strings.TrimPrefix(c.Param("objectKey"), "/")

	if bucket == "" || objectKey == "" {
		dto.Error(c, http.StatusBadRequest, "Invalid bucket or object key parameter")
		return
	}

	if h.minioClient == nil || h.minioClient.Client == nil {
		dto.Error(c, http.StatusServiceUnavailable, "Object storage client is not available")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	obj, err := h.minioClient.Client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		dto.Error(c, http.StatusNotFound, "Object not found")
		return
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		dto.Error(c, http.StatusNotFound, "Object not found")
		return
	}

	contentType := stat.ContentType
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = mime.TypeByExtension(filepath.Ext(objectKey))
		if contentType == "" {
			contentType = "image/jpeg"
		}
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.DataFromReader(http.StatusOK, stat.Size, contentType, obj, nil)
}
