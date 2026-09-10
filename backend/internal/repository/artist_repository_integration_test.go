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
	link := "https://example.com"
	artist := &domain.Artist{
		ID:         uuid.New(),
		UserID:     user.ID,
		Name:       "Jane Doe",
		ArtistLink: &link,
	}

	got, err := repo.Create(context.Background(), artist)

	require.NoError(t, err)
	assert.Equal(t, artist.ID, got.ID)
	assert.Equal(t, "Jane Doe", got.Name)
	assert.Equal(t, &link, got.ArtistLink)

	var row domain.Artist
	require.NoError(t, tx.First(&row, "id = ?", artist.ID).Error)
	assert.Equal(t, "Jane Doe", row.Name)
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

	got, err := repo.GetByID(context.Background(), artist.ID.String(), user.ID)

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

	_, err := repo.GetByID(context.Background(), artist.ID.String(), user2.ID)

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
	got, err := repo.Update(context.Background(), artist.ID.String(), user.ID, usecase.UpdateArtistParams{
		Name: "Jane Smith",
	})

	require.NoError(t, err)
	assert.Equal(t, "Jane Smith", got.Name)
}

func TestArtistRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	_, err := repo.Update(context.Background(), uuid.NewString(), user.ID, usecase.UpdateArtistParams{
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
	_, err := repo.Update(context.Background(), bob.ID.String(), user.ID, usecase.UpdateArtistParams{
		Name: "Alice",
	})

	assert.ErrorIs(t, err, usecase.ErrArtistNameConflict)
}

func TestArtistRepository_Update_NoFields_NoOp(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	got, err := repo.Update(context.Background(), artist.ID.String(), user.ID, usecase.UpdateArtistParams{})

	require.NoError(t, err)
	assert.Equal(t, "Jane Doe", got.Name)
}

func TestArtistRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	err := repo.Delete(context.Background(), artist.ID.String(), user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), artist.ID.String(), user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtistRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")

	err := repo.Delete(context.Background(), uuid.NewString(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtistRepository_Delete_NullsImageArtistID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	img := seedImage(t, tx, user.ID)
	require.NoError(t, tx.Model(img).Update("artist_id", artist.ID).Error)

	err := repo.Delete(context.Background(), artist.ID.String(), user.ID)
	require.NoError(t, err)

	var row domain.Image
	require.NoError(t, tx.First(&row, "id = ?", img.ID).Error)
	assert.Nil(t, row.ArtistID)
}

func TestArtistRepository_Delete_NullsPendingUploadArtistID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtistRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	pending := seedPendingUpload(t, tx, user.ID)
	require.NoError(t, tx.Model(pending).Update("artist_id", artist.ID).Error)

	err := repo.Delete(context.Background(), artist.ID.String(), user.ID)
	require.NoError(t, err)

	var row domain.PendingUpload
	require.NoError(t, tx.First(&row, "id = ?", pending.ID).Error)
	assert.Nil(t, row.ArtistID)
}
