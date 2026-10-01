package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- fakes ---

type fakeFileRepository struct {
	files map[uuid.UUID]*domain.File

	lastUpdatedNotesID    uuid.UUID
	lastUpdatedNotes      *string
	lastDeletedID         uuid.UUID
}

func newFakeFileRepository() *fakeFileRepository {
	return &fakeFileRepository{files: make(map[uuid.UUID]*domain.File)}
}

func (f *fakeFileRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error) {
	file, ok := f.files[id]
	if !ok || file.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return file, nil
}

func (f *fakeFileRepository) List(_ context.Context, userID uuid.UUID, unassigned bool) ([]*domain.File, error) {
	var result []*domain.File
	for _, file := range f.files {
		if file.UserID != userID {
			continue
		}
		if unassigned && file.ArtpieceID != nil {
			continue
		}
		result = append(result, file)
	}
	return result, nil
}

func (f *fakeFileRepository) UpdateNotes(_ context.Context, id uuid.UUID, userID uuid.UUID, notes *string) error {
	file, ok := f.files[id]
	if !ok || file.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	f.lastUpdatedNotesID = id
	f.lastUpdatedNotes = notes
	file.Notes = notes
	return nil
}

func (f *fakeFileRepository) Delete(_ context.Context, id uuid.UUID, userID uuid.UUID) error {
	file, ok := f.files[id]
	if !ok || file.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	f.lastDeletedID = id
	delete(f.files, id)
	return nil
}

type fakeFileStorageService struct {
	deleteErr   error
	deletedKeys []string
}

func (f *fakeFileStorageService) DeleteObject(_ context.Context, key string) error {
	f.deletedKeys = append(f.deletedKeys, key)
	return f.deleteErr
}

func newFileUsecase(repo *fakeFileRepository, storage *fakeFileStorageService) *usecase.FileUsecase {
	return usecase.NewFileUsecase(repo, storage, observability.NewTelemetry(nil, nil, nil))
}

// --- tests ---

func TestFileUsecase_GetByID_NotOwned(t *testing.T) {
	repo := newFakeFileRepository()
	uc := newFileUsecase(repo, &fakeFileStorageService{})

	_, err := uc.GetByID(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileUsecase_List_NoFilter(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	repo := newFakeFileRepository()
	f1ID := uuid.New()
	f2ID := uuid.New()
	repo.files[f1ID] = &domain.File{ID: f1ID, UserID: userID}
	repo.files[f2ID] = &domain.File{ID: f2ID, UserID: userID, ArtpieceID: &artpieceID}

	uc := newFileUsecase(repo, &fakeFileStorageService{})

	files, err := uc.List(context.Background(), userID, false)
	require.NoError(t, err)
	assert.Len(t, files, 2)
}

func TestFileUsecase_List_Unassigned(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	repo := newFakeFileRepository()
	unassignedID := uuid.New()
	assignedID := uuid.New()
	repo.files[unassignedID] = &domain.File{ID: unassignedID, UserID: userID}
	repo.files[assignedID] = &domain.File{ID: assignedID, UserID: userID, ArtpieceID: &artpieceID}

	uc := newFileUsecase(repo, &fakeFileStorageService{})

	files, err := uc.List(context.Background(), userID, true)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, unassignedID, files[0].ID)
}

func TestFileUsecase_UpdateNotes_NotOwned(t *testing.T) {
	repo := newFakeFileRepository()
	uc := newFileUsecase(repo, &fakeFileStorageService{})

	_, err := uc.UpdateNotes(context.Background(), uuid.New(), uuid.New(), nil)

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileUsecase_Delete_DeletesR2Object(t *testing.T) {
	userID := uuid.New()
	fileID := uuid.New()
	thumbPath := "users/thumb/file.jpg"
	repo := newFakeFileRepository()
	repo.files[fileID] = &domain.File{
		ID:              fileID,
		UserID:          userID,
		FileR2Path:      "users/original/file.jpg",
		ThumbnailR2Path: &thumbPath,
	}
	storage := &fakeFileStorageService{}
	uc := newFileUsecase(repo, storage)

	err := uc.Delete(context.Background(), fileID, userID)
	require.NoError(t, err)
	assert.Contains(t, storage.deletedKeys, "users/original/file.jpg")
	assert.Contains(t, storage.deletedKeys, thumbPath)
}

func TestFileUsecase_Delete_R2ErrorStillDeletesRecord(t *testing.T) {
	userID := uuid.New()
	fileID := uuid.New()
	repo := newFakeFileRepository()
	repo.files[fileID] = &domain.File{
		ID:         fileID,
		UserID:     userID,
		FileR2Path: "users/original/file.jpg",
	}
	storage := &fakeFileStorageService{deleteErr: errors.New("r2 unavailable")}
	uc := newFileUsecase(repo, storage)

	err := uc.Delete(context.Background(), fileID, userID)
	require.NoError(t, err)
	_, err = repo.GetByIDAndUserID(context.Background(), fileID, userID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
