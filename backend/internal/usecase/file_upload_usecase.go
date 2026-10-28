package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	bookmime "github.com/devi/booklet/pkg/mime"
	"github.com/devi/booklet/pkg/filesize"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type InitiateFileUploadParams struct {
	UserID        uuid.UUID
	MimeType      string
	ArtpieceID    *uuid.UUID
	Name          *string
	Notes         *string
	Extension     *string
	FileSizeBytes int64
}

type InitiateFileUploadResult struct {
	ID                 uuid.UUID
	UploadURL          string
	ThumbnailUploadURL *string
	ThumbnailR2Key     *string
	ExpiresAt          time.Time
}

type CompleteFileUploadParams struct {
	Width             *int
	Height            *int
	ThumbnailUploaded bool
}

type FileUploadUsecase struct {
	pendingRepo   FileUploadPendingRepository
	artpieceRepo  FileUploadArtpieceRepository
	fileRepo      FileUploadFileRepository
	metadataRepo  FileUploadImageMetadataRepository
	artpieceCover ArtpieceCoverUpdater
	storage       StorageService
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
	tel *observability.Telemetry,
) *FileUploadUsecase {
	return &FileUploadUsecase{
		pendingRepo:   pendingRepo,
		artpieceRepo:  artpieceRepo,
		fileRepo:      fileRepo,
		metadataRepo:  metadataRepo,
		artpieceCover: artpieceCover,
		storage:       storage,
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

	var ext string
	if params.Extension != nil && *params.Extension != "" {
		ext = *params.Extension
	} else {
		ext = strings.TrimPrefix(bookmime.MimeTypeToExt(params.MimeType), ".")
	}
	if params.FileSizeBytes > filesize.MaxBytesForExtension(ext) {
		return nil, ErrFileTooLarge
	}

	id := uuid.New()
	r2Ext := bookmime.MimeTypeToExt(params.MimeType)
	r2Key := fmt.Sprintf("users/%s/files/%s%s", params.UserID.String(), id.String(), r2Ext)
	expiresAt := time.Now().Add(PresignTTL)

	uploadURL, err := u.storage.GeneratePresignedPutURL(ctx, r2Key, params.MimeType, params.FileSizeBytes, PresignTTL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var thumbnailUploadURL *string
	var thumbnailR2Key *string
	if bookmime.IsImage(params.MimeType) {
		thumbExt := bookmime.ThumbnailExt(params.MimeType)
		thumbContentType := bookmime.ThumbnailContentType(params.MimeType)
		thumbKey := fmt.Sprintf("users/%s/thumbnails/%s%s", params.UserID.String(), id.String(), thumbExt)
		thumbURL, err := u.storage.GeneratePresignedPutURL(ctx, thumbKey, thumbContentType, 0, PresignTTL)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		thumbnailUploadURL = &thumbURL
		thumbnailR2Key = &thumbKey
	}

	pending := &domain.PendingFileUpload{
		ID:         id,
		UserID:     params.UserID,
		R2Key:      r2Key,
		MimeType:   params.MimeType,
		ArtpieceID: params.ArtpieceID,
		Name:       params.Name,
		Notes:      params.Notes,
		Extension:  params.Extension,
	}
	if _, err := u.pendingRepo.Create(ctx, pending); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &InitiateFileUploadResult{
		ID:                 id,
		UploadURL:          uploadURL,
		ThumbnailUploadURL: thumbnailUploadURL,
		ThumbnailR2Key:     thumbnailR2Key,
		ExpiresAt:          expiresAt,
	}, nil
}

func (u *FileUploadUsecase) CompleteUpload(ctx context.Context, id uuid.UUID, userID uuid.UUID, params CompleteFileUploadParams) (*domain.File, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CompleteFileUpload")
	defer span.End()

	pending, err := u.pendingRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	file := &domain.File{
		ID:         pending.ID,
		UserID:     pending.UserID,
		ArtpieceID: pending.ArtpieceID,
		FileR2Path: pending.R2Key,
		MimeType:   pending.MimeType,
		Name:       pending.Name,
		Notes:      pending.Notes,
		Extension:  pending.Extension,
	}
	if bookmime.IsImage(pending.MimeType) && params.ThumbnailUploaded {
		thumbKey := fmt.Sprintf("users/%s/thumbnails/%s%s", pending.UserID.String(), pending.ID.String(), bookmime.ThumbnailExt(pending.MimeType))
		file.ThumbnailR2Path = &thumbKey
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
		width, height := 0, 0
		if params.Width != nil {
			width = *params.Width
		}
		if params.Height != nil {
			height = *params.Height
		}
		meta := &domain.ImageMetadata{
			FileID: created.ID,
			Width:  width,
			Height: height,
		}
		if err := u.metadataRepo.CreateImageMetadata(ctx, meta); err != nil {
			observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to insert image_metadata",
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
