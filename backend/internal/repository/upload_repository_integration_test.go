package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedPendingUpload(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.PendingUpload {
	t.Helper()
	p := &domain.PendingUpload{
		ID:           uuid.New(),
		UserID:       userID,
		R2Key:        fmt.Sprintf("users/%s/images/test.jpg", userID.String()),
		MimeType:     "image/jpeg",
		CharacterIDs: []uuid.UUID{},
	}
	require.NoError(t, tx.Create(p).Error)
	return p
}

func TestUploadRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := &domain.PendingUpload{
		ID:           uuid.New(),
		UserID:       user.ID,
		R2Key:        fmt.Sprintf("users/%s/images/abc.jpg", user.ID.String()),
		MimeType:     "image/jpeg",
		CharacterIDs: []uuid.UUID{},
	}

	got, err := repo.Create(context.Background(), p)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, user.ID, got.UserID)
	assert.Equal(t, p.R2Key, got.R2Key)
	assert.Equal(t, "image/jpeg", got.MimeType)

	var fromDB domain.PendingUpload
	require.NoError(t, tx.First(&fromDB, "id = ?", p.ID).Error)
	assert.Equal(t, p.ID, fromDB.ID)
}

func TestUploadRepository_GetByID_Success(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingUpload(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), p.ID, user.ID)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, p.R2Key, got.R2Key)
}

func TestUploadRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	p := seedPendingUpload(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), p.ID, user2.ID)

	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got %v", err)
}

func TestUploadRepository_GetByID_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	_, err := repo.GetByID(context.Background(), uuid.New(), user.ID)

	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got %v", err)
}

func TestUploadRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingUpload(t, tx, user.ID)

	err := repo.Delete(context.Background(), p.ID)

	require.NoError(t, err)
	var fromDB domain.PendingUpload
	assert.ErrorIs(t, tx.First(&fromDB, "id = ?", p.ID).Error, gorm.ErrRecordNotFound)
}

func TestUploadRepository_ListStale_ReturnsOlderThanCutoff(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	stale := &domain.PendingUpload{
		ID:           uuid.New(),
		UserID:       user.ID,
		R2Key:        fmt.Sprintf("users/%s/images/stale.jpg", user.ID.String()),
		MimeType:     "image/jpeg",
		CharacterIDs: []uuid.UUID{},
	}
	require.NoError(t, tx.Create(stale).Error)
	require.NoError(t, tx.Model(stale).Update("created_at", time.Now().Add(-2*time.Hour)).Error)

	cutoff := time.Now().Add(-time.Hour)
	got, err := repo.ListStale(context.Background(), cutoff)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, stale.ID, got[0].ID)
}

func TestUploadRepository_ListStale_ExcludesNewerThanCutoff(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedPendingUpload(t, tx, user.ID)

	cutoff := time.Now().Add(-time.Hour)
	got, err := repo.ListStale(context.Background(), cutoff)

	require.NoError(t, err)
	assert.Empty(t, got)
}
