package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

const presignTTL = 15 * time.Minute

type StorageService interface {
	GeneratePresignedPutURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
}

type UploadCharacterRepository interface {
	GetByIDsAndUserID(ctx context.Context, ids []string, userID string) ([]domain.Character, error)
}

type InitialUploadParams struct {
	UserID       string
	MimeType     string
	Title        *string
	ArtistName   *string
	ArtistLink   *string
	Notes        *string
	CharacterIDs []string
}

type InitialUploadResult struct {
	ID        uuid.UUID
	UploadURL string
	ExpiresAt time.Time
}

type uploadUsecase struct {
	uploadRepo UploadRepository
	storage    StorageService
	charLookup UploadCharacterRepository
	tel        *observability.Telemetry
}

func NewUploadUsecase(uploadRepo UploadRepository, storage StorageService, charLookup UploadCharacterRepository, tel *observability.Telemetry) *uploadUsecase {
	return &uploadUsecase{
		uploadRepo: uploadRepo,
		storage:    storage,
		charLookup: charLookup,
		tel:        tel,
	}
}

func (u *uploadUsecase) InitialUpload(ctx context.Context, params InitialUploadParams) (*InitialUploadResult, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.InitialUpload")
	defer span.End()

	id := uuid.New()
	ext := mimeTypeToExt(params.MimeType)
	r2Key := fmt.Sprintf("users/%s/images/%s%s", params.UserID, id.String(), ext)
	expiresAt := time.Now().Add(presignTTL)

	uploadURL, err := u.storage.GeneratePresignedPutURL(ctx, r2Key, params.MimeType, presignTTL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pending := &domain.PendingUpload{
		ID:           id,
		UserID:       params.UserID,
		R2Key:        r2Key,
		MimeType:     params.MimeType,
		Title:        params.Title,
		ArtistName:   params.ArtistName,
		ArtistLink:   params.ArtistLink,
		Notes:        params.Notes,
		CharacterIDs: params.CharacterIDs,
	}
	if err := u.uploadRepo.CreatePendingUpload(ctx, pending); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &InitialUploadResult{ID: id, UploadURL: uploadURL, ExpiresAt: expiresAt}, nil
}

func (u *uploadUsecase) CompleteUpload(ctx context.Context, pendingID, userID string) (*domain.Image, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CompleteUpload")
	defer span.End()

	pending, err := u.uploadRepo.GetPendingUpload(ctx, pendingID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	validCharIDs := make([]uuid.UUID, 0)
	if len(pending.CharacterIDs) > 0 {
		chars, err := u.charLookup.GetByIDsAndUserID(ctx, pending.CharacterIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		for _, c := range chars {
			validCharIDs = append(validCharIDs, c.ID)
		}
	}

	image, err := u.uploadRepo.CompleteUpload(ctx, pendingID, userID, validCharIDs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// TODO: trigger thumbnail generation job

	return image, nil
}

func mimeTypeToExt(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}
