package repository

import (
	"context"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository_GetOrCreate_Success(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	user, err := repo.GetOrCreate(context.Background(), "kp_abc123")

	require.NoError(t, err)
	assert.Equal(t, "kp_abc123", user.ID)
}

func TestUserRepository_DeleteAllUserData_TombstonesUserRow(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)
	userID := "kp_tombstone_test"
	seedUser(t, tx, userID)

	_, err := repo.DeleteAllUserData(context.Background(), userID)

	require.NoError(t, err)

	var user domain.User
	require.NoError(t, tx.Where("id = ?", userID).First(&user).Error)
	assert.Equal(t, domain.AccountStatePurged, user.AccountState)
	assert.NotNil(t, user.PurgedAt)
}

func TestUserRepository_DeleteAllUserData_NotFound_ReturnsError(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	_, err := repo.DeleteAllUserData(context.Background(), "nonexistent-user")

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_DeleteExpiredTombstones_RemovesOldRows(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	oldPurgedAt := time.Now().Add(-25 * time.Hour)
	oldUser := &domain.User{ID: "kp_old_tombstone", AccountState: domain.AccountStatePurged, PurgedAt: &oldPurgedAt}
	require.NoError(t, tx.Create(oldUser).Error)

	recentPurgedAt := time.Now().Add(-1 * time.Hour)
	recentUser := &domain.User{ID: "kp_recent_tombstone", AccountState: domain.AccountStatePurged, PurgedAt: &recentPurgedAt}
	require.NoError(t, tx.Create(recentUser).Error)

	err := repo.DeleteExpiredTombstones(context.Background())

	require.NoError(t, err)

	var gone domain.User
	assert.ErrorIs(t, tx.Unscoped().Where("id = ?", oldUser.ID).First(&gone).Error, gorm.ErrRecordNotFound)

	var kept domain.User
	assert.NoError(t, tx.Unscoped().Where("id = ?", recentUser.ID).First(&kept).Error)
}
