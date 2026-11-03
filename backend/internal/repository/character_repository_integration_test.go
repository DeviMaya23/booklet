package repository

import (
	"context"
	"testing"
	"time"

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

	got, err := repo.GetByID(context.Background(), c.ID, user.ID)
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

	got, err := repo.GetByID(context.Background(), c.ID, user.ID)

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

	got, err := repo.GetByID(context.Background(), c.ID, user.ID)

	require.NoError(t, err)
	assert.Empty(t, got.Folders)
}

func TestCharacterRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	c := seedCharacter(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), c.ID, user2.ID)

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
	got, err := repo.Update(context.Background(), c.ID, user.ID, usecase.UpdateCharacterParams{
		Name: "Updated Name",
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
}

func TestCharacterRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	_, err := repo.Update(context.Background(), uuid.New(), user.ID, usecase.UpdateCharacterParams{
		Name: "x",
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Update_NotFound_WithFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	_, err := repo.Update(context.Background(), uuid.New(), user.ID, usecase.UpdateCharacterParams{
		Name:      "x",
		FolderIDs: []uuid.UUID{uuid.New()},
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

	got, err := repo.Update(context.Background(), c.ID, user.ID, usecase.UpdateCharacterParams{
		Name: "Updated Name",
		Folders: []domain.CharacterFolder{
			{CharacterID: c.ID, FolderID: newFolder, FolderName: "New Folder"},
		},
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

	got, err := repo.Update(context.Background(), c.ID, user.ID, usecase.UpdateCharacterParams{
		Name:      "Updated Name",
		FolderIDs: []uuid.UUID{},
	})

	require.NoError(t, err)
	assert.Empty(t, got.Folders)

	var count int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}


func TestCharacterRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCharacter(t, tx, user.ID)

	err := repo.Delete(context.Background(), c.ID, user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), c.ID, user.ID)
	assert.Error(t, err)
}

func TestCharacterRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	err := repo.Delete(context.Background(), uuid.New(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Delete_CascadesFolders(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	folderID := uuid.New()
	c := seedCharacterWithFolders(t, tx, user.ID, []uuid.UUID{folderID})

	err := repo.Delete(context.Background(), c.ID, user.ID)
	require.NoError(t, err)

	var count int64
	require.NoError(t, tx.Model(&domain.CharacterFolder{}).Where("character_id = ?", c.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestCharacterRepository_List_VerseIDFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := seedVerse(t, tx, user.ID, "Solaris")

	withVerse := &domain.Character{ID: uuid.New(), UserID: user.ID, Name: "With Verse", VerseID: &verse.ID}
	withoutVerse := &domain.Character{ID: uuid.New(), UserID: user.ID, Name: "Without Verse"}
	require.NoError(t, tx.Create(withVerse).Error)
	require.NoError(t, tx.Create(withoutVerse).Error)

	got, err := repo.List(context.Background(), user.ID, usecase.ListCharacterFilters{VerseID: &verse.ID})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, withVerse.ID, got[0].ID)
	require.NotNil(t, got[0].VerseID)
	assert.Equal(t, verse.ID, *got[0].VerseID)
}

func seedCharacterRecency(t *testing.T, tx *gorm.DB, userID uuid.UUID, name string, createdAt time.Time, lastUsedAt *time.Time) *domain.Character {
	t.Helper()
	c := &domain.Character{ID: uuid.New(), UserID: userID, Name: name, CreatedAt: createdAt, LastUsedAt: lastUsedAt}
	require.NoError(t, tx.Create(c).Error)
	return c
}

func charNames(chars []*domain.Character) []string {
	names := make([]string, len(chars))
	for i, c := range chars {
		names[i] = c.Name
	}
	return names
}

func TestCharacterRepository_List_SortRecent(t *testing.T) {
	recent := "recent"
	base := time.Now().UTC().Truncate(time.Second)
	ts := func(d time.Duration) *time.Time { v := base.Add(d); return &v }

	t.Run("orders by last_used_at desc, nulls last, name tiebreak", func(t *testing.T) {
		tx := testutil.NewTestTx(t, testDB)
		repo := NewCharacterRepository(tx)
		user := seedUser(t, tx, "user_1")
		// created order is the reverse of recency order, so the default order differs
		seedCharacterRecency(t, tx, user.ID, "NullUsed", base.Add(4*time.Hour), nil)
		seedCharacterRecency(t, tx, user.ID, "Old", base.Add(3*time.Hour), ts(-2*time.Hour))
		seedCharacterRecency(t, tx, user.ID, "TieB", base.Add(2*time.Hour), ts(-time.Hour))
		seedCharacterRecency(t, tx, user.ID, "TieA", base.Add(time.Hour), ts(-time.Hour))
		seedCharacterRecency(t, tx, user.ID, "Newest", base, ts(0))

		got, err := repo.List(context.Background(), user.ID, usecase.ListCharacterFilters{Sort: &recent})

		require.NoError(t, err)
		assert.Equal(t, []string{"Newest", "TieA", "TieB", "Old", "NullUsed"}, charNames(got))
	})

	t.Run("default order ignores last_used_at", func(t *testing.T) {
		tx := testutil.NewTestTx(t, testDB)
		repo := NewCharacterRepository(tx)
		user := seedUser(t, tx, "user_1")
		seedCharacterRecency(t, tx, user.ID, "CreatedFirst", base, ts(time.Hour))
		seedCharacterRecency(t, tx, user.ID, "CreatedLast", base.Add(time.Hour), ts(-time.Hour))

		got, err := repo.List(context.Background(), user.ID, usecase.ListCharacterFilters{})

		require.NoError(t, err)
		assert.Equal(t, []string{"CreatedLast", "CreatedFirst"}, charNames(got))
	})

	t.Run("combines with q and verse_id filters", func(t *testing.T) {
		tx := testutil.NewTestTx(t, testDB)
		repo := NewCharacterRepository(tx)
		user := seedUser(t, tx, "user_1")
		verse := &domain.Verse{ID: uuid.New(), UserID: user.ID, Name: "Verse"}
		require.NoError(t, tx.Create(verse).Error)

		mk := func(name string, verseID *uuid.UUID, lastUsed time.Duration) {
			c := &domain.Character{ID: uuid.New(), UserID: user.ID, Name: name, VerseID: verseID, LastUsedAt: ts(lastUsed)}
			require.NoError(t, tx.Create(c).Error)
		}
		mk("Aria Old", &verse.ID, -time.Hour)
		mk("Aria New", &verse.ID, 0)
		mk("Aria NoVerse", nil, time.Hour)
		mk("Bob", &verse.ID, time.Hour)

		q := "aria"
		got, err := repo.List(context.Background(), user.ID, usecase.ListCharacterFilters{Sort: &recent, Q: &q, VerseID: &verse.ID})

		require.NoError(t, err)
		assert.Equal(t, []string{"Aria New", "Aria Old"}, charNames(got))
	})
}

func TestCharacterRepository_UpdateLastUsedAt(t *testing.T) {
	t.Run("bumps all given ids", func(t *testing.T) {
		tx := testutil.NewTestTx(t, testDB)
		repo := NewCharacterRepository(tx)
		user := seedUser(t, tx, "user_1")
		a := seedCharacter(t, tx, user.ID)
		b := seedCharacter(t, tx, user.ID)
		untouched := seedCharacter(t, tx, user.ID)

		before := time.Now().Add(-time.Second)
		require.NoError(t, repo.UpdateLastUsedAt(context.Background(), []uuid.UUID{a.ID, b.ID}, user.ID))

		for _, id := range []uuid.UUID{a.ID, b.ID} {
			got, err := repo.GetByID(context.Background(), id, user.ID)
			require.NoError(t, err)
			require.NotNil(t, got.LastUsedAt)
			assert.True(t, got.LastUsedAt.After(before))
		}
		got, err := repo.GetByID(context.Background(), untouched.ID, user.ID)
		require.NoError(t, err)
		assert.Nil(t, got.LastUsedAt)
	})

	t.Run("does not bump another user's characters", func(t *testing.T) {
		tx := testutil.NewTestTx(t, testDB)
		repo := NewCharacterRepository(tx)
		owner := seedUser(t, tx, "user_1")
		other := seedUser(t, tx, "user_2")
		c := seedCharacter(t, tx, owner.ID)

		require.NoError(t, repo.UpdateLastUsedAt(context.Background(), []uuid.UUID{c.ID}, other.ID))

		got, err := repo.GetByID(context.Background(), c.ID, owner.ID)
		require.NoError(t, err)
		assert.Nil(t, got.LastUsedAt)
	})
}
