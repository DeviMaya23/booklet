package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/devi/booklet/internal/worker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- fakes ---

type fakeArtpieceRepository struct {
	artpieces map[uuid.UUID]*domain.Artpiece
	files     map[uuid.UUID][]*domain.File

	lastCoverID    *uuid.UUID
	lastCoverArtID uuid.UUID
}

func newFakeArtpieceRepository() *fakeArtpieceRepository {
	return &fakeArtpieceRepository{
		artpieces: make(map[uuid.UUID]*domain.Artpiece),
		files:     make(map[uuid.UUID][]*domain.File),
	}
}

func (f *fakeArtpieceRepository) Create(_ context.Context, a *domain.Artpiece) (*domain.Artpiece, error) {
	f.artpieces[a.ID] = a
	return a, nil
}

func (f *fakeArtpieceRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	a, ok := f.artpieces[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

func (f *fakeArtpieceRepository) List(_ context.Context, userID uuid.UUID, _ usecase.ListArtpieceFilters) ([]*domain.Artpiece, error) {
	var result []*domain.Artpiece
	for _, a := range f.artpieces {
		if a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (f *fakeArtpieceRepository) Update(_ context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateArtpieceParams) (*domain.Artpiece, error) {
	a, ok := f.artpieces[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	a.Title = params.Title
	a.Notes = params.Notes
	a.ArtistID = params.ArtistID
	return a, nil
}

func (f *fakeArtpieceRepository) Delete(_ context.Context, id uuid.UUID, userID uuid.UUID) error {
	a, ok := f.artpieces[id]
	if !ok || a.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	delete(f.artpieces, id)
	return nil
}

func (f *fakeArtpieceRepository) UpdateCover(_ context.Context, artpieceID uuid.UUID, coverFileID *uuid.UUID) error {
	f.lastCoverArtID = artpieceID
	f.lastCoverID = coverFileID
	if a, ok := f.artpieces[artpieceID]; ok {
		a.CoverFileID = coverFileID
	}
	return nil
}

func (f *fakeArtpieceRepository) GetFilesForArtpiece(_ context.Context, artpieceID uuid.UUID) ([]*domain.File, error) {
	return f.files[artpieceID], nil
}

func (f *fakeArtpieceRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	a, ok := f.artpieces[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

func (f *fakeArtpieceRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.Artpiece, error) {
	var result []*domain.Artpiece
	for _, id := range ids {
		if a, ok := f.artpieces[id]; ok && a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (f *fakeArtpieceRepository) BulkUpdateCommissionID(_ context.Context, artpieceIDs []uuid.UUID, commissionID *uuid.UUID) error {
	for _, id := range artpieceIDs {
		if a, ok := f.artpieces[id]; ok {
			a.CommissionID = commissionID
		}
	}
	return nil
}

func (f *fakeArtpieceRepository) GetArtpiecesForCommission(_ context.Context, commissionID uuid.UUID) ([]*domain.Artpiece, error) {
	var result []*domain.Artpiece
	for _, a := range f.artpieces {
		if a.CommissionID != nil && *a.CommissionID == commissionID {
			result = append(result, a)
		}
	}
	return result, nil
}

type fakeArtpieceArtistRepository struct {
	artists map[uuid.UUID]*domain.Artist
}

func newFakeArtpieceArtistRepository() *fakeArtpieceArtistRepository {
	return &fakeArtpieceArtistRepository{artists: make(map[uuid.UUID]*domain.Artist)}
}

func (f *fakeArtpieceArtistRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error) {
	a, ok := f.artists[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

type fakeArtpieceCharacterRepository struct {
	characters map[uuid.UUID]*domain.Character
}

func newFakeArtpieceCharacterRepository() *fakeArtpieceCharacterRepository {
	return &fakeArtpieceCharacterRepository{characters: make(map[uuid.UUID]*domain.Character)}
}

func (f *fakeArtpieceCharacterRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Character, error) {
	var result []domain.Character
	for _, id := range ids {
		if c, ok := f.characters[id]; ok && c.UserID == userID {
			result = append(result, *c)
		}
	}
	return result, nil
}

type fakeArtpieceFileRepository struct {
	files map[uuid.UUID]*domain.File
}

func newFakeArtpieceFileRepository() *fakeArtpieceFileRepository {
	return &fakeArtpieceFileRepository{files: make(map[uuid.UUID]*domain.File)}
}

func (f *fakeArtpieceFileRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error) {
	file, ok := f.files[id]
	if !ok || file.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return file, nil
}

func (f *fakeArtpieceFileRepository) UpdateArtpieceID(_ context.Context, fileID uuid.UUID, artpieceID *uuid.UUID) error {
	if file, ok := f.files[fileID]; ok {
		file.ArtpieceID = artpieceID
	}
	return nil
}

func (f *fakeArtpieceFileRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.File, error) {
	var result []*domain.File
	for _, id := range ids {
		file, ok := f.files[id]
		if ok && file.UserID == userID {
			result = append(result, file)
		}
	}
	return result, nil
}

func (f *fakeArtpieceFileRepository) BulkUpdateArtpieceID(_ context.Context, fileIDs []uuid.UUID, artpieceID *uuid.UUID) error {
	for _, id := range fileIDs {
		if file, ok := f.files[id]; ok {
			file.ArtpieceID = artpieceID
		}
	}
	return nil
}

func (f *fakeArtpieceFileRepository) GetFilesForArtpiece(_ context.Context, artpieceID uuid.UUID) ([]*domain.File, error) {
	var result []*domain.File
	for _, file := range f.files {
		if file.ArtpieceID != nil && *file.ArtpieceID == artpieceID {
			result = append(result, file)
		}
	}
	return result, nil
}

func (f *fakeArtpieceFileRepository) BulkDelete(_ context.Context, ids []uuid.UUID, _ uuid.UUID) error {
	for _, id := range ids {
		delete(f.files, id)
	}
	return nil
}

func newArtpieceUsecase(
	artpieceRepo *fakeArtpieceRepository,
	artistRepo *fakeArtpieceArtistRepository,
	charRepo *fakeArtpieceCharacterRepository,
	fileRepo *fakeArtpieceFileRepository,
) *usecase.ArtpieceUsecase {
	return newArtpieceUsecaseWithJobInserter(artpieceRepo, artistRepo, charRepo, fileRepo, &spyJobInserter{})
}

func newArtpieceUsecaseWithJobInserter(
	artpieceRepo *fakeArtpieceRepository,
	artistRepo *fakeArtpieceArtistRepository,
	charRepo *fakeArtpieceCharacterRepository,
	fileRepo *fakeArtpieceFileRepository,
	jobInserter *spyJobInserter,
) *usecase.ArtpieceUsecase {
	return usecase.NewArtpieceUsecase(artpieceRepo, artistRepo, charRepo, fileRepo, &fakeTransactor{}, jobInserter, observability.NewTelemetry(nil, nil, nil))
}

// --- create tests ---

func TestCreateArtpiece_ArtistNotOwned(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	artistRepo := newFakeArtpieceArtistRepository()
	// artist not seeded — GetByIDAndUserID will return not found

	uc := newArtpieceUsecase(artpieceRepo, artistRepo, newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository())
	artistID := uuid.New()
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtpieceParams{ArtistID: &artistID})

	require.ErrorIs(t, err, usecase.ErrArtistNotOwned)
}

func TestCreateArtpiece_CharacterNotOwned(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	artistRepo := newFakeArtpieceArtistRepository()
	charRepo := newFakeArtpieceCharacterRepository()
	// character not seeded

	uc := newArtpieceUsecase(artpieceRepo, artistRepo, charRepo, newFakeArtpieceFileRepository())
	charID := uuid.New()
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtpieceParams{CharacterIDs: []uuid.UUID{charID}})

	require.ErrorIs(t, err, usecase.ErrCharacterNotOwned)
}

// --- update tests ---

func TestUpdateArtpiece_ClearArtist(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	artistID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID, ArtistID: &artistID}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository())

	updated, err := uc.Update(context.Background(), artpieceID, userID, usecase.UpdateArtpieceParams{ArtistID: nil})

	require.NoError(t, err)
	require.Nil(t, updated.ArtistID)
}

// --- detach tests ---

func TestDetachFile_CoverReassignedToImage(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()

	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	coverFileID := uuid.New()
	otherImageID := uuid.New()

	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{
		ID: artpieceID, UserID: userID, CoverFileID: &coverFileID,
	}
	coverFile := &domain.File{ID: coverFileID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/jpeg"}
	otherFile := &domain.File{ID: otherImageID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/png"}
	fileRepo.files[coverFileID] = coverFile
	fileRepo.files[otherImageID] = otherFile
	artpieceRepo.files[artpieceID] = []*domain.File{coverFile, otherFile}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)

	_, err := uc.DetachFile(context.Background(), artpieceID, coverFileID, userID)
	require.NoError(t, err)
	require.NotNil(t, artpieceRepo.lastCoverID)
	require.Equal(t, otherImageID, *artpieceRepo.lastCoverID)
}

func TestDetachFile_CoverClearedWhenNoFilesRemain(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()

	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	fileID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{
		ID: artpieceID, UserID: userID, CoverFileID: &fileID,
	}
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/jpeg"}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)

	_, err := uc.DetachFile(context.Background(), artpieceID, fileID, userID)
	require.NoError(t, err)
	require.True(t, artpieceRepo.lastCoverArtID == artpieceID)
	require.Nil(t, artpieceRepo.lastCoverID)
}

// --- create with file IDs tests ---

func TestCreateArtpiece_WithFileIDs_Success(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()
	fileID := uuid.New()
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, MimeType: "image/jpeg"}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtpieceParams{FileIDs: []uuid.UUID{fileID}})

	require.NoError(t, err)
	require.NotNil(t, artpieceRepo.lastCoverID)
	assert.Equal(t, fileID, *artpieceRepo.lastCoverID)
}

func TestCreateArtpiece_WithFileIDs_Empty(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository())
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtpieceParams{FileIDs: []uuid.UUID{}})

	require.NoError(t, err)
	assert.Nil(t, artpieceRepo.lastCoverID)
}

func TestCreateArtpiece_WithFileIDs_FileNotOwned(t *testing.T) {
	userID := uuid.New()
	fileRepo := newFakeArtpieceFileRepository()
	// file not seeded for this user

	uc := newArtpieceUsecase(newFakeArtpieceRepository(), newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtpieceParams{FileIDs: []uuid.UUID{uuid.New()}})

	require.ErrorIs(t, err, usecase.ErrFileNotOwned)
}

func TestCreateArtpiece_WithFileIDs_FileAlreadyAttached(t *testing.T) {
	userID := uuid.New()
	otherArtpieceID := uuid.New()
	fileRepo := newFakeArtpieceFileRepository()
	fileID := uuid.New()
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &otherArtpieceID}

	uc := newArtpieceUsecase(newFakeArtpieceRepository(), newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtpieceParams{FileIDs: []uuid.UUID{fileID}})

	require.ErrorIs(t, err, usecase.ErrFileAlreadyAttached)
}

// --- replace files tests ---

func TestReplaceFiles_FullReplace(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	oldFileID := uuid.New()
	newFileID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}
	fileRepo.files[oldFileID] = &domain.File{ID: oldFileID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "video/mp4"}
	fileRepo.files[newFileID] = &domain.File{ID: newFileID, UserID: userID, MimeType: "image/jpeg"}
	artpieceRepo.files[artpieceID] = []*domain.File{fileRepo.files[oldFileID]}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.ReplaceFiles(context.Background(), artpieceID, userID, []uuid.UUID{newFileID})

	require.NoError(t, err)
	assert.Nil(t, fileRepo.files[oldFileID].ArtpieceID)
	require.NotNil(t, fileRepo.files[newFileID].ArtpieceID)
	assert.Equal(t, artpieceID, *fileRepo.files[newFileID].ArtpieceID)
	require.NotNil(t, artpieceRepo.lastCoverID)
	assert.Equal(t, newFileID, *artpieceRepo.lastCoverID)
}

func TestReplaceFiles_EmptySet_ClearsAll(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	fileID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID, CoverFileID: &fileID}
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/jpeg"}
	artpieceRepo.files[artpieceID] = []*domain.File{fileRepo.files[fileID]}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.ReplaceFiles(context.Background(), artpieceID, userID, []uuid.UUID{})

	require.NoError(t, err)
	assert.Nil(t, artpieceRepo.lastCoverID)
	assert.Equal(t, artpieceID, artpieceRepo.lastCoverArtID)
}

func TestReplaceFiles_CurrentCoverRemainsInSet(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	coverID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID, CoverFileID: &coverID}
	fileRepo.files[coverID] = &domain.File{ID: coverID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/jpeg"}
	artpieceRepo.files[artpieceID] = []*domain.File{fileRepo.files[coverID]}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.ReplaceFiles(context.Background(), artpieceID, userID, []uuid.UUID{coverID})

	require.NoError(t, err)
	// UpdateCover should not have been called since cover remains
	assert.Equal(t, uuid.Nil, artpieceRepo.lastCoverArtID)
}

func TestReplaceFiles_CoverReassignedWhenRemoved(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	oldCoverID := uuid.New()
	newFileID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID, CoverFileID: &oldCoverID}
	fileRepo.files[oldCoverID] = &domain.File{ID: oldCoverID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "video/mp4"}
	fileRepo.files[newFileID] = &domain.File{ID: newFileID, UserID: userID, MimeType: "image/jpeg"}
	artpieceRepo.files[artpieceID] = []*domain.File{fileRepo.files[oldCoverID]}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.ReplaceFiles(context.Background(), artpieceID, userID, []uuid.UUID{newFileID})

	require.NoError(t, err)
	require.NotNil(t, artpieceRepo.lastCoverID)
	assert.Equal(t, newFileID, *artpieceRepo.lastCoverID)
}

