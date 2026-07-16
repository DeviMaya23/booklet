package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"go.opentelemetry.io/otel/codes"
)

type imageUsecase struct {
	imageRepo ImageRepository
	tel       *observability.Telemetry
}

func NewImageUsecase(imageRepo ImageRepository, tel *observability.Telemetry) *imageUsecase {
	return &imageUsecase{
		imageRepo: imageRepo,
		tel:       tel,
	}
}

func (u *imageUsecase) GetByID(ctx context.Context, id, userID string) (*domain.Image, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetImageByID")
	defer span.End()

	res, err := u.imageRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *imageUsecase) List(ctx context.Context, userID string) ([]*domain.Image, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListImages")
	defer span.End()

	res, err := u.imageRepo.List(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *imageUsecase) Update(ctx context.Context, id, userID string, params UpdateImageParams) (*domain.Image, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateImage")
	defer span.End()

	res, err := u.imageRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *imageUsecase) Delete(ctx context.Context, id, userID string) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteImage")
	defer span.End()

	err := u.imageRepo.Delete(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
