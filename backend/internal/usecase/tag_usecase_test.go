package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateTag_NormalizesName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercases", "WIP", "wip"},
		{"trims", "  sketch  ", "sketch"},
		{"trims and lowercases", "  Final Draft ", "final draft"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeTagRepository()
			uc := usecase.NewTagUsecase(repo, observability.NewTelemetry(nil, nil, nil))
			userID := uuid.New()

			got, err := uc.Create(context.Background(), userID, usecase.CreateTagParams{Name: tt.input})

			require.NoError(t, err)
			require.Equal(t, tt.want, got.Name)
			require.Equal(t, userID, got.UserID)
			require.Equal(t, tt.want, repo.tags[got.ID].Name)
		})
	}
}

func TestCreateTag_StampsLastUsedAt(t *testing.T) {
	repo := newFakeTagRepository()
	uc := usecase.NewTagUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	before := time.Now()
	got, err := uc.Create(context.Background(), uuid.New(), usecase.CreateTagParams{Name: "wip"})

	require.NoError(t, err)
	require.NotNil(t, got.LastUsedAt)
	require.WithinDuration(t, before, *got.LastUsedAt, 5*time.Second)
}

func TestCreateTag_BlankNameInvalid(t *testing.T) {
	for _, name := range []string{"", "   ", "\t\n"} {
		repo := newFakeTagRepository()
		uc := usecase.NewTagUsecase(repo, observability.NewTelemetry(nil, nil, nil))

		got, err := uc.Create(context.Background(), uuid.New(), usecase.CreateTagParams{Name: name})

		require.ErrorIs(t, err, usecase.ErrTagNameInvalid)
		require.Nil(t, got)
		require.Empty(t, repo.tags)
	}
}

func TestUpdateTag_NormalizesName(t *testing.T) {
	repo := newFakeTagRepository()
	uc := usecase.NewTagUsecase(repo, observability.NewTelemetry(nil, nil, nil))
	userID := uuid.New()
	tag := repo.seed(userID, "wip")

	got, err := uc.Update(context.Background(), tag.ID, userID, usecase.UpdateTagParams{Name: "  Finished "})

	require.NoError(t, err)
	require.Equal(t, "finished", got.Name)
	require.Equal(t, "finished", repo.tags[tag.ID].Name)
}

func TestUpdateTag_BlankNameInvalid(t *testing.T) {
	repo := newFakeTagRepository()
	uc := usecase.NewTagUsecase(repo, observability.NewTelemetry(nil, nil, nil))
	userID := uuid.New()
	tag := repo.seed(userID, "wip")

	got, err := uc.Update(context.Background(), tag.ID, userID, usecase.UpdateTagParams{Name: "  "})

	require.ErrorIs(t, err, usecase.ErrTagNameInvalid)
	require.Nil(t, got)
	require.Equal(t, "wip", repo.tags[tag.ID].Name)
}
