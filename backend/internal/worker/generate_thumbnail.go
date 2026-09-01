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

type thumbnailImageRepository interface {
	GetByIDForWorker(ctx context.Context, id uuid.UUID) (*domain.Image, error)
	UpdateThumbnailPath(ctx context.Context, id uuid.UUID, r2Path string) error
}

type thumbnailStorageService interface {
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	PutObject(ctx context.Context, key string, body io.Reader, contentType string) error
}

type GenerateThumbnailArgs struct {
	ImageID uuid.UUID `json:"image_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (GenerateThumbnailArgs) Kind() string { return "generate_thumbnail" }

type GenerateThumbnailWorker struct {
	river.WorkerDefaults[GenerateThumbnailArgs]
	imageRepo thumbnailImageRepository
	storage   thumbnailStorageService
}

func NewGenerateThumbnailWorker(imageRepo thumbnailImageRepository, storage thumbnailStorageService) *GenerateThumbnailWorker {
	return &GenerateThumbnailWorker{imageRepo: imageRepo, storage: storage}
}

func (w *GenerateThumbnailWorker) Work(ctx context.Context, job *river.Job[GenerateThumbnailArgs]) error {
	imageID := job.Args.ImageID
	userID := job.Args.UserID

	img, err := w.imageRepo.GetByIDForWorker(ctx, imageID)
	if err != nil {
		return fmt.Errorf("get image: %w", err)
	}

	rc, err := w.storage.GetObject(ctx, img.ImageR2Path)
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

	thumbnailKey := fmt.Sprintf("users/%s/thumbnails/%s.jpg", userID.String(), imageID.String())
	if err := w.storage.PutObject(ctx, thumbnailKey, &buf, "image/jpeg"); err != nil {
		return fmt.Errorf("upload thumbnail: %w", err)
	}

	if err := w.imageRepo.UpdateThumbnailPath(ctx, imageID, thumbnailKey); err != nil {
		return fmt.Errorf("update thumbnail_r2_path: %w", err)
	}

	return nil
}
