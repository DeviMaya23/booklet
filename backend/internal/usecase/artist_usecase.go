package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type CreateArtistParams struct {
	Name       string
	Notes      *string
	ArtistLink *string
}

type artistUsecase struct {
	artistRepo ArtistRepository
	tel        *observability.Telemetry
}

func NewArtistUsecase(artistRepo ArtistRepository, tel *observability.Telemetry) *artistUsecase {
	return &artistUsecase{
		artistRepo: artistRepo,
		tel:        tel,
	}
}

func (u *artistUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateArtistParams) (*domain.Artist, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateArtist")
	defer span.End()

	artist := &domain.Artist{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       params.Name,
		Notes:      params.Notes,
		ArtistLink: params.ArtistLink,
	}
	res, err := u.artistRepo.Create(ctx, artist)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *artistUsecase) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Artist, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetArtistByID")
	defer span.End()

	res, err := u.artistRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *artistUsecase) List(ctx context.Context, userID uuid.UUID) ([]*domain.Artist, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListArtists")
	defer span.End()

	res, err := u.artistRepo.List(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *artistUsecase) Update(ctx context.Context, id string, userID uuid.UUID, params UpdateArtistParams) (*domain.Artist, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateArtist")
	defer span.End()

	res, err := u.artistRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *artistUsecase) Delete(ctx context.Context, id string, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteArtist")
	defer span.End()

	err := u.artistRepo.Delete(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
