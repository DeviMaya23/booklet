package worker_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"strings"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/worker"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
)

// --- spies ---

type spyThumbnailImageRepo struct {
	imageToReturn       *domain.Image
	lastUpdatedID       uuid.UUID
	lastUpdatedPath     string
}

func (s *spyThumbnailImageRepo) GetByIDForWorker(_ context.Context, id uuid.UUID) (*domain.Image, error) {
	return s.imageToReturn, nil
}

func (s *spyThumbnailImageRepo) UpdateThumbnailPath(_ context.Context, id uuid.UUID, r2Path string) error {
	s.lastUpdatedID = id
	s.lastUpdatedPath = r2Path
	return nil
}

type spyThumbnailStorage struct {
	getErr      error
	getImage    image.Image // if nil, returns a default 100x200 image
	putKey      string
	putBody     []byte
	putMimeType string
}

func makeJPEGBody(img image.Image) io.ReadCloser {
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	return io.NopCloser(&buf)
}

func (s *spyThumbnailStorage) GetObject(_ context.Context, _ string) (io.ReadCloser, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	src := s.getImage
	if src == nil {
		src = image.NewRGBA(image.Rect(0, 0, 100, 200))
		for y := 0; y < 200; y++ {
			for x := 0; x < 100; x++ {
				src.(*image.RGBA).Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
			}
		}
	}
	return makeJPEGBody(src), nil
}

func (s *spyThumbnailStorage) PutObject(_ context.Context, key string, body io.Reader, contentType string) error {
	s.putKey = key
	data, _ := io.ReadAll(body)
	s.putBody = data
	s.putMimeType = contentType
	return nil
}

func makeJob(imageID, userID uuid.UUID) *river.Job[worker.GenerateThumbnailArgs] {
	return &river.Job[worker.GenerateThumbnailArgs]{
		Args: worker.GenerateThumbnailArgs{ImageID: imageID, UserID: userID},
	}
}

// --- tests ---

func TestGenerateThumbnailWorker_UpdatesThumbnailPath(t *testing.T) {
	imageID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyThumbnailImageRepo{
		imageToReturn: &domain.Image{
			ID:          imageID,
			ImageR2Path: "users/" + userID.String() + "/images/" + imageID.String() + ".jpg",
		},
	}
	storageSpy := &spyThumbnailStorage{}

	w := worker.NewGenerateThumbnailWorker(repoSpy, storageSpy)
	err := w.Work(context.Background(), makeJob(imageID, userID))

	require.NoError(t, err)
	expectedKey := "users/" + userID.String() + "/thumbnails/" + imageID.String() + ".jpg"
	require.Equal(t, expectedKey, storageSpy.putKey)
	require.Equal(t, expectedKey, repoSpy.lastUpdatedPath)
	require.Equal(t, imageID, repoSpy.lastUpdatedID)
	require.True(t, strings.HasSuffix(storageSpy.putKey, ".jpg"))
	require.Equal(t, "image/jpeg", storageSpy.putMimeType)
}

func TestGenerateThumbnailWorker_ReturnsErrorOnR2GetFailure(t *testing.T) {
	imageID := uuid.New()
	userID := uuid.New()

	repoSpy := &spyThumbnailImageRepo{
		imageToReturn: &domain.Image{
			ID:          imageID,
			ImageR2Path: "users/" + userID.String() + "/images/" + imageID.String() + ".jpg",
		},
	}
	storageSpy := &spyThumbnailStorage{getErr: errors.New("r2 unavailable")}

	w := worker.NewGenerateThumbnailWorker(repoSpy, storageSpy)
	err := w.Work(context.Background(), makeJob(imageID, userID))

	require.Error(t, err)
	require.Contains(t, err.Error(), "get original from r2")
}

func TestGenerateThumbnailWorker_ResizesLargeImageTo600pxLongEdge(t *testing.T) {
	imageID := uuid.New()
	userID := uuid.New()

	// 1200×900: long edge is 1200px, should become 600px; short edge 450px
	src := image.NewRGBA(image.Rect(0, 0, 1200, 900))
	repoSpy := &spyThumbnailImageRepo{
		imageToReturn: &domain.Image{
			ID:          imageID,
			ImageR2Path: "users/" + userID.String() + "/images/" + imageID.String() + ".jpg",
		},
	}
	storageSpy := &spyThumbnailStorage{getImage: src}

	w := worker.NewGenerateThumbnailWorker(repoSpy, storageSpy)
	err := w.Work(context.Background(), makeJob(imageID, userID))

	require.NoError(t, err)
	decoded, decErr := imaging.Decode(bytes.NewReader(storageSpy.putBody))
	require.NoError(t, decErr)
	bounds := decoded.Bounds()
	longEdge := max(bounds.Dx(), bounds.Dy())
	require.Equal(t, 600, longEdge, "long edge should be exactly 600px")
}

func TestGenerateThumbnailWorker_DoesNotUpscaleSmallImage(t *testing.T) {
	imageID := uuid.New()
	userID := uuid.New()

	// 100×200: both edges < 600px, should not be upscaled
	src := image.NewRGBA(image.Rect(0, 0, 100, 200))
	repoSpy := &spyThumbnailImageRepo{
		imageToReturn: &domain.Image{
			ID:          imageID,
			ImageR2Path: "users/" + userID.String() + "/images/" + imageID.String() + ".jpg",
		},
	}
	storageSpy := &spyThumbnailStorage{getImage: src}

	w := worker.NewGenerateThumbnailWorker(repoSpy, storageSpy)
	err := w.Work(context.Background(), makeJob(imageID, userID))

	require.NoError(t, err)
	decoded, decErr := imaging.Decode(bytes.NewReader(storageSpy.putBody))
	require.NoError(t, decErr)
	bounds := decoded.Bounds()
	require.LessOrEqual(t, bounds.Dx(), 100, "width should not exceed original")
	require.LessOrEqual(t, bounds.Dy(), 200, "height should not exceed original")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
