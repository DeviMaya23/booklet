package usecase_test

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
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

type fakeArtpieceArtistRepository struct {
	artists map[uuid.UUID]*domain.Artist
}

func newFakeArtpieceArtistRepository() *fakeArtpieceArtistRepository {
	return &fakeArtpieceArtistRepository{artists: make(map[uuid.UUID]*domain.Artist)}
}

func (f *fakeArtpieceArtistRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error) {
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

func (f *fakeArtpieceFileRepository) GetFilesForArtpiece(_ context.Context, artpieceID uuid.UUID) ([]*domain.File, error) {
	var result []*domain.File
	for _, file := range f.files {
		if file.ArtpieceID != nil && *file.ArtpieceID == artpieceID {
			result = append(result, file)
		}
	}
	return result, nil
}

func newArtpieceUsecase(
	artpieceRepo *fakeArtpieceRepository,
	artistRepo *fakeArtpieceArtistRepository,
	charRepo *fakeArtpieceCharacterRepository,
	fileRepo *fakeArtpieceFileRepository,
) *usecase.ArtpieceUsecase {
	return usecase.NewArtpieceUsecase(artpieceRepo, artistRepo, charRepo, fileRepo, observability.NewTelemetry(nil, nil, nil))
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
	fileRepo.files[coverFileID] = &domain.File{ID: coverFileID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/jpeg"}
	fileRepo.files[otherImageID] = &domain.File{ID: otherImageID, UserID: userID, ArtpieceID: &artpieceID, MimeType: "image/png"}

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
