package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func newFileUploadTestUsecase(
	pendingRepo *fakeFileUploadPendingRepository,
	artpieceRepo *fakeFileUploadArtpieceRepository,
	fileRepo *fakeFileUploadFileRepository,
	metaRepo *fakeImageMetadataRepository,
	coverUpdater *spyCoverUpdater,
	storage *spyStorageService,
) *usecase.FileUploadUsecase {
	return usecase.NewFileUploadUsecase(
		pendingRepo,
		artpieceRepo,
		fileRepo,
		metaRepo,
		coverUpdater,
		storage,
		observability.NewTelemetry(nil, nil, nil),
	)
}

func TestInitiateFileUpload_ImageType_ReturnsThumbnailUploadURL(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	result, err := uc.InitiateUpload(context.Background(), usecase.InitiateFileUploadParams{
		UserID:   userID,
		MimeType: "image/jpeg",
	})
	require.NoError(t, err)
	assert.NotNil(t, result.ThumbnailUploadURL, "image upload should return a thumbnail_upload_url")
	assert.NotNil(t, result.ThumbnailR2Key)
	assert.Contains(t, *result.ThumbnailR2Key, ".jpg")
}

func TestInitiateFileUpload_NonImageType_NoThumbnailUploadURL(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	result, err := uc.InitiateUpload(context.Background(), usecase.InitiateFileUploadParams{
		UserID:   userID,
		MimeType: "video/mp4",
	})
	require.NoError(t, err)
	assert.Nil(t, result.ThumbnailUploadURL, "non-image upload should not return a thumbnail_upload_url")
	assert.Nil(t, result.ThumbnailR2Key)
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

	artpieceRepo.artpieces[artpieceID] = &domain.Artpiece{ID: artpieceID, UserID: userID}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:         fileID,
		UserID:     userID,
		R2Key:      "users/x/files/f.jpg",
		MimeType:   "image/jpeg",
		ArtpieceID: &artpieceID,
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID, usecase.CompleteFileUploadParams{})
	require.NoError(t, err)
	require.True(t, coverUpdater.called, "cover should have been set")
	require.Equal(t, artpieceID, coverUpdater.lastArtpieceID)
	require.Equal(t, file.ID, *coverUpdater.lastCoverID)
}

func TestCompleteFileUpload_ExtensionCopiedFromPending(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	ext := "clip"
	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:        fileID,
		UserID:    userID,
		R2Key:     "users/x/files/f.bin",
		MimeType:  "application/octet-stream",
		Extension: &ext,
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID, usecase.CompleteFileUploadParams{})
	require.NoError(t, err)
	require.NotNil(t, file.Extension)
	require.Equal(t, "clip", *file.Extension)
}

func TestCompleteFileUpload_NonImageType_NoThumbnailAndNoMetadata(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:       fileID,
		UserID:   userID,
		R2Key:    "users/x/files/f.mp4",
		MimeType: "video/mp4",
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID, usecase.CompleteFileUploadParams{})
	require.NoError(t, err)
	assert.Nil(t, file.ThumbnailR2Path, "non-image file should have no thumbnail path")
	assert.Empty(t, metaRepo.created, "no image_metadata row should be created for non-image")
}

func TestCompleteFileUpload_ImageType_SetsThumbnailPathAndMetadata(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:       fileID,
		UserID:   userID,
		R2Key:    "users/x/files/f.jpg",
		MimeType: "image/jpeg",
	}

	w, h := 1920, 1080
	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID, usecase.CompleteFileUploadParams{Width: &w, Height: &h, ThumbnailUploaded: true})
	require.NoError(t, err)
	require.NotNil(t, file.ThumbnailR2Path)
	assert.Equal(t, "users/"+userID.String()+"/thumbnails/"+fileID.String()+".jpg", *file.ThumbnailR2Path)
	require.Len(t, metaRepo.created, 1)
	assert.Equal(t, 1920, metaRepo.created[0].Width)
	assert.Equal(t, 1080, metaRepo.created[0].Height)
}

func TestCompleteFileUpload_AlphaImageType_SetsPNGThumbnailPath(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:       fileID,
		UserID:   userID,
		R2Key:    "users/x/files/f.png",
		MimeType: "image/png",
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID, usecase.CompleteFileUploadParams{ThumbnailUploaded: true})
	require.NoError(t, err)
	require.NotNil(t, file.ThumbnailR2Path)
	assert.Equal(t, "users/"+userID.String()+"/thumbnails/"+fileID.String()+".png", *file.ThumbnailR2Path)
}

func TestCompleteFileUpload_ImageType_ThumbnailNotUploaded_NoThumbnailPath(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	fileID := uuid.New()
	pendingRepo.pending[fileID] = &domain.PendingFileUpload{
		ID:       fileID,
		UserID:   userID,
		R2Key:    "users/x/files/f.heic",
		MimeType: "image/heic",
	}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	file, err := uc.CompleteUpload(context.Background(), fileID, userID, usecase.CompleteFileUploadParams{ThumbnailUploaded: false})
	require.NoError(t, err)
	assert.Nil(t, file.ThumbnailR2Path, "thumbnail path should not be set when thumbnail was not uploaded")
}

func TestInitiateUpload_FileSizeExceedsLimit_ReturnsErrFileTooLarge(t *testing.T) {
	userID := uuid.New()

	pendingRepo := newFakeFileUploadPendingRepository()
	artpieceRepo := newFakeFileUploadArtpieceRepository()
	fileRepo := newFakeFileUploadFileRepository()
	metaRepo := &fakeImageMetadataRepository{}
	coverUpdater := &spyCoverUpdater{}
	storage := &spyStorageService{}

	uc := newFileUploadTestUsecase(pendingRepo, artpieceRepo, fileRepo, metaRepo, coverUpdater, storage)

	// image limit is 50 MB; 51 MB exceeds it
	ext := "png"
	_, err := uc.InitiateUpload(context.Background(), usecase.InitiateFileUploadParams{
		UserID:        userID,
		MimeType:      "image/png",
		Extension:     &ext,
		FileSizeBytes: 51 * 1024 * 1024,
	})

	require.ErrorIs(t, err, usecase.ErrFileTooLarge)
}
