package usecase_test

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeCharacterRepository struct {
	characters map[uuid.UUID]*domain.Character

	lastCreated          *domain.Character
	lastUpdated          usecase.UpdateCharacterParams
	lastUpdatedAvatarID  string
	lastUpdatedAvatarKey string
	lastListFilters      usecase.ListCharacterFilters
}

func newFakeCharacterRepository() *fakeCharacterRepository {
	return &fakeCharacterRepository{
		characters: make(map[uuid.UUID]*domain.Character),
	}
}

func (f *fakeCharacterRepository) Create(_ context.Context, character *domain.Character) error {
	f.characters[character.ID] = character
	f.lastCreated = character
	return nil
}

func (f *fakeCharacterRepository) GetByID(_ context.Context, id string, userID uuid.UUID) (*domain.Character, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}
	c, ok := f.characters[parsed]
	if !ok {
		return nil, fmt.Errorf("get character: %w", gorm.ErrRecordNotFound)
	}
	if c.UserID != userID {
		return nil, fmt.Errorf("get character: %w", gorm.ErrRecordNotFound)
	}
	return c, nil
}

func (f *fakeCharacterRepository) List(_ context.Context, userID uuid.UUID, filters usecase.ListCharacterFilters) ([]*domain.Character, error) {
	f.lastListFilters = filters
	var result []*domain.Character
	for _, c := range f.characters {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (f *fakeCharacterRepository) Update(_ context.Context, id string, _ uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	f.lastUpdated = params
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("update character: %w", err)
	}
	c, ok := f.characters[parsed]
	if !ok {
		return nil, fmt.Errorf("update character: %w", gorm.ErrRecordNotFound)
	}
	if params.Name != nil {
		c.Name = *params.Name
	}
	if params.AvatarR2Path != nil {
		c.AvatarR2Path = params.AvatarR2Path
	}
	if params.Biography != nil {
		c.Biography = params.Biography
	}
	if params.IsPublic != nil {
		c.IsPublic = *params.IsPublic
	}
	if params.FolderIDs != nil {
		folders := make([]domain.CharacterFolder, len(*params.FolderIDs))
		for i, id := range *params.FolderIDs {
			folders[i] = domain.CharacterFolder{CharacterID: parsed, FolderID: id}
		}
		c.Folders = folders
	}
	return c, nil
}

func (f *fakeCharacterRepository) Delete(_ context.Context, id string, _ uuid.UUID) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("delete character: %w", err)
	}
	if _, ok := f.characters[parsed]; !ok {
		return fmt.Errorf("delete character: %w", gorm.ErrRecordNotFound)
	}
	delete(f.characters, parsed)
	return nil
}

func (f *fakeCharacterRepository) ClearAvatarR2Path(_ context.Context, id string, userID uuid.UUID) (string, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("clear avatar_r2_path: %w", err)
	}
	c, ok := f.characters[parsed]
	if !ok || c.UserID != userID {
		return "", gorm.ErrRecordNotFound
	}
	if c.AvatarR2Path == nil {
		return "", nil
	}
	oldKey := *c.AvatarR2Path
	c.AvatarR2Path = nil
	return oldKey, nil
}

func (f *fakeCharacterRepository) UpdateAvatarR2Path(_ context.Context, id string, _ uuid.UUID, r2Key string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("update avatar_r2_path: %w", err)
	}
	c, ok := f.characters[parsed]
	if !ok {
		return fmt.Errorf("update avatar_r2_path: %w", gorm.ErrRecordNotFound)
	}
	f.lastUpdatedAvatarID = id
	f.lastUpdatedAvatarKey = r2Key
	c.AvatarR2Path = &r2Key
	return nil
}

type fakeCharacterAvatarUploadRepository struct {
	pending map[uuid.UUID]*domain.PendingCharacterAvatarUpload
	stale   []*domain.PendingCharacterAvatarUpload

	lastDeleted uuid.UUID
	deletedIDs  []uuid.UUID
}

func newFakeCharacterAvatarUploadRepository() *fakeCharacterAvatarUploadRepository {
	return &fakeCharacterAvatarUploadRepository{
		pending: make(map[uuid.UUID]*domain.PendingCharacterAvatarUpload),
	}
}

func (f *fakeCharacterAvatarUploadRepository) Create(_ context.Context, p *domain.PendingCharacterAvatarUpload) (*domain.PendingCharacterAvatarUpload, error) {
	f.pending[p.ID] = p
	return p, nil
}

func (f *fakeCharacterAvatarUploadRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingCharacterAvatarUpload, error) {
	p, ok := f.pending[id]
	if !ok {
		return nil, fmt.Errorf("select pending character avatar upload: %w", gorm.ErrRecordNotFound)
	}
	if p.UserID != userID {
		return nil, fmt.Errorf("select pending character avatar upload: %w", gorm.ErrRecordNotFound)
	}
	return p, nil
}

func (f *fakeCharacterAvatarUploadRepository) Delete(_ context.Context, id uuid.UUID) error {
	f.lastDeleted = id
	f.deletedIDs = append(f.deletedIDs, id)
	delete(f.pending, id)
	return nil
}

func (f *fakeCharacterAvatarUploadRepository) ListStale(_ context.Context, _ time.Time) ([]*domain.PendingCharacterAvatarUpload, error) {
	return f.stale, nil
}

type fakeStorageService struct {
	lastKey        string
	deletedKeys    []string
	presignURL     string
	presignErr     error
	deleteErr      error
}

func (f *fakeStorageService) GeneratePresignedPutURL(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	f.lastKey = key
	if f.presignErr != nil {
		return "", f.presignErr
	}
	url := f.presignURL
	if url == "" {
		url = "https://example.com/presigned"
	}
	return url, nil
}

func (f *fakeStorageService) DeleteObject(_ context.Context, key string) error {
	f.deletedKeys = append(f.deletedKeys, key)
	return f.deleteErr
}

type fakeTransactor struct{}

func (f *fakeTransactor) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type fakeImageRepository struct {
	images            []*domain.Image
	listByCharacterID []*domain.Image
	lastListFilters   usecase.ListImageFilters
}

func (f *fakeImageRepository) GetByID(_ context.Context, _ string, _ uuid.UUID) (*domain.Image, error) {
	return nil, fmt.Errorf("get image: %w", gorm.ErrRecordNotFound)
}

func (f *fakeImageRepository) List(_ context.Context, _ uuid.UUID, filters usecase.ListImageFilters) ([]*domain.Image, error) {
	f.lastListFilters = filters
	return f.images, nil
}

func (f *fakeImageRepository) Update(_ context.Context, _ string, _ uuid.UUID, _ usecase.UpdateImageParams) (*domain.Image, error) {
	return nil, fmt.Errorf("update image: %w", gorm.ErrRecordNotFound)
}

func (f *fakeImageRepository) Delete(_ context.Context, _ string, _ uuid.UUID) error {
	return nil
}

func (f *fakeImageRepository) ListByCharacterID(_ context.Context, _ uuid.UUID, _ uuid.UUID) ([]*domain.Image, error) {
	return f.listByCharacterID, nil
}
