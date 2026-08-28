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

	got, err := repo.List(context.Background(), user1.ID)

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, img := range got {
		assert.Equal(t, user1.ID, img.UserID)
	}
}

func TestImageRepository_Update_CharacterAssociation(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewImageRepository(tx)

	user := seedUser(t, tx, "user_1")
	img := seedImage(t, tx, user.ID)
	char := seedCharacter(t, tx, user.ID)

	charIDs := []string{char.ID.String()}
	got, err := repo.Update(context.Background(), img.ID.String(), user.ID, usecase.UpdateImageParams{
		CharacterIDs: &charIDs,
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
		CharacterIDs: &charIDs,
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
