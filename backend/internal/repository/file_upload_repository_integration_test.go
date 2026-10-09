package repository

import (
	"context"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedPendingFileUpload(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.PendingFileUpload {
	t.Helper()
	p := &domain.PendingFileUpload{
		ID:       uuid.New(),
		UserID:   userID,
		R2Key:    "uploads/" + uuid.New().String(),
		MimeType: "image/png",
	}
	require.NoError(t, tx.Create(p).Error)
	return p
}

func TestPendingFileUploadRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewPendingFileUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	name := "artwork.png"
	p := &domain.PendingFileUpload{
		ID:       uuid.New(),
		UserID:   user.ID,
		R2Key:    "uploads/artwork.png",
		MimeType: "image/png",
		Name:     &name,
	}

	got, err := repo.Create(context.Background(), p)

	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, "uploads/artwork.png", got.R2Key)
	assert.Equal(t, "image/png", got.MimeType)
	require.NotNil(t, got.Name)
	assert.Equal(t, "artwork.png", *got.Name)
}

func TestPendingFileUploadRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewPendingFileUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingFileUpload(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), p.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)
}

func TestPendingFileUploadRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewPendingFileUploadRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	p := seedPendingFileUpload(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), p.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestPendingFileUploadRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewPendingFileUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingFileUpload(t, tx, user.ID)

	err := repo.Delete(context.Background(), p.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), p.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestPendingFileUploadRepository_ListStale(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewPendingFileUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedPendingFileUpload(t, tx, user.ID)

	future := time.Now().Add(time.Hour)
	got, err := repo.ListStale(context.Background(), future)

	require.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestPendingFileUploadRepository_ListStale_ExcludesRecent(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewPendingFileUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedPendingFileUpload(t, tx, user.ID)

	past := time.Now().Add(-time.Hour)
	got, err := repo.ListStale(context.Background(), past)

	require.NoError(t, err)
	assert.Empty(t, got)
}
