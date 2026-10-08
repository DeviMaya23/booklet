package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

var ErrSavedFilterNotFound = errors.New("saved filter not found")

// Patch[T] is a three-state PATCH field: absent (not in request), null (explicit clear), or a value.
type Patch[T any] struct {
	Set   bool
	Value *T
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	p.Set = true
	if string(data) == "null" {
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	p.Value = &v
	return nil
}

type CreateSavedFilterParams struct {
	Name            string
	ThumbnailR2Path *string
	FilterPayload   json.RawMessage
}

type UpdateSavedFilterParams struct {
	Name            Patch[string]
	ThumbnailR2Path Patch[string]
	FilterPayload   *json.RawMessage
}

type SavedFilterRepository interface {
	Create(ctx context.Context, sf *domain.SavedFilter) (*domain.SavedFilter, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.SavedFilter, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.SavedFilter, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateSavedFilterParams) (*domain.SavedFilter, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
