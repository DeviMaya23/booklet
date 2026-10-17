package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	_ "golang.org/x/image/webp"
	"go.uber.org/zap"
)

const thumbnailCacheControl = "private, max-age=172800"

type thumbnailFormat struct {
	encoding    imaging.Format
	ext         string
	contentType string
}

func formatForMimeType(mimeType string) thumbnailFormat {
	switch mimeType {
	case "image/png", "image/gif", "image/webp":
		return thumbnailFormat{imaging.PNG, ".png", "image/png"}
	default:
		return thumbnailFormat{imaging.JPEG, ".jpg", "image/jpeg"}
	}
}

type fileThumbnailFileRepository interface {
	GetByIDForWorker(ctx context.Context, id uuid.UUID) (*domain.File, error)
	UpdateThumbnailPath(ctx context.Context, id uuid.UUID, r2Path string) error
	UpdateThumbnailGenState(ctx context.Context, id uuid.UUID, state string) error
	UpdateImageMetadataDimensions(ctx context.Context, fileID uuid.UUID, width, height int) error
}

type fileThumbnailStorageService interface {
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	PutObject(ctx context.Context, key string, body io.Reader, contentType, cacheControl string) error
}

type GenerateFileThumbnailArgs struct {
	FileID uuid.UUID `json:"file_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (GenerateFileThumbnailArgs) Kind() string { return "generate_file_thumbnail" }

type GenerateFileThumbnailWorker struct {
	river.WorkerDefaults[GenerateFileThumbnailArgs]
	fileRepo      fileThumbnailFileRepository
	storage       fileThumbnailStorageService
	logger        *zap.Logger
	decodeTimeout time.Duration
}

func NewGenerateFileThumbnailWorker(fileRepo fileThumbnailFileRepository, storage fileThumbnailStorageService, logger *zap.Logger) *GenerateFileThumbnailWorker {
	return NewGenerateFileThumbnailWorkerWithTimeout(fileRepo, storage, logger, 30*time.Second)
}

func NewGenerateFileThumbnailWorkerWithTimeout(fileRepo fileThumbnailFileRepository, storage fileThumbnailStorageService, logger *zap.Logger, decodeTimeout time.Duration) *GenerateFileThumbnailWorker {
	return &GenerateFileThumbnailWorker{
		fileRepo:      fileRepo,
		storage:       storage,
		logger:        logger,
		decodeTimeout: decodeTimeout,
	}
}

func (w *GenerateFileThumbnailWorker) NextRetry(job *river.Job[GenerateFileThumbnailArgs]) time.Time {
	return time.Now().Add(15 * time.Second)
}

type decodeResizeResult struct {
	img    image.Image
	width  int
	height int
	err    error
}

func (w *GenerateFileThumbnailWorker) Work(ctx context.Context, job *river.Job[GenerateFileThumbnailArgs]) error {
	fileID := job.Args.FileID
	userID := job.Args.UserID

	file, err := w.fileRepo.GetByIDForWorker(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get file: %w", err)
	}

	setFailed := func() error {
		if err := w.fileRepo.UpdateThumbnailGenState(context.Background(), fileID, "failed"); err != nil {
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

	decodeCtx, cancelDecode := context.WithTimeout(ctx, w.decodeTimeout)
	defer cancelDecode()

	ch := make(chan decodeResizeResult, 1)
	go func() {
		decoded, decodeErr := imaging.Decode(rc)
		if decodeErr != nil {
			ch <- decodeResizeResult{err: decodeErr}
			return
		}
		bounds := decoded.Bounds()
		thumbnail := imaging.Fit(decoded, 600, 600, imaging.Lanczos)
		ch <- decodeResizeResult{img: thumbnail, width: bounds.Dx(), height: bounds.Dy()}
	}()

	var res decodeResizeResult
	select {
	case res = <-ch:
	case <-decodeCtx.Done():
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("decode timeout: %w", decodeCtx.Err())
	}

	if res.err != nil {
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("decode image: %w", res.err)
	}

	if err := w.fileRepo.UpdateImageMetadataDimensions(ctx, fileID, res.width, res.height); err != nil {
		w.logger.Error("failed to update image_metadata dimensions",
			zap.String("file_id", fileID.String()),
			zap.Error(err),
		)
	}

	format := formatForMimeType(file.MimeType)

	var encodeOpts []imaging.EncodeOption
	if format.encoding == imaging.JPEG {
		encodeOpts = append(encodeOpts, imaging.JPEGQuality(85))
	}

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, res.img, format.encoding, encodeOpts...); err != nil {
		if job.Attempt >= job.MaxAttempts {
			return setFailed()
		}
		return fmt.Errorf("encode thumbnail: %w", err)
	}

	thumbnailKey := fmt.Sprintf("users/%s/thumbnails/%s%s", userID.String(), fileID.String(), format.ext)
	if err := w.storage.PutObject(ctx, thumbnailKey, &buf, format.contentType, thumbnailCacheControl); err != nil {
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
