package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

type savedFilterUsecase struct {
	repo         SavedFilterRepository
	artpieceRepo ArtpieceThumbnailRepository
	tel          *observability.Telemetry
}

func NewSavedFilterUsecase(repo SavedFilterRepository, artpieceRepo ArtpieceThumbnailRepository, tel *observability.Telemetry) *savedFilterUsecase {
	return &savedFilterUsecase{repo: repo, artpieceRepo: artpieceRepo, tel: tel}
}

func (u *savedFilterUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateSavedFilterParams) (*domain.SavedFilter, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateSavedFilter")
	defer span.End()

	var thumbnailR2Path *string
	if params.ThumbnailArtpieceID != nil {
		artpiece, err := u.artpieceRepo.GetByID(ctx, *params.ThumbnailArtpieceID, userID)
		if err == nil && artpiece.CoverFile != nil {
			thumbnailR2Path = artpiece.CoverFile.ThumbnailR2Path
		}
	}

	sf := &domain.SavedFilter{
		ID:              uuid.New(),
		UserID:          userID,
		Name:            params.Name,
		ThumbnailR2Path: thumbnailR2Path,
		FilterPayload:   params.FilterPayload,
	}
	res, err := u.repo.Create(ctx, sf)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, ErrDuplicateSavedFilterName) {
			return nil, ErrDuplicateSavedFilterName
		}
		return nil, err
	}
	return res, nil
}

func (u *savedFilterUsecase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.SavedFilter, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetSavedFilterByID")
	defer span.End()

	res, err := u.repo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSavedFilterNotFound
		}
		return nil, err
	}
	return res, nil
}

func (u *savedFilterUsecase) List(ctx context.Context, userID uuid.UUID) ([]*domain.SavedFilter, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListSavedFilters")
	defer span.End()

	res, err := u.repo.List(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *savedFilterUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateSavedFilterParams) (*domain.SavedFilter, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateSavedFilter")
	defer span.End()

	if params.ThumbnailArtpieceID != nil {
		artpiece, err := u.artpieceRepo.GetByID(ctx, *params.ThumbnailArtpieceID, userID)
		var r2Path *string
		if err == nil && artpiece.CoverFile != nil {
			r2Path = artpiece.CoverFile.ThumbnailR2Path
		}
		params.ThumbnailR2Path = Patch[string]{Set: true, Value: r2Path}
	}

	res, err := u.repo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSavedFilterNotFound
		}
		return nil, err
	}
	return res, nil
}

func (u *savedFilterUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteSavedFilter")
	defer span.End()

	err := u.repo.Delete(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSavedFilterNotFound
		}
		return err
	}
	return nil
}