func TestReplaceFiles_FileNotOwned(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository())
	_, err := uc.ReplaceFiles(context.Background(), artpieceID, userID, []uuid.UUID{uuid.New()})

	require.ErrorIs(t, err, usecase.ErrFileNotOwned)
}

func TestReplaceFiles_FileAttachedToOtherArtpiece(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	otherArtpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	fileID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &otherArtpieceID}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)
	_, err := uc.ReplaceFiles(context.Background(), artpieceID, userID, []uuid.UUID{fileID})

	require.ErrorIs(t, err, usecase.ErrFileAlreadyAttached)
}

func TestReplaceFiles_ArtpieceNotFound(t *testing.T) {
	userID := uuid.New()
	uc := newArtpieceUsecase(newFakeArtpieceRepository(), newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository())

	_, err := uc.ReplaceFiles(context.Background(), uuid.New(), userID, []uuid.UUID{})

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

// --- set cover tests ---

func TestSetCover_FileNotInArtpiece(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	otherArtpieceID := uuid.New()

	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	fileID := uuid.New()
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &otherArtpieceID}

	uc := newArtpieceUsecase(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo)

	_, err := uc.SetCover(context.Background(), artpieceID, fileID, userID)
	require.ErrorIs(t, err, usecase.ErrFileNotInArtpiece)
}

