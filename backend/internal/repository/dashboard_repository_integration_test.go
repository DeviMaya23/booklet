package repository

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardRepository_GetRecentArtpieces_ScopedToUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedArtpiece(t, tx, user1.ID)
	seedArtpiece(t, tx, user2.ID)

	got, err := repo.GetRecentArtpieces(context.Background(), user1.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, user1.ID, got[0].UserID)
}

func TestDashboardRepository_GetRecentArtpieces_LimitsFive(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	for i := 0; i < 7; i++ {
		seedArtpiece(t, tx, user.ID)
	}

	got, err := repo.GetRecentArtpieces(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Len(t, got, 5)
}

func TestDashboardRepository_GetRecentArtpieces_PreloadsArtistLinks(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	require.NoError(t, tx.Create(&domain.ArtistLink{
		ID:        uuid.New(),
		ArtistID:  artist.ID,
		URL:       "https://example.com",
		IsPrimary: true,
	}).Error)

	ap := &domain.Artpiece{
		ID:       uuid.New(),
		UserID:   user.ID,
		ArtistID: &artist.ID,
	}
	require.NoError(t, tx.Create(ap).Error)

	got, err := repo.GetRecentArtpieces(context.Background(), user.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Artist)
	require.Len(t, got[0].Artist.Links, 1)
	assert.Equal(t, "https://example.com", got[0].Artist.Links[0].URL)
	assert.True(t, got[0].Artist.Links[0].IsPrimary)
}

func TestDashboardRepository_GetInProgressCommissions_ScopedToUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedCommission(t, tx, user1.ID, "wip")
	seedCommission(t, tx, user2.ID, "wip")

	got, err := repo.GetInProgressCommissions(context.Background(), user1.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, user1.ID, got[0].UserID)
}

func TestDashboardRepository_GetInProgressCommissions_OnlyWaitlistAndWip(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedCommission(t, tx, user.ID, "waitlist")
	seedCommission(t, tx, user.ID, "wip")
	seedCommission(t, tx, user.ID, "done")

	got, err := repo.GetInProgressCommissions(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, c := range got {
		assert.Contains(t, []string{"waitlist", "wip"}, c.Status)
	}
}

func TestDashboardRepository_GetInProgressCommissions_PreloadsArtistLinks(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	require.NoError(t, tx.Create(&domain.ArtistLink{
		ID:        uuid.New(),
		ArtistID:  artist.ID,
		URL:       "https://portfolio.example.com",
		IsPrimary: true,
	}).Error)
	seedCommissionWithArtist(t, tx, user.ID, artist, "wip")

	got, err := repo.GetInProgressCommissions(context.Background(), user.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Artist)
	require.Len(t, got[0].Artist.Links, 1)
	assert.Equal(t, "https://portfolio.example.com", got[0].Artist.Links[0].URL)
	assert.True(t, got[0].Artist.Links[0].IsPrimary)
}

func TestDashboardRepository_GetCommissionsNoArtist(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")
	seedCommission(t, tx, user.ID, "waitlist")
	seedCommissionWithArtist(t, tx, user.ID, artist, "wip")

	got, err := repo.GetCommissionsNoArtist(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestDashboardRepository_GetArtpiecesNoArtist(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Jane Doe")

	seedArtpiece(t, tx, user.ID)
	apWithArtist := &domain.Artpiece{
		ID:       uuid.New(),
		UserID:   user.ID,
		ArtistID: &artist.ID,
	}
	require.NoError(t, tx.Create(apWithArtist).Error)

	got, err := repo.GetArtpiecesNoArtist(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestDashboardRepository_GetDoneCommissionsNoArtpieces(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	done := seedCommission(t, tx, user.ID, "done")
	doneWithArtpiece := seedCommission(t, tx, user.ID, "done")
	ap := &domain.Artpiece{
		ID:           uuid.New(),
		UserID:       user.ID,
		CommissionID: &doneWithArtpiece.ID,
	}
	require.NoError(t, tx.Create(ap).Error)
	_ = done

	got, err := repo.GetDoneCommissionsNoArtpieces(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, done.ID, got[0].ID)
}

func TestDashboardRepository_GetArtpiecesNoFiles(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user := seedUser(t, tx, "user_1")
	seedArtpiece(t, tx, user.ID)

	got, err := repo.GetArtpiecesNoFiles(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestDashboardRepository_HousekeepingItems_ScopedToUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewDashboardRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedCommission(t, tx, user1.ID, "waitlist")
	seedCommission(t, tx, user2.ID, "waitlist")

	got, err := repo.GetCommissionsNoArtist(context.Background(), user1.ID)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}
