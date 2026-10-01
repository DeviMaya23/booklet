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

func seedArtpiece(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.Artpiece {
	t.Helper()
	a := &domain.Artpiece{
		ID:     uuid.New(),
		UserID: userID,
	}
	require.NoError(t, tx.Create(a).Error)
	return a
}

func TestArtpieceRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpiece(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), artpiece.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, artpiece.ID, got.ID)
	assert.Empty(t, got.Characters)
	assert.Nil(t, got.CoverFile)
}

func TestArtpieceRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	artpiece := seedArtpiece(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), artpiece.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtpieceRepository_List_FilterByArtist(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Picasso")

	a1 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, ArtistID: &artist.ID}
	a2 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID}
	require.NoError(t, tx.Create(a1).Error)
	require.NoError(t, tx.Create(a2).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		ArtistIDs: []uuid.UUID{artist.ID},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, a1.ID, results[0].ID)
}

func TestArtpieceRepository_List_FilterByCharacter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	char := seedCharacter(t, tx, user.ID)

	a1 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, Characters: []domain.Character{{ID: char.ID}}}
	a2 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID}
	require.NoError(t, tx.Create(a1).Error)
	require.NoError(t, tx.Create(a2).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		CharacterIDs: []uuid.UUID{char.ID},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, a1.ID, results[0].ID)
}

func TestArtpieceRepository_Delete_DoesNotCascadeToFiles(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	artpieceRepo := NewArtpieceRepository(tx)
	fileRepo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpiece(t, tx, user.ID)

	f := &domain.File{
		ID:         uuid.New(),
		UserID:     user.ID,
		ArtpieceID: &artpiece.ID,
		FileR2Path: "users/x/files/abc.jpg",
		MimeType:   "image/jpeg",
	}
	created, err := fileRepo.Create(context.Background(), f)
	require.NoError(t, err)

	err = artpieceRepo.Delete(context.Background(), artpiece.ID, user.ID)
	require.NoError(t, err)

	// file should still exist with artpiece_id NULL
	got, err := fileRepo.GetByIDAndUserID(context.Background(), created.ID, user.ID)
	require.NoError(t, err)
	assert.Nil(t, got.ArtpieceID, "file should have artpiece_id set to NULL after artpiece deletion")
}
