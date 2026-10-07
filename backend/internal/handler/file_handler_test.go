package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyFileUsecase struct {
	bulkDeleteErr error
	bulkDeleteIDs []uuid.UUID
	getByIDFile   *domain.File
	getByIDErr    error
}

func (s *spyFileUsecase) GetByID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.File, error) {
	return s.getByIDFile, s.getByIDErr
}

func (s *spyFileUsecase) List(_ context.Context, _ uuid.UUID, _ bool) ([]*domain.File, error) {
	return nil, nil
}

func (s *spyFileUsecase) Update(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ *string, _ *string) (*domain.File, error) {
	return nil, nil
}

func (s *spyFileUsecase) Delete(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

func (s *spyFileUsecase) BulkDelete(_ context.Context, ids []uuid.UUID, _ uuid.UUID) error {
	s.bulkDeleteIDs = ids
	return s.bulkDeleteErr
}

func newFileHandler(spy *spyFileUsecase) *handler.FileHandler {
	return handler.NewFileHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))
}

// --- BulkDeleteFiles ---

func TestBulkDeleteFiles_HappyPath(t *testing.T) {
	id1 := uuid.New()
	spy := &spyFileUsecase{}
	h := newFileHandler(spy)
	e := setupEcho(testUserID)
	e.DELETE("/files", h.BulkDeleteFiles)

	body := map[string]any{"ids": []string{id1.String()}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodDelete, "/files", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	require.Len(t, spy.bulkDeleteIDs, 1)
	assert.Equal(t, id1, spy.bulkDeleteIDs[0])
}

func TestBulkDeleteFiles_EmptyIDs(t *testing.T) {
	h := newFileHandler(&spyFileUsecase{})
	e := setupEcho(testUserID)
	e.DELETE("/files", h.BulkDeleteFiles)

	body := map[string]any{"ids": []string{}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodDelete, "/files", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteFiles_OwnershipViolation(t *testing.T) {
	spy := &spyFileUsecase{bulkDeleteErr: usecase.ErrFileOwnershipViolation}
	h := newFileHandler(spy)
	e := setupEcho(testUserID)
	e.DELETE("/files", h.BulkDeleteFiles)

	body := map[string]any{"ids": []string{uuid.New().String()}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodDelete, "/files", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteFiles_BindError(t *testing.T) {
	h := newFileHandler(&spyFileUsecase{})
	e := setupEcho(testUserID)
	e.DELETE("/files", h.BulkDeleteFiles)

	req := httptest.NewRequest(http.MethodDelete, "/files", strings.NewReader("not json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- GetDownloadURL ---

func TestGetDownloadURL_Success(t *testing.T) {
	fileID := uuid.New()
	name := "my-file"
	spy := &spyFileUsecase{
		getByIDFile: &domain.File{
			ID:         fileID,
			UserID:     testUserID,
			FileR2Path: "users/files/my-file.psd",
			Name:       &name,
		},
	}
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/my-file.psd?sig=dl"}
	h := handler.NewFileHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))
	e := setupEcho(testUserID)
	e.GET("/files/:id/download", h.GetDownloadURL)

	req := httptest.NewRequest(http.MethodGet, "/files/"+fileID.String()+"/download", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "https://cdn.example.com/my-file.psd?sig=dl", resp["download_url"])
}

func TestGetDownloadURL_InvalidUUID(t *testing.T) {
	h := handler.NewFileHandler(&spyFileUsecase{}, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))
	e := setupEcho(testUserID)
	e.GET("/files/:id/download", h.GetDownloadURL)

	req := httptest.NewRequest(http.MethodGet, "/files/not-a-uuid/download", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetDownloadURL_NotFound(t *testing.T) {
	spy := &spyFileUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := handler.NewFileHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))
	e := setupEcho(testUserID)
	e.GET("/files/:id/download", h.GetDownloadURL)

	req := httptest.NewRequest(http.MethodGet, "/files/"+uuid.New().String()+"/download", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// suppress unused import
var _ = errors.New
