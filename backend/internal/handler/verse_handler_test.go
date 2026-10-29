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

type verseResponseBody struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type spyVerseUsecase struct {
	createResult *domain.Verse
	createErr    error

	getByIDResult *domain.Verse
	getByIDErr    error

	listResult []*domain.Verse
	listErr    error

	updateResult *domain.Verse
	updateErr    error

	deleteErr  error
	reorderErr error
}

func (s *spyVerseUsecase) Create(_ context.Context, _ uuid.UUID, _ usecase.CreateVerseParams) (*domain.Verse, error) {
	return s.createResult, s.createErr
}

func (s *spyVerseUsecase) GetByID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Verse, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *spyVerseUsecase) List(_ context.Context, _ uuid.UUID, _ usecase.ListVerseFilters) ([]*domain.Verse, error) {
	return s.listResult, s.listErr
}

func (s *spyVerseUsecase) Update(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ usecase.UpdateVerseParams) (*domain.Verse, error) {
	return s.updateResult, s.updateErr
}

func (s *spyVerseUsecase) Delete(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return s.deleteErr
}

func (s *spyVerseUsecase) Reorder(_ context.Context, _ uuid.UUID, _ []uuid.UUID) error {
	return s.reorderErr
}

func makeVerse() *domain.Verse {
	return &domain.Verse{
		ID:        uuid.New(),
		UserID:    testUserID,
		Name:      "Solaris",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newVerseHandler(spy *spyVerseUsecase) *handler.VerseHandler {
	return handler.NewVerseHandler(spy, observability.NewTelemetry(nil, nil, nil))
}

func setupEchoUnauthenticated() *echo.Echo {
	e := echo.New()
	e.Validator = handler.NewEchoValidator()
	return e
}

// --- CreateVerse ---

func TestCreateVerse_HappyPath(t *testing.T) {
	verse := makeVerse()
	spy := &spyVerseUsecase{createResult: verse}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses", h.CreateVerse)

	req := httptest.NewRequest(http.MethodPost, "/verses", strings.NewReader(`{"name":"Solaris"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got verseResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, verse.ID.String(), got.ID)
	require.Equal(t, "Solaris", got.Name)
}

func TestCreateVerse_MissingName(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses", h.CreateVerse)

	req := httptest.NewRequest(http.MethodPost, "/verses", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateVerse_EmptyName(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses", h.CreateVerse)

	req := httptest.NewRequest(http.MethodPost, "/verses", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateVerse_Conflict(t *testing.T) {
	spy := &spyVerseUsecase{createErr: usecase.ErrVerseNameConflict}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses", h.CreateVerse)

	req := httptest.NewRequest(http.MethodPost, "/verses", strings.NewReader(`{"name":"Solaris"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateVerse_MalformedJSON(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses", h.CreateVerse)

	req := httptest.NewRequest(http.MethodPost, "/verses", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- ListVerses ---

func TestListVerses_HappyPath(t *testing.T) {
	verses := []*domain.Verse{makeVerse(), makeVerse()}
	spy := &spyVerseUsecase{listResult: verses}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/verses", h.ListVerses)

	req := httptest.NewRequest(http.MethodGet, "/verses", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []verseResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 2)
}

// --- GetVerseByID ---

func TestGetVerseByID_HappyPath(t *testing.T) {
	verse := makeVerse()
	spy := &spyVerseUsecase{getByIDResult: verse}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/verses/:id", h.GetVerseByID)

	req := httptest.NewRequest(http.MethodGet, "/verses/"+verse.ID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got verseResponseBody
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, verse.ID.String(), got.ID)
}

func TestGetVerseByID_NotFound(t *testing.T) {
	spy := &spyVerseUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/verses/:id", h.GetVerseByID)

	req := httptest.NewRequest(http.MethodGet, "/verses/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetVerseByID_InvalidID(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/verses/:id", h.GetVerseByID)

	req := httptest.NewRequest(http.MethodGet, "/verses/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- UpdateVerse ---

func TestUpdateVerse_HappyPath(t *testing.T) {
	verse := makeVerse()
	spy := &spyVerseUsecase{updateResult: verse}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.PATCH("/verses/:id", h.UpdateVerse)

	req := httptest.NewRequest(http.MethodPatch, "/verses/"+verse.ID.String(), strings.NewReader(`{"name":"Terran"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateVerse_EmptyName(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.PATCH("/verses/:id", h.UpdateVerse)

	req := httptest.NewRequest(http.MethodPatch, "/verses/"+uuid.New().String(), strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateVerse_Conflict(t *testing.T) {
	spy := &spyVerseUsecase{updateErr: usecase.ErrVerseNameConflict}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.PATCH("/verses/:id", h.UpdateVerse)

	req := httptest.NewRequest(http.MethodPatch, "/verses/"+uuid.New().String(), strings.NewReader(`{"name":"Solaris"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestUpdateVerse_NotFound(t *testing.T) {
	spy := &spyVerseUsecase{updateErr: gorm.ErrRecordNotFound}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.PATCH("/verses/:id", h.UpdateVerse)

	req := httptest.NewRequest(http.MethodPatch, "/verses/"+uuid.New().String(), strings.NewReader(`{"name":"Terran"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

// --- DeleteVerse ---

func TestDeleteVerse_HappyPath(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/verses/:id", h.DeleteVerse)

	req := httptest.NewRequest(http.MethodDelete, "/verses/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteVerse_NotFound(t *testing.T) {
	spy := &spyVerseUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/verses/:id", h.DeleteVerse)

	req := httptest.NewRequest(http.MethodDelete, "/verses/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

// --- ReorderVerses ---

func TestReorderVerses_HappyPath(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses/reorder", h.ReorderVerses)

	body := `{"ids":["` + uuid.New().String() + `","` + uuid.New().String() + `"]}`
	req := httptest.NewRequest(http.MethodPost, "/verses/reorder", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestReorderVerses_IDsMismatch_Returns422(t *testing.T) {
	spy := &spyVerseUsecase{reorderErr: usecase.ErrVerseReorderIDsMismatch}
	h := newVerseHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/verses/reorder", h.ReorderVerses)

	body := `{"ids":["` + uuid.New().String() + `"]}`
	req := httptest.NewRequest(http.MethodPost, "/verses/reorder", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestReorderVerses_Unauthorized(t *testing.T) {
	spy := &spyVerseUsecase{}
	h := newVerseHandler(spy)

	e := setupEchoUnauthenticated()
	e.POST("/verses/reorder", h.ReorderVerses)

	body := `{"ids":["` + uuid.New().String() + `"]}`
	req := httptest.NewRequest(http.MethodPost, "/verses/reorder", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
