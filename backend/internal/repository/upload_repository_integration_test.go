package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedPendingUpload(t *testing.T, tx *gorm.DB, userID string) *domain.PendingUpload {
	t.Helper()
	seedUser(t, tx, userID)
	p := &domain.PendingUpload{
		ID:       uuid.New(),
		UserID:   userID,
		R2Key:    fmt.Sprintf("users/%s/images/test.jpg", userID),
		MimeType: "image/jpeg",
	}
	require.NoError(t, tx.Create(p).Error)
	return p
}

func TestUploadRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	seedUser(t, tx, "user_1")
	p := &domain.PendingUpload{
		ID:       uuid.New(),
		UserID:   "user_1",
		R2Key:    "users/user_1/images/abc.jpg",
		MimeType: "image/jpeg",
	}

	got, err := repo.Create(context.Background(), p)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, "user_1", got.UserID)
	assert.Equal(t, "users/user_1/images/abc.jpg", got.R2Key)
	assert.Equal(t, "image/jpeg", got.MimeType)

	var fromDB domain.PendingUpload
	require.NoError(t, tx.First(&fromDB, "id = ?", p.ID).Error)
	assert.Equal(t, p.ID, fromDB.ID)
}

func TestUploadRepository_GetByID_Success(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	p := seedPendingUpload(t, tx, "user_1")

	got, err := repo.GetByID(context.Background(), p.ID, "user_1")

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, p.R2Key, got.R2Key)
}

func TestUploadRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	p := seedPendingUpload(t, tx, "user_1")

	_, err := repo.GetByID(context.Background(), p.ID, "user_2")

	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got %v", err)
}

func TestUploadRepository_GetByID_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	_, err := repo.GetByID(context.Background(), uuid.New(), "user_1")

	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got %v", err)
}

func TestUploadRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	p := seedPendingUpload(t, tx, "user_1")

	err := repo.Delete(context.Background(), p.ID)

	require.NoError(t, err)
	var fromDB domain.PendingUpload
	assert.ErrorIs(t, tx.First(&fromDB, "id = ?", p.ID).Error, gorm.ErrRecordNotFound)
}
