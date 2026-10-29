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

func seedVerse(t *testing.T, tx *gorm.DB, userID uuid.UUID, name string) *domain.Verse {
	t.Helper()
	v := &domain.Verse{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}
	require.NoError(t, tx.Create(v).Error)
	return v
}

func TestVerseRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := &domain.Verse{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Solaris",
	}

	got, err := repo.Create(context.Background(), verse)

	require.NoError(t, err)
	assert.Equal(t, verse.ID, got.ID)
	assert.Equal(t, "Solaris", got.Name)
	assert.Equal(t, user.ID, got.UserID)
}

func TestVerseRepository_Create_DuplicateNameSameUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedVerse(t, tx, user.ID, "Solaris")

	_, err := repo.Create(context.Background(), &domain.Verse{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Solaris",
	})

	assert.ErrorIs(t, err, usecase.ErrVerseNameConflict)
}

func TestVerseRepository_Create_SameNameDifferentUsers(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedVerse(t, tx, user1.ID, "Solaris")

	_, err := repo.Create(context.Background(), &domain.Verse{
		ID:     uuid.New(),
		UserID: user2.ID,
		Name:   "Solaris",
	})

	assert.NoError(t, err)
}

func TestVerseRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := seedVerse(t, tx, user.ID, "Solaris")

	got, err := repo.GetByID(context.Background(), verse.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, verse.ID, got.ID)
	assert.Equal(t, "Solaris", got.Name)
}

func TestVerseRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	verse := seedVerse(t, tx, user1.ID, "Solaris")

	_, err := repo.GetByID(context.Background(), verse.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestVerseRepository_List_OrderedByRecencyThenName(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")

	older := time.Now().Add(-2 * time.Hour)
	newer := time.Now().Add(-1 * time.Hour)

	v1 := &domain.Verse{ID: uuid.New(), UserID: user.ID, Name: "Zeta", LastUsedAt: &newer}
	v2 := &domain.Verse{ID: uuid.New(), UserID: user.ID, Name: "Alpha", LastUsedAt: &older}
	v3 := &domain.Verse{ID: uuid.New(), UserID: user.ID, Name: "Beta"}
	require.NoError(t, tx.Create(v1).Error)
	require.NoError(t, tx.Create(v2).Error)
	require.NoError(t, tx.Create(v3).Error)

	got, err := repo.List(context.Background(), user.ID, usecase.ListVerseFilters{})

	require.NoError(t, err)
	require.Len(t, got, 3)
	assert.Equal(t, "Zeta", got[0].Name)  // most recently used
	assert.Equal(t, "Alpha", got[1].Name) // older last_used_at
	assert.Equal(t, "Beta", got[2].Name)  // null last_used_at sorts last
}

func TestVerseRepository_List_IsolatesByUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedVerse(t, tx, user1.ID, "Solaris")
	seedVerse(t, tx, user2.ID, "Other")

	got, err := repo.List(context.Background(), user1.ID, usecase.ListVerseFilters{})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Solaris", got[0].Name)
}

func TestVerseRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := seedVerse(t, tx, user.ID, "Solaris")

	got, err := repo.Update(context.Background(), verse.ID, user.ID, usecase.UpdateVerseParams{Name: "Terran"})

	require.NoError(t, err)
	assert.Equal(t, "Terran", got.Name)
}

func TestVerseRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	_, err := repo.Update(context.Background(), uuid.New(), user.ID, usecase.UpdateVerseParams{Name: "x"})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestVerseRepository_Update_DuplicateName(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedVerse(t, tx, user.ID, "Solaris")
	v2 := seedVerse(t, tx, user.ID, "Terran")

	_, err := repo.Update(context.Background(), v2.ID, user.ID, usecase.UpdateVerseParams{Name: "Solaris"})

	assert.ErrorIs(t, err, usecase.ErrVerseNameConflict)
}

func TestVerseRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := seedVerse(t, tx, user.ID, "Solaris")

	err := repo.Delete(context.Background(), verse.ID, user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), verse.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestVerseRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")

	err := repo.Delete(context.Background(), uuid.New(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestVerseRepository_Delete_SetsVerseIDNullOnCharacters(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	verseRepo := NewVerseRepository(tx)
	charRepo := NewCharacterRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := seedVerse(t, tx, user.ID, "Solaris")

	verseID := verse.ID
	char := &domain.Character{
		ID:      uuid.New(),
		UserID:  user.ID,
		Name:    "Aria",
		VerseID: &verseID,
	}
	require.NoError(t, charRepo.Create(context.Background(), char))

	err := verseRepo.Delete(context.Background(), verse.ID, user.ID)
	require.NoError(t, err)

	got, err := charRepo.GetByID(context.Background(), char.ID, user.ID)
	require.NoError(t, err)
	assert.Nil(t, got.VerseID)
}

func TestVerseRepository_UpdateLastUsedAt(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewVerseRepository(tx)

	user := seedUser(t, tx, "user_1")
	verse := seedVerse(t, tx, user.ID, "Solaris")
	assert.Nil(t, verse.LastUsedAt)

	err := repo.UpdateLastUsedAt(context.Background(), verse.ID, user.ID)

	require.NoError(t, err)
	got, err := repo.GetByID(context.Background(), verse.ID, user.ID)
	require.NoError(t, err)
	assert.NotNil(t, got.LastUsedAt)
}
