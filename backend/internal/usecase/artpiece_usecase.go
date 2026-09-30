package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	bookmime "github.com/devi/booklet/pkg/mime"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type ArtpieceUsecase struct {
	artpieceRepo  ArtpieceRepository
	artistRepo    ArtpieceArtistRepository
	characterRepo ArtpieceCharacterRepository
	fileRepo      ArtpieceFileRepository
	tel           *observability.Telemetry
}

func NewArtpieceUsecase(
	artpieceRepo ArtpieceRepository,
	artistRepo ArtpieceArtistRepository,
	characterRepo ArtpieceCharacterRepository,
	fileRepo ArtpieceFileRepository,
	tel *observability.Telemetry,
) *ArtpieceUsecase {
	return &ArtpieceUsecase{
		artpieceRepo:  artpieceRepo,
		artistRepo:    artistRepo,
		characterRepo: characterRepo,
		fileRepo:      fileRepo,
		tel:           tel,
	}
}

func (u *ArtpieceUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateArtpieceParams) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateArtpiece")
	defer span.End()

	if params.ArtistID != nil {
		if _, err := u.artistRepo.GetByIDAndUserID(ctx, *params.ArtistID, userID); err != nil {
			return nil, ErrArtistNotOwned
		}
	}

	if len(params.CharacterIDs) > 0 {
		chars, err := u.characterRepo.GetByIDsAndUserID(ctx, params.CharacterIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if len(chars) != len(params.CharacterIDs) {
			return nil, ErrCharacterNotOwned
		}
	}

	characters := make([]domain.Character, len(params.CharacterIDs))
	for i, cid := range params.CharacterIDs {
		characters[i] = domain.Character{ID: cid}
	}

	artpiece := &domain.Artpiece{
		ID:         uuid.New(),
		UserID:     userID,
		Title:      params.Title,
		ArtistID:   params.ArtistID,
		Notes:      params.Notes,
		Characters: characters,
	}

	res, err := u.artpieceRepo.Create(ctx, artpiece)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *ArtpieceUsecase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetArtpieceByID")
	defer span.End()

	res, err := u.artpieceRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *ArtpieceUsecase) List(ctx context.Context, userID uuid.UUID, filters ListArtpieceFilters) ([]*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListArtpieces")
	defer span.End()

	res, err := u.artpieceRepo.List(ctx, userID, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *ArtpieceUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateArtpieceParams) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateArtpiece")
	defer span.End()

	res, err := u.artpieceRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *ArtpieceUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteArtpiece")
	defer span.End()

	if err := u.artpieceRepo.Delete(ctx, id, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (u *ArtpieceUsecase) AttachFile(ctx context.Context, artpieceID uuid.UUID, fileID uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.AttachFile")
	defer span.End()

	artpiece, err := u.artpieceRepo.GetByID(ctx, artpieceID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	file, err := u.fileRepo.GetByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return nil, ErrFileNotOwned
	}

	if file.ArtpieceID != nil && *file.ArtpieceID != artpieceID {
		return nil, ErrFileAlreadyAttached
	}

	if err := u.fileRepo.UpdateArtpieceID(ctx, fileID, &artpieceID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if artpiece.CoverFileID == nil {
		file.ArtpieceID = &artpieceID
		if err := u.maybeAutoSetCover(ctx, artpiece, file); err != nil {
			span.RecordError(err)
		}
	}

	return u.artpieceRepo.GetByID(ctx, artpieceID, userID)
}

func (u *ArtpieceUsecase) DetachFile(ctx context.Context, artpieceID uuid.UUID, fileID uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DetachFile")
	defer span.End()

	artpiece, err := u.artpieceRepo.GetByID(ctx, artpieceID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	file, err := u.fileRepo.GetByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return nil, ErrFileNotOwned
	}

	if file.ArtpieceID == nil || *file.ArtpieceID != artpieceID {
		return nil, ErrFileNotInArtpiece
	}

	if err := u.fileRepo.UpdateArtpieceID(ctx, fileID, nil); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if artpiece.CoverFileID != nil && *artpiece.CoverFileID == fileID {
		if err := u.reassignCover(ctx, artpieceID, fileID); err != nil {
			span.RecordError(err)
		}
	}

	return u.artpieceRepo.GetByID(ctx, artpieceID, userID)
}

func (u *ArtpieceUsecase) SetCover(ctx context.Context, artpieceID uuid.UUID, fileID uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.SetCover")
	defer span.End()

	_, err := u.artpieceRepo.GetByID(ctx, artpieceID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	file, err := u.fileRepo.GetByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return nil, ErrFileNotOwned
	}

	if file.ArtpieceID == nil || *file.ArtpieceID != artpieceID {
		return nil, ErrFileNotInArtpiece
	}

	if err := u.artpieceRepo.UpdateCover(ctx, artpieceID, &fileID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return u.artpieceRepo.GetByID(ctx, artpieceID, userID)
}

func (u *ArtpieceUsecase) UpdateCover(ctx context.Context, artpieceID uuid.UUID, coverFileID *uuid.UUID) error {
	return u.artpieceRepo.UpdateCover(ctx, artpieceID, coverFileID)
}

func (u *ArtpieceUsecase) maybeAutoSetCover(ctx context.Context, artpiece *domain.Artpiece, newFile *domain.File) error {
	if bookmime.IsImage(newFile.MimeType) {
		coverID := newFile.ID
		return u.artpieceRepo.UpdateCover(ctx, artpiece.ID, &coverID)
	}

	// New file is not an image; only set as cover if no image files already exist.
	files, err := u.artpieceRepo.GetFilesForArtpiece(ctx, artpiece.ID)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.ID != newFile.ID && bookmime.IsImage(f.MimeType) {
			return nil
		}
	}

	coverID := newFile.ID
	return u.artpieceRepo.UpdateCover(ctx, artpiece.ID, &coverID)
}

func (u *ArtpieceUsecase) reassignCover(ctx context.Context, artpieceID uuid.UUID, detachedFileID uuid.UUID) error {
	files, err := u.artpieceRepo.GetFilesForArtpiece(ctx, artpieceID)
	if err != nil {
		return err
	}

	var remaining []*domain.File
	for _, f := range files {
		if f.ID != detachedFileID {
			remaining = append(remaining, f)
		}
	}

	if len(remaining) == 0 {
		return u.artpieceRepo.UpdateCover(ctx, artpieceID, nil)
	}

	// prefer an image file
	for _, f := range remaining {
		if bookmime.IsImage(f.MimeType) {
			coverID := f.ID
			return u.artpieceRepo.UpdateCover(ctx, artpieceID, &coverID)
		}
	}

	coverID := remaining[0].ID
	return u.artpieceRepo.UpdateCover(ctx, artpieceID, &coverID)
}
