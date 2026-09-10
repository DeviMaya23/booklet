package usecase

import (
	"context"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type UpdateCharacterParams struct {
	Name         string
	AvatarR2Path *string
	Notes        *string
	IsPublic     bool
	FolderIDs    []uuid.UUID
}

type ListCharacterFilters struct {
	Q *string `query:"q"`
}

type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) error
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Character, error)
	List(ctx context.Context, userID uuid.UUID, filters ListCharacterFilters) ([]*domain.Character, error)
	Update(ctx context.Context, id string, userID uuid.UUID, params UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id string, userID uuid.UUID) error
	UpdateAvatarR2Path(ctx context.Context, id string, userID uuid.UUID, r2Key string) error
	ClearAvatarR2Path(ctx context.Context, id string, userID uuid.UUID) (oldKey string, err error)
}

type CharacterAvatarUploadRepository interface {
	Create(ctx context.Context, p *domain.PendingCharacterAvatarUpload) (*domain.PendingCharacterAvatarUpload, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingCharacterAvatarUpload, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingCharacterAvatarUpload, error)
}
