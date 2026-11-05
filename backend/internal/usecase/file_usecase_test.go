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

	lastDeletedID    uuid.UUID
	bulkDeleteCalled bool
	bulkDeletedIDs   []uuid.UUID
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

func (f *fakeFileRepository) Update(_ context.Context, id uuid.UUID, userID uuid.UUID, name *string, notes *string, isVisible *bool) error {
	file, ok := f.files[id]
	if !ok || file.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	if isVisible == nil || name != nil {
		file.Name = name
	}
	if isVisible == nil || notes != nil {
		file.Notes = notes
	}
	if isVisible != nil {
		file.IsVisible = *isVisible
	}
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

type spyCoverReassigner struct {
	calls []reassignCall
	err   error
}

type reassignCall struct {
	artpieceID uuid.UUID
	fileID     uuid.UUID
}

func (s *spyCoverReassigner) ReassignCoverIfCurrent(_ context.Context, artpieceID uuid.UUID, fileID uuid.UUID, _ uuid.UUID) error {
	s.calls = append(s.calls, reassignCall{artpieceID: artpieceID, fileID: fileID})
	return s.err
}

func newFileUsecase(repo *fakeFileRepository, storage *fakeFileStorageService) *usecase.FileUsecase {
	return newFileUsecaseWithReassigner(repo, &spyCoverReassigner{}, storage)
}

func newFileUsecaseWithReassigner(repo *fakeFileRepository, reassigner *spyCoverReassigner, storage *fakeFileStorageService) *usecase.FileUsecase {
	return usecase.NewFileUsecase(repo, reassigner, storage, &fakeJobInserter{}, observability.NewTelemetry(nil, nil, nil))
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

func TestFileUsecase_Update_NotOwned(t *testing.T) {
	repo := newFakeFileRepository()
	uc := newFileUsecase(repo, &fakeFileStorageService{})

	_, err := uc.Update(context.Background(), uuid.New(), uuid.New(), nil, nil, nil)

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileUsecase_Update_IsVisible(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()
	boolPtr := func(b bool) *bool { return &b }

	t.Run("returned file reflects provided is_visible", func(t *testing.T) {
		fileID := uuid.New()
		repo := newFakeFileRepository()
		repo.files[fileID] = &domain.File{ID: fileID, UserID: userID, IsVisible: true}
		uc := newFileUsecase(repo, &fakeFileStorageService{})

		got, err := uc.Update(context.Background(), fileID, userID, nil, nil, boolPtr(false))

		require.NoError(t, err)
		assert.False(t, got.IsVisible)
	})

	t.Run("is_visible unchanged when omitted", func(t *testing.T) {
		fileID := uuid.New()
		repo := newFakeFileRepository()
		repo.files[fileID] = &domain.File{ID: fileID, UserID: userID, IsVisible: true}
		name := "renamed"
		uc := newFileUsecase(repo, &fakeFileStorageService{})

		got, err := uc.Update(context.Background(), fileID, userID, &name, nil, nil)

		require.NoError(t, err)
		assert.True(t, got.IsVisible)
	})

	t.Run("is_visible only leaves name and notes untouched", func(t *testing.T) {
		fileID := uuid.New()
		name, notes := "keep", "keep notes"
		repo := newFakeFileRepository()
		repo.files[fileID] = &domain.File{ID: fileID, UserID: userID, IsVisible: true, Name: &name, Notes: &notes}
		uc := newFileUsecase(repo, &fakeFileStorageService{})

		got, err := uc.Update(context.Background(), fileID, userID, nil, nil, boolPtr(false))

		require.NoError(t, err)
		require.NotNil(t, got.Name)
		require.NotNil(t, got.Notes)
		assert.Equal(t, "keep", *got.Name)
		assert.Equal(t, "keep notes", *got.Notes)
	})

	t.Run("hiding an attached file triggers cover reassignment check", func(t *testing.T) {
		fileID := uuid.New()
		repo := newFakeFileRepository()
		repo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &artpieceID, IsVisible: true}
		reassigner := &spyCoverReassigner{}
		uc := newFileUsecaseWithReassigner(repo, reassigner, &fakeFileStorageService{})

		_, err := uc.Update(context.Background(), fileID, userID, nil, nil, boolPtr(false))

		require.NoError(t, err)
		require.Len(t, reassigner.calls, 1)
		assert.Equal(t, reassignCall{artpieceID: artpieceID, fileID: fileID}, reassigner.calls[0])
	})

	t.Run("no reassignment when making visible, omitted, or unattached", func(t *testing.T) {
		attachedID, inboxID := uuid.New(), uuid.New()
		repo := newFakeFileRepository()
		repo.files[attachedID] = &domain.File{ID: attachedID, UserID: userID, ArtpieceID: &artpieceID}
		repo.files[inboxID] = &domain.File{ID: inboxID, UserID: userID, IsVisible: true}
		reassigner := &spyCoverReassigner{}
		uc := newFileUsecaseWithReassigner(repo, reassigner, &fakeFileStorageService{})

		_, err := uc.Update(context.Background(), attachedID, userID, nil, nil, boolPtr(true))
		require.NoError(t, err)
		name := "n"
		_, err = uc.Update(context.Background(), attachedID, userID, &name, nil, nil)
		require.NoError(t, err)
		_, err = uc.Update(context.Background(), inboxID, userID, nil, nil, boolPtr(false))
		require.NoError(t, err)

		assert.Empty(t, reassigner.calls)
	})

	t.Run("reassignment failure is wrapped and returned", func(t *testing.T) {
		fileID := uuid.New()
		repo := newFakeFileRepository()
		repo.files[fileID] = &domain.File{ID: fileID, UserID: userID, ArtpieceID: &artpieceID, IsVisible: true}
		boom := errors.New("boom")
		uc := newFileUsecaseWithReassigner(repo, &spyCoverReassigner{err: boom}, &fakeFileStorageService{})

		_, err := uc.Update(context.Background(), fileID, userID, nil, nil, boolPtr(false))

		require.ErrorIs(t, err, boom)
		require.ErrorContains(t, err, "reassign cover after hiding file")
	})
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
	uc := usecase.NewFileUsecase(repo, &spyCoverReassigner{}, &fakeFileStorageService{}, inserter, observability.NewTelemetry(nil, nil, nil))

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
	uc := usecase.NewFileUsecase(repo, &spyCoverReassigner{}, &fakeFileStorageService{}, inserter, observability.NewTelemetry(nil, nil, nil))

	err := uc.BulkDelete(context.Background(), []uuid.UUID{f1ID}, userID)
	require.NoError(t, err)
	assert.True(t, repo.bulkDeleteCalled)
}
