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
	"gorm.io/gorm"
)

type artistResponseBody struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Notes      *string `json:"notes"`
	ArtistLink *string `json:"artist_link"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type spyArtistUsecase struct {
	createResult *domain.Artist
	createErr    error

	getByIDResult *domain.Artist
	getByIDErr    error

	listResult []*domain.Artist
	listErr    error

	updateResult *domain.Artist
	updateErr    error

	deleteErr error
}

func (s *spyArtistUsecase) Create(_ context.Context, _ uuid.UUID, _ usecase.CreateArtistParams) (*domain.Artist, error) {
	return s.createResult, s.createErr
}

func (s *spyArtistUsecase) GetByID(_ context.Context, _ string, _ uuid.UUID) (*domain.Artist, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *spyArtistUsecase) List(_ context.Context, _ uuid.UUID, _ usecase.ListArtistFilters) ([]*domain.Artist, error) {
	return s.listResult, s.listErr
}

func (s *spyArtistUsecase) Update(_ context.Context, _ string, _ uuid.UUID, _ usecase.UpdateArtistParams) (*domain.Artist, error) {
	return s.updateResult, s.updateErr
}

func (s *spyArtistUsecase) Delete(_ context.Context, _ string, _ uuid.UUID) error {
	return s.deleteErr
}

func makeArtist() *domain.Artist {
	return &domain.Artist{
		ID:        uuid.New(),
		UserID:    testUserID,
		Name:      "Jane Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// --- CreateArtist ---

func TestCreateArtist_HappyPath(t *testing.T) {
	artist := makeArtist()
	spy := &spyArtistUsecase{createResult: artist}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artists", h.CreateArtist)

	body := `{"name":"Jane Doe"}`
	req := httptest.NewRequest(http.MethodPost, "/artists", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got artistResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, artist.ID.String(), got.ID)
	require.Equal(t, "Jane Doe", got.Name)
}

func TestCreateArtist_MissingName(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artists", h.CreateArtist)

	req := httptest.NewRequest(http.MethodPost, "/artists", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateArtist_Conflict(t *testing.T) {
	spy := &spyArtistUsecase{createErr: usecase.ErrArtistNameConflict}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artists", h.CreateArtist)

	req := httptest.NewRequest(http.MethodPost, "/artists", strings.NewReader(`{"name":"Jane Doe"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateArtist_MalformedJSON(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artists", h.CreateArtist)

	req := httptest.NewRequest(http.MethodPost, "/artists", strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateArtist_InvalidArtistLink(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artists", h.CreateArtist)

	req := httptest.NewRequest(http.MethodPost, "/artists", strings.NewReader(`{"name":"Jane","artist_link":"not-a-url"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// --- ListArtists ---

func TestListArtists_HappyPath(t *testing.T) {
	artists := []*domain.Artist{makeArtist(), makeArtist()}
	spy := &spyArtistUsecase{listResult: artists}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artists", h.ListArtists)

	req := httptest.NewRequest(http.MethodGet, "/artists", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []artistResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 2)
}

func TestListArtists_QFilterBoundToUsecase(t *testing.T) {
	artists := []*domain.Artist{makeArtist()}
	spy := &spyArtistUsecase{listResult: artists}

	var capturedFilters usecase.ListArtistFilters
	spy2 := &captureArtistListSpy{inner: spy, capture: &capturedFilters}
	h := handler.NewArtistHandler(spy2, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artists", h.ListArtists)

	req := httptest.NewRequest(http.MethodGet, "/artists?q=jane", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, capturedFilters.Q)
	require.Equal(t, "jane", *capturedFilters.Q)
}

type captureArtistListSpy struct {
	inner   *spyArtistUsecase
	capture *usecase.ListArtistFilters
}

func (s *captureArtistListSpy) Create(ctx context.Context, userID uuid.UUID, params usecase.CreateArtistParams) (*domain.Artist, error) {
	return s.inner.Create(ctx, userID, params)
}
func (s *captureArtistListSpy) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Artist, error) {
	return s.inner.GetByID(ctx, id, userID)
}
func (s *captureArtistListSpy) List(ctx context.Context, userID uuid.UUID, filters usecase.ListArtistFilters) ([]*domain.Artist, error) {
	*s.capture = filters
	return s.inner.List(ctx, userID, filters)
}
func (s *captureArtistListSpy) Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateArtistParams) (*domain.Artist, error) {
	return s.inner.Update(ctx, id, userID, params)
}
func (s *captureArtistListSpy) Delete(ctx context.Context, id string, userID uuid.UUID) error {
	return s.inner.Delete(ctx, id, userID)
}

// --- GetArtistByID ---

func TestGetArtistByID_HappyPath(t *testing.T) {
	artist := makeArtist()
	spy := &spyArtistUsecase{getByIDResult: artist}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artists/:id", h.GetArtistByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/artists/%s", artist.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got artistResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, artist.ID.String(), got.ID)
}

func TestGetArtistByID_NotFound(t *testing.T) {
	spy := &spyArtistUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artists/:id", h.GetArtistByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/artists/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetArtistByID_InvalidUUID(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artists/:id", h.GetArtistByID)

	req := httptest.NewRequest(http.MethodGet, "/artists/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- UpdateArtist ---

func TestUpdateArtist_HappyPath(t *testing.T) {
	artist := makeArtist()
	spy := &spyArtistUsecase{updateResult: artist}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/artists/:id", h.UpdateArtist)

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/artists/%s", artist.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got artistResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, artist.ID.String(), got.ID)
}

func TestUpdateArtist_Conflict(t *testing.T) {
	spy := &spyArtistUsecase{updateErr: usecase.ErrArtistNameConflict}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/artists/:id", h.UpdateArtist)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/artists/%s", uuid.New()), strings.NewReader(`{"name":"Taken Name"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestUpdateArtist_NotFound(t *testing.T) {
	spy := &spyArtistUsecase{updateErr: gorm.ErrRecordNotFound}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/artists/:id", h.UpdateArtist)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/artists/%s", uuid.New()), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateArtist_MalformedJSON(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/artists/:id", h.UpdateArtist)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/artists/%s", uuid.New()), strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateArtist_InvalidArtistLink(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/artists/:id", h.UpdateArtist)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/artists/%s", uuid.New()), strings.NewReader(`{"artist_link":"not-a-url"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// --- DeleteArtist ---

func TestDeleteArtist_HappyPath(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artists/:id", h.DeleteArtist)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/artists/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteArtist_NotFound(t *testing.T) {
	spy := &spyArtistUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artists/:id", h.DeleteArtist)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/artists/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteArtist_InvalidUUID(t *testing.T) {
	spy := &spyArtistUsecase{}
	h := handler.NewArtistHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artists/:id", h.DeleteArtist)

	req := httptest.NewRequest(http.MethodDelete, "/artists/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
