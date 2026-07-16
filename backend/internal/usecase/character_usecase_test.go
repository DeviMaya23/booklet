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

