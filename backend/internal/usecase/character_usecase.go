package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type CreateCharacterParams struct {
	Name            string
	HeroImageR2Path *string
	Biography       *string
	IsPublic        bool
}

type characterUsecase struct {
	characterRepo CharacterRepository
	tel           *observability.Telemetry
}

func NewCharacterUsecase(characterRepo CharacterRepository, tel *observability.Telemetry) *characterUsecase {
	return &characterUsecase{
		characterRepo: characterRepo,
		tel:           tel,
	}
}

func (u *characterUsecase) Create(ctx context.Context, userID string, params CreateCharacterParams) (*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateCharacter")
	defer span.End()

	character := &domain.Character{
		ID:              uuid.New(),
		UserID:          userID,
		Name:            params.Name,
		HeroImageR2Path: params.HeroImageR2Path,
		Biography:       params.Biography,
		IsPublic:        params.IsPublic,
	}
	if err := u.characterRepo.Create(ctx, character); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return character, nil
}

func (u *characterUsecase) GetByID(ctx context.Context, id, userID string) (*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetCharacterByID")
	defer span.End()

	res, err := u.characterRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *characterUsecase) List(ctx context.Context, userID string) ([]*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListCharacters")
	defer span.End()

	res, err := u.characterRepo.List(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *characterUsecase) Update(ctx context.Context, id, userID string, params UpdateCharacterParams) (*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateCharacter")
	defer span.End()

	res, err := u.characterRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *characterUsecase) Delete(ctx context.Context, id, userID string) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteCharacter")
	defer span.End()

	err := u.characterRepo.Delete(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
