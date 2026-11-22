package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	appnotif "teman-mancing/backend/internal/application/notification"
	appspot "teman-mancing/backend/internal/application/fishing_spot"
	"teman-mancing/backend/internal/delivery/http/dto"
	"teman-mancing/backend/internal/delivery/http/middleware"
	domainspot "teman-mancing/backend/internal/domain/fishing_spot"
	domainnotif "teman-mancing/backend/internal/domain/notification"
	"teman-mancing/backend/internal/infrastructure/storage"
)

type SpotHandler struct {
	spotService  *appspot.Service
	minioClient  *storage.MinIOClient
	notifService *appnotif.Service
}

func NewSpotHandler(spotService *appspot.Service, minioClient *storage.MinIOClient, notifServices ...*appnotif.Service) *SpotHandler {
	var ns *appnotif.Service
	if len(notifServices) > 0 {
		ns = notifServices[0]
	}
	return &SpotHandler{
		spotService:  spotService,
		minioClient:  minioClient,
		notifService: ns,
	}
}

type CreateSpotRequest struct {
	Name             string                     `json:"name" binding:"required"`
	Description      string                     `json:"description" binding:"required"`
	Address          string                     `json:"address" binding:"required"`
	Latitude         float64                    `json:"latitude" binding:"required"`
	Longitude        float64                    `json:"longitude" binding:"required"`
	FishTypes        []string                   `json:"fish_types"`
	Techniques       []string                   `json:"techniques"`
	RecommendedBaits []domainspot.RecommendedBait `json:"recommended_baits"`
}

type UpdateSpotRequest struct {
	Name             *string                    `json:"name"`
	Description      *string                    `json:"description"`
	Address          *string                    `json:"address"`
	Latitude         *float64                   `json:"latitude"`
	Longitude        *float64                   `json:"longitude"`
	FishTypes        []string                   `json:"fish_types"`
	Techniques       []string                   `json:"techniques"`
	RecommendedBaits []domainspot.RecommendedBait `json:"recommended_baits"`
}

type CreateCommentRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	PhotoURL *string    `json:"photo_url"`
	Comment  string     `json:"comment"`
}

// GET /spots
func (h *SpotHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	query := c.Query("q")
	fishType := c.Query("fish_type")
	technique := c.Query("technique")
	sort := c.Query("sort")

	var filterUserID *uuid.UUID
	if uidStr := c.Query("user_id"); uidStr != "" {
		if uid, err := uuid.Parse(uidStr); err == nil {
			filterUserID = &uid
		}
	}

	spots, total, err := h.spotService.ListSpots(c.Request.Context(), appspot.ListSpotsFilter{
		Page:      page,
		Limit:     limit,
		Query:     query,
		UserID:    filterUserID,
		FishType:  fishType,
		Technique: technique,
		Sort:      sort,
	})
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve fishing spots")
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	dto.SuccessWithMeta(c, http.StatusOK, "Fishing spots retrieved", spots, gin.H{
		"page":        page,
		"limit":       limit,
		"total_items": total,
		"total_pages": totalPages,
	})
}

// GET /spots/:id
func (h *SpotHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	spotID, err := uuid.Parse(idStr)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID format")
		return
	}

	spot, err := h.spotService.GetSpot(c.Request.Context(), spotID)
	if err != nil {
		if errors.Is(err, appspot.ErrSpotNotFound) {
			dto.Error(c, http.StatusNotFound, "Fishing spot not found")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve spot")
		return
	}

	dto.Success(c, http.StatusOK, "Spot retrieved successfully", spot)
}

// POST /spots
func (h *SpotHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	var req CreateSpotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid input data", dto.FieldError{Field: "validation", Message: err.Error()})
		return
	}

	spot, err := h.spotService.CreateSpot(c.Request.Context(), appspot.CreateSpotInput{
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		Address:          req.Address,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		FishTypeNames:    req.FishTypes,
		TechniqueNames:   req.Techniques,
		RecommendedBaits: req.RecommendedBaits,
	})
	if err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(c, http.StatusCreated, "Fishing spot created successfully", spot)
}

