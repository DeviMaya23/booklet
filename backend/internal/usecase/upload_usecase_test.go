package usecase_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// --- spies ---

type spyStorageService struct {
	lastKey string
}

func (s *spyStorageService) GeneratePresignedPutURL(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	s.lastKey = key
	return "https://example.com/presigned", nil
}

type spyUploadRepository struct {
	pendingToReturn *domain.PendingUpload
	lastDeletedID   uuid.UUID
}

func (s *spyUploadRepository) Create(_ context.Context, p *domain.PendingUpload) (*domain.PendingUpload, error) {
	return p, nil
}

func (s *spyUploadRepository) GetByID(_ context.Context, _ uuid.UUID, _ string) (*domain.PendingUpload, error) {
	return s.pendingToReturn, nil
}

func (s *spyUploadRepository) Delete(_ context.Context, id uuid.UUID) error {
	s.lastDeletedID = id
	return nil
}

type spyUploadCharacterRepository struct {
	charsToReturn []domain.Character
	lastIDs       []uuid.UUID
}

func (s *spyUploadCharacterRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, _ string) ([]domain.Character, error) {
	s.lastIDs = ids
	return s.charsToReturn, nil
}

type spyUploadImageRepository struct {
	lastImage *domain.Image
}

func (s *spyUploadImageRepository) Create(_ context.Context, image *domain.Image) (*domain.Image, error) {
	s.lastImage = image
	return image, nil
}

type spyTransactor struct{}

func (s *spyTransactor) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

// --- tests ---

func TestInitialUpload_R2KeyFormat(t *testing.T) {
	storageSpy := &spyStorageService{}
	repoSpy := &spyUploadRepository{}
	charSpy := &spyUploadCharacterRepository{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, imageSpy, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	result, err := uc.InitialUpload(context.Background(), usecase.InitialUploadParams{
		UserID:   "user-1",
		MimeType: "image/jpeg",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	expectedPrefix := "users/user-1/images/"
	expectedSuffix := ".jpg"
	require.True(t, strings.HasPrefix(storageSpy.lastKey, expectedPrefix),
		"expected key to start with %q, got %q", expectedPrefix, storageSpy.lastKey)
	require.True(t, strings.HasSuffix(storageSpy.lastKey, expectedSuffix),
		"expected key to end with %q, got %q", expectedSuffix, storageSpy.lastKey)

	middle := storageSpy.lastKey[len(expectedPrefix) : len(storageSpy.lastKey)-len(expectedSuffix)]
	_, err = uuid.Parse(middle)
	require.NoError(t, err, "expected UUID segment in key, got %q", middle)
}

func TestCompleteUpload_SomeCharsValid(t *testing.T) {
	charID1, charID2 := uuid.New(), uuid.New()

	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:           uuid.New(),
			CharacterIDs: []uuid.UUID{charID1, charID2},
		},
	}
	charSpy := &spyUploadCharacterRepository{
		charsToReturn: []domain.Character{{ID: charID1}},
	}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, charSpy, imageSpy, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), "user-1")

	require.NoError(t, err)
	require.NotNil(t, imageSpy.lastImage)
	require.Len(t, charSpy.lastIDs, 2)
	require.Len(t, imageSpy.lastImage.Characters, 1)
	require.Equal(t, charID1, imageSpy.lastImage.Characters[0].ID)
}

func TestCompleteUpload_AllCharsInvalid(t *testing.T) {
	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:           uuid.New(),
			CharacterIDs: []uuid.UUID{uuid.New()},
		},
	}
	charSpy := &spyUploadCharacterRepository{
		charsToReturn: []domain.Character{},
	}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, charSpy, imageSpy, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), "user-1")

	require.NoError(t, err)
	require.Empty(t, imageSpy.lastImage.Characters)
}
