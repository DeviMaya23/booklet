package usecase_test

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestListImages_PassesFiltersToRepo(t *testing.T) {
	img := &domain.Image{ID: uuid.New(), UserID: uuid.New(), Characters: []domain.Character{}}
	repo := &fakeImageRepository{images: []*domain.Image{img}}
	uc := usecase.NewImageUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	q := "sunset"
	charID := uuid.New().String()
	artistID := uuid.New().String()
	filters := usecase.ListImageFilters{
		Q:            &q,
		CharacterIDs: []string{charID},
		ArtistIDs:    []string{artistID},
	}
	_, err := uc.List(context.Background(), img.UserID, filters)

	require.NoError(t, err)
	require.NotNil(t, repo.lastListFilters.Q)
	require.Equal(t, "sunset", *repo.lastListFilters.Q)
	require.Equal(t, []string{charID}, repo.lastListFilters.CharacterIDs)
	require.Equal(t, []string{artistID}, repo.lastListFilters.ArtistIDs)
}
