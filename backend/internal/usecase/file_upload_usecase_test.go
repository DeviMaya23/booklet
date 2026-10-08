package usecase_test

import (
	"context"
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

// --- fakes ---

type fakeFileUploadPendingRepository struct {
	pending map[uuid.UUID]*domain.PendingFileUpload
}

func newFakeFileUploadPendingRepository() *fakeFileUploadPendingRepository {
	return &fakeFileUploadPendingRepository{pending: make(map[uuid.UUID]*domain.PendingFileUpload)}
}

func (f *fakeFileUploadPendingRepository) Create(_ context.Context, p *domain.PendingFileUpload) (*domain.PendingFileUpload, error) {
	f.pending[p.ID] = p
	return p, nil
}

func (f *fakeFileUploadPendingRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingFileUpload, error) {
	p, ok := f.pending[id]
	if !ok || p.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}

func (f *fakeFileUploadPendingRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(f.pending, id)
	return nil
}

func (f *fakeFileUploadPendingRepository) ListStale(_ context.Context, _ time.Time) ([]*domain.PendingFileUpload, error) {
	return nil, nil
}

type fakeFileUploadFileRepository struct {
	files             map[uuid.UUID]*domain.File
	lastArtpieceID    *uuid.UUID
	lastUpdatedFileID uuid.UUID
}

func newFakeFileUploadFileRepository() *fakeFileUploadFileRepository {
	return &fakeFileUploadFileRepository{files: make(map[uuid.UUID]*domain.File)}
}

func (f *fakeFileUploadFileRepository) Create(_ context.Context, file *domain.File) (*domain.File, error) {
	f.files[file.ID] = file
	return file, nil
}

func (f *fakeFileUploadFileRepository) GetFilesForArtpiece(_ context.Context, artpieceID uuid.UUID) ([]*domain.File, error) {
	var result []*domain.File
	for _, file := range f.files {
		if file.ArtpieceID != nil && *file.ArtpieceID == artpieceID {
			result = append(result, file)
		}
	}
	return result, nil
}

func (f *fakeFileUploadFileRepository) UpdateArtpieceID(_ context.Context, fileID uuid.UUID, artpieceID *uuid.UUID) error {
	f.lastUpdatedFileID = fileID
	f.lastArtpieceID = artpieceID
	if file, ok := f.files[fileID]; ok {
		file.ArtpieceID = artpieceID
	}
	return nil
}

type fakeFileUploadArtpieceRepository struct {
	artpieces map[uuid.UUID]*domain.Artpiece
}

func newFakeFileUploadArtpieceRepository() *fakeFileUploadArtpieceRepository {
	return &fakeFileUploadArtpieceRepository{artpieces: make(map[uuid.UUID]*domain.Artpiece)}
}

func (f *fakeFileUploadArtpieceRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	a, ok := f.artpieces[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

type fakeImageMetadataRepository struct {
	created []*domain.ImageMetadata
}

func (f *fakeImageMetadataRepository) CreateImageMetadata(_ context.Context, m *domain.ImageMetadata) error {
	f.created = append(f.created, m)
	return nil
}

type spyCoverUpdater struct {
	lastArtpieceID uuid.UUID
	lastCoverID    *uuid.UUID
	called         bool
}

func (s *spyCoverUpdater) UpdateCover(_ context.Context, artpieceID uuid.UUID, coverFileID *uuid.UUID) error {
	s.called = true
	s.lastArtpieceID = artpieceID
	s.lastCoverID = coverFileID
	return nil
}

type spyFileJobInserter struct {
	insertedArgs []worker.GenerateFileThumbnailArgs
}

func (s *spyFileJobInserter) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	if a, ok := args.(worker.GenerateFileThumbnailArgs); ok {
		s.insertedArgs = append(s.insertedArgs, a)
	}
	return &rivertype.JobInsertResult{}, nil
}

func newFileUploadTestUsecase(
	pendingRepo *fakeFileUploadPendingRepository,
	artpieceRepo *fakeFileUploadArtpieceRepository,
	fileRepo *fakeFileUploadFileRepository,
	metaRepo *fakeImageMetadataRepository,
	coverUpdater *spyCoverUpdater,
	storage *spyStorageService,
	jobInserter *spyFileJobInserter,
) *usecase.FileUploadUsecase {
	return usecase.NewFileUploadUsecase(
		pendingRepo,
		artpieceRepo,
		fileRepo,
		metaRepo,
		coverUpdater,
		storage,
		jobInserter,
		observability.NewTelemetry(nil, nil, nil),
	)
}

func TestCompleteFileUpload_AutoSetsCoverOnFirstImageAttach(t *testing.T) {
	userID := uuid.New()
	artpieceID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}
	jobInserter := &spyFileJobInserter{}

	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:         fileID,
		UserID:     userID,
		R2Key:      "users/x/files/f.jpg",
		MimeType:   "image/jpeg",
		ArtpieceID: &artpieceID,
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage, jobInserter)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID)
	require.NoError(t, err)
	require.True(t, coverUpdater.called, "cover should have been set")
	require.Equal(t, artpieceID, coverUpdater.lastArtpieceID)
	require.Equal(t, file.ID, *coverUpdater.lastCoverID)
}

func TestCompleteFileUpload_NoThumbnailEnqueuedForNonImageType(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}
	jobInserter := &spyFileJobInserter{}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:       fileID,
		UserID:   userID,
		R2Key:    "users/x/files/f.mp4",
		MimeType: "video/mp4",
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage, jobInserter)

	_, err := uc.CompleteUpload(context.Background(), fileID, userID)
	require.NoError(t, err)
	require.Empty(t, jobInserter.insertedArgs, "no thumbnail job should be enqueued for non-image")
	require.Empty(t, metaRepo.created, "no image_metadata row should be created for non-image")
}
