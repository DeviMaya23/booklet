package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/worker"
	bookmime "github.com/devi/booklet/pkg/mime"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type InitiateFileUploadParams struct {
	UserID     uuid.UUID
	MimeType   string
	ArtpieceID *uuid.UUID
	Name       *string
	Notes      *string
}

type InitiateFileUploadResult struct {
	ID        uuid.UUID
	UploadURL string
	ExpiresAt time.Time
}

type FileUploadUsecase struct {
	pendingRepo   FileUploadPendingRepository
	artpieceRepo  FileUploadArtpieceRepository
	fileRepo      FileUploadFileRepository
	metadataRepo  FileUploadImageMetadataRepository
	artpieceCover ArtpieceCoverUpdater
	storage       StorageService
	jobInserter   JobInserter
	tel           *observability.Telemetry
}

type ArtpieceCoverUpdater interface {
	UpdateCover(ctx context.Context, artpieceID uuid.UUID, coverFileID *uuid.UUID) error
}

func NewFileUploadUsecase(
	pendingRepo FileUploadPendingRepository,
	artpieceRepo FileUploadArtpieceRepository,
	fileRepo FileUploadFileRepository,
	metadataRepo FileUploadImageMetadataRepository,
	artpieceCover ArtpieceCoverUpdater,
	storage StorageService,
	jobInserter JobInserter,
	tel *observability.Telemetry,
) *FileUploadUsecase {
	return &FileUploadUsecase{
		pendingRepo:   pendingRepo,
		artpieceRepo:  artpieceRepo,
		fileRepo:      fileRepo,
		metadataRepo:  metadataRepo,
		artpieceCover: artpieceCover,
		storage:       storage,
		jobInserter:   jobInserter,
		tel:           tel,
	}
}

func (u *FileUploadUsecase) InitiateUpload(ctx context.Context, params InitiateFileUploadParams) (*InitiateFileUploadResult, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.InitiateFileUpload")
	defer span.End()

	if params.ArtpieceID != nil {
		if _, err := u.artpieceRepo.GetByIDAndUserID(ctx, *params.ArtpieceID, params.UserID); err != nil {
			return nil, ErrArtpieceNotOwned
		}
	}

	id := uuid.New()
	ext := bookmime.MimeTypeToExt(params.MimeType)
	r2Key := fmt.Sprintf("users/%s/files/%s%s", params.UserID.String(), id.String(), ext)
	expiresAt := time.Now().Add(PresignTTL)

	uploadURL, err := u.storage.GeneratePresignedPutURL(ctx, r2Key, params.MimeType, PresignTTL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pending := &domain.PendingFileUpload{
		ID:         id,
		UserID:     params.UserID,
		R2Key:      r2Key,
		MimeType:   params.MimeType,
		ArtpieceID: params.ArtpieceID,
		Name:       params.Name,
		Notes:      params.Notes,
	}
	if _, err := u.pendingRepo.Create(ctx, pending); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &InitiateFileUploadResult{ID: id, UploadURL: uploadURL, ExpiresAt: expiresAt}, nil
}

func (u *FileUploadUsecase) CompleteUpload(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CompleteFileUpload")
	defer span.End()

	pending, err := u.pendingRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	thumbnailGenState := "not_applicable"
	if bookmime.IsImage(pending.MimeType) {
		thumbnailGenState = "pending"
	}

	file := &domain.File{
		ID:                pending.ID,
		UserID:            pending.UserID,
		ArtpieceID:        pending.ArtpieceID,
		FileR2Path:        pending.R2Key,
		MimeType:          pending.MimeType,
		ThumbnailGenState: &thumbnailGenState,
		Name:              pending.Name,
		Notes:             pending.Notes,
	}

	created, err := u.fileRepo.Create(ctx, file)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if err := u.pendingRepo.Delete(ctx, pending.ID); err != nil {
		observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to delete pending file upload",
			zap.String("pending_id", pending.ID.String()),
			zap.Error(err),
		)
	}

	if bookmime.IsImage(pending.MimeType) {
		meta := &domain.ImageMetadata{
			FileID: created.ID,
			Width:  0,
			Height: 0,
		}
		if err := u.metadataRepo.CreateImageMetadata(ctx, meta); err != nil {
			observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to insert image_metadata",
				zap.String("file_id", created.ID.String()),
				zap.Error(err),
			)
		}

		if _, err := u.jobInserter.Insert(ctx, worker.GenerateFileThumbnailArgs{FileID: created.ID, UserID: created.UserID}, nil); err != nil {
			observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to enqueue file thumbnail job",
				zap.String("file_id", created.ID.String()),
				zap.Error(err),
			)
		}
	}

	if pending.ArtpieceID != nil {
		if err := u.maybeSetCover(ctx, created, *pending.ArtpieceID); err != nil {
			observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to auto-set cover",
				zap.String("artpiece_id", pending.ArtpieceID.String()),
				zap.Error(err),
			)
		}
	}

	return created, nil
}

func (u *FileUploadUsecase) maybeSetCover(ctx context.Context, file *domain.File, artpieceID uuid.UUID) error {
	if bookmime.IsImage(file.MimeType) {
		coverID := file.ID
		return u.artpieceCover.UpdateCover(ctx, artpieceID, &coverID)
	}

	// New file is not an image; only set as cover if no image files already exist.
	files, err := u.fileRepo.GetFilesForArtpiece(ctx, artpieceID)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.ID != file.ID && bookmime.IsImage(f.MimeType) {
			return nil
		}
	}

	coverID := file.ID
	return u.artpieceCover.UpdateCover(ctx, artpieceID, &coverID)
}

func (u *FileUploadUsecase) CleanupStaleUploads(ctx context.Context, threshold time.Duration) error {
	cutoff := time.Now().Add(-threshold)
	records, err := u.pendingRepo.ListStale(ctx, cutoff)
	if err != nil {
		return err
	}

	for _, p := range records {
		if err := u.storage.DeleteObject(ctx, p.R2Key); err != nil {
			return err
		}
		if err := u.pendingRepo.Delete(ctx, p.ID); err != nil {
			return err
		}
	}
	return nil
}