// --- delete tests ---

func TestDeleteArtpiece_DeleteFilesfalse_NoFileDeletion(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}
	fileID := uuid.New()
	fileRepo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &artpieceID}

	spy := &spyJobInserter{}
	uc := newArtpieceUsecaseWithJobInserter(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo, spy)

	err := uc.Delete(context.Background(), artpieceID, userID, false)

	require.NoError(t, err)
	assert.Nil(t, artpieceRepo.artpieces[artpieceID], "artpiece should be deleted")
	assert.NotNil(t, fileRepo.files[fileID], "file should remain")
	assert.Nil(t, spy.lastArgs, "no job should be enqueued")
}

func TestDeleteArtpiece_DeleteFilesTrue_WithFiles(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	thumbPath := "thumb/path"
	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}
	fileID := uuid.New()
	fileRepo.files[fileID] = &domain.File{
		ID: fileID, UserID: userID, ArtpieceID: &artpieceID,
		FileR2Path: "file/path", ThumbnailR2Path: &thumbPath,
	}

	spy := &spyJobInserter{}
	uc := newArtpieceUsecaseWithJobInserter(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo, spy)

	err := uc.Delete(context.Background(), artpieceID, userID, true)

	require.NoError(t, err)
	assert.Nil(t, artpieceRepo.artpieces[artpieceID], "artpiece should be deleted")
	assert.Nil(t, fileRepo.files[fileID], "file should be deleted")
	require.NotNil(t, spy.lastArgs)
	purgeArgs, ok := spy.lastArgs.(worker.PurgeR2ObjectsArgs)
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"file/path", "thumb/path"}, purgeArgs.R2Keys)
}

