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

func seedImage(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.Image {
	t.Helper()
	img := &domain.Image{
		ID:          uuid.New(),
		UserID:      userID,
		ImageR2Path: "images/test.jpg",
		MimeType:    "image/jpeg",
	}
	require.NoError(t, tx.Create(img).Error)
	return img
}

func TestImageRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	img := seedImage(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), img.ID.String(), user.ID)

	require.NoError(t, err)
	assert.Equal(t, img.ID, got.ID)
	assert.Empty(t, got.Characters)
}

func TestImageRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	img := seedImage(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), img.ID.String(), user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestImageRepository_List(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedImage(t, tx, user1.ID)
	seedImage(t, tx, user1.ID)
	seedImage(t, tx, user2.ID)

	got, err := repo.List(context.Background(), user1.ID, usecase.ListImageFilters{})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, img := range got {
		assert.Equal(t, user1.ID, img.UserID)
	}
}

func TestImageRepository_List_QFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	sunset := seedImage(t, tx, user.ID)
	title := "Sunset at Sea"
	require.NoError(t, tx.Model(sunset).Update("title", title).Error)
	other := seedImage(t, tx, user.ID)
	otherTitle := "Mountain View"
	require.NoError(t, tx.Model(other).Update("title", otherTitle).Error)

	q := "sunset"
	got, err := repo.List(context.Background(), user.ID, usecase.ListImageFilters{Q: &q})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, sunset.ID, got[0].ID)
}

func TestImageRepository_List_CharacterIDsFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	charA := seedCharacter(t, tx, user.ID)
	charB := seedCharacter(t, tx, user.ID)
	img1 := seedImage(t, tx, user.ID)
	img2 := seedImage(t, tx, user.ID)
	// tag img1 with charA, img2 with charB
	require.NoError(t, tx.Model(img1).Association("Characters").Replace([]domain.Character{{ID: charA.ID}}))
	require.NoError(t, tx.Model(img2).Association("Characters").Replace([]domain.Character{{ID: charB.ID}}))

	got, err := repo.List(context.Background(), user.ID, usecase.ListImageFilters{CharacterIDs: []string{charA.ID.String()}})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, img1.ID, got[0].ID)
}

func TestImageRepository_List_CharacterIDsFilter_NoDuplicates(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	charA := seedCharacter(t, tx, user.ID)
	charB := seedCharacter(t, tx, user.ID)
	img := seedImage(t, tx, user.ID)
	// tag img with both characters
	require.NoError(t, tx.Model(img).Association("Characters").Replace([]domain.Character{{ID: charA.ID}, {ID: charB.ID}}))

	got, err := repo.List(context.Background(), user.ID, usecase.ListImageFilters{CharacterIDs: []string{charA.ID.String(), charB.ID.String()}})

	require.NoError(t, err)
	require.Len(t, got, 1, "image should appear once despite matching multiple character IDs")
	assert.Equal(t, img.ID, got[0].ID)
}

func TestImageRepository_List_ArtistIDsFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	img1 := seedImage(t, tx, user.ID)
	img2 := seedImage(t, tx, user.ID)
	require.NoError(t, tx.Model(img1).Update("artist_id", artist.ID).Error)

	got, err := repo.List(context.Background(), user.ID, usecase.ListImageFilters{ArtistIDs: []string{artist.ID.String()}})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, img1.ID, got[0].ID)
	_ = img2
}

func TestImageRepository_List_CombinedCharacterAndArtistFilter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	char := seedCharacter(t, tx, user.ID)
	artistA := seedArtist(t, tx, user.ID, "Artist A")
	artistB := seedArtist(t, tx, user.ID, "Artist B")

	// img1: tagged with char, under artistA — should match both filters
	img1 := seedImage(t, tx, user.ID)
	require.NoError(t, tx.Model(img1).Update("artist_id", artistA.ID).Error)
	require.NoError(t, tx.Model(img1).Association("Characters").Replace([]domain.Character{{ID: char.ID}}))

	// img2: tagged with char but under artistB — matches character filter only
	img2 := seedImage(t, tx, user.ID)
	require.NoError(t, tx.Model(img2).Update("artist_id", artistB.ID).Error)
	require.NoError(t, tx.Model(img2).Association("Characters").Replace([]domain.Character{{ID: char.ID}}))

	got, err := repo.List(context.Background(), user.ID, usecase.ListImageFilters{
		CharacterIDs: []string{char.ID.String()},
		ArtistIDs:    []string{artistA.ID.String()},
	})

	require.NoError(t, err)
	require.Len(t, got, 1, "only image matching both character and artist filter should be returned")
	assert.Equal(t, img1.ID, got[0].ID)
}

func TestImageRepository_Update_CharacterAssociation(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	img := seedImage(t, tx, user.ID)
	char := seedCharacter(t, tx, user.ID)

	charIDs := []string{char.ID.String()}
	got, err := repo.Update(context.Background(), img.ID.String(), user.ID, usecase.UpdateImageParams{
		CharacterIDs: charIDs,
	})

	require.NoError(t, err)
	require.Len(t, got.Characters, 1)
	assert.Equal(t, char.ID, got.Characters[0].ID)
}

func TestImageRepository_Update_CharacterNotOwned(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	img := seedImage(t, tx, user1.ID)
	char := seedCharacter(t, tx, user2.ID)

	charIDs := []string{char.ID.String()}
	_, err := repo.Update(context.Background(), img.ID.String(), user1.ID, usecase.UpdateImageParams{
		CharacterIDs: charIDs,
	})

	assert.ErrorIs(t, err, usecase.ErrCharacterNotOwned)
}

func TestImageRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	img := seedImage(t, tx, user.ID)

	err := repo.Delete(context.Background(), img.ID.String(), user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), img.ID.String(), user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestImageRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	err := repo.Delete(context.Background(), uuid.NewString(), user.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
