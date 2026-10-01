package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyFileUploadUsecase struct {
	initiateResult *usecase.InitiateFileUploadResult
	initiateErr    error

	completeResult *domain.File
	completeErr    error
}

func (s *spyFileUploadUsecase) InitiateUpload(_ context.Context, _ usecase.InitiateFileUploadParams) (*usecase.InitiateFileUploadResult, error) {
	return s.initiateResult, s.initiateErr
}

func (s *spyFileUploadUsecase) CompleteUpload(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.File, error) {
	return s.completeResult, s.completeErr
}

func makeFile() *domain.File {
	return &domain.File{
		ID:         uuid.New(),
		UserID:     testUserID,
		FileR2Path: "users/x/files/abc.mp4",
		MimeType:   "video/mp4",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// --- InitiateUpload ---

func TestInitiateFileUpload_HappyPath(t *testing.T) {
	pendingID := uuid.New()
	spy := &spyFileUploadUsecase{
		initiateResult: &usecase.InitiateFileUploadResult{
			ID:        pendingID,
			UploadURL: "https://example.com/presigned",
			ExpiresAt: time.Now().Add(15 * time.Minute),
		},
	}
	h := handler.NewFileUploadHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/files", h.InitiateUpload)

	body := `{"mime_type":"video/mp4"}`
	req := httptest.NewRequest(http.MethodPost, "/files", strings.NewReader(body))
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

func TestInitiateFileUpload_MissingMimeType(t *testing.T) {
	spy := &spyFileUploadUsecase{}
	h := handler.NewFileUploadHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/files", h.InitiateUpload)

	req := httptest.NewRequest(http.MethodPost, "/files", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestInitiateFileUpload_ArtpieceNotOwned(t *testing.T) {
	spy := &spyFileUploadUsecase{initiateErr: usecase.ErrArtpieceNotOwned}
	h := handler.NewFileUploadHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/files", h.InitiateUpload)

	body := `{"mime_type":"video/mp4","artpiece_id":"00000000-0000-0000-0000-000000000099"}`
	req := httptest.NewRequest(http.MethodPost, "/files", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// --- CompleteUpload ---

func TestCompleteFileUpload_HappyPath(t *testing.T) {
	f := makeFile()
	spy := &spyFileUploadUsecase{completeResult: f}
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/signed"}
	h := handler.NewFileUploadHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/files/:id/complete", h.CompleteUpload)

	req := httptest.NewRequest(http.MethodPost, "/files/"+f.ID.String()+"/complete", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, f.ID.String(), got["id"])
	require.Equal(t, "video/mp4", got["mime_type"])
}

func TestCompleteFileUpload_InvalidUUID(t *testing.T) {
	spy := &spyFileUploadUsecase{}
	h := handler.NewFileUploadHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/files/:id/complete", h.CompleteUpload)

	req := httptest.NewRequest(http.MethodPost, "/files/not-a-uuid/complete", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCompleteFileUpload_NotFound(t *testing.T) {
	spy := &spyFileUploadUsecase{completeErr: gorm.ErrRecordNotFound}
	h := handler.NewFileUploadHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/files/:id/complete", h.CompleteUpload)

	req := httptest.NewRequest(http.MethodPost, "/files/"+uuid.New().String()+"/complete", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
