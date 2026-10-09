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

func seedCommission(t *testing.T, tx *gorm.DB, userID uuid.UUID, status string) *domain.Commission {
	t.Helper()
	c := &domain.Commission{
		ID:     uuid.New(),
		UserID: userID,
		Status: status,
	}
	require.NoError(t, tx.Create(c).Error)
	return c
}

func seedCommissionWithArtist(t *testing.T, tx *gorm.DB, userID uuid.UUID, artist *domain.Artist, status string) *domain.Commission {
	t.Helper()
	c := &domain.Commission{
		ID:       uuid.New(),
		UserID:   userID,
		ArtistID: &artist.ID,
		Status:   status,
	}
	require.NoError(t, tx.Create(c).Error)
	return c
}

func TestCommissionRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	title := "My Commission"
	c := &domain.Commission{
		ID:     uuid.New(),
		UserID: user.ID,
		Status: "waitlist",
		Title:  &title,
	}

	got, err := repo.Create(context.Background(), c)

	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	require.NotNil(t, got.Title)
	assert.Equal(t, "My Commission", *got.Title)
	assert.Equal(t, "waitlist", got.Status)
}

func TestCommissionRepository_Create_PreloadsArtistLinks(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	require.NoError(t, tx.Create(&domain.ArtistLink{
		ID:        uuid.New(),
		ArtistID:  artist.ID,
		URL:       "https://example.com",
		IsPrimary: true,
	}).Error)

	c := &domain.Commission{
		ID:       uuid.New(),
		UserID:   user.ID,
		ArtistID: &artist.ID,
		Status:   "waitlist",
	}

	got, err := repo.Create(context.Background(), c)

	require.NoError(t, err)
	require.NotNil(t, got.Artist)
	require.Len(t, got.Artist.Links, 1)
	assert.Equal(t, "https://example.com", got.Artist.Links[0].URL)
	assert.True(t, got.Artist.Links[0].IsPrimary)
}

func TestCommissionRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCommission(t, tx, user.ID, "wip")

	got, err := repo.GetByID(context.Background(), c.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	assert.Equal(t, "wip", got.Status)
}

func TestCommissionRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	c := seedCommission(t, tx, user1.ID, "wip")

	_, err := repo.GetByID(context.Background(), c.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCommissionRepository_List_ScopedToUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedCommission(t, tx, user1.ID, "waitlist")
	seedCommission(t, tx, user1.ID, "wip")
	seedCommission(t, tx, user2.ID, "waitlist")

	got, err := repo.List(context.Background(), user1.ID)

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, c := range got {
		assert.Equal(t, user1.ID, c.UserID)
	}
}

func TestCommissionRepository_List_PreloadsArtistLinks(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	require.NoError(t, tx.Create(&domain.ArtistLink{
		ID:        uuid.New(),
		ArtistID:  artist.ID,
		URL:       "https://portfolio.example.com",
		IsPrimary: true,
	}).Error)
	seedCommissionWithArtist(t, tx, user.ID, artist, "wip")

	got, err := repo.List(context.Background(), user.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Artist)
	require.Len(t, got[0].Artist.Links, 1)
	assert.Equal(t, "https://portfolio.example.com", got[0].Artist.Links[0].URL)
}

func TestCommissionRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCommission(t, tx, user.ID, "waitlist")
	newTitle := "Updated Title"

	got, err := repo.Update(context.Background(), c.ID, user.ID, usecase.UpdateCommissionParams{
		Title:  &newTitle,
		Status: "wip",
	})

	require.NoError(t, err)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Updated Title", *got.Title)
	assert.Equal(t, "wip", got.Status)
}

func TestCommissionRepository_Update_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	c := seedCommission(t, tx, user1.ID, "waitlist")

	_, err := repo.Update(context.Background(), c.ID, user2.ID, usecase.UpdateCommissionParams{Status: "wip"})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCommissionRepository_Patch_Status(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCommission(t, tx, user.ID, "waitlist")
	newStatus := "done"

	got, err := repo.Patch(context.Background(), c.ID, user.ID, usecase.PatchCommissionParams{
		Status: &newStatus,
	})

	require.NoError(t, err)
	assert.Equal(t, "done", got.Status)
}

func TestCommissionRepository_Patch_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	c := seedCommission(t, tx, user1.ID, "waitlist")
	newStatus := "done"

	_, err := repo.Patch(context.Background(), c.ID, user2.ID, usecase.PatchCommissionParams{
		Status: &newStatus,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCommissionRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user := seedUser(t, tx, "user_1")
	c := seedCommission(t, tx, user.ID, "waitlist")

	err := repo.Delete(context.Background(), c.ID, user.ID)

	require.NoError(t, err)
	_, err = repo.GetByID(context.Background(), c.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCommissionRepository_Delete_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewCommissionRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	c := seedCommission(t, tx, user1.ID, "waitlist")

	err := repo.Delete(context.Background(), c.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
