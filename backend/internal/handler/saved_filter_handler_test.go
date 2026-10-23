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

func makeSavedFilterWithThumbnail(r2Path string) *domain.SavedFilter {
	sf := makeSavedFilter()
	sf.ThumbnailR2Path = &r2Path
	return sf
}

// --- CreateSavedFilter ---

func TestCreateSavedFilter_MissingName_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

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
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

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

func TestCreateSavedFilter_DuplicateName_Returns422WithNameFieldError(t *testing.T) {
	spy := &spySavedFilterUsecase{createErr: usecase.ErrDuplicateSavedFilterName}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	body := `{"name":"Duplicate","filter_payload":{}}`
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	errors, ok := got["errors"].([]interface{})
	require.True(t, ok)
	require.Len(t, errors, 1)
	field := errors[0].(map[string]interface{})
	require.Equal(t, "name", field["field"])
}

func TestCreateSavedFilter_WithThumbnail_ResponseHasThumbnailURL(t *testing.T) {
	r2Path := "thumbnails/artpieces/abc.webp"
	sf := makeSavedFilterWithThumbnail(r2Path)
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/thumb?sig=xyz"}
	spy := &spySavedFilterUsecase{createResult: sf}
	h := handler.NewSavedFilterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	body := fmt.Sprintf(`{"name":"My Filter","thumbnail_artpiece_id":"%s","filter_payload":{}}`, uuid.New())
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Contains(t, got, "thumbnail_url")
	require.NotContains(t, got, "thumbnail_r2_path")
	require.NotNil(t, got["thumbnail_url"])
}

func TestCreateSavedFilter_NonObjectPayload_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	body := `{"name":"My Filter","filter_payload":[1,2,3]}`
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	errs, ok := got["errors"].([]interface{})
	require.True(t, ok)
	require.Len(t, errs, 1)
	require.Equal(t, "filter_payload", errs[0].(map[string]interface{})["field"])
}

func TestCreateSavedFilter_NullPayload_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	body := `{"name":"My Filter","filter_payload":null}`
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	errs, ok := got["errors"].([]interface{})
	require.True(t, ok)
	require.Len(t, errs, 1)
	require.Equal(t, "filter_payload", errs[0].(map[string]interface{})["field"])
}

func TestCreateSavedFilter_OversizedPayload_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/saved_filters", h.CreateSavedFilter)

	// build a payload object > 4096 bytes
	large := `{"name":"My Filter","filter_payload":{"k":"` + strings.Repeat("x", 4100) + `"}}`
	req := httptest.NewRequest(http.MethodPost, "/saved_filters", strings.NewReader(large))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	errs, ok := got["errors"].([]interface{})
	require.True(t, ok)
	require.Len(t, errs, 1)
	require.Equal(t, "filter_payload", errs[0].(map[string]interface{})["field"])
}

// --- ListSavedFilters ---

func TestListSavedFilters_ResponseHasThumbnailURL(t *testing.T) {
	r2Path := "thumbnails/artpieces/list.webp"
	sf := makeSavedFilterWithThumbnail(r2Path)
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/list-thumb?sig=xyz"}
	spy := &spySavedFilterUsecase{listResult: []*domain.SavedFilter{sf}}
	h := handler.NewSavedFilterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/saved_filters", h.ListSavedFilters)

	req := httptest.NewRequest(http.MethodGet, "/saved_filters", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	require.Contains(t, got[0], "thumbnail_url")
	require.NotContains(t, got[0], "thumbnail_r2_path")
	require.NotNil(t, got[0]["thumbnail_url"])
}

// --- GetSavedFilter ---

func TestGetSavedFilter_HappyPath_ResponseHasThumbnailURL(t *testing.T) {
	r2Path := "thumbnails/artpieces/get.webp"
	sf := makeSavedFilterWithThumbnail(r2Path)
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/get-thumb?sig=xyz"}
	spy := &spySavedFilterUsecase{getByIDResult: sf}
	h := handler.NewSavedFilterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/saved_filters/:id", h.GetSavedFilter)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/saved_filters/%s", sf.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Contains(t, got, "thumbnail_url")
	require.NotContains(t, got, "thumbnail_r2_path")
	require.NotNil(t, got["thumbnail_url"])
}

func TestGetSavedFilter_NotFound_Returns404(t *testing.T) {
	spy := &spySavedFilterUsecase{getByIDErr: usecase.ErrSavedFilterNotFound}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

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
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/saved_filters/:id", h.PatchSavedFilter)

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/saved_filters/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPatchSavedFilter_WithThumbnailArtpieceID_Returns200WithThumbnailURL(t *testing.T) {
	r2Path := "thumbnails/artpieces/patch.webp"
	sf := makeSavedFilterWithThumbnail(r2Path)
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/patch-thumb?sig=xyz"}
	spy := &spySavedFilterUsecase{updateResult: sf}
	h := handler.NewSavedFilterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/saved_filters/:id", h.PatchSavedFilter)

	body := fmt.Sprintf(`{"thumbnail_artpiece_id":"%s"}`, uuid.New())
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/saved_filters/%s", sf.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Contains(t, got, "thumbnail_url")
	require.NotNil(t, got["thumbnail_url"])
}

func TestPatchSavedFilter_WithThumbnailArtpieceID_ArtpieceHasNoThumbnail_Returns200WithNullURL(t *testing.T) {
	sf := makeSavedFilter()
	spy := &spySavedFilterUsecase{updateResult: sf}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/saved_filters/:id", h.PatchSavedFilter)

	body := fmt.Sprintf(`{"thumbnail_artpiece_id":"%s"}`, uuid.New())
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/saved_filters/%s", sf.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Contains(t, got, "thumbnail_url")
	require.Nil(t, got["thumbnail_url"])
}

func TestPatchSavedFilter_NonObjectPayload_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/saved_filters/:id", h.PatchSavedFilter)

	body := `{"filter_payload":"not an object"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/saved_filters/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	errs, ok := got["errors"].([]interface{})
	require.True(t, ok)
	require.Len(t, errs, 1)
	require.Equal(t, "filter_payload", errs[0].(map[string]interface{})["field"])
}

func TestPatchSavedFilter_OversizedPayload_Returns422(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/saved_filters/:id", h.PatchSavedFilter)

	large := `{"filter_payload":{"k":"` + strings.Repeat("x", 4100) + `"}}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/saved_filters/%s", uuid.New()), strings.NewReader(large))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	errs, ok := got["errors"].([]interface{})
	require.True(t, ok)
	require.Len(t, errs, 1)
	require.Equal(t, "filter_payload", errs[0].(map[string]interface{})["field"])
}

// --- DeleteSavedFilter ---

func TestDeleteSavedFilter_NotFound_Returns404(t *testing.T) {
	spy := &spySavedFilterUsecase{deleteErr: usecase.ErrSavedFilterNotFound}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/saved_filters/:id", h.DeleteSavedFilter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/saved_filters/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteSavedFilter_HappyPath_Returns204(t *testing.T) {
	spy := &spySavedFilterUsecase{}
	h := handler.NewSavedFilterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/saved_filters/:id", h.DeleteSavedFilter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/saved_filters/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}
