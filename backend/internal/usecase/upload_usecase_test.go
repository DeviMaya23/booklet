package usecase_test

import (
	"context"
	"fmt"
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
	pendingToReturn  *domain.PendingUpload
	lastValidCharIDs []uuid.UUID
}

func (s *spyUploadRepository) CreatePendingUpload(_ context.Context, p *domain.PendingUpload) error {
	return nil
}

func (s *spyUploadRepository) GetPendingUpload(_ context.Context, _, _ string) (*domain.PendingUpload, error) {
	return s.pendingToReturn, nil
}

func (s *spyUploadRepository) CompleteUpload(_ context.Context, _, _ string, validCharIDs []uuid.UUID) (*domain.Image, error) {
	s.lastValidCharIDs = validCharIDs
	return &domain.Image{ID: uuid.New(), Characters: []domain.Character{}}, nil
}

type spyUploadCharacterRepository struct {
	charsToReturn []domain.Character
}

func (s *spyUploadCharacterRepository) GetByIDsAndUserID(_ context.Context, _ []string, _ string) ([]domain.Character, error) {
	return s.charsToReturn, nil
}

// --- tests ---

func TestInitialUpload_R2KeyFormat(t *testing.T) {
	storageSpy := &spyStorageService{}
	repoSpy := &spyUploadRepository{}
	charSpy := &spyUploadCharacterRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, observability.NewTelemetry(nil, nil, nil))

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

	// The middle segment should be a valid UUID
	middle := storageSpy.lastKey[len(expectedPrefix) : len(storageSpy.lastKey)-len(expectedSuffix)]
	_, err = uuid.Parse(middle)
	require.NoError(t, err, "expected UUID segment in key, got %q", middle)
}

func TestCompleteUpload_SomeCharsValid(t *testing.T) {
	charID1 := uuid.New()
	charID2 := uuid.New()

	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:           uuid.New(),
			CharacterIDs: []string{charID1.String(), charID2.String()},
		},
	}
	charSpy := &spyUploadCharacterRepository{
		// Only charID1 is valid/owned
		charsToReturn: []domain.Character{
			{ID: charID1, UserID: "user-1", Name: fmt.Sprintf("Char %s", charID1)},
		},
	}
	storageSpy := &spyStorageService{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.CompleteUpload(context.Background(), uuid.NewString(), "user-1")

	require.NoError(t, err)
	require.Len(t, repoSpy.lastValidCharIDs, 1)
	require.Equal(t, charID1, repoSpy.lastValidCharIDs[0])
}

func TestCompleteUpload_AllCharsInvalid(t *testing.T) {
	charID1 := uuid.New()

	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:           uuid.New(),
			CharacterIDs: []string{charID1.String()},
		},
	}
	charSpy := &spyUploadCharacterRepository{
		// No valid chars returned
		charsToReturn: []domain.Character{},
	}
	storageSpy := &spyStorageService{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.CompleteUpload(context.Background(), uuid.NewString(), "user-1")

	require.NoError(t, err)
	require.Empty(t, repoSpy.lastValidCharIDs)
}
