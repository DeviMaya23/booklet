package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/worker"
	bookmime "github.com/devi/booklet/pkg/mime"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type ArtpieceUsecase struct {
	artpieceRepo  ArtpieceRepository
	artistRepo    ArtpieceArtistRepository
	characterRepo ArtpieceCharacterRepository
	fileRepo      ArtpieceFileRepository
	transactor    Transactor
	jobInserter   JobInserter
	tel           *observability.Telemetry
}

func NewArtpieceUsecase(
	artpieceRepo ArtpieceRepository,
	artistRepo ArtpieceArtistRepository,
	characterRepo ArtpieceCharacterRepository,
	fileRepo ArtpieceFileRepository,
	transactor Transactor,
	jobInserter JobInserter,
	tel *observability.Telemetry,
) *ArtpieceUsecase {
	return &ArtpieceUsecase{
		artpieceRepo:  artpieceRepo,
		artistRepo:    artistRepo,
		characterRepo: characterRepo,
		fileRepo:      fileRepo,
		transactor:    transactor,
		jobInserter:   jobInserter,
		tel:           tel,
	}
}

func (u *ArtpieceUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateArtpieceParams) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateArtpiece")
	defer span.End()

	if params.ArtistID != nil {
		if _, err := u.artistRepo.GetByID(ctx, *params.ArtistID, userID); err != nil {
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

	fileIDs := dedupUUIDs(params.FileIDs)
	if len(fileIDs) > 0 {
		files, err := u.fileRepo.GetByIDsAndUserID(ctx, fileIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if len(files) != len(fileIDs) {
			return nil, ErrFileNotOwned
		}
		for _, f := range files {
			if f.ArtpieceID != nil {
				return nil, ErrFileAlreadyAttached
			}
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

	created, err := u.artpieceRepo.Create(ctx, artpiece)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(fileIDs) > 0 {
		err = u.transactor.InTransaction(ctx, func(ctx context.Context) error {
			if err := u.fileRepo.BulkUpdateArtpieceID(ctx, fileIDs, &artpiece.ID); err != nil {
				return err
			}
			files, err := u.fileRepo.GetFilesForArtpiece(ctx, artpiece.ID)
			if err != nil {
				return err
			}
			return u.setCoverFromFileSet(ctx, artpiece.ID, nil, files)
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	return u.artpieceRepo.GetByID(ctx, created.ID, userID)
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

func (u *ArtpieceUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, deleteFiles bool) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteArtpiece")
	defer span.End()

	if !deleteFiles {
		if err := u.artpieceRepo.Delete(ctx, id, userID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		return nil
	}

	if _, err := u.artpieceRepo.GetByID(ctx, id, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	files, err := u.fileRepo.GetFilesForArtpiece(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	fileIDs := make([]uuid.UUID, len(files))
	var r2Keys []string
	for i, f := range files {
		fileIDs[i] = f.ID
		r2Keys = append(r2Keys, f.FileR2Path)
		if f.ThumbnailR2Path != nil {
			r2Keys = append(r2Keys, *f.ThumbnailR2Path)
		}
	}

	err = u.transactor.InTransaction(ctx, func(ctx context.Context) error {
		if len(fileIDs) > 0 {
			if err := u.fileRepo.BulkDelete(ctx, fileIDs, userID); err != nil {
				return err
			}
		}
		return u.artpieceRepo.Delete(ctx, id, userID)
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if len(r2Keys) > 0 {
		if _, err := u.jobInserter.Insert(ctx, worker.PurgeR2ObjectsArgs{R2Keys: r2Keys}, nil); err != nil {
			observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to enqueue purge_r2_objects job after artpiece delete with files",
				zap.String("artpiece_id", id.String()),
				zap.Error(err),
			)
		}
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

// ReplaceFiles replaces the full file set for an artpiece in a single transaction.
func (u *ArtpieceUsecase) ReplaceFiles(ctx context.Context, artpieceID uuid.UUID, userID uuid.UUID, fileIDs []uuid.UUID) (*domain.Artpiece, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ReplaceFiles")
	defer span.End()

	fileIDs = dedupUUIDs(fileIDs)

	artpiece, err := u.artpieceRepo.GetByID(ctx, artpieceID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(fileIDs) > 0 {
		owned, err := u.fileRepo.GetByIDsAndUserID(ctx, fileIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if len(owned) != len(fileIDs) {
			return nil, ErrFileNotOwned
		}
		for _, f := range owned {
			if f.ArtpieceID != nil && *f.ArtpieceID != artpieceID {
				return nil, ErrFileAlreadyAttached
			}
		}
	}

	currentFiles, err := u.artpieceRepo.GetFilesForArtpiece(ctx, artpieceID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	newSet := make(map[uuid.UUID]struct{}, len(fileIDs))
	for _, id := range fileIDs {
		newSet[id] = struct{}{}
	}
	currentSet := make(map[uuid.UUID]struct{}, len(currentFiles))
	for _, f := range currentFiles {
		currentSet[f.ID] = struct{}{}
	}

	var toAttach, toDetach []uuid.UUID
	for id := range newSet {
		if _, exists := currentSet[id]; !exists {
			toAttach = append(toAttach, id)
		}
	}
	for id := range currentSet {
		if _, exists := newSet[id]; !exists {
			toDetach = append(toDetach, id)
		}
	}

	err = u.transactor.InTransaction(ctx, func(ctx context.Context) error {
		if len(toAttach) > 0 {
			if err := u.fileRepo.BulkUpdateArtpieceID(ctx, toAttach, &artpieceID); err != nil {
				return err
			}
		}
		if len(toDetach) > 0 {
			if err := u.fileRepo.BulkUpdateArtpieceID(ctx, toDetach, nil); err != nil {
				return err
			}
		}

		finalFiles, err := u.fileRepo.GetFilesForArtpiece(ctx, artpieceID)
		if err != nil {
			return err
		}
		return u.setCoverFromFileSet(ctx, artpieceID, artpiece.CoverFileID, finalFiles)
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return u.artpieceRepo.GetByID(ctx, artpieceID, userID)
}

func dedupUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	out := ids[:0:0]
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			out = append(out, id)
		}
	}
	return out
}

// setCoverFromFileSet applies cover assignment rules given a known file set.
// currentCoverID is nil if the artpiece has no cover.
func (u *ArtpieceUsecase) setCoverFromFileSet(ctx context.Context, artpieceID uuid.UUID, currentCoverID *uuid.UUID, files []*domain.File) error {
	if len(files) == 0 {
		return u.artpieceRepo.UpdateCover(ctx, artpieceID, nil)
	}

	// If the current cover is still in the set, leave it unchanged.
	if currentCoverID != nil {
		for _, f := range files {
			if f.ID == *currentCoverID {
				return nil
			}
		}
	}

	// Cover is absent or was removed — pick the best candidate.
	for _, f := range files {
		if bookmime.IsImage(f.MimeType) {
			coverID := f.ID
			return u.artpieceRepo.UpdateCover(ctx, artpieceID, &coverID)
		}
	}
	coverID := files[0].ID
	return u.artpieceRepo.UpdateCover(ctx, artpieceID, &coverID)
}
