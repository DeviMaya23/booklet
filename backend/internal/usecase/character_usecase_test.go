package usecase_test

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreate_AssemblesCharacter(t *testing.T) {
	repo := newFakeCharacterRepository()
	uc := usecase.NewCharacterUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	got, err := uc.Create(context.Background(), "user-1", usecase.CreateCharacterParams{
		Name: "Aria Stormweaver",
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, "user-1", got.UserID)
	require.Equal(t, "Aria Stormweaver", got.Name)
}

func TestUpdate_PartialFields(t *testing.T) {
	newName := "Updated Name"
	falseVal := false
	truVal := true
	bio := "A new biography."

	tests := []struct {
		name         string
		params       usecase.UpdateCharacterParams
		expectName   *string
		expectPublic *bool
		expectBio    *string
	}{
		{
			name:       "name only",
			params:     usecase.UpdateCharacterParams{Name: &newName},
			expectName: &newName,
		},
		{
			name:         "is_public false only",
			params:       usecase.UpdateCharacterParams{IsPublic: &falseVal},
			expectPublic: &falseVal,
		},
		{
			name:         "multiple fields",
			params:       usecase.UpdateCharacterParams{Name: &newName, IsPublic: &truVal, Biography: &bio},
			expectName:   &newName,
			expectPublic: &truVal,
			expectBio:    &bio,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeCharacterRepository()
			uc := usecase.NewCharacterUsecase(repo, observability.NewTelemetry(nil, nil, nil))

			_, _ = uc.Create(context.Background(), "user-1", usecase.CreateCharacterParams{Name: "Original"})

			_, err := uc.Update(context.Background(), repo.lastCreated.ID.String(), "user-1", tc.params)
			require.NoError(t, err)

			got := repo.lastUpdated
			require.Equal(t, tc.expectName, got.Name)
			require.Equal(t, tc.expectPublic, got.IsPublic)
			require.Equal(t, tc.expectBio, got.Biography)
		})
	}
}
