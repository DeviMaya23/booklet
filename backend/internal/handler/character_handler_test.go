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
	"time"

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
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	AvatarURL *string  `json:"avatar_url"`
	Biography *string  `json:"biography"`
	IsPublic  bool     `json:"is_public"`
	FolderIDs []string `json:"folder_ids"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type spyPresigner struct {
	presignedURL string
	presignErr   error
	calls        []string
}

func (s *spyPresigner) GeneratePresignedGetURL(_ context.Context, key string, _ time.Duration) (string, error) {
	s.calls = append(s.calls, key)
	if s.presignErr != nil {
		return "", s.presignErr
	}
	url := s.presignedURL
	if url == "" {
		url = "https://cdn.example.com/presigned?sig=abc"
	}
	return url, nil
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

	initAvatarUploadResult *usecase.AvatarUploadResult
	initAvatarUploadErr    error

	completeAvatarUploadErr error

	deleteAvatarErr error

	getCharacterImagesResult []*domain.Image
	getCharacterImagesErr    error

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

func (s *spyCharacterUsecase) List(_ context.Context, _ uuid.UUID, _ usecase.ListCharacterFilters) ([]*domain.Character, error) {
	return s.listResult, s.listErr
}

func (s *spyCharacterUsecase) Update(_ context.Context, _ string, _ uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	s.lastUpdateParams = params
	return s.updateResult, s.updateErr
}

func (s *spyCharacterUsecase) Delete(_ context.Context, _ string, _ uuid.UUID) error {
	return s.deleteErr
}

func (s *spyCharacterUsecase) InitAvatarUpload(_ context.Context, _ uuid.UUID, _ string, _ string) (*usecase.AvatarUploadResult, error) {
	return s.initAvatarUploadResult, s.initAvatarUploadErr
}

func (s *spyCharacterUsecase) CompleteAvatarUpload(_ context.Context, _ uuid.UUID, _ string, _ uuid.UUID) error {
	return s.completeAvatarUploadErr
}

func (s *spyCharacterUsecase) DeleteAvatar(_ context.Context, _ uuid.UUID, _ string) error {
	return s.deleteAvatarErr
}

func (s *spyCharacterUsecase) GetCharacterImages(_ context.Context, _ uuid.UUID, _ uuid.UUID) ([]*domain.Image, error) {
	return s.getCharacterImagesResult, s.getCharacterImagesErr
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

func newCharacterHandler(spy *spyCharacterUsecase) *handler.CharacterHandler {
	return handler.NewCharacterHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))
}

// --- CreateCharacter ---

func TestCreateCharacter_HappyPath(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{createResult: character}
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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

func TestGetCharacterByID_AvatarURLPresigned(t *testing.T) {
	avatarKey := "users/123/files/avatar.jpg"
	character := makeCharacter()
	character.AvatarR2Path = &avatarKey
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/avatar?sig=xyz"}
	spy := &spyCharacterUsecase{getByIDResult: character}
	h := handler.NewCharacterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s", character.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.NotNil(t, got.AvatarURL)
	require.Equal(t, "https://cdn.example.com/avatar?sig=xyz", *got.AvatarURL)
	require.Contains(t, presigner.calls, avatarKey)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &raw))
	_, hasRawPath := raw["avatar_r2_path"]
	require.False(t, hasRawPath)
}

func TestGetCharacterByID_AvatarURLNullWhenNoAvatar(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{getByIDResult: character}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s", character.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Nil(t, got.AvatarURL)
}

func TestGetCharacterByID_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/characters/:id", h.GetCharacterByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCharacterByID_InvalidUUID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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

func TestListCharacters_AvatarURLsPresigned(t *testing.T) {
	key1 := "users/1/files/a.jpg"
	key2 := "users/1/files/b.jpg"
	c1 := makeCharacter()
	c1.AvatarR2Path = &key1
	c2 := makeCharacter()
	c2.AvatarR2Path = &key2
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/signed"}
	spy := &spyCharacterUsecase{listResult: []*domain.Character{c1, c2}}
	h := handler.NewCharacterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/characters", h.ListCharacters)

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []characterResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 2)
	require.NotNil(t, got[0].AvatarURL)
	require.Equal(t, "https://cdn.example.com/signed", *got[0].AvatarURL)
	require.NotNil(t, got[1].AvatarURL)
	require.Contains(t, presigner.calls, key1)
	require.Contains(t, presigner.calls, key2)
}

func TestListCharacters_QFilterBoundToUsecase(t *testing.T) {
	characters := []*domain.Character{makeCharacter()}
	spy := &spyCharacterUsecase{listResult: characters}

	var capturedFilters usecase.ListCharacterFilters
	spy2 := &captureCharacterListSpy{inner: spy, capture: &capturedFilters}
	h := handler.NewCharacterHandler(spy2, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/characters", h.ListCharacters)

	req := httptest.NewRequest(http.MethodGet, "/characters?q=aria", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, capturedFilters.Q)
	require.Equal(t, "aria", *capturedFilters.Q)
}

type captureCharacterListSpy struct {
	inner   *spyCharacterUsecase
	capture *usecase.ListCharacterFilters
}

func (s *captureCharacterListSpy) Create(ctx context.Context, userID uuid.UUID, params usecase.CreateCharacterParams) (*domain.Character, error) {
	return s.inner.Create(ctx, userID, params)
}
func (s *captureCharacterListSpy) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Character, error) {
	return s.inner.GetByID(ctx, id, userID)
}
func (s *captureCharacterListSpy) List(ctx context.Context, userID uuid.UUID, filters usecase.ListCharacterFilters) ([]*domain.Character, error) {
	*s.capture = filters
	return s.inner.List(ctx, userID, filters)
}
func (s *captureCharacterListSpy) Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	return s.inner.Update(ctx, id, userID, params)
}
func (s *captureCharacterListSpy) Delete(ctx context.Context, id string, userID uuid.UUID) error {
	return s.inner.Delete(ctx, id, userID)
}
func (s *captureCharacterListSpy) InitAvatarUpload(ctx context.Context, userID uuid.UUID, characterID string, mimeType string) (*usecase.AvatarUploadResult, error) {
	return s.inner.InitAvatarUpload(ctx, userID, characterID, mimeType)
}
func (s *captureCharacterListSpy) CompleteAvatarUpload(ctx context.Context, userID uuid.UUID, characterID string, uploadID uuid.UUID) error {
	return s.inner.CompleteAvatarUpload(ctx, userID, characterID, uploadID)
}
func (s *captureCharacterListSpy) DeleteAvatar(ctx context.Context, userID uuid.UUID, characterID string) error {
	return s.inner.DeleteAvatar(ctx, userID, characterID)
}
func (s *captureCharacterListSpy) GetCharacterImages(ctx context.Context, characterID uuid.UUID, userID uuid.UUID) ([]*domain.Image, error) {
	return s.inner.GetCharacterImages(ctx, characterID, userID)
}

// --- UpdateCharacter ---

func TestUpdateCharacter_HappyPath(t *testing.T) {
	character := makeCharacter()
	spy := &spyCharacterUsecase{updateResult: character}
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/characters/:id", h.DeleteCharacter)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteCharacter_NotFound(t *testing.T) {
	spy := &spyCharacterUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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
	h := newCharacterHandler(spy)

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

// --- InitAvatarUpload ---

func TestInitAvatarUpload_CharacterNotFound(t *testing.T) {
	spy := &spyCharacterUsecase{initAvatarUploadErr: usecase.ErrCharacterNotFound}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/init", h.InitAvatarUpload)

	body := `{"mime_type":"image/jpeg"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/init", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestInitAvatarUpload_MissingMimeType(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/init", h.InitAvatarUpload)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/init", uuid.New()), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestInitAvatarUpload_UnsupportedMimeType(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/init", h.InitAvatarUpload)

	body := `{"mime_type":"image/gif"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/init", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestInitAvatarUpload_Success(t *testing.T) {
	pendingID := uuid.New()
	spy := &spyCharacterUsecase{
		initAvatarUploadResult: &usecase.AvatarUploadResult{
			ID:        pendingID,
			UploadURL: "https://example.com/presigned",
			ExpiresAt: time.Now().Add(15 * time.Minute),
		},
	}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/init", h.InitAvatarUpload)

	body := `{"mime_type":"image/jpeg"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/init", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, pendingID.String(), got["id"])
	require.NotEmpty(t, got["upload_url"])
	require.NotEmpty(t, got["expires_at"])
}

