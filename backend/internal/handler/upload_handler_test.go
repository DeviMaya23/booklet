package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyUploadUsecase struct {
	initialUploadResult *usecase.InitialUploadResult
	initialUploadErr    error

	completeUploadErr error
}

func (s *spyUploadUsecase) InitialUpload(_ context.Context, _ usecase.InitialUploadParams) (*usecase.InitialUploadResult, error) {
	return s.initialUploadResult, s.initialUploadErr
}

func (s *spyUploadUsecase) CompleteUpload(_ context.Context, _ uuid.UUID, _ string) error {
	return s.completeUploadErr
}

// --- InitialUpload ---

func TestInitialUpload_HappyPath(t *testing.T) {
	pendingID := uuid.New()
	spy := &spyUploadUsecase{
		initialUploadResult: &usecase.InitialUploadResult{
			ID:        pendingID,
			UploadURL: "https://example.com/presigned",
			ExpiresAt: time.Now().Add(15 * time.Minute),
		},
	}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images", h.InitialUpload)

	body := `{"mime_type":"image/jpeg","title":"Test Image"}`
	req := httptest.NewRequest(http.MethodPost, "/images", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.NotEmpty(t, got["id"])
	require.NotEmpty(t, got["upload_url"])
	require.NotEmpty(t, got["expires_at"])
}

func TestInitialUpload_MissingMimeType(t *testing.T) {
	spy := &spyUploadUsecase{}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images", h.InitialUpload)

	body := `{"title":"Test Image"}`
	req := httptest.NewRequest(http.MethodPost, "/images", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestInitialUpload_InvalidMimeType(t *testing.T) {
	spy := &spyUploadUsecase{}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images", h.InitialUpload)

	body := `{"mime_type":"image/gif"}`
	req := httptest.NewRequest(http.MethodPost, "/images", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestInitialUpload_MalformedJSON(t *testing.T) {
	spy := &spyUploadUsecase{}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images", h.InitialUpload)

	req := httptest.NewRequest(http.MethodPost, "/images", strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- CompleteUpload ---

func TestCompleteUpload_HappyPath(t *testing.T) {
	spy := &spyUploadUsecase{}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images/:id/complete", h.CompleteUpload)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/images/%s/complete", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Empty(t, rec.Body.Bytes())
}

func TestCompleteUpload_NotFound(t *testing.T) {
	spy := &spyUploadUsecase{completeUploadErr: gorm.ErrRecordNotFound}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images/:id/complete", h.CompleteUpload)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/images/%s/complete", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCompleteUpload_InvalidUUID(t *testing.T) {
	spy := &spyUploadUsecase{}
	h := handler.NewUploadHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/images/:id/complete", h.CompleteUpload)

	req := httptest.NewRequest(http.MethodPost, "/images/not-a-uuid/complete", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
