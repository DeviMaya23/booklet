package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/devi/booklet/internal/worker"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	rivertype "github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- spies ---

type spyStorageService struct {
	lastKey           string
	lastDeletedKey    string
	deleteObjectCalls []string
}

func (s *spyStorageService) GeneratePresignedPutURL(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	s.lastKey = key
	return "https://example.com/presigned", nil
}

func (s *spyStorageService) DeleteObject(_ context.Context, key string) error {
	s.lastDeletedKey = key
	s.deleteObjectCalls = append(s.deleteObjectCalls, key)
	return nil
}

type spyUploadRepository struct {
	pendingToReturn *domain.PendingUpload
	staleToReturn   []*domain.PendingUpload
	lastDeletedID   uuid.UUID
	deletedIDs      []uuid.UUID
	listStaleCalled bool
}

func (s *spyUploadRepository) Create(_ context.Context, p *domain.PendingUpload) (*domain.PendingUpload, error) {
	return p, nil
}

func (s *spyUploadRepository) GetByID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.PendingUpload, error) {
	return s.pendingToReturn, nil
}

func (s *spyUploadRepository) Delete(_ context.Context, id uuid.UUID) error {
	s.lastDeletedID = id
	s.deletedIDs = append(s.deletedIDs, id)
	return nil
}

func (s *spyUploadRepository) ListStale(_ context.Context, _ time.Time) ([]*domain.PendingUpload, error) {
	s.listStaleCalled = true
	return s.staleToReturn, nil
}

type spyUploadCharacterRepository struct {
	charsToReturn []domain.Character
	lastIDs       []uuid.UUID
}

func (s *spyUploadCharacterRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, _ uuid.UUID) ([]domain.Character, error) {
	s.lastIDs = ids
	return s.charsToReturn, nil
}

type spyUploadArtistRepository struct {
	artistToReturn *domain.Artist
	returnErr      error
}

func (s *spyUploadArtistRepository) GetByIDAndUserID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Artist, error) {
	return s.artistToReturn, s.returnErr
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

type spyJobInserter struct {
	lastArgs river.JobArgs
	returnErr error
}

func (s *spyJobInserter) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	s.lastArgs = args
	return &rivertype.JobInsertResult{}, s.returnErr
}

// --- tests ---

func TestInitialUpload_R2KeyFormat(t *testing.T) {
	storageSpy := &spyStorageService{}
	repoSpy := &spyUploadRepository{}
	charSpy := &spyUploadCharacterRepository{}
	artistSpy := &spyUploadArtistRepository{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	userID := uuid.New()
	result, err := uc.InitialUpload(context.Background(), usecase.InitialUploadParams{
		UserID:   userID,
		MimeType: "image/jpeg",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	expectedPrefix := fmt.Sprintf("users/%s/images/", userID.String())
	expectedSuffix := ".jpg"
	require.True(t, strings.HasPrefix(storageSpy.lastKey, expectedPrefix),
		"expected key to start with %q, got %q", expectedPrefix, storageSpy.lastKey)
	require.True(t, strings.HasSuffix(storageSpy.lastKey, expectedSuffix),
		"expected key to end with %q, got %q", expectedSuffix, storageSpy.lastKey)

	middle := storageSpy.lastKey[len(expectedPrefix) : len(storageSpy.lastKey)-len(expectedSuffix)]
	_, err = uuid.Parse(middle)
	require.NoError(t, err, "expected UUID segment in key, got %q", middle)
}

func TestInitialUpload_ArtistIDNotOwnedReturnsError(t *testing.T) {
	artistSpy := &spyUploadArtistRepository{
		returnErr: fmt.Errorf("get artist: %w", gorm.ErrRecordNotFound),
	}
	artistID := uuid.New()

	uc := usecase.NewUploadUsecase(&spyUploadRepository{}, &spyStorageService{}, &spyUploadCharacterRepository{}, artistSpy, &spyUploadImageRepository{}, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.InitialUpload(context.Background(), usecase.InitialUploadParams{
		UserID:   uuid.New(),
		MimeType: "image/jpeg",
		ArtistID: &artistID,
	})

	require.ErrorIs(t, err, usecase.ErrArtistNotOwned)
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
	artistSpy := &spyUploadArtistRepository{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, charSpy, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), uuid.New())

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
	artistSpy := &spyUploadArtistRepository{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, charSpy, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
	require.Empty(t, imageSpy.lastImage.Characters)
}

func TestCompleteUpload_ArtistIDCarriedThroughIfValid(t *testing.T) {
	artistID := uuid.New()
	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:       uuid.New(),
			ArtistID: &artistID,
		},
	}
	artistSpy := &spyUploadArtistRepository{
		artistToReturn: &domain.Artist{ID: artistID},
	}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, &spyUploadCharacterRepository{}, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
	require.NotNil(t, imageSpy.lastImage.ArtistID)
	require.Equal(t, artistID, *imageSpy.lastImage.ArtistID)
}

func TestCompleteUpload_ArtistIDNulledIfNotFound(t *testing.T) {
	artistID := uuid.New()
	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:       uuid.New(),
			ArtistID: &artistID,
		},
	}
	artistSpy := &spyUploadArtistRepository{
		returnErr: fmt.Errorf("get artist: %w", gorm.ErrRecordNotFound),
	}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, &spyUploadCharacterRepository{}, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
	require.Nil(t, imageSpy.lastImage.ArtistID)
}

