package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/disintegration/imaging"
	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

type fileThumbnailFileRepository interface {
	GetByIDForWorker(ctx context.Context, id uuid.UUID) (*domain.File, error)
	UpdateThumbnailPath(ctx context.Context, id uuid.UUID, r2Path string) error
	UpdateThumbnailGenState(ctx context.Context, id uuid.UUID, state string) error
	UpdateImageMetadataDimensions(ctx context.Context, fileID uuid.UUID, width, height int) error
}

type fileThumbnailStorageService interface {
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	PutObject(ctx context.Context, key string, body io.Reader, contentType string) error
}

type GenerateFileThumbnailArgs struct {
	FileID uuid.UUID `json:"file_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (GenerateFileThumbnailArgs) Kind() string { return "generate_file_thumbnail" }

type GenerateFileThumbnailWorker struct {
	river.WorkerDefaults[GenerateFileThumbnailArgs]
	fileRepo fileThumbnailFileRepository
	storage  fileThumbnailStorageService
	logger   *zap.Logger
}

func NewGenerateFileThumbnailWorker(fileRepo fileThumbnailFileRepository, storage fileThumbnailStorageService, logger *zap.Logger) *GenerateFileThumbnailWorker {
	return &GenerateFileThumbnailWorker{fileRepo: fileRepo, storage: storage, logger: logger}
}

func (w *GenerateFileThumbnailWorker) Work(ctx context.Context, job *river.Job[GenerateFileThumbnailArgs]) error {
	fileID := job.Args.FileID
	userID := job.Args.UserID

	file, err := w.fileRepo.GetByIDForWorker(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get file: %w", err)
	}

	setFailed := func() error {
		if err := w.fileRepo.UpdateThumbnailGenState(ctx, fileID, "failed"); err != nil {
			w.logger.Error("failed to set thumbnail_gen_state to failed",
				zap.String("file_id", fileID.String()),
				zap.Error(err),
			)
		}
		return nil
	}

	rc, err := w.storage.GetObject(ctx, file.FileR2Path)
	if err != nil {
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("get original from r2: %w", err)
	}
	defer func() { _ = rc.Close() }()

	decoded, err := imaging.Decode(rc)
	if err != nil {
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("decode image: %w", err)
	}

	bounds := decoded.Bounds()
	if err := w.fileRepo.UpdateImageMetadataDimensions(ctx, fileID, bounds.Dx(), bounds.Dy()); err != nil {
		w.logger.Error("failed to update image_metadata dimensions",
			zap.String("file_id", fileID.String()),
			zap.Error(err),
		)
	}

	thumbnail := imaging.Fit(decoded, 600, 600, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumbnail, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("encode thumbnail: %w", err)
	}

	thumbnailKey := fmt.Sprintf("users/%s/thumbnails/%s.jpg", userID.String(), fileID.String())
	if err := w.storage.PutObject(ctx, thumbnailKey, &buf, "image/jpeg"); err != nil {
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("upload thumbnail: %w", err)
	}

	if err := w.fileRepo.UpdateThumbnailPath(ctx, fileID, thumbnailKey); err != nil {
		return fmt.Errorf("update thumbnail_r2_path: %w", err)
	}

	if err := w.fileRepo.UpdateThumbnailGenState(ctx, fileID, "done"); err != nil {
		return fmt.Errorf("update thumbnail_gen_state: %w", err)
	}

	return nil
}
