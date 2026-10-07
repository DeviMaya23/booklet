package worker_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"testing"

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

func (s *spyFileThumbnailRepo) UpdateThumbnailGenState(_ context.Context, _ uuid.UUID, state string) error {
	s.updatedGenState = state
	return nil
}

func (s *spyFileThumbnailRepo) UpdateImageMetadataDimensions(_ context.Context, _ uuid.UUID, width, height int) error {
	s.updatedWidth = width
	s.updatedHeight = height
	return nil
}

type spyFileThumbnailStorage struct {
	getErr   error
	getImage image.Image
	putKey   string
}

func (s *spyFileThumbnailStorage) GetObject(_ context.Context, _ string) (io.ReadCloser, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	src := s.getImage
	if src == nil {
		src = image.NewRGBA(image.Rect(0, 0, 100, 100))
		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				src.(*image.RGBA).Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
			}
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, src, nil)
	return io.NopCloser(&buf), nil
}

func (s *spyFileThumbnailStorage) PutObject(_ context.Context, key string, _ io.Reader, _ string) error {
	s.putKey = key
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

// --- tests ---

func TestGenerateFileThumbnailWorker_SuccessSetsStateAndPath(t *testing.T) {
	fileID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyFileThumbnailRepo{
		fileToReturn: &domain.File{
			ID:          fileID,
			FileR2Path:  "users/" + userID.String() + "/files/" + fileID.String() + ".jpg",
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
		},
	}
	storageSpy := &spyFileThumbnailStorage{getImage: src}

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
