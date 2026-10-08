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

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type spySavedFilterUsecase struct {
	createResult *domain.SavedFilter
	createErr    error

	getByIDResult *domain.SavedFilter
	getByIDErr    error

	listResult []*domain.SavedFilter
	listErr    error

	updateResult *domain.SavedFilter
	updateErr    error

	deleteErr error
}

func (s *spySavedFilterUsecase) Create(_ context.Context, _ uuid.UUID, _ usecase.CreateSavedFilterParams) (*domain.SavedFilter, error) {
	return s.createResult, s.createErr
}

func (s *spySavedFilterUsecase) GetByID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.SavedFilter, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *spySavedFilterUsecase) List(_ context.Context, _ uuid.UUID) ([]*domain.SavedFilter, error) {
	return s.listResult, s.listErr
}

func (s *spySavedFilterUsecase) Update(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ usecase.UpdateSavedFilterParams) (*domain.SavedFilter, error) {
	return s.updateResult, s.updateErr
}

func (s *spySavedFilterUsecase) Delete(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return s.deleteErr
}

func makeSavedFilter() *domain.SavedFilter {
	return &domain.SavedFilter{
		ID:            uuid.New(),
		UserID:        testUserID,
		Name:          "My Filter",
		FilterPayload: json.RawMessage(`{"artists":[]}`),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// --- CreateSavedFilter ---

func TestCreateSavedFilter_MissingName_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	body := `{"filter_payload":{"artists":[]}}`
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateSavedFilter_HappyPath_Returns201WithBody(t *testing.T) {
	sf := makeSavedFilter()
	spy := &spySavedFilterUsecase{createResult: sf}
	h := handler.NewSavedFilterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	body := `{"name":"My Filter","filter_payload":{"artists":[]}}`
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, sf.ID.String(), got["id"])
	require.Equal(t, "My Filter", got["name"])
}

// --- GetSavedFilter ---

func TestGetSavedFilter_NotFound_Returns404(t *testing.T) {
	spy := &spySavedFilterUsecase{getByIDErr: usecase.ErrSavedFilterNotFound}
	h := handler.NewSavedFilterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/saved_filters/:id", h.GetSavedFilter)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/saved_filters/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

// --- PatchSavedFilter ---

func TestPatchSavedFilter_NotFound_Returns404(t *testing.T) {
	spy := &spySavedFilterUsecase{updateErr: usecase.ErrSavedFilterNotFound}
	h := handler.NewSavedFilterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/saved_filters/:id", h.PatchSavedFilter)

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/saved_filters/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

// --- DeleteSavedFilter ---

func TestDeleteSavedFilter_NotFound_Returns404(t *testing.T) {
	spy := &spySavedFilterUsecase{deleteErr: usecase.ErrSavedFilterNotFound}
	h := handler.NewSavedFilterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/saved_filters/:id", h.DeleteSavedFilter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/saved_filters/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteSavedFilter_HappyPath_Returns204(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/saved_filters/:id", h.DeleteSavedFilter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/saved_filters/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}
