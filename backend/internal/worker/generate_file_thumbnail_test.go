package worker_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/worker"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// --- spies ---

type spyFileThumbnailRepo struct {
	fileToReturn         *domain.File
	updatedThumbnailPath string
	updatedGenState      string
	updatedWidth         int
	updatedHeight        int
}

func (s *spyFileThumbnailRepo) GetByIDForWorker(_ context.Context, _ uuid.UUID) (*domain.File, error) {
	return s.fileToReturn, nil
}

func (s *spyFileThumbnailRepo) UpdateThumbnailPath(_ context.Context, _ uuid.UUID, r2Path string) error {
	s.updatedThumbnailPath = r2Path
	return nil
}

func (s *spyFileThumbnailRepo) UpdateThumbnailGenState(ctx context.Context, _ uuid.UUID, state string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	s.updatedGenState = state
	return nil
}

func (s *spyFileThumbnailRepo) UpdateImageMetadataDimensions(_ context.Context, _ uuid.UUID, width, height int) error {
	s.updatedWidth = width
	s.updatedHeight = height
	return nil
}

type spyFileThumbnailStorage struct {
	getErr          error
	getReader       func() io.ReadCloser
	putKey          string
	putContentType  string
	putCacheControl string
}

func (s *spyFileThumbnailStorage) GetObject(_ context.Context, _ string) (io.ReadCloser, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.getReader != nil {
		return s.getReader(), nil
	}
	// default: a small opaque JPEG
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			src.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, src, nil)
	return io.NopCloser(&buf), nil
}

func (s *spyFileThumbnailStorage) PutObject(_ context.Context, key string, _ io.Reader, contentType, cacheControl string) error {
	s.putKey = key
	s.putContentType = contentType
	s.putCacheControl = cacheControl
	return nil
}

func makeFileJob(fileID, userID uuid.UUID, attempt, maxAttempts int) *river.Job[worker.GenerateFileThumbnailArgs] {
	return &river.Job[worker.GenerateFileThumbnailArgs]{
		JobRow: &rivertype.JobRow{
			Attempt:     attempt,
			MaxAttempts: maxAttempts,
		},
		Args: worker.GenerateFileThumbnailArgs{FileID: fileID, UserID: userID},
	}
}

func jpegReader(img image.Image) io.ReadCloser {
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	return io.NopCloser(&buf)
}

func pngReader(img image.Image) io.ReadCloser {
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return io.NopCloser(&buf)
}

// --- tests ---

func TestGenerateFileThumbnailWorker_SuccessSetsStateAndPath(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
			MimeType:   "image/jpeg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	err := w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25))

	require.NoError(t, err)
	expectedKey := "users/" + userID.String() + "/thumbnails/" + fileID.String() + ".jpg"
	require.Equal(t, expectedKey, repoSpy.updatedThumbnailPath)
	require.Equal(t, "done", repoSpy.updatedGenState)
}

func TestGenerateFileThumbnailWorker_WritesImageDimensions(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	src := image.NewRGBA(image.Rect(0, 0, 800, 600))
	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
			MimeType:   "image/jpeg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getReader: func() io.ReadCloser { return jpegReader(src) }}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	err := w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25))

	require.NoError(t, err)
	require.Equal(t, 800, repoSpy.updatedWidth)
	require.Equal(t, 600, repoSpy.updatedHeight)
}

func TestGenerateFileThumbnailWorker_FinalAttemptFailureSetsFailedAndReturnsNil(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getErr: errors.New("r2 unavailable")}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	err := w.Work(context.Background(), makeFileJob(fileID, userID, 25, 25))

	require.NoError(t, err)
	require.Equal(t, "failed", repoSpy.updatedGenState)
}

func TestGenerateFileThumbnailWorker_NonFinalAttemptFailureReturnsError(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getErr: errors.New("r2 unavailable")}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	err := w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25))

	require.Error(t, err)
	require.Contains(t, err.Error(), "get original from r2")
	require.Empty(t, repoSpy.updatedGenState)
}

func TestGenerateFileThumbnailWorker_UploadsThumbnailWithCacheControl(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
			MimeType:   "image/jpeg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	require.NoError(t, w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25)))

	require.Equal(t, "private, max-age=172800", storageSpy.putCacheControl)
}

func TestGenerateFileThumbnailWorker_JPEGSourceProducesJPEGThumbnail(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
			MimeType:   "image/jpeg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getReader: func() io.ReadCloser { return jpegReader(src) }}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	require.NoError(t, w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25)))

	require.True(t, strings.HasSuffix(storageSpy.putKey, ".jpg"), "key should end in .jpg, got: %s", storageSpy.putKey)
	require.Equal(t, "image/jpeg", storageSpy.putContentType)
}

func TestGenerateFileThumbnailWorker_PNGSourceProducesPNGThumbnail(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	src := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			src.Set(x, y, color.NRGBA{R: 0, G: 128, B: 255, A: 128}) // semi-transparent
		}
	}
	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".png",
			MimeType:   "image/png",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getReader: func() io.ReadCloser { return pngReader(src) }}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	require.NoError(t, w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25)))

	require.True(t, strings.HasSuffix(storageSpy.putKey, ".png"), "key should end in .png, got: %s", storageSpy.putKey)
	require.Equal(t, "image/png", storageSpy.putContentType)
}

func TestGenerateFileThumbnailWorker_FinalAttemptWithCancelledCtxSetsFailedState(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getErr: errors.New("r2 unavailable")}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	err := w.Work(ctx, makeFileJob(fileID, userID, 25, 25))

	require.NoError(t, err)
	require.Equal(t, "failed", repoSpy.updatedGenState)
}

func TestGenerateFileThumbnailWorker_DecodeTimeoutOnFinalAttemptSetsFailedAndReturnsNil(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	pr, _ := io.Pipe() // reader blocks; writer never writes
	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getReader: func() io.ReadCloser { return pr }}

	w := worker.NewGenerateFileThumbnailWorkerWithTimeout(repoSpy, storageSpy, zap.NewNop(), 1*time.Millisecond)
	err := w.Work(context.Background(), makeFileJob(fileID, userID, 25, 25))

	require.NoError(t, err)
	require.Equal(t, "failed", repoSpy.updatedGenState)
}

func TestGenerateFileThumbnailWorker_DecodeTimeoutOnNonFinalAttemptReturnsError(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	pr, _ := io.Pipe()
	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getReader: func() io.ReadCloser { return pr }}

	w := worker.NewGenerateFileThumbnailWorkerWithTimeout(repoSpy, storageSpy, zap.NewNop(), 1*time.Millisecond)
	err := w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25))

	require.Error(t, err)
	require.ErrorContains(t, err, "decode timeout")
}

func TestGenerateFileThumbnailWorker_WebPSourceGeneratesThumbnail(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	webpData, err := os.ReadFile("testdata/small.webp")
	require.NoError(t, err)

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:         fileID,
			FileR2Path: "users/" + userID.String() + "/files/" + fileID.String() + ".webp",
			MimeType:   "image/webp",
		},
	}
	storageSpy := &spyFileThumbnailStorage{getReader: func() io.ReadCloser {
		return io.NopCloser(bytes.NewReader(webpData))
	}}

	w := worker.NewGenerateFileThumbnailWorker(repoSpy, storageSpy, zap.NewNop())
	require.NoError(t, w.Work(context.Background(), makeFileJob(fileID, userID, 1, 25)))

	require.Equal(t, "done", repoSpy.updatedGenState)
	require.True(t, strings.HasSuffix(storageSpy.putKey, ".png"), "key should end in .png, got: %s", storageSpy.putKey)
}
