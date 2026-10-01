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
	"github.com/riverqueue/river"
	rivertype "github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- fakes ---

type fakeFileRepository struct {
	files map[uuid.UUID]*domain.File

	lastUpdatedNotesID uuid.UUID
	lastUpdatedNotes   *string
	lastDeletedID      uuid.UUID
	bulkDeleteCalled   bool
	bulkDeletedIDs     []uuid.UUID
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

func (f *fakeFileRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.File, error) {
	var result []*domain.File
	for _, id := range ids {
		file, ok := f.files[id]
		if ok && file.UserID == userID {
			result = append(result, file)
		}
	}
	return result, nil
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

func (f *fakeFileRepository) BulkDelete(_ context.Context, ids []uuid.UUID, _ uuid.UUID) error {
	f.bulkDeleteCalled = true
	f.bulkDeletedIDs = ids
	for _, id := range ids {
		delete(f.files, id)
	}
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

type fakeJobInserter struct {
	insertErr    error
	insertedArgs []any
}

func (f *fakeJobInserter) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	f.insertedArgs = append(f.insertedArgs, args)
	return &rivertype.JobInsertResult{}, f.insertErr
}

func newFileUsecase(repo *fakeFileRepository, storage *fakeFileStorageService) *usecase.FileUsecase {
	return usecase.NewFileUsecase(repo, storage, &fakeJobInserter{}, observability.NewTelemetry(nil, nil, nil))
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

func TestFileUsecase_BulkDelete_AllOwned(t *testing.T) {
	userID := uuid.New()
	f1ID, f2ID := uuid.New(), uuid.New()
	thumbPath := "users/thumb/f2.jpg"
	repo := newFakeFileRepository()
	repo.files[f1ID] = &domain.File{ID: f1ID, UserID: userID, FileR2Path: "users/files/f1.jpg"}
	repo.files[f2ID] = &domain.File{ID: f2ID, UserID: userID, FileR2Path: "users/files/f2.jpg", ThumbnailR2Path: &thumbPath}
	inserter := &fakeJobInserter{}
	uc := usecase.NewFileUsecase(repo, &fakeFileStorageService{}, inserter, observability.NewTelemetry(nil, nil, nil))

	err := uc.BulkDelete(context.Background(), []uuid.UUID{f1ID, f2ID}, userID)
	require.NoError(t, err)
	assert.True(t, repo.bulkDeleteCalled)
	require.Len(t, inserter.insertedArgs, 1)
	args, ok := inserter.insertedArgs[0].(worker.PurgeR2ObjectsArgs)
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"users/files/f1.jpg", "users/files/f2.jpg", thumbPath}, args.R2Keys)
}

func TestFileUsecase_BulkDelete_OwnershipViolation(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	f1ID, f2ID := uuid.New(), uuid.New()
	repo := newFakeFileRepository()
	repo.files[f1ID] = &domain.File{ID: f1ID, UserID: userID, FileR2Path: "users/files/f1.jpg"}
	repo.files[f2ID] = &domain.File{ID: f2ID, UserID: otherUserID, FileR2Path: "users/files/f2.jpg"}
	uc := newFileUsecase(repo, &fakeFileStorageService{})

	err := uc.BulkDelete(context.Background(), []uuid.UUID{f1ID, f2ID}, userID)
	require.ErrorIs(t, err, usecase.ErrFileOwnershipViolation)
	assert.False(t, repo.bulkDeleteCalled)
}

func TestFileUsecase_BulkDelete_JobEnqueueFailureStillSucceeds(t *testing.T) {
	userID := uuid.New()
	f1ID := uuid.New()
	repo := newFakeFileRepository()
	repo.files[f1ID] = &domain.File{ID: f1ID, UserID: userID, FileR2Path: "users/files/f1.jpg"}
	inserter := &fakeJobInserter{insertErr: errors.New("river unavailable")}
	uc := usecase.NewFileUsecase(repo, &fakeFileStorageService{}, inserter, observability.NewTelemetry(nil, nil, nil))

	err := uc.BulkDelete(context.Background(), []uuid.UUID{f1ID}, userID)
	require.NoError(t, err)
	assert.True(t, repo.bulkDeleteCalled)
}
