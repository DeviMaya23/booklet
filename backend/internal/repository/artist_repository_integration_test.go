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

func seedArtist(t *testing.T, tx *gorm.DB, userID uuid.UUID, name string) *domain.Artist {
	t.Helper()
	a := &domain.Artist{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}
	require.NoError(t, tx.Create(a).Error)
	return a
}

func TestArtistRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := &domain.Artist{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Jane Doe",
		Links: []domain.ArtistLink{
			{ID: uuid.New(), URL: "https://example.com", IsPrimary: true},
		},
	}

	got, err := repo.Create(context.Background(), artist)

	require.NoError(t, err)
	assert.Equal(t, artist.ID, got.ID)
	assert.Equal(t, "Jane Doe", got.Name)
	require.Len(t, got.Links, 1)
	assert.Equal(t, "https://example.com", got.Links[0].URL)
	assert.True(t, got.Links[0].IsPrimary)
}

func TestArtistRepository_Update_ReplacesLinks(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	updated, err := repo.Update(context.Background(), artist.ID, user.ID, usecase.UpdateArtistParams{
		Name: "Jane Doe",
		Links: []usecase.ArtistLinkInput{
			{URL: "https://first.example.com", IsPrimary: true},
			{URL: "https://second.example.com", IsPrimary: false},
		},
	})

	require.NoError(t, err)
	require.Len(t, updated.Links, 2)
	primaryCount := 0
	for _, l := range updated.Links {
		if l.IsPrimary {
			primaryCount++
		}
	}
	assert.Equal(t, 1, primaryCount)
}

func TestArtistRepository_Create_DuplicateNameSameUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedArtist(t, tx, user.ID, "Jane Doe")

	_, err := repo.Create(context.Background(), &domain.Artist{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Jane Doe",
	})

	assert.ErrorIs(t, err, usecase.ErrArtistNameConflict)
}

func TestArtistRepository_Create_SameNameDifferentUsers(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedArtist(t, tx, user1.ID, "Jane Doe")

	_, err := repo.Create(context.Background(), &domain.Artist{
		ID:     uuid.New(),
		UserID: user2.ID,
		Name:   "Jane Doe",
	})

	assert.NoError(t, err)
}

func TestArtistRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	got, err := repo.GetByID(context.Background(), artist.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, artist.ID, got.ID)
	assert.Equal(t, "Jane Doe", got.Name)
}

func TestArtistRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	artist := seedArtist(t, tx, user1.ID, "Jane Doe")

	_, err := repo.GetByID(context.Background(), artist.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtistRepository_List(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedArtist(t, tx, user1.ID, "Zara")
	seedArtist(t, tx, user1.ID, "Alice")
	seedArtist(t, tx, user2.ID, "Other")

	got, err := repo.List(context.Background(), user1.ID, usecase.ListArtistFilters{})

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Alice", got[0].Name)
	assert.Equal(t, "Zara", got[1].Name)
}

func TestArtistRepository_List_QFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedArtist(t, tx, user.ID, "Jane Doe")
	seedArtist(t, tx, user.ID, "Bob Smith")

	q := "jane"
	got, err := repo.List(context.Background(), user.ID, usecase.ListArtistFilters{Q: &q})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Jane Doe", got[0].Name)
}

func TestArtistRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	got, err := repo.Update(context.Background(), artist.ID, user.ID, usecase.UpdateArtistParams{
		Name: "Jane Smith",
	})

	require.NoError(t, err)
	assert.Equal(t, "Jane Smith", got.Name)
}

func TestArtistRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	_, err := repo.Update(context.Background(), uuid.New(), user.ID, usecase.UpdateArtistParams{
		Name: "x",
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtistRepository_Update_DuplicateName(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedArtist(t, tx, user.ID, "Alice")
	bob := seedArtist(t, tx, user.ID, "Bob")
	_, err := repo.Update(context.Background(), bob.ID, user.ID, usecase.UpdateArtistParams{
		Name: "Alice",
	})

	assert.ErrorIs(t, err, usecase.ErrArtistNameConflict)
}

func TestArtistRepository_Update_NoFields_NoOp(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	got, err := repo.Update(context.Background(), artist.ID, user.ID, usecase.UpdateArtistParams{})

	require.NoError(t, err)
	assert.Equal(t, "Jane Doe", got.Name)
}

func TestArtistRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	err := repo.Delete(context.Background(), artist.ID, user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), artist.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtistRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")

	err := repo.Delete(context.Background(), uuid.New(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

