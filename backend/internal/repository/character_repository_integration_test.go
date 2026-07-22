package repository

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedUser(t *testing.T, tx *gorm.DB, userID string) {
	t.Helper()
	require.NoError(t, tx.FirstOrCreate(&domain.User{ID: userID}).Error)
}

func seedCharacter(t *testing.T, tx *gorm.DB, userID string) *domain.Character {
	t.Helper()
	seedUser(t, tx, userID)
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: userID,
		Name:   "Test Character",
	}
	require.NoError(t, tx.Create(c).Error)
	return c
}

func TestCharacterRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	seedUser(t, tx, "user_1")
	c := &domain.Character{
		ID:     uuid.New(),
		UserID: "user_1",
		Name:   "Aria",
	}
	err := repo.Create(context.Background(), c)

	require.NoError(t, err)

	var got domain.Character
	require.NoError(t, tx.First(&got, "id = ?", c.ID).Error)
	assert.Equal(t, "Aria", got.Name)
	assert.Equal(t, "user_1", got.UserID)
}

func TestCharacterRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")

	got, err := repo.GetByID(context.Background(), c.ID.String(), "user_1")

	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
}

func TestCharacterRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")

	_, err := repo.GetByID(context.Background(), c.ID.String(), "user_2")

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_List(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	seedCharacter(t, tx, "user_1")
	seedCharacter(t, tx, "user_1")
	seedCharacter(t, tx, "user_2")

	got, err := repo.List(context.Background(), "user_1")

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, c := range got {
		assert.Equal(t, "user_1", c.UserID)
	}
}

func TestCharacterRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")
	newName := "Updated Name"

	got, err := repo.Update(context.Background(), c.ID.String(), "user_1", usecase.UpdateCharacterParams{
		Name: &newName,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
}

func TestCharacterRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	newName := "x"
	_, err := repo.Update(context.Background(), uuid.NewString(), "user_1", usecase.UpdateCharacterParams{
		Name: &newName,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCharacterRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	c := seedCharacter(t, tx, "user_1")

	err := repo.Delete(context.Background(), c.ID.String(), "user_1")

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), c.ID.String(), "user_1")
	assert.Error(t, err)
}

func TestCharacterRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCharacterRepository(tx)

	err := repo.Delete(context.Background(), uuid.NewString(), "user_1")

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
