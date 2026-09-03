package usecase_test

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreate_AssemblesCharacter(t *testing.T) {
	repo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	got, err := uc.Create(context.Background(), userID, usecase.CreateCharacterParams{
		Name: "Aria Stormweaver",
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "Aria Stormweaver", got.Name)
	require.Empty(t, repo.lastCreated.Folders)
}

func TestCreate_AssemblesCharacter_WithFolders(t *testing.T) {
	repo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	folderID := uuid.New()
	folderIDs := []uuid.UUID{folderID}

	got, err := uc.Create(context.Background(), userID, usecase.CreateCharacterParams{
		Name:      "Aria Stormweaver",
		FolderIDs: &folderIDs,
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Len(t, repo.lastCreated.Folders, 1)
	require.Equal(t, folderID, repo.lastCreated.Folders[0].FolderID)
	require.Equal(t, got.ID, repo.lastCreated.Folders[0].CharacterID)
}

// --- InitAvatarUpload ---

func TestInitAvatarUpload_CharacterNotOwnedReturnsErrCharacterNotFound(t *testing.T) {
	charRepo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.InitAvatarUpload(context.Background(), uuid.New(), uuid.New().String(), "image/jpeg")

	require.ErrorIs(t, err, usecase.ErrCharacterNotFound)
}

func TestInitAvatarUpload_SuccessReturnsResult(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}
	avatarRepo := newFakeCharacterAvatarUploadRepository()
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, avatarRepo, &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	result, err := uc.InitAvatarUpload(context.Background(), userID, charID.String(), "image/jpeg")

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, result.ID)
	require.NotEmpty(t, result.UploadURL)
	require.False(t, result.ExpiresAt.IsZero())
	require.Len(t, avatarRepo.pending, 1)
}

// --- CompleteAvatarUpload ---

func TestCompleteAvatarUpload_PendingNotFoundReturnsError(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteAvatarUpload(context.Background(), userID, charID.String(), uuid.New())

	require.ErrorIs(t, err, usecase.ErrPendingUploadNotFound)
}

func TestCompleteAvatarUpload_SuccessNoPriorAvatar(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	uploadID := uuid.New()
	newKey := userID.String() + "/files/" + uploadID.String()

	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}

	avatarRepo := newFakeCharacterAvatarUploadRepository()
	avatarRepo.pending[uploadID] = &domain.PendingCharacterAvatarUpload{
		ID:          uploadID,
		UserID:      userID,
		CharacterID: charID,
		R2Key:       newKey,
	}

	storage := &fakeStorageService{}
	uc := usecase.NewCharacterUsecase(charRepo, storage, avatarRepo, &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteAvatarUpload(context.Background(), userID, charID.String(), uploadID)

	require.NoError(t, err)
	require.Equal(t, newKey, charRepo.lastUpdatedAvatarKey)
	require.Empty(t, storage.deletedKeys)
	require.NotContains(t, avatarRepo.pending, uploadID)
}

// --- DeleteAvatar ---

func TestDeleteAvatar_CharacterNotFoundReturnsErrCharacterNotFound(t *testing.T) {
	charRepo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.DeleteAvatar(context.Background(), uuid.New(), uuid.New().String())

	require.ErrorIs(t, err, usecase.ErrCharacterNotFound)
}

func TestDeleteAvatar_NoAvatarReturnsNilWithoutR2Call(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}

	storage := &fakeStorageService{}
	uc := usecase.NewCharacterUsecase(charRepo, storage, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.DeleteAvatar(context.Background(), userID, charID.String())

	require.NoError(t, err)
	require.Empty(t, storage.deletedKeys)
}

func TestDeleteAvatar_ExistingAvatarClearsPathAndDeletesR2Object(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	oldKey := "old/avatar/key"
	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria", AvatarR2Path: &oldKey}

	storage := &fakeStorageService{}
	uc := usecase.NewCharacterUsecase(charRepo, storage, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.DeleteAvatar(context.Background(), userID, charID.String())

	require.NoError(t, err)
	require.Nil(t, charRepo.characters[charID].AvatarR2Path)
	require.Contains(t, storage.deletedKeys, oldKey)
}

func TestCompleteAvatarUpload_SuccessReplacingExistingAvatar(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	uploadID := uuid.New()
	oldKey := "old/key/avatar"
	newKey := userID.String() + "/files/" + uploadID.String()

	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria", AvatarR2Path: &oldKey}

	avatarRepo := newFakeCharacterAvatarUploadRepository()
	avatarRepo.pending[uploadID] = &domain.PendingCharacterAvatarUpload{
		ID:          uploadID,
		UserID:      userID,
		CharacterID: charID,
		R2Key:       newKey,
	}

	storage := &fakeStorageService{}
	uc := usecase.NewCharacterUsecase(charRepo, storage, avatarRepo, &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteAvatarUpload(context.Background(), userID, charID.String(), uploadID)

	require.NoError(t, err)
	require.Equal(t, newKey, charRepo.lastUpdatedAvatarKey)
	require.Contains(t, storage.deletedKeys, oldKey)
	require.NotContains(t, avatarRepo.pending, uploadID)
}

// --- GetCharacterImages ---

func TestGetCharacterImages_ReturnsImages(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	imageRepo := &fakeImageRepository{
		listByCharacterID: []*domain.Image{
			{ID: uuid.New(), UserID: userID, ImageR2Path: "images/a.jpg", MimeType: "image/jpeg", Characters: []domain.Character{}},
			{ID: uuid.New(), UserID: userID, ImageR2Path: "images/b.jpg", MimeType: "image/jpeg", Characters: []domain.Character{}},
		},
	}
	uc := usecase.NewCharacterUsecase(newFakeCharacterRepository(), &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), imageRepo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	got, err := uc.GetCharacterImages(context.Background(), charID, userID)

	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestGetCharacterImages_ReturnsEmptySliceWhenNone(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	imageRepo := &fakeImageRepository{listByCharacterID: []*domain.Image{}}
	uc := usecase.NewCharacterUsecase(newFakeCharacterRepository(), &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), imageRepo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	got, err := uc.GetCharacterImages(context.Background(), charID, userID)

	require.NoError(t, err)
	require.Empty(t, got)
}
