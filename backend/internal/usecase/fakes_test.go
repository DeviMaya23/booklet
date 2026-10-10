package usecase_test

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	rivertype "github.com/riverqueue/river/rivertype"
	"gorm.io/gorm"
)

type fakeCharacterRepository struct {
	characters map[uuid.UUID]*domain.Character

	lastCreated          *domain.Character
	lastUpdated          usecase.UpdateCharacterParams
	lastUpdatedAvatarID  uuid.UUID
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

func (f *fakeCharacterRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Character, error) {
	c, ok := f.characters[id]
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

func (f *fakeCharacterRepository) Update(_ context.Context, id uuid.UUID, _ uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	f.lastUpdated = params
	c, ok := f.characters[id]
	if !ok {
		return nil, fmt.Errorf("update character: %w", gorm.ErrRecordNotFound)
	}
	c.Name = params.Name
	c.IsPublic = params.IsPublic
	if params.AvatarR2Path != nil {
		c.AvatarR2Path = params.AvatarR2Path
	}
	c.Notes = params.Notes
	c.Folders = params.Folders
	return c, nil
}

func (f *fakeCharacterRepository) Delete(_ context.Context, id uuid.UUID, _ uuid.UUID) error {
	if _, ok := f.characters[id]; !ok {
		return fmt.Errorf("delete character: %w", gorm.ErrRecordNotFound)
	}
	delete(f.characters, id)
	return nil
}

func (f *fakeCharacterRepository) ClearAvatarR2Path(_ context.Context, id uuid.UUID, userID uuid.UUID) (string, error) {
	c, ok := f.characters[id]
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

func (f *fakeCharacterRepository) UpdateAvatarR2Path(_ context.Context, id uuid.UUID, _ uuid.UUID, r2Key string) error {
	c, ok := f.characters[id]
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
	lastKey     string
	deletedKeys []string
	presignURL  string
	presignErr  error
	deleteErr   error
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

type fakeBookleafClient struct {
	folderList *bookleaf.FolderList
	err        error
	called     bool
	lastUserID string
}

func (f *fakeBookleafClient) GetPublicFolders(_ context.Context, userID string) (*bookleaf.FolderList, error) {
	f.called = true
	f.lastUserID = userID
	return f.folderList, f.err
}

func (f *fakeBookleafClient) GetFolderImages(_ context.Context, _, _ string) (*bookleaf.FolderImageList, error) {
	return &bookleaf.FolderImageList{Images: []bookleaf.FolderImage{}}, nil
}

func (f *fakeBookleafClient) DeleteAccount(_ context.Context, _ string) error {
	return nil
}

type spyTransactor struct{}

func (s *spyTransactor) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type spyJobInserter struct {
	lastArgs  river.JobArgs
	returnErr error
}

func (s *spyJobInserter) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	s.lastArgs = args
	return &rivertype.JobInsertResult{}, s.returnErr
}

type spyStorageService struct{}

func (s *spyStorageService) GeneratePresignedPutURL(_ context.Context, _ string, _ string, _ time.Duration) (string, error) {
	return "https://presigned.example.com/put", nil
}

func (s *spyStorageService) DeleteObject(_ context.Context, _ string) error {
	return nil
}
