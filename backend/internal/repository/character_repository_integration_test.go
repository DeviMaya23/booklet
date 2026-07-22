package repository

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedUser(t *testing.T, tx *gorm.DB, userID string) {
	t.Helper()
	require.NoError(t, tx.FirstOrCreate(&domain.User{ID: userID}).Error)
}

func seedCharacter(t *testing.T, tx *gorm.DB, userID string) *domain.Character {
	t.Helper()
	seedUser(t, tx, userID)
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: userID,
		Name:   "Test Character",
	}
	require.NoError(t, tx.Create(c).Error)
	return c
}

func seedCharacterWithFolders(t *testing.T, tx *gorm.DB, userID string, folderIDs []uuid.UUID) *domain.Character {
	t.Helper()
	c := seedCharacter(t, tx, userID)
	for _, fid := range folderIDs {
		require.NoError(t, tx.Create(&domain.CharacterFolder{CharacterID: c.ID, FolderID: fid}).Error)
	}
	return c
}

func TestCharacterRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	seedUser(t, tx, "user_1")
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: "user_1",
		Name:   "Aria",
	}
	err := repo.Create(context.Background(), c)

	require.NoError(t, err)

	var got domain.Character
	require.NoError(t, tx.First(&got, "id = ?", c.ID).Error)
	assert.Equal(t, "Aria", got.Name)
	assert.Equal(t, "user_1", got.UserID)

	var folderCount int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&folderCount).Error)
	assert.Equal(t, int64(0), folderCount)
}

func TestCharacterRepository_Create_WithFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	seedUser(t, tx, "user_1")
	folderID1 := uuid.New()
	folderID2 := uuid.New()
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: "user_1",
		Name:   "Aria",
		Folders: []domain.CharacterFolder{
			{FolderID: folderID1},
			{FolderID: folderID2},
		},
	}
	// CharacterID on folders is set by usecase; simulate that here
	c.Folders[0].CharacterID = c.ID
	c.Folders[1].CharacterID = c.ID

	err := repo.Create(context.Background(), c)
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), c.ID.String(), "user_1")
	require.NoError(t, err)
	require.Len(t, got.Folders, 2)
	folderIDsGot := []uuid.UUID{got.Folders[0].FolderID, got.Folders[1].FolderID}
	assert.Contains(t, folderIDsGot, folderID1)
	assert.Contains(t, folderIDsGot, folderID2)
}

func TestCharacterRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, "user_1", []uuid.UUID{folderID})

	got, err := repo.GetByID(context.Background(), c.ID.String(), "user_1")

	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	require.Len(t, got.Folders, 1)
	assert.Equal(t, folderID, got.Folders[0].FolderID)
}

func TestCharacterRepository_GetByID_NoFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")

	got, err := repo.GetByID(context.Background(), c.ID.String(), "user_1")

	require.NoError(t, err)
	assert.Empty(t, got.Folders)
}

func TestCharacterRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")

	_, err := repo.GetByID(context.Background(), c.ID.String(), "user_2")

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_List(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	folderID := uuid.New()
	seedCharacterWithFolders(t, tx, "user_1", []uuid.UUID{folderID})
	seedCharacter(t, tx, "user_1")
	seedCharacter(t, tx, "user_2")

	got, err := repo.List(context.Background(), "user_1")

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, c := range got {
		assert.Equal(t, "user_1", c.UserID)
		assert.NotNil(t, c.Folders)
	}
}

func TestCharacterRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")
	newName := "Updated Name"

	got, err := repo.Update(context.Background(), c.ID.String(), "user_1", usecase.UpdateCharacterParams{
		Name: &newName,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
}

func TestCharacterRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	newName := "x"
	_, err := repo.Update(context.Background(), uuid.NewString(), "user_1", usecase.UpdateCharacterParams{
		Name: &newName,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Update_NotFound_WithFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	folderIDs := []uuid.UUID{uuid.New()}
	_, err := repo.Update(context.Background(), uuid.NewString(), "user_1", usecase.UpdateCharacterParams{
		FolderIDs: &folderIDs,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Update_ReplaceFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	oldFolderA := uuid.New()
	oldFolderB := uuid.New()
	c := seedCharacterWithFolders(t, tx, "user_1", []uuid.UUID{oldFolderA, oldFolderB})

	newFolder := uuid.New()
	newFolders := []uuid.UUID{newFolder}

	got, err := repo.Update(context.Background(), c.ID.String(), "user_1", usecase.UpdateCharacterParams{
		FolderIDs: &newFolders,
	})

	require.NoError(t, err)
	require.Len(t, got.Folders, 1)
	assert.Equal(t, newFolder, got.Folders[0].FolderID)

	var count int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ? AND folder_id IN ?", c.ID, []uuid.UUID{oldFolderA, oldFolderB}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestCharacterRepository_Update_ClearFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, "user_1", []uuid.UUID{folderID})

	emptyFolders := []uuid.UUID{}

	got, err := repo.Update(context.Background(), c.ID.String(), "user_1", usecase.UpdateCharacterParams{
		FolderIDs: &emptyFolders,
	})

	require.NoError(t, err)
	assert.Empty(t, got.Folders)

	var count int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestCharacterRepository_Update_NilFolders_NoOp(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, "user_1", []uuid.UUID{folderID})

	newName := "New Name"
	got, err := repo.Update(context.Background(), c.ID.String(), "user_1", usecase.UpdateCharacterParams{
		Name:      &newName,
		FolderIDs: nil,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
	require.Len(t, got.Folders, 1)
	assert.Equal(t, folderID, got.Folders[0].FolderID)
}

func TestCharacterRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")

	err := repo.Delete(context.Background(), c.ID.String(), "user_1")

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), c.ID.String(), "user_1")
	assert.Error(t, err)
}

func TestCharacterRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	err := repo.Delete(context.Background(), uuid.NewString(), "user_1")

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Delete_CascadesFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, "user_1", []uuid.UUID{folderID})

	err := repo.Delete(context.Background(), c.ID.String(), "user_1")
	require.NoError(t, err)

	var count int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}
