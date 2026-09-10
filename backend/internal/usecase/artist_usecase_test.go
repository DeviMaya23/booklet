package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// fakeArtistRepository is an in-memory fake for ArtistRepository.
type fakeArtistRepository struct {
	artists         map[uuid.UUID]*domain.Artist
	conflictName    string // if non-empty, any Create/Update with this name returns ErrArtistNameConflict
	lastListFilters usecase.ListArtistFilters
}

func newFakeArtistRepository() *fakeArtistRepository {
	return &fakeArtistRepository{
		artists: make(map[uuid.UUID]*domain.Artist),
	}
}

func (f *fakeArtistRepository) Create(_ context.Context, artist *domain.Artist) (*domain.Artist, error) {
	if f.conflictName != "" && artist.Name == f.conflictName {
		return nil, usecase.ErrArtistNameConflict
	}
	f.artists[artist.ID] = artist
	return artist, nil
}

func (f *fakeArtistRepository) GetByID(_ context.Context, id string, _ uuid.UUID) (*domain.Artist, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("get artist: %w", err)
	}
	a, ok := f.artists[parsed]
	if !ok {
		return nil, fmt.Errorf("get artist: %w", gorm.ErrRecordNotFound)
	}
	return a, nil
}

func (f *fakeArtistRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, _ uuid.UUID) (*domain.Artist, error) {
	a, ok := f.artists[id]
	if !ok {
		return nil, fmt.Errorf("get artist: %w", gorm.ErrRecordNotFound)
	}
	return a, nil
}

func (f *fakeArtistRepository) List(_ context.Context, userID uuid.UUID, filters usecase.ListArtistFilters) ([]*domain.Artist, error) {
	f.lastListFilters = filters
	var result []*domain.Artist
	for _, a := range f.artists {
		if a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (f *fakeArtistRepository) Update(_ context.Context, id string, _ uuid.UUID, params usecase.UpdateArtistParams) (*domain.Artist, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("update artist: %w", err)
	}
	a, ok := f.artists[parsed]
	if !ok {
		return nil, fmt.Errorf("update artist: %w", gorm.ErrRecordNotFound)
	}
	if f.conflictName != "" && params.Name == f.conflictName {
		return nil, usecase.ErrArtistNameConflict
	}
	a.Name = params.Name
	a.Notes = params.Notes
	a.ArtistLink = params.ArtistLink
	return a, nil
}

func (f *fakeArtistRepository) Delete(_ context.Context, id string, _ uuid.UUID) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("delete artist: %w", err)
	}
	if _, ok := f.artists[parsed]; !ok {
		return fmt.Errorf("delete artist: %w", gorm.ErrRecordNotFound)
	}
	delete(f.artists, parsed)
	return nil
}

// --- tests ---

func TestCreateArtist_AssemblesArtist(t *testing.T) {
	repo := newFakeArtistRepository()
	uc := usecase.NewArtistUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	got, err := uc.Create(context.Background(), userID, usecase.CreateArtistParams{
		Name: "Jane Doe",
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "Jane Doe", got.Name)
}

func TestCreateArtist_DuplicateName_ReturnsConflict(t *testing.T) {
	repo := newFakeArtistRepository()
	repo.conflictName = "Jane Doe"
	uc := usecase.NewArtistUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.Create(context.Background(), uuid.New(), usecase.CreateArtistParams{
		Name: "Jane Doe",
	})

	require.ErrorIs(t, err, usecase.ErrArtistNameConflict)
}

func TestGetArtistByID_NotFound(t *testing.T) {
	repo := newFakeArtistRepository()
	uc := usecase.NewArtistUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.GetByID(context.Background(), uuid.New().String(), uuid.New())

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUpdateArtist_DuplicateName_ReturnsConflict(t *testing.T) {
	repo := newFakeArtistRepository()
	uc := usecase.NewArtistUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	existing, err := uc.Create(context.Background(), userID, usecase.CreateArtistParams{Name: "Aria"})
	require.NoError(t, err)

	repo.conflictName = "Jane Doe"
	_, err = uc.Update(context.Background(), existing.ID.String(), userID, usecase.UpdateArtistParams{
		Name: "Jane Doe",
	})

	require.ErrorIs(t, err, usecase.ErrArtistNameConflict)
}

func TestDeleteArtist_RemovesArtist(t *testing.T) {
	repo := newFakeArtistRepository()
	uc := usecase.NewArtistUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	artist, err := uc.Create(context.Background(), userID, usecase.CreateArtistParams{Name: "Aria"})
	require.NoError(t, err)

	err = uc.Delete(context.Background(), artist.ID.String(), userID)
	require.NoError(t, err)

	_, err = uc.GetByID(context.Background(), artist.ID.String(), userID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestListArtists_PassesFiltersToRepo(t *testing.T) {
	repo := newFakeArtistRepository()
	uc := usecase.NewArtistUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	_, err := uc.Create(context.Background(), userID, usecase.CreateArtistParams{Name: "Jane"})
	require.NoError(t, err)

	q := "jane"
	filters := usecase.ListArtistFilters{Q: &q}
	_, err = uc.List(context.Background(), userID, filters)

	require.NoError(t, err)
	require.NotNil(t, repo.lastListFilters.Q)
	require.Equal(t, "jane", *repo.lastListFilters.Q)
}