func TestInitAvatarUpload_InvalidCharacterUUID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/init", h.InitAvatarUpload)

	body := `{"mime_type":"image/jpeg"}`
	req := httptest.NewRequest(http.MethodPost, "/characters/not-a-uuid/avatar/init", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- CompleteAvatarUpload ---

func TestCompleteAvatarUpload_PendingNotFound(t *testing.T) {
	spy := &spyCharacterUsecase{completeAvatarUploadErr: usecase.ErrPendingUploadNotFound}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/:uploadID/complete", h.CompleteAvatarUpload)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/%s/complete", uuid.New(), uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCompleteAvatarUpload_Success(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/:uploadID/complete", h.CompleteAvatarUpload)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/%s/complete", uuid.New(), uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCompleteAvatarUpload_InvalidUploadID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.POST("/characters/:id/avatar/:uploadID/complete", h.CompleteAvatarUpload)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/characters/%s/avatar/not-a-uuid/complete", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- DeleteAvatar ---

func TestDeleteAvatar_CharacterNotFound(t *testing.T) {
	spy := &spyCharacterUsecase{deleteAvatarErr: usecase.ErrCharacterNotFound}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/characters/:id/avatar", h.DeleteAvatar)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s/avatar", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteAvatar_Success(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/characters/:id/avatar", h.DeleteAvatar)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/characters/%s/avatar", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

// --- GetCharacterImages ---

func TestGetCharacterImages_ReturnsPressignedThumbnails(t *testing.T) {
	thumbKey := "users/1/thumbnails/img.jpg"
	img := &domain.Image{
		ID:              uuid.New(),
		UserID:          testUserID,
		ImageR2Path:     "users/1/images/img.jpg",
		MimeType:        "image/jpeg",
		ThumbnailR2Path: &thumbKey,
		Characters:      []domain.Character{},
	}
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/thumb?sig=xyz"}
	spy := &spyCharacterUsecase{getCharacterImagesResult: []*domain.Image{img}}
	h := handler.NewCharacterHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/characters/:id/images", h.GetCharacterImages)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s/images", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	require.Equal(t, img.ID.String(), got[0]["image_id"])
	require.Equal(t, "https://cdn.example.com/thumb?sig=xyz", got[0]["thumbnail_url"])
	require.Contains(t, presigner.calls, thumbKey)
}

func TestGetCharacterImages_ReturnsEmptyArray(t *testing.T) {
	spy := &spyCharacterUsecase{getCharacterImagesResult: []*domain.Image{}}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/characters/:id/images", h.GetCharacterImages)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/characters/%s/images", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Empty(t, got)
}

func TestGetCharacterImages_InvalidUUID(t *testing.T) {
	spy := &spyCharacterUsecase{}
	h := newCharacterHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/characters/:id/images", h.GetCharacterImages)

	req := httptest.NewRequest(http.MethodGet, "/characters/not-a-uuid/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
