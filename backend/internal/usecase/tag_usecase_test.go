package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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
			uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))
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
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	before := time.Now()
	got, err := uc.Create(context.Background(), uuid.New(), usecase.CreateTagParams{Name: "wip"})

	require.NoError(t, err)
	require.NotNil(t, got.LastUsedAt)
	require.WithinDuration(t, before, *got.LastUsedAt, 5*time.Second)
}

func TestCreateTag_BlankNameInvalid(t *testing.T) {
	for _, name := range []string{"", "   ", "\t\n"} {
		repo := newFakeTagRepository()
		uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

		got, err := uc.Create(context.Background(), uuid.New(), usecase.CreateTagParams{Name: name})

		require.ErrorIs(t, err, usecase.ErrTagNameInvalid)
		require.Nil(t, got)
		require.Empty(t, repo.tags)
	}
}

func TestUpdateTag_NormalizesName(t *testing.T) {
	repo := newFakeTagRepository()
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))
	userID := uuid.New()
	tag := repo.seed(userID, "wip")

	got, err := uc.Update(context.Background(), tag.ID, userID, usecase.UpdateTagParams{Name: "  Finished "})

	require.NoError(t, err)
	require.Equal(t, "finished", got.Name)
	require.Equal(t, "finished", repo.tags[tag.ID].Name)
}

func TestUpdateTag_BlankNameInvalid(t *testing.T) {
	repo := newFakeTagRepository()
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))
	userID := uuid.New()
	tag := repo.seed(userID, "wip")

	got, err := uc.Update(context.Background(), tag.ID, userID, usecase.UpdateTagParams{Name: "  "})

	require.ErrorIs(t, err, usecase.ErrTagNameInvalid)
	require.Nil(t, got)
	require.Equal(t, "wip", repo.tags[tag.ID].Name)
}

// --- MergeTag ---

func newMergeFixture() (*fakeTagRepository, uuid.UUID) {
	return newFakeTagRepository(), uuid.New()
}

func TestMergeTag_SameTagRejected(t *testing.T) {
	repo, userID := newMergeFixture()
	tag := repo.seed(userID, "sketch")
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.Merge(context.Background(), tag.ID, userID, usecase.MergeTagParams{TargetID: tag.ID})

	require.ErrorIs(t, err, usecase.ErrTagMergeSameTag)
	require.Contains(t, repo.tags, tag.ID)
}

func TestMergeTag_MissingOrForeignTagNotFound(t *testing.T) {
	tests := []struct {
		name  string
		setup func(repo *fakeTagRepository, userID uuid.UUID) (source, target uuid.UUID)
	}{
		{"target missing", func(repo *fakeTagRepository, userID uuid.UUID) (uuid.UUID, uuid.UUID) {
			return repo.seed(userID, "a").ID, uuid.New()
		}},
		{"source missing", func(repo *fakeTagRepository, userID uuid.UUID) (uuid.UUID, uuid.UUID) {
			return uuid.New(), repo.seed(userID, "b").ID
		}},
		{"target owned by another user", func(repo *fakeTagRepository, userID uuid.UUID) (uuid.UUID, uuid.UUID) {
			return repo.seed(userID, "a").ID, repo.seed(uuid.New(), "b").ID
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, userID := newMergeFixture()
			source, target := tt.setup(repo, userID)
			artpiece := uuid.New()
			repo.assign(source, artpiece)
			uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

			err := uc.Merge(context.Background(), source, userID, usecase.MergeTagParams{TargetID: target})

			require.ErrorIs(t, err, gorm.ErrRecordNotFound)
			require.Equal(t, map[uuid.UUID]bool{artpiece: true}, repo.assignments[source])
			require.Empty(t, repo.assignments[target])
		})
	}
}

func TestMergeTag_MovesArtpiecesAndDeletesSource(t *testing.T) {
	repo, userID := newMergeFixture()
	source := repo.seed(userID, "sketches")
	target := repo.seed(userID, "sketch")
	onlySource, both := uuid.New(), uuid.New()
	repo.assign(source.ID, onlySource)
	repo.assign(source.ID, both)
	repo.assign(target.ID, both)
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.Merge(context.Background(), source.ID, userID, usecase.MergeTagParams{TargetID: target.ID})

	require.NoError(t, err)
	require.NotContains(t, repo.tags, source.ID)
	require.Equal(t, map[uuid.UUID]bool{onlySource: true, both: true}, repo.assignments[target.ID])
	require.Equal(t, "sketch", repo.tags[target.ID].Name)
	require.Nil(t, repo.tags[target.ID].LastUsedAt)
}

func TestMergeTag_ReassignFailureKeepsSource(t *testing.T) {
	repo, userID := newMergeFixture()
	source := repo.seed(userID, "sketches")
	target := repo.seed(userID, "sketch")
	boom := errors.New("boom")
	repo.reassignErr = boom
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.Merge(context.Background(), source.ID, userID, usecase.MergeTagParams{TargetID: target.ID})

	require.ErrorIs(t, err, boom)
	require.Contains(t, repo.tags, source.ID)
}

func TestMergeTag_DeleteFailureIsReturned(t *testing.T) {
	repo, userID := newMergeFixture()
	source := repo.seed(userID, "sketches")
	target := repo.seed(userID, "sketch")
	boom := errors.New("boom")
	repo.deleteErr = boom
	uc := usecase.NewTagUsecase(repo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.Merge(context.Background(), source.ID, userID, usecase.MergeTagParams{TargetID: target.ID})

	require.ErrorIs(t, err, boom)
	require.Contains(t, repo.tags, source.ID)
}
