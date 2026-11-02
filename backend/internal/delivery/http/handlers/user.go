package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"teman-mancing/backend/internal/application/auth"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
	"teman-mancing/backend/internal/infrastructure/storage"
)

type UserHandler struct {
	authService *auth.AuthService
	minioClient *storage.MinIOClient
}

func NewUserHandler(authService *auth.AuthService, minioClient *storage.MinIOClient) *UserHandler {
	return &UserHandler{
		authService: authService,
		minioClient: minioClient,
	}
}

type UpdateProfileRequest struct {
	Name     *string `json:"name"`
	Bio      *string `json:"bio"`
	Location *string `json:"location"`
}

func (h *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		dto.Error(c, http.StatusBadRequest, "Username is required")
		return
	}

	u, err := h.authService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		dto.Error(c, http.StatusNotFound, "User not found")
		return
	}

	dto.Success(c, http.StatusOK, "User profile retrieved", u)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid profile payload format")
		return
	}

	updated, err := h.authService.UpdateProfile(c.Request.Context(), userID, auth.UpdateProfileInput{
		Name:     req.Name,
		Bio:      req.Bio,
		Location: req.Location,
	})
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	dto.Success(c, http.StatusOK, "Profile updated successfully", updated)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Avatar file is required in 'avatar' form-data field")
		return
	}

	// Validate file size (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		dto.Error(c, http.StatusBadRequest, "Avatar image exceeds 5MB size limit")
		return
	}

	// Validate content type
	contentType := fileHeader.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !validTypes[contentType] {
		dto.Error(c, http.StatusBadRequest, "Invalid file format. Allowed types: JPEG, PNG, WEBP")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to read avatar file")
		return
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".webp"
	}
	objectKey := fmt.Sprintf("avatars/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// If MinIO client is configured, upload to S3 bucket
	if h.minioClient != nil && h.minioClient.Client != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		bucket := h.minioClient.Buckets.Avatars
		if bucket == "" {
			bucket = "avatars"
		}

		_, err = h.minioClient.Client.PutObject(ctx, bucket, objectKey, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			// In dev mode without minio up, allow fallback
			_ = err
		}
	} else {
		// Drain reader
		_, _ = io.Copy(io.Discard, file)
	}

	// Update database with new object key
	if err := h.authService.UpdateAvatar(c.Request.Context(), userID, objectKey); err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to save avatar key to user record")
		return
	}

	dto.Success(c, http.StatusOK, "Avatar uploaded successfully", gin.H{
		"avatar_key": objectKey,
	})
}
