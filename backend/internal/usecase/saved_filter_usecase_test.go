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

// fakeSavedFilterRepository implements SavedFilterRepository for tests.
type fakeSavedFilterRepository struct {
	filters   map[uuid.UUID]*domain.SavedFilter
	createErr error
}

func newFakeSavedFilterRepository() *fakeSavedFilterRepository {
	return &fakeSavedFilterRepository{
		filters: make(map[uuid.UUID]*domain.SavedFilter),
	}
}

func (f *fakeSavedFilterRepository) Create(_ context.Context, sf *domain.SavedFilter) (*domain.SavedFilter, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
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

// spyArtpieceThumbnailRepository implements ArtpieceThumbnailRepository for tests.
type spyArtpieceThumbnailRepository struct {
	result *domain.Artpiece
	err    error
}

func (s *spyArtpieceThumbnailRepository) GetByID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Artpiece, error) {
	return s.result, s.err
}

func TestCreateSavedFilter_AssemblesRecord(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	artpieceRepo := &spyArtpieceThumbnailRepository{}
	uc := usecase.NewSavedFilterUsecase(repo, artpieceRepo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	payload := json.RawMessage(`{"artist_id":null,"character_ids":[],"match_mode":"all"}`)
	got, err := uc.Create(context.Background(), userID, usecase.CreateSavedFilterParams{
		Name:          "My Filter",
		FilterPayload: payload,
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "My Filter", got.Name)
	require.Nil(t, got.ThumbnailR2Path)
	require.JSONEq(t, `{"artist_id":null,"character_ids":[],"match_mode":"all"}`, string(got.FilterPayload))
}

func TestCreateSavedFilter_WithThumbnailArtpieceID_CopiesR2Path(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	r2Path := "thumbnails/artpieces/abc.webp"
	artpieceRepo := &spyArtpieceThumbnailRepository{
		result: &domain.Artpiece{
			CoverFile: &domain.File{ThumbnailR2Path: &r2Path},
		},
	}
	uc := usecase.NewSavedFilterUsecase(repo, artpieceRepo, observability.NewTelemetry(nil, nil, nil))

	artpieceID := uuid.New()
	got, err := uc.Create(context.Background(), uuid.New(), usecase.CreateSavedFilterParams{
		Name:                "Filtered",
		ThumbnailArtpieceID: &artpieceID,
		FilterPayload:       json.RawMessage(`{}`),
	})

	require.NoError(t, err)
	require.NotNil(t, got.ThumbnailR2Path)
	require.Equal(t, r2Path, *got.ThumbnailR2Path)
}

func TestCreateSavedFilter_ThumbnailArtpieceHasNoFile_StoresNilPath(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	artpieceRepo := &spyArtpieceThumbnailRepository{
		result: &domain.Artpiece{CoverFile: nil},
	}
	uc := usecase.NewSavedFilterUsecase(repo, artpieceRepo, observability.NewTelemetry(nil, nil, nil))

	artpieceID := uuid.New()
	got, err := uc.Create(context.Background(), uuid.New(), usecase.CreateSavedFilterParams{
		Name:                "No Thumb",
		ThumbnailArtpieceID: &artpieceID,
		FilterPayload:       json.RawMessage(`{}`),
	})

	require.NoError(t, err)
	require.Nil(t, got.ThumbnailR2Path)
}

func TestCreateSavedFilter_DuplicateName_ReturnsErrDuplicateSavedFilterName(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	repo.createErr = usecase.ErrDuplicateSavedFilterName
	artpieceRepo := &spyArtpieceThumbnailRepository{}
	uc := usecase.NewSavedFilterUsecase(repo, artpieceRepo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.Create(context.Background(), uuid.New(), usecase.CreateSavedFilterParams{
		Name:          "Duplicate",
		FilterPayload: json.RawMessage(`{}`),
	})

	require.ErrorIs(t, err, usecase.ErrDuplicateSavedFilterName)
}

func TestGetSavedFilterByID_NotFound_ReturnsErrSavedFilterNotFound(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, &spyArtpieceThumbnailRepository{}, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.GetByID(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, usecase.ErrSavedFilterNotFound)
}

func TestUpdateSavedFilter_WithThumbnailArtpieceID_CopiesR2Path(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	r2Path := "thumbnails/artpieces/update.webp"
	artpieceRepo := &spyArtpieceThumbnailRepository{
		result: &domain.Artpiece{
			CoverFile: &domain.File{ThumbnailR2Path: &r2Path},
		},
	}
	uc := usecase.NewSavedFilterUsecase(repo, artpieceRepo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	existing := &domain.SavedFilter{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          "My Filter",
		FilterPayload: json.RawMessage(`{}`),
	}
	repo.filters[existing.ID] = existing

	artpieceID := uuid.New()
	got, err := uc.Update(context.Background(), existing.ID, userID, usecase.UpdateSavedFilterParams{
		ThumbnailArtpieceID: &artpieceID,
	})

	require.NoError(t, err)
	require.NotNil(t, got.ThumbnailR2Path)
	require.Equal(t, r2Path, *got.ThumbnailR2Path)
}

func TestUpdateSavedFilter_WithThumbnailArtpieceID_ArtpieceHasNoFile_StoresNilPath(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	artpieceRepo := &spyArtpieceThumbnailRepository{
		result: &domain.Artpiece{CoverFile: nil},
	}
	uc := usecase.NewSavedFilterUsecase(repo, artpieceRepo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	existing := &domain.SavedFilter{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          "My Filter",
		FilterPayload: json.RawMessage(`{}`),
	}
	repo.filters[existing.ID] = existing

	artpieceID := uuid.New()
	got, err := uc.Update(context.Background(), existing.ID, userID, usecase.UpdateSavedFilterParams{
		ThumbnailArtpieceID: &artpieceID,
	})

	require.NoError(t, err)
	require.Nil(t, got.ThumbnailR2Path)
}

func TestUpdateSavedFilter_NotFound_ReturnsErrSavedFilterNotFound(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, &spyArtpieceThumbnailRepository{}, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.Update(context.Background(), uuid.New(), uuid.New(), usecase.UpdateSavedFilterParams{})

	require.ErrorIs(t, err, usecase.ErrSavedFilterNotFound)
}

func TestDeleteSavedFilter_NotFound_ReturnsErrSavedFilterNotFound(t *testing.T) {
	repo := newFakeSavedFilterRepository()
	uc := usecase.NewSavedFilterUsecase(repo, &spyArtpieceThumbnailRepository{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.Delete(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, usecase.ErrSavedFilterNotFound)
}
