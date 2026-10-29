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

type fakeVerseRepository struct {
	verses       map[uuid.UUID]*domain.Verse
	conflictName string
}

func newFakeVerseRepository() *fakeVerseRepository {
	return &fakeVerseRepository{verses: make(map[uuid.UUID]*domain.Verse)}
}

func (f *fakeVerseRepository) Create(_ context.Context, verse *domain.Verse) (*domain.Verse, error) {
	if f.conflictName != "" && verse.Name == f.conflictName {
		return nil, usecase.ErrVerseNameConflict
	}
	f.verses[verse.ID] = verse
	return verse, nil
}

func (f *fakeVerseRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Verse, error) {
	v, ok := f.verses[id]
	if !ok || v.UserID != userID {
		return nil, fmt.Errorf("get verse: %w", gorm.ErrRecordNotFound)
	}
	return v, nil
}

func (f *fakeVerseRepository) List(_ context.Context, userID uuid.UUID, _ usecase.ListVerseFilters) ([]*domain.Verse, error) {
	var result []*domain.Verse
	for _, v := range f.verses {
		if v.UserID == userID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (f *fakeVerseRepository) Update(_ context.Context, id uuid.UUID, _ uuid.UUID, params usecase.UpdateVerseParams) (*domain.Verse, error) {
	v, ok := f.verses[id]
	if !ok {
		return nil, fmt.Errorf("update verse: %w", gorm.ErrRecordNotFound)
	}
	if f.conflictName != "" && params.Name == f.conflictName {
		return nil, usecase.ErrVerseNameConflict
	}
	v.Name = params.Name
	return v, nil
}

func (f *fakeVerseRepository) Delete(_ context.Context, id uuid.UUID, _ uuid.UUID) error {
	if _, ok := f.verses[id]; !ok {
		return fmt.Errorf("delete verse: %w", gorm.ErrRecordNotFound)
	}
	delete(f.verses, id)
	return nil
}

func (f *fakeVerseRepository) UpdateLastUsedAt(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

func newFakeVerseRepositoryWithVerse(repo *fakeVerseRepository, userID uuid.UUID, name string) (*domain.Verse, error) {
	v := &domain.Verse{ID: uuid.New(), UserID: userID, Name: name}
	repo.verses[v.ID] = v
	return v, nil
}

// --- tests ---

func TestCreateVerse_AssemblesVerse(t *testing.T) {
	repo := newFakeVerseRepository()
	uc := usecase.NewVerseUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	got, err := uc.Create(context.Background(), userID, usecase.CreateVerseParams{Name: "Solaris"})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "Solaris", got.Name)
}

func TestCreateVerse_DuplicateName_ReturnsConflict(t *testing.T) {
	repo := newFakeVerseRepository()
	repo.conflictName = "Solaris"
	uc := usecase.NewVerseUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.Create(context.Background(), uuid.New(), usecase.CreateVerseParams{Name: "Solaris"})

	require.ErrorIs(t, err, usecase.ErrVerseNameConflict)
}

func TestUpdateVerse_NotFound(t *testing.T) {
	repo := newFakeVerseRepository()
	uc := usecase.NewVerseUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.Update(context.Background(), uuid.New(), uuid.New(), usecase.UpdateVerseParams{Name: "x"})

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUpdateVerse_DuplicateName_ReturnsConflict(t *testing.T) {
	repo := newFakeVerseRepository()
	uc := usecase.NewVerseUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	existing, err := uc.Create(context.Background(), userID, usecase.CreateVerseParams{Name: "Solaris"})
	require.NoError(t, err)

	repo.conflictName = "Terran"
	_, err = uc.Update(context.Background(), existing.ID, userID, usecase.UpdateVerseParams{Name: "Terran"})

	require.ErrorIs(t, err, usecase.ErrVerseNameConflict)
}

func TestDeleteVerse_NotFound(t *testing.T) {
	repo := newFakeVerseRepository()
	uc := usecase.NewVerseUsecase(repo, observability.NewTelemetry(nil, nil, nil))

	err := uc.Delete(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
