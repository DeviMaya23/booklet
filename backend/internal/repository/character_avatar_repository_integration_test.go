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

func seedPendingCharacterAvatarUpload(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.PendingCharacterAvatarUpload {
	t.Helper()
	char := seedCharacter(t, tx, userID)
	p := &domain.PendingCharacterAvatarUpload{
		ID:          uuid.New(),
		UserID:      userID,
		CharacterID: char.ID,
		R2Key:       "uploads/" + uuid.New().String(),
		MimeType:    "image/png",
	}
	require.NoError(t, tx.Create(p).Error)
	return p
}

func TestCharacterAvatarRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterAvatarRepository(tx)

	user := seedUser(t, tx, "user_1")
	char := seedCharacter(t, tx, user.ID)
	p := &domain.PendingCharacterAvatarUpload{
		ID:          uuid.New(),
		UserID:      user.ID,
		CharacterID: char.ID,
		R2Key:       "uploads/test-key",
		MimeType:    "image/jpeg",
	}

	got, err := repo.Create(context.Background(), p)

	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)
	assert.Equal(t, "uploads/test-key", got.R2Key)
	assert.Equal(t, "image/jpeg", got.MimeType)
	assert.Equal(t, char.ID, got.CharacterID)
}

func TestCharacterAvatarRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterAvatarRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingCharacterAvatarUpload(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), p.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)
}

func TestCharacterAvatarRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterAvatarRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	p := seedPendingCharacterAvatarUpload(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), p.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterAvatarRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterAvatarRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingCharacterAvatarUpload(t, tx, user.ID)

	err := repo.Delete(context.Background(), p.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), p.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterAvatarRepository_ListStale(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterAvatarRepository(tx)

	user := seedUser(t, tx, "user_1")
	p := seedPendingCharacterAvatarUpload(t, tx, user.ID)
	_ = p

	future := time.Now().Add(time.Hour)
	got, err := repo.ListStale(context.Background(), future)

	require.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestCharacterAvatarRepository_ListStale_ExcludesRecent(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterAvatarRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedPendingCharacterAvatarUpload(t, tx, user.ID)

	past := time.Now().Add(-time.Hour)
	got, err := repo.ListStale(context.Background(), past)

	require.NoError(t, err)
	assert.Empty(t, got)
}
