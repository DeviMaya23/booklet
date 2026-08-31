package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

const PresignTTL = 15 * time.Minute

type InitialUploadParams struct {
	UserID       uuid.UUID
	MimeType     string
	Title        *string
	ArtistID     *uuid.UUID
	Notes        *string
	CharacterIDs []uuid.UUID
}

type InitialUploadResult struct {
	ID        uuid.UUID
	UploadURL string
	ExpiresAt time.Time
}

type uploadUsecase struct {
	uploadRepo    UploadRepository
	storage       StorageService
	characterRepo UploadCharacterRepository
	artistRepo    UploadArtistRepository
	imageRepo     UploadImageRepository
	transactor    Transactor
	tel           *observability.Telemetry
	uploadCount   metric.Int64Counter
}

func NewUploadUsecase(
	uploadRepo UploadRepository,
	storage StorageService,
	characterRepo UploadCharacterRepository,
	artistRepo UploadArtistRepository,
	imageRepo UploadImageRepository,
	transactor Transactor,
	tel *observability.Telemetry,
) *uploadUsecase {
	uploadCount, _ := tel.Meter.Int64Counter(
		"r2.upload.count",
		metric.WithDescription("Total number of upload completion requests"),
	)
	return &uploadUsecase{
		uploadRepo:    uploadRepo,
		storage:       storage,
		characterRepo: characterRepo,
		artistRepo:    artistRepo,
		imageRepo:     imageRepo,
		transactor:    transactor,
		tel:           tel,
		uploadCount:   uploadCount,
	}
}

func (u *uploadUsecase) InitialUpload(ctx context.Context, params InitialUploadParams) (*InitialUploadResult, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.InitialUpload")
	defer span.End()

	id := uuid.New()
	ext := mimeTypeToExt(params.MimeType)
	r2Key := fmt.Sprintf("users/%s/images/%s%s", params.UserID.String(), id.String(), ext)
	expiresAt := time.Now().Add(PresignTTL)

	observability.LoggerFromContext(ctx, u.tel.Logger).Info("upload initiated",
		zap.String("event", "r2.upload.started"),
		zap.String("image_id", id.String()),
		zap.String("user_id", params.UserID.String()),
		zap.String("mime_type", params.MimeType),
		zap.String("r2_key", r2Key),
	)

	uploadURL, err := u.storage.GeneratePresignedPutURL(ctx, r2Key, params.MimeType, PresignTTL)
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
		ArtistID:     params.ArtistID,
		Notes:        params.Notes,
		CharacterIDs: params.CharacterIDs,
	}
	if _, err := u.uploadRepo.Create(ctx, pending); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &InitialUploadResult{ID: id, UploadURL: uploadURL, ExpiresAt: expiresAt}, nil
}

func (u *uploadUsecase) CompleteUpload(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CompleteUpload")
	defer span.End()

	start := time.Now()

	pending, err := u.uploadRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	validCharIDs := make([]uuid.UUID, 0)
	if len(pending.CharacterIDs) > 0 {
		chars, err := u.characterRepo.GetByIDsAndUserID(ctx, pending.CharacterIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		for _, c := range chars {
			validCharIDs = append(validCharIDs, c.ID)
		}
	}

	var resolvedArtistID *uuid.UUID
	if pending.ArtistID != nil {
		_, err := u.artistRepo.GetByIDAndUserID(ctx, *pending.ArtistID, userID)
		if err == nil {
			resolvedArtistID = pending.ArtistID
		}
		// if not found, silently leave resolvedArtistID as nil
	}

	u.uploadCount.Add(ctx, 1, metric.WithAttributes(attribute.String("r2.status", "success")))
	observability.LoggerFromContext(ctx, u.tel.Logger).Info("upload completed",
		zap.String("event", "r2.upload.completed"),
		zap.String("image_id", id.String()),
		zap.String("user_id", userID.String()),
		zap.Float64("duration_ms", float64(time.Since(start).Milliseconds())),
	)

	img := &domain.Image{
		ID:          pending.ID,
		UserID:      pending.UserID,
		ImageR2Path: pending.R2Key,
		MimeType:    pending.MimeType,
		Title:       pending.Title,
		ArtistID:    resolvedArtistID,
		Notes:       pending.Notes,
		Characters:  make([]domain.Character, len(validCharIDs)),
	}
	for i, cid := range validCharIDs {
		img.Characters[i] = domain.Character{ID: cid}
	}

	if err := u.transactor.InTransaction(ctx, func(txCtx context.Context) error {
		if err := u.uploadRepo.Delete(txCtx, pending.ID); err != nil {
			return err
		}
		_, err := u.imageRepo.Create(txCtx, img)
		return err
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (u *uploadUsecase) CleanupStaleUploads(ctx context.Context, threshold time.Duration) error {
	cutoff := time.Now().Add(-threshold)
	records, err := u.uploadRepo.ListStale(ctx, cutoff)
	if err != nil {
		return err
	}

	deleted := 0
	for _, p := range records {
		if err := u.storage.DeleteObject(ctx, p.R2Key); err != nil {
			return err
		}
		if err := u.uploadRepo.Delete(ctx, p.ID); err != nil {
			return err
		}
		deleted++
	}

	observability.LoggerFromContext(ctx, u.tel.Logger).Info("stale uploads purged",
		zap.String("event", "r2.upload.purged"),
		zap.Int("count", deleted),
	)
	return nil
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
