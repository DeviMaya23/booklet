package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreate_AssemblesCharacter(t *testing.T) {
	repo := newFakeCharacterRepository()
	bl := &fakeBookleafClient{}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	userID := uuid.New()
	got, err := uc.Create(context.Background(), userID, usecase.CreateCharacterParams{
		Name: "Aria Stormweaver",
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "Aria Stormweaver", got.Name)
	require.Empty(t, repo.lastCreated.Folders)
	require.False(t, bl.called, "Bookleaf should not be called when no folder IDs")
}

// --- InitAvatarUpload ---

func TestInitAvatarUpload_CharacterNotOwnedReturnsErrCharacterNotFound(t *testing.T) {
	charRepo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

	_, err := uc.InitAvatarUpload(context.Background(), uuid.New(), uuid.New().String(), "image/jpeg")

	require.ErrorIs(t, err, usecase.ErrCharacterNotFound)
}

func TestInitAvatarUpload_SuccessReturnsResult(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}
	avatarRepo := newFakeCharacterAvatarUploadRepository()
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, avatarRepo, &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

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
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

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
	uc := usecase.NewCharacterUsecase(charRepo, storage, avatarRepo, &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

	err := uc.CompleteAvatarUpload(context.Background(), userID, charID.String(), uploadID)

	require.NoError(t, err)
	require.Equal(t, newKey, charRepo.lastUpdatedAvatarKey)
	require.Empty(t, storage.deletedKeys)
	require.NotContains(t, avatarRepo.pending, uploadID)
}

// --- DeleteAvatar ---

func TestDeleteAvatar_CharacterNotFoundReturnsErrCharacterNotFound(t *testing.T) {
	charRepo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(charRepo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

	err := uc.DeleteAvatar(context.Background(), uuid.New(), uuid.New().String())

	require.ErrorIs(t, err, usecase.ErrCharacterNotFound)
}

func TestDeleteAvatar_NoAvatarReturnsNilWithoutR2Call(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	charRepo := newFakeCharacterRepository()
	charRepo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}

	storage := &fakeStorageService{}
	uc := usecase.NewCharacterUsecase(charRepo, storage, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

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
	uc := usecase.NewCharacterUsecase(charRepo, storage, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

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
	uc := usecase.NewCharacterUsecase(charRepo, storage, avatarRepo, &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

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
	uc := usecase.NewCharacterUsecase(newFakeCharacterRepository(), &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), imageRepo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

	got, err := uc.GetCharacterImages(context.Background(), charID, userID)

	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestGetCharacterImages_ReturnsEmptySliceWhenNone(t *testing.T) {
	userID := uuid.New()
	charID := uuid.New()
	imageRepo := &fakeImageRepository{listByCharacterID: []*domain.Image{}}
	uc := usecase.NewCharacterUsecase(newFakeCharacterRepository(), &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), imageRepo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

	got, err := uc.GetCharacterImages(context.Background(), charID, userID)

	require.NoError(t, err)
	require.Empty(t, got)
}

func TestListCharacters_PassesFiltersToRepo(t *testing.T) {
	repo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), &fakeBookleafClient{})

	userID := uuid.New()
	q := "aria"
	filters := usecase.ListCharacterFilters{Q: &q}
	_, err := uc.List(context.Background(), userID, filters)

	require.NoError(t, err)
	require.NotNil(t, repo.lastListFilters.Q)
	require.Equal(t, "aria", *repo.lastListFilters.Q)
}

// --- Bookleaf folder validation: Create ---

func TestCreate_BookleafAllValid_EnrichesAndPersists(t *testing.T) {
	repo := newFakeCharacterRepository()
	folderID := uuid.New()
	bl := &fakeBookleafClient{
		folderList: &bookleaf.FolderList{
			FolderList: []bookleaf.Folder{
				{FolderID: folderID.String(), FolderName: "My Folder"},
			},
		},
	}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	userID := uuid.New()
	folderIDs := []uuid.UUID{folderID}
	got, err := uc.Create(context.Background(), userID, usecase.CreateCharacterParams{
		Name:       "Aria",
		FolderIDs:  &folderIDs,
		IDPSubject: "kp_user_1",
	})

	require.NoError(t, err)
	require.Len(t, repo.lastCreated.Folders, 1)
	require.Equal(t, folderID, repo.lastCreated.Folders[0].FolderID)
	require.Equal(t, "My Folder", repo.lastCreated.Folders[0].FolderName)
	require.Equal(t, got.ID, repo.lastCreated.Folders[0].CharacterID)
}

func TestCreate_BookleafSomeAbsent_DropsAbsentIDs(t *testing.T) {
	repo := newFakeCharacterRepository()
	validID := uuid.New()
	unknownID := uuid.New()
	bl := &fakeBookleafClient{
		folderList: &bookleaf.FolderList{
			FolderList: []bookleaf.Folder{
				{FolderID: validID.String(), FolderName: "Valid Folder"},
			},
		},
	}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	folderIDs := []uuid.UUID{validID, unknownID}
	_, err := uc.Create(context.Background(), uuid.New(), usecase.CreateCharacterParams{
		Name:       "Aria",
		FolderIDs:  &folderIDs,
		IDPSubject: "kp_user_1",
	})

	require.NoError(t, err)
	require.Len(t, repo.lastCreated.Folders, 1)
	require.Equal(t, validID, repo.lastCreated.Folders[0].FolderID)
	require.Equal(t, "Valid Folder", repo.lastCreated.Folders[0].FolderName)
}

func TestCreate_BookleafFails_ReturnsError(t *testing.T) {
	repo := newFakeCharacterRepository()
	folderID := uuid.New()
	bookleafErr := errors.New("bookleaf unavailable")
	bl := &fakeBookleafClient{err: bookleafErr}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	folderIDs := []uuid.UUID{folderID}
	_, err := uc.Create(context.Background(), uuid.New(), usecase.CreateCharacterParams{
		Name:       "Aria",
		FolderIDs:  &folderIDs,
		IDPSubject: "kp_user_1",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "bookleaf unavailable")
	require.Nil(t, repo.lastCreated, "character should not have been persisted")
}

func TestCreate_EmptyFolderIDs_SkipsBookleafAndPersists(t *testing.T) {
	repo := newFakeCharacterRepository()
	bl := &fakeBookleafClient{}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	emptyFolders := []uuid.UUID{}
	_, err := uc.Create(context.Background(), uuid.New(), usecase.CreateCharacterParams{
		Name:       "Aria",
		FolderIDs:  &emptyFolders,
		IDPSubject: "kp_user_1",
	})

	require.NoError(t, err)
	require.False(t, bl.called, "Bookleaf should not be called for empty folder_ids")
	require.Empty(t, repo.lastCreated.Folders)
}

// --- Bookleaf folder validation: Update ---

func TestUpdate_BookleafAllValid_EnrichesAndPersists(t *testing.T) {
	repo := newFakeCharacterRepository()
	userID := uuid.New()
	charID := uuid.New()
	repo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}

	folderID := uuid.New()
	bl := &fakeBookleafClient{
		folderList: &bookleaf.FolderList{
			FolderList: []bookleaf.Folder{
				{FolderID: folderID.String(), FolderName: "My Folder"},
			},
		},
	}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	got, err := uc.Update(context.Background(), charID.String(), userID, usecase.UpdateCharacterParams{
		Name:       "Aria",
		FolderIDs:  []uuid.UUID{folderID},
		IDPSubject: "kp_user_1",
	})

	require.NoError(t, err)
	require.Len(t, got.Folders, 1)
	require.Equal(t, folderID, got.Folders[0].FolderID)
	require.Equal(t, "My Folder", got.Folders[0].FolderName)
}

func TestUpdate_BookleafSomeAbsent_DropsAbsentIDs(t *testing.T) {
	repo := newFakeCharacterRepository()
	userID := uuid.New()
	charID := uuid.New()
	repo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}

	validID := uuid.New()
	unknownID := uuid.New()
	bl := &fakeBookleafClient{
		folderList: &bookleaf.FolderList{
			FolderList: []bookleaf.Folder{
				{FolderID: validID.String(), FolderName: "Valid Folder"},
			},
		},
	}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	got, err := uc.Update(context.Background(), charID.String(), userID, usecase.UpdateCharacterParams{
		Name:       "Aria",
		FolderIDs:  []uuid.UUID{validID, unknownID},
		IDPSubject: "kp_user_1",
	})

	require.NoError(t, err)
	require.Len(t, got.Folders, 1)
	require.Equal(t, validID, got.Folders[0].FolderID)
}

func TestUpdate_BookleafFails_ReturnsError(t *testing.T) {
	repo := newFakeCharacterRepository()
	userID := uuid.New()
	charID := uuid.New()
	repo.characters[charID] = &domain.Character{ID: charID, UserID: userID, Name: "Aria"}

	bookleafErr := errors.New("bookleaf unavailable")
	bl := &fakeBookleafClient{err: bookleafErr}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	_, err := uc.Update(context.Background(), charID.String(), userID, usecase.UpdateCharacterParams{
		Name:       "Aria",
		FolderIDs:  []uuid.UUID{uuid.New()},
		IDPSubject: "kp_user_1",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "bookleaf unavailable")
}

func TestUpdate_EmptyFolderIDs_SkipsBookleafAndClearsAssignments(t *testing.T) {
	existingFolderID := uuid.New()
	repo := newFakeCharacterRepository()
	userID := uuid.New()
	charID := uuid.New()
	repo.characters[charID] = &domain.Character{
		ID:     charID,
		UserID: userID,
		Name:   "Aria",
		Folders: []domain.CharacterFolder{
			{CharacterID: charID, FolderID: existingFolderID, FolderName: "Old Folder"},
		},
	}

	bl := &fakeBookleafClient{}
	uc := usecase.NewCharacterUsecase(repo, &fakeStorageService{}, newFakeCharacterAvatarUploadRepository(), &fakeImageRepository{}, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil), bl)

	got, err := uc.Update(context.Background(), charID.String(), userID, usecase.UpdateCharacterParams{
		Name:       "Aria",
		FolderIDs:  []uuid.UUID{},
		IDPSubject: "kp_user_1",
	})

	require.NoError(t, err)
	require.Empty(t, got.Folders)
	require.False(t, bl.called, "Bookleaf should not be called for empty folder_ids")
}