func TestCleanupStaleUploads_NoStaleRecords(t *testing.T) {
	repoSpy := &spyUploadRepository{staleToReturn: nil}
	storageSpy := &spyStorageService{}
	charSpy := &spyUploadCharacterRepository{}
	artistSpy := &spyUploadArtistRepository{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CleanupStaleUploads(context.Background(), 15*time.Minute)

	require.NoError(t, err)
	require.Empty(t, storageSpy.deleteObjectCalls)
	require.Empty(t, repoSpy.deletedIDs)
}

func TestCompleteUpload_EnqueuesThumbnailJobOnSuccess(t *testing.T) {
	imageID := uuid.New()
	userID := uuid.New()
	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{
			ID:     imageID,
			UserID: userID,
		},
	}
	jobSpy := &spyJobInserter{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, &spyUploadCharacterRepository{}, &spyUploadArtistRepository{}, imageSpy, &spyTransactor{}, jobSpy, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), imageID, userID)

	require.NoError(t, err)
	require.NotNil(t, jobSpy.lastArgs, "expected thumbnail job to be enqueued")
	args, ok := jobSpy.lastArgs.(worker.GenerateThumbnailArgs)
	require.True(t, ok, "expected GenerateThumbnailArgs")
	require.Equal(t, imageID, args.ImageID)
	require.Equal(t, userID, args.UserID)
}

func TestCompleteUpload_EnqueueFailureDoesNotFail(t *testing.T) {
	repoSpy := &spyUploadRepository{
		pendingToReturn: &domain.PendingUpload{ID: uuid.New()},
	}
	jobSpy := &spyJobInserter{returnErr: errors.New("river unavailable")}

	uc := usecase.NewUploadUsecase(repoSpy, &spyStorageService{}, &spyUploadCharacterRepository{}, &spyUploadArtistRepository{}, &spyUploadImageRepository{}, &spyTransactor{}, jobSpy, observability.NewTelemetry(nil, nil, nil))

	err := uc.CompleteUpload(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
}

func TestCleanupStaleUploads_StaleRecordsExist(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	stale := []*domain.PendingUpload{
		{ID: id1, R2Key: "users/u1/images/a.jpg"},
		{ID: id2, R2Key: "users/u1/images/b.jpg"},
	}
	repoSpy := &spyUploadRepository{staleToReturn: stale}
	storageSpy := &spyStorageService{}
	charSpy := &spyUploadCharacterRepository{}
	artistSpy := &spyUploadArtistRepository{}
	imageSpy := &spyUploadImageRepository{}

	uc := usecase.NewUploadUsecase(repoSpy, storageSpy, charSpy, artistSpy, imageSpy, &spyTransactor{}, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CleanupStaleUploads(context.Background(), 15*time.Minute)

	require.NoError(t, err)
	require.Equal(t, []string{"users/u1/images/a.jpg", "users/u1/images/b.jpg"}, storageSpy.deleteObjectCalls)
	require.Equal(t, []uuid.UUID{id1, id2}, repoSpy.deletedIDs)
}
