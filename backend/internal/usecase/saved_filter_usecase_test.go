package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeSavedFilterRepository struct {
	filters map[uuid.UUID]*domain.SavedFilter
}

func newFakeSavedFilterRepository() *fakeSavedFilterRepository {
	return &fakeSavedFilterRepository{
		filters: make(map[uuid.UUID]*domain.SavedFilter),
	}
}

func (f *fakeSavedFilterRepository) Create(_ context.Context, sf *domain.SavedFilter) (*domain.SavedFilter, error) {
	f.filters[sf.ID] = sf
	return sf, nil
}

func (f *fakeSavedFilterRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.SavedFilter, error) {
	sf, ok := f.filters[id]
	if !ok || sf.UserID != userID {
		return nil, fmt.Errorf("get saved filter: %w", gorm.ErrRecordNotFound)
	}
	return sf, nil
}

func (f *fakeSavedFilterRepository) List(_ context.Context, userID uuid.UUID) ([]*domain.SavedFilter, error) {
	var result []*domain.SavedFilter
	for _, sf := range f.filters {
		if sf.UserID == userID {
			result = append(result, sf)
		}
	}
	return result, nil
}

func (f *fakeSavedFilterRepository) Update(_ context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateSavedFilterParams) (*domain.SavedFilter, error) {
	sf, ok := f.filters[id]
	if !ok || sf.UserID != userID {
		return nil, fmt.Errorf("update saved filter: %w", gorm.ErrRecordNotFound)
	}
	if params.Name.Set && params.Name.Value != nil {
		sf.Name = *params.Name.Value
	}
	if params.ThumbnailR2Path.Set {
		sf.ThumbnailR2Path = params.ThumbnailR2Path.Value
	}
	if params.FilterPayload != nil {
		sf.FilterPayload = *params.FilterPayload
	}
	return sf, nil
}

func (f *fakeSavedFilterRepository) Delete(_ context.Context, id uuid.UUID, userID uuid.UUID) error {
	sf, ok := f.filters[id]
	if !ok || sf.UserID != userID {
		return fmt.Errorf("delete saved filter: %w", gorm.ErrRecordNotFound)
	}
	delete(f.filters, id)
	return nil
}

func TestCreateSavedFilter_AssemblesRecord(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	payload := json.RawMessage(`{"artists":["abc"]}`)
	got, err := uc.Create(context.Background(), userID, usecase.CreateSavedFilterParams{
		Name:          "My Filter",
		FilterPayload: payload,
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "My Filter", got.Name)
	require.JSONEq(t, `{"artists":["abc"]}`, string(got.FilterPayload))
}

func TestGetSavedFilterByID_NotFound_ReturnsErrSavedFilterNotFound(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.GetByID(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, usecase.ErrSavedFilterNotFound)
}

func TestUpdateSavedFilter_NotFound_ReturnsErrSavedFilterNotFound(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.Update(context.Background(), uuid.New(), uuid.New(), usecase.UpdateSavedFilterParams{})

	require.ErrorIs(t, err, usecase.ErrSavedFilterNotFound)
}

func TestDeleteSavedFilter_NotFound_ReturnsErrSavedFilterNotFound(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	err := uc.Delete(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, usecase.ErrSavedFilterNotFound)
}
