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

func TestUploadRepository_CreatePendingUpload(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	seedUser(t, tx, "user_1")
	p := &domain.PendingUpload{
		ID:       uuid.New(),
		UserID:   "user_1",
		R2Key:    "users/user_1/images/abc.jpg",
		MimeType: "image/jpeg",
	}

	err := repo.CreatePendingUpload(context.Background(), p)

	require.NoError(t, err)
	var got domain.PendingUpload
	require.NoError(t, tx.First(&got, "id = ?", p.ID).Error)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, "user_1", got.UserID)
	assert.Equal(t, "users/user_1/images/abc.jpg", got.R2Key)
	assert.Equal(t, "image/jpeg", got.MimeType)
}

func TestUploadRepository_CompleteUpload_Success(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	char := seedCharacter(t, tx, "user_1")
	p := seedPendingUpload(t, tx, "user_1")

	image, err := repo.CompleteUpload(context.Background(), p.ID.String(), "user_1", []uuid.UUID{char.ID})

	require.NoError(t, err)
	require.NotNil(t, image)

	// pending_upload should be deleted
	var pending domain.PendingUpload
	err = tx.First(&pending, "id = ?", p.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	// image should exist with correct fields
	assert.Equal(t, "user_1", image.UserID)
	assert.Equal(t, p.R2Key, image.ImageR2Path)
	assert.Equal(t, "image/jpeg", image.MimeType)

	// character should be associated
	require.Len(t, image.Characters, 1)
	assert.Equal(t, char.ID, image.Characters[0].ID)
}

func TestUploadRepository_CompleteUpload_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	p := seedPendingUpload(t, tx, "user_1")

	_, err := repo.CompleteUpload(context.Background(), p.ID.String(), "user_2", []uuid.UUID{})

	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got %v", err)
}

func TestUploadRepository_CompleteUpload_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUploadRepository(tx)

	_, err := repo.CompleteUpload(context.Background(), uuid.NewString(), "user_1", []uuid.UUID{})

	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got %v", err)
}