func TestDeleteArtpiece_DeleteFilesTrue_NoFiles(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()

	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}

	spy := &spyJobInserter{}
	uc := newArtpieceUsecaseWithJobInserter(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository(), spy)

	err := uc.Delete(context.Background(), artpieceID, userID, true)

	require.NoError(t, err)
	assert.Nil(t, artpieceRepo.artpieces[artpieceID], "artpiece should be deleted")
	assert.Nil(t, spy.lastArgs, "no job should be enqueued when there are no files")
}

func TestDeleteArtpiece_DeleteFilesTrue_ArtpieceNotFound(t *testing.T) {
	userID := uuid.New()

	uc := newArtpieceUsecaseWithJobInserter(newFakeArtpieceRepository(), newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), newFakeArtpieceFileRepository(), &spyJobInserter{})

	err := uc.Delete(context.Background(), uuid.New(), userID, true)

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteArtpiece_DeleteFilesTrue_JobEnqueueFails_ReturnsNil(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	artpieceRepo := newFakeArtpieceRepository()
	fileRepo := newFakeArtpieceFileRepository()

	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}
	fileID := uuid.New()
	fileRepo.files[fileID] = &domain.File{
		ID: fileID, UserID: userID, ArtpieceID: &artpieceID, FileR2Path: "file/path",
	}

	spy := &spyJobInserter{returnErr: errors.New("river down")}
	uc := newArtpieceUsecaseWithJobInserter(artpieceRepo, newFakeArtpieceArtistRepository(), newFakeArtpieceCharacterRepository(), fileRepo, spy)

	err := uc.Delete(context.Background(), artpieceID, userID, true)

	require.NoError(t, err, "enqueue failure should not surface as an error")
	assert.Nil(t, artpieceRepo.artpieces[artpieceID], "artpiece should still be deleted")
	assert.Nil(t, fileRepo.files[fileID], "file DB row should still be deleted")
}