// PATCH /spots/:id
func (h *SpotHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	idStr := c.Param("id")
	spotID, err := uuid.Parse(idStr)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID format")
		return
	}

	var req UpdateSpotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid update payload")
		return
	}

	spot, err := h.spotService.UpdateSpot(c.Request.Context(), spotID, userID, appspot.UpdateSpotInput{
		Name:             req.Name,
		Description:      req.Description,
		Address:          req.Address,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		FishTypeNames:    req.FishTypes,
		TechniqueNames:   req.Techniques,
		RecommendedBaits: req.RecommendedBaits,
	})
	if err != nil {
		if errors.Is(err, appspot.ErrSpotNotFound) {
			dto.Error(c, http.StatusNotFound, "Fishing spot not found")
			return
		}
		if errors.Is(err, appspot.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You do not have permission to edit this spot")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to update spot")
		return
	}

	dto.Success(c, http.StatusOK, "Spot updated successfully", spot)
}

// DELETE /spots/:id
func (h *SpotHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	idStr := c.Param("id")
	spotID, err := uuid.Parse(idStr)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID format")
		return
	}

	err = h.spotService.DeleteSpot(c.Request.Context(), spotID, userID)
	if err != nil {
		if errors.Is(err, appspot.ErrSpotNotFound) {
			dto.Error(c, http.StatusNotFound, "Fishing spot not found")
			return
		}
		if errors.Is(err, appspot.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You do not have permission to delete this spot")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to delete spot")
		return
	}

	dto.Success(c, http.StatusOK, "Spot deleted successfully", nil)
}

// POST /spots/:id/photos
func (h *SpotHandler) UploadPhoto(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	spotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID format")
		return
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Missing photo file in form-data")
		return
	}

	if fileHeader.Size > 5*1024*1024 {
		dto.Error(c, http.StatusBadRequest, "Photo file exceeds 5MB limit")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to read photo file")
		return
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".webp"
	}
	objectKey := fmt.Sprintf("fishing-spots/%s/%s%s", spotID.String(), uuid.New().String(), ext)

	if h.minioClient != nil && h.minioClient.Client != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		bucket := h.minioClient.Buckets.Spots
		if bucket == "" {
			bucket = "fishing-spots"
		}
		_, _ = h.minioClient.Client.PutObject(ctx, bucket, objectKey, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
	}

	isPrimary := c.DefaultPostForm("is_primary", "false") == "true"
	photo, err := h.spotService.AddPhoto(c.Request.Context(), spotID, userID, objectKey, isPrimary)
	if err != nil {
		if errors.Is(err, appspot.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You do not own this spot")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to register photo")
		return
	}

	dto.Success(c, http.StatusCreated, "Photo uploaded successfully", photo)
}

// DELETE /spots/:id/photos/:photoId
func (h *SpotHandler) DeletePhoto(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	spotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID")
		return
	}

	photoID, err := uuid.Parse(c.Param("photoId"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	if err := h.spotService.DeletePhoto(c.Request.Context(), spotID, photoID, userID); err != nil {
		if errors.Is(err, appspot.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You do not own this spot")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to delete photo")
		return
	}

	dto.Success(c, http.StatusOK, "Photo deleted successfully", nil)
}

// GET /spots/:id/comments
func (h *SpotHandler) ListComments(c *gin.Context) {
	spotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	comments, total, err := h.spotService.ListComments(c.Request.Context(), spotID, page, limit)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to retrieve comments")
		return
	}

	dto.SuccessWithMeta(c, http.StatusOK, "Comments retrieved", comments, gin.H{
		"page":        page,
		"limit":       limit,
		"total_items": total,
	})
}

