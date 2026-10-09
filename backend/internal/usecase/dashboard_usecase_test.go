package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- fake ---

type fakeDashboardRepository struct {
	recentArtpieces      []*domain.Artpiece
	inProgress           []*domain.Commission
	commissionsNoArtist  []usecase.DashboardHousekeepingItem
	artpiecesNoArtist    []usecase.DashboardHousekeepingItem
	doneNoArtpieces      []usecase.DashboardHousekeepingItem
	artpiecesNoFiles     []usecase.DashboardHousekeepingItem

	errRecentArtpieces     error
	errInProgress          error
	errCommissionsNoArtist error
	errArtpiecesNoArtist   error
	errDoneNoArtpieces     error
	errArtpiecesNoFiles    error
}

func (f *fakeDashboardRepository) GetRecentArtpieces(_ context.Context, _ uuid.UUID) ([]*domain.Artpiece, error) {
	return f.recentArtpieces, f.errRecentArtpieces
}

func (f *fakeDashboardRepository) GetInProgressCommissions(_ context.Context, _ uuid.UUID) ([]*domain.Commission, error) {
	return f.inProgress, f.errInProgress
}

func (f *fakeDashboardRepository) GetCommissionsNoArtist(_ context.Context, _ uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return f.commissionsNoArtist, f.errCommissionsNoArtist
}

func (f *fakeDashboardRepository) GetArtpiecesNoArtist(_ context.Context, _ uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return f.artpiecesNoArtist, f.errArtpiecesNoArtist
}

func (f *fakeDashboardRepository) GetDoneCommissionsNoArtpieces(_ context.Context, _ uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return f.doneNoArtpieces, f.errDoneNoArtpieces
}

func (f *fakeDashboardRepository) GetArtpiecesNoFiles(_ context.Context, _ uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return f.artpiecesNoFiles, f.errArtpiecesNoFiles
}

func newDashboardUsecase(repo *fakeDashboardRepository) *usecase.DashboardUsecase {
	return usecase.NewDashboardUsecase(repo, observability.NewTelemetry(nil, nil, nil))
}

// --- tests ---

func TestDashboardUsecase_GetDashboard_Success(t *testing.T) {
	title := "Chibi"
	artpieceID := uuid.New()
	commissionID := uuid.New()

	repo := &fakeDashboardRepository{
		recentArtpieces: []*domain.Artpiece{{ID: artpieceID, Title: &title}},
		inProgress:      []*domain.Commission{{ID: commissionID, Status: "wip"}},
		commissionsNoArtist: []usecase.DashboardHousekeepingItem{
			{ID: uuid.New(), Title: &title},
		},
		artpiecesNoArtist:   []usecase.DashboardHousekeepingItem{},
		doneNoArtpieces:     []usecase.DashboardHousekeepingItem{},
		artpiecesNoFiles:    []usecase.DashboardHousekeepingItem{},
	}

	u := newDashboardUsecase(repo)
	result, err := u.GetDashboard(context.Background(), uuid.New())

	require.NoError(t, err)
	require.Len(t, result.RecentArtpieces, 1)
	assert.Equal(t, artpieceID, result.RecentArtpieces[0].ID)
	require.Len(t, result.InProgress, 1)
	assert.Equal(t, commissionID, result.InProgress[0].ID)
	require.Len(t, result.CommissionsNoArtist, 1)
	assert.Empty(t, result.ArtpiecesNoArtist)
	assert.Empty(t, result.DoneNoArtpieces)
	assert.Empty(t, result.ArtpiecesNoFiles)
}

func TestDashboardUsecase_GetDashboard_RecentArtpiecesError(t *testing.T) {
	repo := &fakeDashboardRepository{
		errRecentArtpieces: errors.New("db error"),
	}
	u := newDashboardUsecase(repo)
	_, err := u.GetDashboard(context.Background(), uuid.New())
	assert.ErrorContains(t, err, "db error")
}

func TestDashboardUsecase_GetDashboard_InProgressError(t *testing.T) {
	repo := &fakeDashboardRepository{
		recentArtpieces: []*domain.Artpiece{},
		errInProgress:   errors.New("db error"),
	}
	u := newDashboardUsecase(repo)
	_, err := u.GetDashboard(context.Background(), uuid.New())
	assert.ErrorContains(t, err, "db error")
}

func TestDashboardUsecase_GetDashboard_CommissionsNoArtistError(t *testing.T) {
	repo := &fakeDashboardRepository{
		recentArtpieces:        []*domain.Artpiece{},
		inProgress:             []*domain.Commission{},
		errCommissionsNoArtist: errors.New("db error"),
	}
	u := newDashboardUsecase(repo)
	_, err := u.GetDashboard(context.Background(), uuid.New())
	assert.ErrorContains(t, err, "db error")
}

func TestDashboardUsecase_GetDashboard_ArtpiecesNoArtistError(t *testing.T) {
	repo := &fakeDashboardRepository{
		recentArtpieces:      []*domain.Artpiece{},
		inProgress:           []*domain.Commission{},
		commissionsNoArtist:  []usecase.DashboardHousekeepingItem{},
		errArtpiecesNoArtist: errors.New("db error"),
	}
	u := newDashboardUsecase(repo)
	_, err := u.GetDashboard(context.Background(), uuid.New())
	assert.ErrorContains(t, err, "db error")
}

func TestDashboardUsecase_GetDashboard_DoneNoArtpiecesError(t *testing.T) {
	repo := &fakeDashboardRepository{
		recentArtpieces:     []*domain.Artpiece{},
		inProgress:          []*domain.Commission{},
		commissionsNoArtist: []usecase.DashboardHousekeepingItem{},
		artpiecesNoArtist:   []usecase.DashboardHousekeepingItem{},
		errDoneNoArtpieces:  errors.New("db error"),
	}
	u := newDashboardUsecase(repo)
	_, err := u.GetDashboard(context.Background(), uuid.New())
	assert.ErrorContains(t, err, "db error")
}

func TestDashboardUsecase_GetDashboard_ArtpiecesNoFilesError(t *testing.T) {
	repo := &fakeDashboardRepository{
		recentArtpieces:     []*domain.Artpiece{},
		inProgress:          []*domain.Commission{},
		commissionsNoArtist: []usecase.DashboardHousekeepingItem{},
		artpiecesNoArtist:   []usecase.DashboardHousekeepingItem{},
		doneNoArtpieces:     []usecase.DashboardHousekeepingItem{},
		errArtpiecesNoFiles: errors.New("db error"),
	}
	u := newDashboardUsecase(repo)
	_, err := u.GetDashboard(context.Background(), uuid.New())
	assert.ErrorContains(t, err, "db error")
}
