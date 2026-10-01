package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type FileUsecase struct {
	fileRepo FileRepository
	storage  FileStorageService
	tel      *observability.Telemetry
}

func NewFileUsecase(fileRepo FileRepository, storage FileStorageService, tel *observability.Telemetry) *FileUsecase {
	return &FileUsecase{fileRepo: fileRepo, storage: storage, tel: tel}
}

func (u *FileUsecase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetFileByID")
	defer span.End()

	file, err := u.fileRepo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return file, nil
}

func (u *FileUsecase) List(ctx context.Context, userID uuid.UUID, unassigned bool) ([]*domain.File, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListFiles")
	defer span.End()

	files, err := u.fileRepo.List(ctx, userID, unassigned)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return files, nil
}

func (u *FileUsecase) UpdateNotes(ctx context.Context, id uuid.UUID, userID uuid.UUID, notes *string) (*domain.File, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateFileNotes")
	defer span.End()

	if err := u.fileRepo.UpdateNotes(ctx, id, userID, notes); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return u.fileRepo.GetByIDAndUserID(ctx, id, userID)
}

func (u *FileUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteFile")
	defer span.End()

	file, err := u.fileRepo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := u.fileRepo.Delete(ctx, id, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	logger := observability.LoggerFromContext(ctx, u.tel.Logger)
	if err := u.storage.DeleteObject(ctx, file.FileR2Path); err != nil {
		logger.Error("failed to delete file from R2",
			zap.String("file_id", id.String()),
			zap.String("r2_path", file.FileR2Path),
			zap.Error(err),
		)
	}
	if file.ThumbnailR2Path != nil {
		if err := u.storage.DeleteObject(ctx, *file.ThumbnailR2Path); err != nil {
			logger.Error("failed to delete file thumbnail from R2",
				zap.String("file_id", id.String()),
				zap.String("r2_path", *file.ThumbnailR2Path),
				zap.Error(err),
			)
		}
	}

	return nil
}
