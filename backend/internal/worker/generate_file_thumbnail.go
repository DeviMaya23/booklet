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
)

type fileThumbnailFileRepository interface {
	GetByIDForWorker(ctx context.Context, id uuid.UUID) (*domain.File, error)
	UpdateThumbnailPath(ctx context.Context, id uuid.UUID, r2Path string) error
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
}

func NewGenerateFileThumbnailWorker(fileRepo fileThumbnailFileRepository, storage fileThumbnailStorageService) *GenerateFileThumbnailWorker {
	return &GenerateFileThumbnailWorker{fileRepo: fileRepo, storage: storage}
}

func (w *GenerateFileThumbnailWorker) Work(ctx context.Context, job *river.Job[GenerateFileThumbnailArgs]) error {
	fileID := job.Args.FileID
	userID := job.Args.UserID

	file, err := w.fileRepo.GetByIDForWorker(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get file: %w", err)
	}

	rc, err := w.storage.GetObject(ctx, file.FileR2Path)
	if err != nil {
		return fmt.Errorf("get original from r2: %w", err)
	}
	defer func() { _ = rc.Close() }()

	decoded, err := imaging.Decode(rc)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	thumbnail := imaging.Fit(decoded, 600, 600, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumbnail, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
		return fmt.Errorf("encode thumbnail: %w", err)
	}

	thumbnailKey := fmt.Sprintf("users/%s/thumbnails/%s.jpg", userID.String(), fileID.String())
	if err := w.storage.PutObject(ctx, thumbnailKey, &buf, "image/jpeg"); err != nil {
		return fmt.Errorf("upload thumbnail: %w", err)
	}

	if err := w.fileRepo.UpdateThumbnailPath(ctx, fileID, thumbnailKey); err != nil {
		return fmt.Errorf("update thumbnail_r2_path: %w", err)
	}

	return nil
}
