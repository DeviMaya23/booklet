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

func seedUser(t *testing.T, tx *gorm.DB, idpSubject string) *domain.User {
	t.Helper()
	u := &domain.User{ID: uuid.New(), IDPSubject: idpSubject}
	require.NoError(t, tx.FirstOrCreate(u, "idp_subject = ?", idpSubject).Error)
	return u
}

func seedCharacter(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.Character {
	t.Helper()
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: userID,
		Name:   "Test Character",
	}
	require.NoError(t, tx.Create(c).Error)
	return c
}

func seedCharacterWithFolders(t *testing.T, tx *gorm.DB, userID uuid.UUID, folderIDs []uuid.UUID) *domain.Character {
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

	user := seedUser(t, tx, "user_1")
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Aria",
	}
	err := repo.Create(context.Background(), c)

	require.NoError(t, err)

	var got domain.Character
	require.NoError(t, tx.First(&got, "id = ?", c.ID).Error)
	assert.Equal(t, "Aria", got.Name)
	assert.Equal(t, user.ID, got.UserID)

	var folderCount int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&folderCount).Error)
	assert.Equal(t, int64(0), folderCount)
}

func TestCharacterRepository_Create_WithFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	folderID1 := uuid.New()
	folderID2 := uuid.New()
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Aria",
		Folders: []domain.CharacterFolder{
			{FolderID: folderID1},
			{FolderID: folderID2},
		},
	}
	c.Folders[0].CharacterID = c.ID
	c.Folders[1].CharacterID = c.ID

	err := repo.Create(context.Background(), c)
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), c.ID.String(), user.ID)
	require.NoError(t, err)
	require.Len(t, got.Folders, 2)
	folderIDsGot := []uuid.UUID{got.Folders[0].FolderID, got.Folders[1].FolderID}
	assert.Contains(t, folderIDsGot, folderID1)
	assert.Contains(t, folderIDsGot, folderID2)
}

func TestCharacterRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, user.ID, []uuid.UUID{folderID})

	got, err := repo.GetByID(context.Background(), c.ID.String(), user.ID)

	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	require.Len(t, got.Folders, 1)
	assert.Equal(t, folderID, got.Folders[0].FolderID)
}

func TestCharacterRepository_GetByID_NoFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCharacter(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), c.ID.String(), user.ID)

	require.NoError(t, err)
	assert.Empty(t, got.Folders)
}

func TestCharacterRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	c := seedCharacter(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), c.ID.String(), user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_List(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	folderID := uuid.New()
	seedCharacterWithFolders(t, tx, user1.ID, []uuid.UUID{folderID})
	seedCharacter(t, tx, user1.ID)
	seedCharacter(t, tx, user2.ID)

	got, err := repo.List(context.Background(), user1.ID, usecase.ListCharacterFilters{})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, c := range got {
		assert.Equal(t, user1.ID, c.UserID)
		assert.NotNil(t, c.Folders)
	}
}

func TestCharacterRepository_List_QFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	aria := &domain.Character{ID: uuid.New(), UserID: user.ID, Name: "Aria Stormweaver"}
	bob := &domain.Character{ID: uuid.New(), UserID: user.ID, Name: "Bob Morrow"}
	require.NoError(t, tx.Create(aria).Error)
	require.NoError(t, tx.Create(bob).Error)

	q := "aria"
	got, err := repo.List(context.Background(), user.ID, usecase.ListCharacterFilters{Q: &q})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Aria Stormweaver", got[0].Name)
}

func TestCharacterRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCharacter(t, tx, user.ID)
	newName := "Updated Name"

	got, err := repo.Update(context.Background(), c.ID.String(), user.ID, usecase.UpdateCharacterParams{
		Name: &newName,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
}

func TestCharacterRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	newName := "x"
	_, err := repo.Update(context.Background(), uuid.NewString(), user.ID, usecase.UpdateCharacterParams{
		Name: &newName,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Update_NotFound_WithFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	folderIDs := []uuid.UUID{uuid.New()}
	_, err := repo.Update(context.Background(), uuid.NewString(), user.ID, usecase.UpdateCharacterParams{
		FolderIDs: &folderIDs,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Update_ReplaceFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	oldFolderA := uuid.New()
	oldFolderB := uuid.New()
	c := seedCharacterWithFolders(t, tx, user.ID, []uuid.UUID{oldFolderA, oldFolderB})

	newFolder := uuid.New()
	newFolders := []uuid.UUID{newFolder}

	got, err := repo.Update(context.Background(), c.ID.String(), user.ID, usecase.UpdateCharacterParams{
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

	user := seedUser(t, tx, "user_1")
	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, user.ID, []uuid.UUID{folderID})

	emptyFolders := []uuid.UUID{}

	got, err := repo.Update(context.Background(), c.ID.String(), user.ID, usecase.UpdateCharacterParams{
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

	user := seedUser(t, tx, "user_1")
	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, user.ID, []uuid.UUID{folderID})

	newName := "New Name"
	got, err := repo.Update(context.Background(), c.ID.String(), user.ID, usecase.UpdateCharacterParams{
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

	user := seedUser(t, tx, "user_1")
	c := seedCharacter(t, tx, user.ID)

	err := repo.Delete(context.Background(), c.ID.String(), user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), c.ID.String(), user.ID)
	assert.Error(t, err)
}

func TestCharacterRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	err := repo.Delete(context.Background(), uuid.NewString(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Delete_CascadesFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, user.ID, []uuid.UUID{folderID})

	err := repo.Delete(context.Background(), c.ID.String(), user.ID)
	require.NoError(t, err)

	var count int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}
