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

	pendingID := uuid.New()
	transactor := NewGormTransactor(testDB)
	uploadRepo := NewUploadRepository(testDB)

	err := transactor.InTransaction(ctx, func(txCtx context.Context) error {
		_, createErr := uploadRepo.Create(txCtx, &domain.PendingUpload{
			ID:       pendingID,
			UserID:   user.ID,
			R2Key:    "users/" + user.ID.String() + "/images/test.jpg",
			MimeType: "image/jpeg",
		})
		if createErr != nil {
			return createErr
		}
		return errors.New("intentional failure")
	})

	require.Error(t, err)

	var count int64
	testDB.Model(&domain.PendingUpload{}).Where("id = ?", pendingID).Count(&count)
	assert.Equal(t, int64(0), count, "pending_upload row must not exist after rollback")
}
