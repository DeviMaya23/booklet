package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type tagUsecase struct {
	tagRepo TagRepository
	tel     *observability.Telemetry
}

func NewTagUsecase(tagRepo TagRepository, tel *observability.Telemetry) *tagUsecase {
	return &tagUsecase{tagRepo: tagRepo, tel: tel}
}

func normalizeTagName(name string) (string, error) {
	n := strings.ToLower(strings.TrimSpace(name))
	if n == "" {
		return "", ErrTagNameInvalid
	}
	return n, nil
}

func (u *tagUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateTagParams) (*domain.Tag, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateTag")
	defer span.End()

	name, err := normalizeTagName(params.Name)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	res, err := u.tagRepo.Create(ctx, &domain.Tag{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       name,
		LastUsedAt: &now,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *tagUsecase) List(ctx context.Context, userID uuid.UUID, filters ListTagFilters) ([]*domain.Tag, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListTags")
	defer span.End()

	res, err := u.tagRepo.List(ctx, userID, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *tagUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateTagParams) (*domain.Tag, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateTag")
	defer span.End()

	name, err := normalizeTagName(params.Name)
	if err != nil {
		return nil, err
	}

	res, err := u.tagRepo.Update(ctx, id, userID, UpdateTagParams{Name: name})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *tagUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteTag")
	defer span.End()

	if err := u.tagRepo.Delete(ctx, id, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
