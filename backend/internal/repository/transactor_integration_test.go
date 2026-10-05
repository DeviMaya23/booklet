package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGormTransactor_RollsBackOnError(t *testing.T) {
	ctx := context.Background()

	tx := testutil.NewTestTx(t, testDB)
	user := seedUser(t, tx, "transactor_rollback_test_user")

	characterID := uuid.New()
	transactor := NewGormTransactor(testDB)
	characterRepo := NewCharacterRepository(testDB)

	err := transactor.InTransaction(ctx, func(txCtx context.Context) error {
		createErr := characterRepo.Create(txCtx, &domain.Character{
			ID:     characterID,
			UserID: user.ID,
			Name:   "Rollback Test Character",
		})
		if createErr != nil {
			return createErr
		}
		return errors.New("intentional failure")
	})

	require.Error(t, err)

	var count int64
	testDB.Model(&domain.Character{}).Where("id = ?", characterID).Count(&count)
	assert.Equal(t, int64(0), count, "character row must not exist after rollback")
}
