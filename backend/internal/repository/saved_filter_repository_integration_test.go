package repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedSavedFilter(t *testing.T, tx *gorm.DB, userID uuid.UUID, name string) *domain.SavedFilter {
	t.Helper()
	sf := &domain.SavedFilter{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          name,
		FilterPayload: json.RawMessage(`{"artists":[]}`),
	}
	require.NoError(t, tx.Create(sf).Error)
	return sf
}

func TestSavedFilterRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")
	payload := json.RawMessage(`{"artists":["abc"]}`)
	sf := &domain.SavedFilter{
		ID:            uuid.New(),
		UserID:        user.ID,
		Name:          "My Filter",
		FilterPayload: payload,
	}

	got, err := repo.Create(context.Background(), sf)

	require.NoError(t, err)
	assert.Equal(t, sf.ID, got.ID)
	assert.Equal(t, "My Filter", got.Name)
	assert.JSONEq(t, `{"artists":["abc"]}`, string(got.FilterPayload))
}

func TestSavedFilterRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")
	sf := seedSavedFilter(t, tx, user.ID, "My Filter")

	got, err := repo.GetByID(context.Background(), sf.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, sf.ID, got.ID)
	assert.Equal(t, "My Filter", got.Name)
}

func TestSavedFilterRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	sf := seedSavedFilter(t, tx, user1.ID, "My Filter")

	_, err := repo.GetByID(context.Background(), sf.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestSavedFilterRepository_List_OrderedByCreatedAtDesc(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")
	first := seedSavedFilter(t, tx, user.ID, "First")
	second := seedSavedFilter(t, tx, user.ID, "Second")

	got, err := repo.List(context.Background(), user.ID)

	require.NoError(t, err)
	require.Len(t, got, 2)
	// created_at DESC — second inserted last so it comes first
	assert.Equal(t, second.ID, got[0].ID)
	assert.Equal(t, first.ID, got[1].ID)
}

func TestSavedFilterRepository_List_Isolation(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedSavedFilter(t, tx, user1.ID, "User1 Filter")
	seedSavedFilter(t, tx, user2.ID, "User2 Filter")

	got, err := repo.List(context.Background(), user1.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "User1 Filter", got[0].Name)
}

func TestSavedFilterRepository_Update_Name(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")
	sf := seedSavedFilter(t, tx, user.ID, "Old Name")
	newName := "New Name"

	got, err := repo.Update(context.Background(), sf.ID, user.ID, usecase.UpdateSavedFilterParams{
		Name: usecase.Patch[string]{Set: true, Value: &newName},
	})

	require.NoError(t, err)
	assert.Equal(t, "New Name", got.Name)
}

func TestSavedFilterRepository_Update_ClearThumbnail(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")
	path := "r2/thumb.jpg"
	sf := &domain.SavedFilter{
		ID:              uuid.New(),
		UserID:          user.ID,
		Name:            "Filter",
		ThumbnailR2Path: &path,
		FilterPayload:   json.RawMessage(`{}`),
	}
	require.NoError(t, tx.Create(sf).Error)

	got, err := repo.Update(context.Background(), sf.ID, user.ID, usecase.UpdateSavedFilterParams{
		ThumbnailR2Path: usecase.Patch[string]{Set: true, Value: nil},
	})

	require.NoError(t, err)
	assert.Nil(t, got.ThumbnailR2Path)
}

func TestSavedFilterRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")

	_, err := repo.Update(context.Background(), uuid.New(), user.ID, usecase.UpdateSavedFilterParams{})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestSavedFilterRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")
	sf := seedSavedFilter(t, tx, user.ID, "My Filter")

	err := repo.Delete(context.Background(), sf.ID, user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), sf.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestSavedFilterRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewSavedFilterRepository(tx)

	user := seedUser(t, tx, "user_1")

	err := repo.Delete(context.Background(), uuid.New(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
