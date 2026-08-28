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

type imageResponse struct {
	ID              string          `json:"id"`
	ImageR2Path     string          `json:"image_r2_path"`
	MimeType        string          `json:"mime_type"`
	Title           *string         `json:"title"`
	ThumbnailR2Path *string         `json:"thumbnail_r2_path"`
	ArtistName      *string         `json:"artist_name"`
	ArtistLink      *string         `json:"artist_link"`
	Notes           *string         `json:"notes"`
	Characters      []characterItem `json:"characters"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

type characterItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type spyImageUsecase struct {
	getByIDResult *domain.Image
	getByIDErr    error

	listResult []*domain.Image
	listErr    error

	updateResult     *domain.Image
	updateErr        error
	lastUpdateParams usecase.UpdateImageParams

	deleteErr error
}

func (s *spyImageUsecase) GetByID(_ context.Context, _ string, _ uuid.UUID) (*domain.Image, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *spyImageUsecase) List(_ context.Context, _ uuid.UUID) ([]*domain.Image, error) {
	return s.listResult, s.listErr
}

func (s *spyImageUsecase) Update(_ context.Context, _ string, _ uuid.UUID, params usecase.UpdateImageParams) (*domain.Image, error) {
	s.lastUpdateParams = params
	return s.updateResult, s.updateErr
}

func (s *spyImageUsecase) Delete(_ context.Context, _ string, _ uuid.UUID) error {
	return s.deleteErr
}

func makeImage() *domain.Image {
	return &domain.Image{
		ID:          uuid.New(),
		UserID:      testUserID,
		ImageR2Path: "images/abc.jpg",
		MimeType:    "image/jpeg",
		Characters:  []domain.Character{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// --- GetImageByID ---

func TestGetImageByID_HappyPath(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{getByIDResult: image}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images/:id", h.GetImageByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/images/%s", image.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, image.ID.String(), got["id"])
	require.Equal(t, image.ImageR2Path, got["image_r2_path"])
	chars, ok := got["characters"].([]interface{})
	require.True(t, ok, "characters should be a JSON array, not null")
	require.Empty(t, chars)
}

func TestGetImageByID_NotFound(t *testing.T) {
	spy := &spyImageUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images/:id", h.GetImageByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/images/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetImageByID_InvalidUUID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images/:id", h.GetImageByID)

	req := httptest.NewRequest(http.MethodGet, "/images/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- ListImages ---

func TestListImages_HappyPath(t *testing.T) {
	images := []*domain.Image{makeImage(), makeImage()}
	spy := &spyImageUsecase{listResult: images}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []imageResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 2)
}

func TestListImages_EmptyCharactersArray(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{listResult: []*domain.Image{image}}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var raw []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &raw))
	require.Len(t, raw, 1)
	chars, ok := raw[0]["characters"].([]interface{})
	require.True(t, ok, "characters should be an array, not null")
	require.Empty(t, chars)
}

// --- UpdateImage ---

func TestUpdateImage_HappyPath(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{updateResult: image}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/images/:id", h.UpdateImage)

	body := `{"artist_name":"Jane Doe"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/images/%s", image.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got imageResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, image.ID.String(), got.ID)
}

func TestUpdateImage_NotFound(t *testing.T) {
	spy := &spyImageUsecase{updateErr: gorm.ErrRecordNotFound}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/images/:id", h.UpdateImage)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateImage_CharacterNotOwned(t *testing.T) {
	spy := &spyImageUsecase{updateErr: usecase.ErrCharacterNotOwned}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/images/:id", h.UpdateImage)

	body := `{"character_ids":["` + uuid.New().String() + `"]}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateImage_MalformedJSON(t *testing.T) {
	spy := &spyImageUsecase{}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/images/:id", h.UpdateImage)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateImage_InvalidUUID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/images/:id", h.UpdateImage)

	req := httptest.NewRequest(http.MethodPatch, "/images/not-a-uuid", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateImage_AbsentCharacterIDs(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{updateResult: image}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PATCH("/images/:id", h.UpdateImage)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/images/%s", image.ID), strings.NewReader(`{"artist_name":"Jane"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, spy.lastUpdateParams.CharacterIDs)
}

// --- DeleteImage ---

func TestDeleteImage_HappyPath(t *testing.T) {
	spy := &spyImageUsecase{}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/images/:id", h.DeleteImage)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/images/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteImage_NotFound(t *testing.T) {
	spy := &spyImageUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/images/:id", h.DeleteImage)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/images/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteImage_InvalidUUID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := handler.NewImageHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/images/:id", h.DeleteImage)

	req := httptest.NewRequest(http.MethodDelete, "/images/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
