package fishing_spot

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	domainspot "teman-mancing/backend/internal/domain/fishing_spot"
)

var (
	ErrSpotNotFound    = errors.New("fishing spot not found")
	ErrCommentNotFound = errors.New("comment not found")
	ErrPhotoNotFound   = errors.New("photo not found")
	ErrUnauthorized    = errors.New("you are not authorized to modify this resource")
)

type CreateSpotInput struct {
	UserID           uuid.UUID                 `json:"user_id"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description"`
	Address          string                    `json:"address"`
	Latitude         float64                   `json:"latitude"`
	Longitude        float64                   `json:"longitude"`
	FishTypeNames    []string                  `json:"fish_types"`
	TechniqueNames   []string                  `json:"techniques"`
	RecommendedBaits []domainspot.RecommendedBait `json:"recommended_baits"`
}

type UpdateSpotInput struct {
	Name             *string                   `json:"name"`
	Description      *string                   `json:"description"`
	Address          *string                   `json:"address"`
	Latitude         *float64                  `json:"latitude"`
	Longitude        *float64                  `json:"longitude"`
	FishTypeNames    []string                  `json:"fish_types"`
	TechniqueNames   []string                  `json:"techniques"`
	RecommendedBaits []domainspot.RecommendedBait `json:"recommended_baits"`
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListSpots(ctx context.Context, filter ListSpotsFilter) ([]*domainspot.FishingSpot, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) GetSpot(ctx context.Context, id uuid.UUID) (*domainspot.FishingSpot, error) {
	spot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if spot == nil {
		return nil, ErrSpotNotFound
	}
	return spot, nil
}

func (s *Service) CreateSpot(ctx context.Context, input CreateSpotInput) (*domainspot.FishingSpot, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("spot name is required")
	}
	if strings.TrimSpace(input.Description) == "" {
		return nil, errors.New("spot description is required")
	}
	if strings.TrimSpace(input.Address) == "" {
		return nil, errors.New("spot address is required")
	}

	spotID := uuid.New()
	now := time.Now().UTC()

	var fishTypes []domainspot.FishType
	for _, fn := range input.FishTypeNames {
		fn = strings.TrimSpace(fn)
		if fn != "" {
			fishTypes = append(fishTypes, domainspot.FishType{
				ID:   uuid.New(),
				Name: fn,
			})
		}
	}

	var techniques []domainspot.FishingTechnique
	for _, tn := range input.TechniqueNames {
		tn = strings.TrimSpace(tn)
		if tn != "" {
			techniques = append(techniques, domainspot.FishingTechnique{
				ID:   uuid.New(),
				Name: tn,
			})
		}
	}

	var baits []domainspot.RecommendedBait
	for _, b := range input.RecommendedBaits {
		name := strings.TrimSpace(b.Name)
		if name != "" {
			baits = append(baits, domainspot.RecommendedBait{
				ID:          uuid.New(),
				SpotID:      spotID,
				Name:        name,
				Description: b.Description,
			})
		}
	}

	spot := &domainspot.FishingSpot{
		ID:               spotID,
		UserID:           input.UserID,
		Name:             input.Name,
		Description:      input.Description,
		Address:          input.Address,
		Latitude:         input.Latitude,
		Longitude:        input.Longitude,
		FishTypes:        fishTypes,
		Techniques:       techniques,
		RecommendedBaits: baits,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, spot); err != nil {
		return nil, err
	}

	return spot, nil
}

func (s *Service) UpdateSpot(ctx context.Context, spotID, actorID uuid.UUID, input UpdateSpotInput) (*domainspot.FishingSpot, error) {
	spot, err := s.repo.GetByID(ctx, spotID)
	if err != nil {
		return nil, err
	}
	if spot == nil {
		return nil, ErrSpotNotFound
	}

	// Ownership validation (README Section 12.4)
	if spot.UserID != actorID {
		return nil, ErrUnauthorized
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		spot.Name = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil && strings.TrimSpace(*input.Description) != "" {
		spot.Description = strings.TrimSpace(*input.Description)
	}
	if input.Address != nil && strings.TrimSpace(*input.Address) != "" {
		spot.Address = strings.TrimSpace(*input.Address)
	}
	if input.Latitude != nil {
		spot.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		spot.Longitude = *input.Longitude
	}

	spot.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, spot); err != nil {
		return nil, err
	}

	return spot, nil
}

func (s *Service) DeleteSpot(ctx context.Context, spotID, actorID uuid.UUID) error {
	spot, err := s.repo.GetByID(ctx, spotID)
	if err != nil {
		return err
	}
	if spot == nil {
		return ErrSpotNotFound
	}

	// Ownership validation
	if spot.UserID != actorID {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, spotID)
}

// Comments Management
func (s *Service) ListComments(ctx context.Context, spotID uuid.UUID, page, limit int) ([]*domainspot.Comment, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListComments(ctx, spotID, page, limit)
}

func (s *Service) AddComment(ctx context.Context, spotID, userID uuid.UUID, userName string, userAvatar *string, commentText string) (*domainspot.Comment, error) {
	cleanText := strings.TrimSpace(commentText)
	if cleanText == "" {
		return nil, errors.New("comment content cannot be empty")
	}

	// Verify spot exists
	spot, err := s.repo.GetByID(ctx, spotID)
	if err != nil {
		return nil, err
	}
	if spot == nil {
		return nil, ErrSpotNotFound
	}

	now := time.Now().UTC()
	c := &domainspot.Comment{
		ID:         uuid.New(),
		SpotID:     spotID,
		UserID:     userID,
		UserName:   userName,
		UserAvatar: userAvatar,
		Comment:    cleanText,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateComment(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) DeleteComment(ctx context.Context, commentID, actorID uuid.UUID) error {
	comment, err := s.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrCommentNotFound
	}

	// Ownership validation: only author can delete comment (README Section 12.4)
	if comment.UserID != actorID {
		return ErrUnauthorized
	}

	return s.repo.DeleteComment(ctx, commentID)
}

// Photos Management
func (s *Service) AddPhoto(ctx context.Context, spotID, actorID uuid.UUID, objectKey string, isPrimary bool) (*domainspot.Photo, error) {
	spot, err := s.repo.GetByID(ctx, spotID)
	if err != nil {
		return nil, err
	}
	if spot == nil {
		return nil, ErrSpotNotFound
	}

	if spot.UserID != actorID {
		return nil, ErrUnauthorized
	}

	photo := &domainspot.Photo{
		ID:        uuid.New(),
		SpotID:    spotID,
		ObjectKey: objectKey,
		IsPrimary: isPrimary,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.AddPhoto(ctx, photo); err != nil {
		return nil, err
	}

	return photo, nil
}

func (s *Service) DeletePhoto(ctx context.Context, spotID, photoID, actorID uuid.UUID) error {
	spot, err := s.repo.GetByID(ctx, spotID)
	if err != nil {
		return err
	}
	if spot == nil {
		return ErrSpotNotFound
	}

	if spot.UserID != actorID {
		return ErrUnauthorized
	}

	return s.repo.DeletePhoto(ctx, photoID)
}

func (s *Service) MasterCatalogs(ctx context.Context) ([]domainspot.FishType, []domainspot.FishingTechnique, error) {
	fishes, err := s.repo.ListFishTypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	techniques, err := s.repo.ListTechniques(ctx)
	if err != nil {
		return nil, nil, err
	}
	return fishes, techniques, nil
}
