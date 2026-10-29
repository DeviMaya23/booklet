package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type verseUsecase struct {
	verseRepo VerseRepository
	tel       *observability.Telemetry
}

func NewVerseUsecase(verseRepo VerseRepository, tel *observability.Telemetry) *verseUsecase {
	return &verseUsecase{verseRepo: verseRepo, tel: tel}
}

func (u *verseUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateVerseParams) (*domain.Verse, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateVerse")
	defer span.End()

	verse := &domain.Verse{
		ID:     uuid.New(),
		UserID: userID,
		Name:   params.Name,
	}
	res, err := u.verseRepo.Create(ctx, verse)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *verseUsecase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Verse, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetVerseByID")
	defer span.End()

	res, err := u.verseRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *verseUsecase) List(ctx context.Context, userID uuid.UUID, filters ListVerseFilters) ([]*domain.Verse, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListVerses")
	defer span.End()

	res, err := u.verseRepo.List(ctx, userID, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *verseUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateVerseParams) (*domain.Verse, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateVerse")
	defer span.End()

	res, err := u.verseRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *verseUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteVerse")
	defer span.End()

	err := u.verseRepo.Delete(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
