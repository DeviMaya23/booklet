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

func TestUserRepository_GetOrCreate_CreatesWithUUIDAndIDPSubject(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	user, err := repo.GetOrCreate(context.Background(), "kp_abc123")

	require.NoError(t, err)
	assert.Equal(t, "kp_abc123", user.IDPSubject)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestUserRepository_GetOrCreate_IdempotentOnSameIDPSubject(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	first, err := repo.GetOrCreate(context.Background(), "kp_idempotent")
	require.NoError(t, err)

	second, err := repo.GetOrCreate(context.Background(), "kp_idempotent")
	require.NoError(t, err)

	assert.Equal(t, first.ID, second.ID)
}

func TestUserRepository_GetByIDPSubject_ReturnsUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	created, err := repo.GetOrCreate(context.Background(), "kp_lookup")
	require.NoError(t, err)

	got, err := repo.GetByIDPSubject(context.Background(), "kp_lookup")
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "kp_lookup", got.IDPSubject)
}

func TestUserRepository_DeleteAllUserData_TombstonesUserRow(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	user := seedUser(t, tx, "kp_tombstone_test")

	_, err := repo.DeleteAllUserData(context.Background(), user.ID)

	require.NoError(t, err)

	var got domain.User
	require.NoError(t, tx.Where("id = ?", user.ID).First(&got).Error)
	assert.Equal(t, domain.AccountStatePurged, got.AccountState)
	assert.NotNil(t, got.PurgedAt)
}

func TestUserRepository_DeleteAllUserData_NotFound_ReturnsError(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	_, err := repo.DeleteAllUserData(context.Background(), uuid.New())

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_DeleteExpiredTombstones_RemovesOldRows(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewUserRepository(tx)

	oldPurgedAt := time.Now().Add(-25 * time.Hour)
	oldUser := &domain.User{ID: uuid.New(), IDPSubject: "kp_old_tombstone", AccountState: domain.AccountStatePurged, PurgedAt: &oldPurgedAt}
	require.NoError(t, tx.Create(oldUser).Error)

	recentPurgedAt := time.Now().Add(-1 * time.Hour)
	recentUser := &domain.User{ID: uuid.New(), IDPSubject: "kp_recent_tombstone", AccountState: domain.AccountStatePurged, PurgedAt: &recentPurgedAt}
	require.NoError(t, tx.Create(recentUser).Error)

	err := repo.DeleteExpiredTombstones(context.Background())

	require.NoError(t, err)

	var gone domain.User
	assert.ErrorIs(t, tx.Unscoped().Where("id = ?", oldUser.ID).First(&gone).Error, gorm.ErrRecordNotFound)

	var kept domain.User
	assert.NoError(t, tx.Unscoped().Where("id = ?", recentUser.ID).First(&kept).Error)
}
