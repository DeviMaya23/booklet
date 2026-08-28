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

var (
	testUserID     = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	testIDPSubject = "kp_test_user_1"
)

type characterResponse struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	HeroImageR2Path *string  `json:"hero_image_r2_path"`
	Biography       *string  `json:"biography"`
	IsPublic        bool     `json:"is_public"`
	FolderIDs       []string `json:"folder_ids"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
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

func (s *spyCharacterUsecase) Create(_ context.Context, _ uuid.UUID, params usecase.CreateCharacterParams) (*domain.Character, error) {
	s.lastCreateParams = params
	return s.createResult, s.createErr
}

func (s *spyCharacterUsecase) GetByID(_ context.Context, _ string, _ uuid.UUID) (*domain.Character, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *spyCharacterUsecase) List(_ context.Context, _ uuid.UUID) ([]*domain.Character, error) {
	return s.listResult, s.listErr
}

func (s *spyCharacterUsecase) Update(_ context.Context, _ string, _ uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	s.lastUpdateParams = params
	return s.updateResult, s.updateErr
}

func (s *spyCharacterUsecase) Delete(_ context.Context, _ string, _ uuid.UUID) error {
	return s.deleteErr
}

func setupEcho(userID uuid.UUID) *echo.Echo {
	e := echo.New()
	e.Validator = handler.NewEchoValidator()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(string(middleware.AuthenticatedUserIDContextKey), userID)
			c.Set(string(middleware.AuthenticatedIDPSubjectContextKey), testIDPSubject)
			return next(c)
		}
	})
	return e
}

func makeCharacter() *domain.Character {
	return &domain.Character{
		ID:     uuid.New(),
		UserID: testUserID,
		Name:   "Aria",
	}
}

// --- CreateCharacter ---

func TestCreateCharacter_HappyPath(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{createResult: character}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
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
	require.Equal(t, []string{}, got.FolderIDs)
}

func TestCreateCharacter_MissingName(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/characters", h.CreateCharacter)

	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got struct {
		Errors []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.Errors, 1)
	require.Equal(t, "name", got.Errors[0].Field)
	require.Equal(t, "name is required", got.Errors[0].Message)
}

func TestCreateCharacter_MalformedJSON(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
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

	e := setupEcho(testUserID)
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

	e := setupEcho(testUserID)
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

	e := setupEcho(testUserID)
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCharacterByID_InvalidUUID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
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

	e := setupEcho(testUserID)
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

	e := setupEcho(testUserID)
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
	require.Equal(t, []string{}, got.FolderIDs)
}

func TestUpdateCharacter_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{updateErr: gorm.ErrRecordNotFound}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
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

	e := setupEcho(testUserID)
	e.PATCH("/characters/:id", h.UpdateCharacter)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", uuid.New().String()), strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateCharacter_EmptyName(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/characters/:id", h.UpdateCharacter)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", uuid.New().String()), strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got struct {
		Errors []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.Errors, 1)
	require.Equal(t, "name", got.Errors[0].Field)
	require.Equal(t, "name must not be empty", got.Errors[0].Message)
}

func TestUpdateCharacter_AbsentName(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{updateResult: character}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/characters/:id", h.UpdateCharacter)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", character.ID), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, spy.lastUpdateParams.Name)
}

// --- DeleteCharacter ---

func TestDeleteCharacter_HappyPath(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/characters/:id", h.DeleteCharacter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteCharacter_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/characters/:id", h.DeleteCharacter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateCharacter_DuplicateFolderIDs_Deduped(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{createResult: character}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/characters", h.CreateCharacter)

	folderID := uuid.New().String()
	body := fmt.Sprintf(`{"name":"Aria","folder_ids":["%s","%s"]}`, folderID, folderID)
	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, spy.lastCreateParams.FolderIDs)
	require.Len(t, *spy.lastCreateParams.FolderIDs, 1)
}

func TestCreateCharacter_InvalidFolderID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/characters", h.CreateCharacter)

	body := `{"name":"Aria","folder_ids":["not-a-uuid"]}`
	req := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got struct {
		Errors []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.NotEmpty(t, got.Errors)
}

func TestUpdateCharacter_InvalidFolderID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := handler.NewCharacterHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/characters/:id", h.UpdateCharacter)

	body := `{"folder_ids":["not-a-uuid"]}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/characters/%s", uuid.New().String()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var got struct {
		Errors []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.NotEmpty(t, got.Errors)
}
