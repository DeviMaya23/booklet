package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type characterResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	HeroImageR2Path *string `json:"hero_image_r2_path"`
	Biography       *string `json:"biography"`
	IsPublic        bool    `json:"is_public"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// spyCharacterUsecase is a value-return spy for CharacterUsecase.
type spyCharacterUsecase struct {
	createResult *domain.Character
	createErr    error

	getByIDResult *domain.Character
	getByIDErr    error

	listResult []*domain.Character
	listErr    error

	updateResult *domain.Character
	updateErr    error

	deleteErr error

	lastCreateParams usecase.CreateCharacterParams
	lastUpdateParams usecase.UpdateCharacterParams
}

func (s *spyCharacterUsecase) Create(_ context.Context, _ string, params usecase.CreateCharacterParams) (*domain.Character, error) {
	s.lastCreateParams = params
	return s.createResult, s.createErr
}

func (s *spyCharacterUsecase) GetByID(_ context.Context, _, _ string) (*domain.Character, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *spyCharacterUsecase) List(_ context.Context, _ string) ([]*domain.Character, error) {
	return s.listResult, s.listErr
}

func (s *spyCharacterUsecase) Update(_ context.Context, _, _ string, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	s.lastUpdateParams = params
	return s.updateResult, s.updateErr
}

func (s *spyCharacterUsecase) Delete(_ context.Context, _, _ string) error {
	return s.deleteErr
}

func setupEcho(userID string) *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(string(middleware.AuthenticatedUserIDContextKey), userID)
			return next(c)
		}
	})
	return e
}

func makeCharacter() *domain.Character {
	return &domain.Character{
		ID:     uuid.New(),
		UserID: "user-1",
		Name:   "Aria",
	}
}

// --- CreateCharacter ---

func TestCreateCharacter_HappyPath(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{createResult: character}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/characters", h.CreateCharacter)

	body := `{"name":"Aria"}`
	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, character.ID.String(), got.ID)
	require.Equal(t, "Aria", spy.lastCreateParams.Name)
}

func TestCreateCharacter_MissingName(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/characters", h.CreateCharacter)

	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCharacter_MalformedJSON(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/characters", h.CreateCharacter)

	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(`{bad json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCharacter_UsecaseError(t *testing.T) {
	spy := &spyCharacterUsecase{createErr: errors.New("db down")}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.POST("/characters", h.CreateCharacter)

	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(`{"name":"Aria"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- GetCharacterByID ---

func TestGetCharacterByID_HappyPath(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{getByIDResult: character}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s", character.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, character.ID.String(), got.ID)
}

func TestGetCharacterByID_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCharacterByID_InvalidUUID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, "/characters/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- ListCharacters ---

func TestListCharacters_HappyPath(t *testing.T) {
	characters := []*domain.Character{makeCharacter(), makeCharacter()}
	spy := &spyCharacterUsecase{listResult: characters}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.GET("/characters", h.ListCharacters)

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 2)
}

// --- UpdateCharacter ---

func TestUpdateCharacter_HappyPath(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{updateResult: character}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.PATCH("/characters/:id", h.UpdateCharacter)

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", character.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, character.ID.String(), got.ID)
}

func TestUpdateCharacter_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{updateErr: gorm.ErrRecordNotFound}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.PATCH("/characters/:id", h.UpdateCharacter)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", uuid.New().String()), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCharacter_MalformedJSON(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.PATCH("/characters/:id", h.UpdateCharacter)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", uuid.New().String()), strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- DeleteCharacter ---

func TestDeleteCharacter_HappyPath(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.DELETE("/characters/:id", h.DeleteCharacter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteCharacter_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.DELETE("/characters/:id", h.DeleteCharacter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
