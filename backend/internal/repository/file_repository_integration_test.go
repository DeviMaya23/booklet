package repository

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedFile(t *testing.T, tx *gorm.DB, userID uuid.UUID, mimeType string) *domain.File {
	t.Helper()
	f := &domain.File{
		ID:         uuid.New(),
		UserID:     userID,
		FileR2Path: "users/" + userID.String() + "/files/test",
		MimeType:   mimeType,
	}
	require.NoError(t, tx.Create(f).Error)
	return f
}

func seedArtpieceForFile(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.Artpiece {
	t.Helper()
	a := &domain.Artpiece{
		ID:     uuid.New(),
		UserID: userID,
	}
	require.NoError(t, tx.Create(a).Error)
	return a
}

func TestFileRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := &domain.File{
		ID:         uuid.New(),
		UserID:     user.ID,
		FileR2Path: "users/" + user.ID.String() + "/files/abc.jpg",
		MimeType:   "image/jpeg",
	}

	created, err := repo.Create(context.Background(), f)

	require.NoError(t, err)
	assert.Equal(t, f.ID, created.ID)
	assert.Equal(t, "image/jpeg", created.MimeType)
}

func TestFileRepository_GetByIDAndUserID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	f := seedFile(t, tx, user1.ID, "image/jpeg")

	_, err := repo.GetByIDAndUserID(context.Background(), f.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileRepository_UpdateThumbnailPath(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := seedFile(t, tx, user.ID, "image/jpeg")

	err := repo.UpdateThumbnailPath(context.Background(), f.ID, "users/x/thumbnails/abc.jpg")
	require.NoError(t, err)

	got, err := repo.GetByIDAndUserID(context.Background(), f.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got.ThumbnailR2Path)
	assert.Equal(t, "users/x/thumbnails/abc.jpg", *got.ThumbnailR2Path)
}

func TestFileRepository_GetFilesForArtpiece(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpieceForFile(t, tx, user.ID)

	f1 := seedFile(t, tx, user.ID, "image/jpeg")
	f2 := seedFile(t, tx, user.ID, "video/mp4")
	artpieceID := artpiece.ID
	require.NoError(t, tx.Model(f1).Update("artpiece_id", artpieceID).Error)
	require.NoError(t, tx.Model(f2).Update("artpiece_id", artpieceID).Error)

	// file belonging to another artpiece — should not appear
	other := seedArtpieceForFile(t, tx, user.ID)
	f3 := seedFile(t, tx, user.ID, "image/png")
	require.NoError(t, tx.Model(f3).Update("artpiece_id", other.ID).Error)

	files, err := repo.GetFilesForArtpiece(context.Background(), artpieceID)
	require.NoError(t, err)
	assert.Len(t, files, 2)
}