// POST /spots/:id/comments
func (h *SpotHandler) CreateComment(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userName, _ := middleware.GetAuthUsername(c)

	spotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid comment format")
		return
	}

	comment, err := h.spotService.AddComment(c.Request.Context(), spotID, userID, userName, nil, req.ParentID, req.PhotoURL, req.Comment)
	if err != nil {
		if errors.Is(err, appspot.ErrSpotNotFound) {
			dto.Error(c, http.StatusNotFound, "Spot not found")
			return
		}
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Trigger notifications for new comments (Issue #24)
	if h.notifService != nil {
		go func() {
			ctx := context.Background()
			spot, err := h.spotService.GetSpot(ctx, spotID)
			if err != nil || spot == nil {
				return
			}

			notifiedUsers := make(map[uuid.UUID]bool)
			notifiedUsers[userID] = true // Never notify the author of the comment

			// 1. If this is a reply to another comment, notify the parent comment's author
			if req.ParentID != nil {
				if parent, err := h.spotService.GetCommentByID(ctx, *req.ParentID); err == nil && parent != nil {
					if !notifiedUsers[parent.UserID] {
						notifiedUsers[parent.UserID] = true
						_, _ = h.notifService.CreateNotification(
							ctx,
							parent.UserID,
							domainnotif.TypeSpotComment,
							"Balasan Komentar Baru",
							fmt.Sprintf("%s membalas komentar Anda di spot \"%s\"", userName, spot.Name),
							&spotID,
						)
					}
				}
			}

			// 2. Notify the spot creator/author if not already notified
			if !notifiedUsers[spot.UserID] {
				notifiedUsers[spot.UserID] = true
				_, _ = h.notifService.CreateNotification(
					ctx,
					spot.UserID,
					domainnotif.TypeSpotComment,
					"Komentar Baru pada Spot Anda",
					fmt.Sprintf("%s menambahkan komentar baru pada spot \"%s\"", userName, spot.Name),
					&spotID,
				)
			}

			// 3. Notify previous commenters of the spot
			commenterIDs, err := h.spotService.GetCommenterIDs(ctx, spotID)
			if err == nil {
				for _, cid := range commenterIDs {
					if !notifiedUsers[cid] {
						notifiedUsers[cid] = true
						_, _ = h.notifService.CreateNotification(
							ctx,
							cid,
							domainnotif.TypeSpotComment,
							"Diskusi Baru di Spot Mancing",
							fmt.Sprintf("%s ikut berdiskusi pada spot \"%s\"", userName, spot.Name),
							&spotID,
						)
					}
				}
			}
		}()
	}

	dto.Success(c, http.StatusCreated, "Comment posted successfully", comment)
}

// POST /spots/:id/comments/photo
func (h *SpotHandler) UploadCommentPhoto(c *gin.Context) {
	_, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	spotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid spot ID format")
		return
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Missing photo file in form-data")
		return
	}

	if fileHeader.Size > 5*1024*1024 {
		dto.Error(c, http.StatusBadRequest, "Photo file exceeds 5MB limit")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, "Failed to read photo file")
		return
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".webp"
	}
	objectKey := fmt.Sprintf("fishing-spots/%s/comments/%s%s", spotID.String(), uuid.New().String(), ext)

	if h.minioClient != nil && h.minioClient.Client != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		bucket := h.minioClient.Buckets.Spots
		if bucket == "" {
			bucket = "fishing-spots"
		}
		_, err = h.minioClient.Client.PutObject(ctx, bucket, objectKey, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
		if err != nil {
			dto.Error(c, http.StatusInternalServerError, "Failed to upload photo to storage")
			return
		}
	}

	dto.Success(c, http.StatusOK, "Photo uploaded successfully", gin.H{
		"photo_url": objectKey,
	})
}

// DELETE /comments/:id
func (h *SpotHandler) DeleteComment(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		dto.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	if err := h.spotService.DeleteComment(c.Request.Context(), commentID, userID); err != nil {
		if errors.Is(err, appspot.ErrCommentNotFound) {
			dto.Error(c, http.StatusNotFound, "Comment not found")
			return
		}
		if errors.Is(err, appspot.ErrUnauthorized) {
			dto.Error(c, http.StatusForbidden, "You can only delete your own comments")
			return
		}
		dto.Error(c, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	dto.Success(c, http.StatusOK, "Comment deleted successfully", nil)
}
