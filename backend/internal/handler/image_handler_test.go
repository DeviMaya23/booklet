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
	ID           string          `json:"id"`
	ImageURL     *string         `json:"image_url"`
	MimeType     string          `json:"mime_type"`
	Title        *string         `json:"title"`
	ThumbnailURL *string         `json:"thumbnail_url"`
	ArtistID     *string         `json:"artist_id"`
	ArtistName   *string         `json:"artist_name"`
	Notes        *string         `json:"notes"`
	Characters   []characterItem `json:"characters"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
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

func (s *spyImageUsecase) List(_ context.Context, _ uuid.UUID, _ usecase.ListImageFilters) ([]*domain.Image, error) {
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

func newImageHandler(spy *spyImageUsecase) *handler.ImageHandler {
	return handler.NewImageHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))
}

// --- GetImageByID ---

func TestGetImageByID_HappyPath(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{getByIDResult: image}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/images/:id", h.GetImageByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/images/%s", image.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, image.ID.String(), got["id"])
	chars, ok := got["characters"].([]interface{})
	require.True(t, ok, "characters should be a JSON array, not null")
	require.Empty(t, chars)
}

func TestGetImageByID_ImageURLAndThumbnailURLPresigned(t *testing.T) {
	thumbKey := "users/1/thumbnails/img.jpg"
	image := makeImage()
	image.ThumbnailR2Path = &thumbKey
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/signed?sig=abc"}
	spy := &spyImageUsecase{getByIDResult: image}
	h := handler.NewImageHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images/:id", h.GetImageByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/images/%s", image.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	imageURL, hasImageURL := got["image_url"]
	require.True(t, hasImageURL)
	require.Equal(t, "https://cdn.example.com/signed?sig=abc", imageURL)

	thumbnailURL, hasThumbnailURL := got["thumbnail_url"]
	require.True(t, hasThumbnailURL)
	require.Equal(t, "https://cdn.example.com/signed?sig=abc", thumbnailURL)

	_, hasRawImagePath := got["image_r2_path"]
	require.False(t, hasRawImagePath)
	_, hasThumbnailPath := got["thumbnail_r2_path"]
	require.False(t, hasThumbnailPath)

	require.Contains(t, presigner.calls, image.ImageR2Path)
	require.Contains(t, presigner.calls, thumbKey)
}

func TestGetImageByID_NotFound(t *testing.T) {
	spy := &spyImageUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/images/:id", h.GetImageByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/images/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetImageByID_InvalidUUID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

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
	h := newImageHandler(spy)

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

func TestListImages_ThumbnailURLPresignedImageURLAbsent(t *testing.T) {
	thumbKey := "users/1/thumbnails/img.jpg"
	image := makeImage()
	image.ThumbnailR2Path = &thumbKey
	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/thumb?sig=abc"}
	spy := &spyImageUsecase{listResult: []*domain.Image{image}}
	h := handler.NewImageHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var raw []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &raw))
	require.Len(t, raw, 1)

	_, hasImageURL := raw[0]["image_url"]
	require.False(t, hasImageURL, "image_url must be absent from list responses")

	thumbnailURL, hasThumbnailURL := raw[0]["thumbnail_url"]
	require.True(t, hasThumbnailURL)
	require.Equal(t, "https://cdn.example.com/thumb?sig=abc", thumbnailURL)

	require.Contains(t, presigner.calls, thumbKey)
}

func TestListImages_EmptyCharactersArray(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{listResult: []*domain.Image{image}}
	h := newImageHandler(spy)

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

func TestListImages_QFilterBoundToUsecase(t *testing.T) {
	images := []*domain.Image{makeImage()}
	spy := &captureImageListSpy{inner: &spyImageUsecase{listResult: images}}
	h := handler.NewImageHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images?q=sunset", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, spy.capturedFilters.Q)
	require.Equal(t, "sunset", *spy.capturedFilters.Q)
}

func TestListImages_CharacterIDsFilterBoundToUsecase(t *testing.T) {
	charID := uuid.New()
	images := []*domain.Image{makeImage()}
	spy := &captureImageListSpy{inner: &spyImageUsecase{listResult: images}}
	h := handler.NewImageHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images?character_ids="+charID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, spy.capturedFilters.CharacterIDs, 1)
	require.Equal(t, charID.String(), spy.capturedFilters.CharacterIDs[0])
}

func TestListImages_ArtistIDsFilterBoundToUsecase(t *testing.T) {
	artistID := uuid.New()
	images := []*domain.Image{makeImage()}
	spy := &captureImageListSpy{inner: &spyImageUsecase{listResult: images}}
	h := handler.NewImageHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images?artist_ids="+artistID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, spy.capturedFilters.ArtistIDs, 1)
	require.Equal(t, artistID.String(), spy.capturedFilters.ArtistIDs[0])
}

func TestListImages_MalformedCharacterID_Returns400(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images?character_ids=not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListImages_MalformedArtistID_Returns400(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.GET("/images", h.ListImages)

	req := httptest.NewRequest(http.MethodGet, "/images?artist_ids=not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

type captureImageListSpy struct {
	inner           *spyImageUsecase
	capturedFilters usecase.ListImageFilters
}

func (s *captureImageListSpy) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Image, error) {
	return s.inner.GetByID(ctx, id, userID)
}
func (s *captureImageListSpy) List(ctx context.Context, userID uuid.UUID, filters usecase.ListImageFilters) ([]*domain.Image, error) {
	s.capturedFilters = filters
	return s.inner.List(ctx, userID, filters)
}
func (s *captureImageListSpy) Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateImageParams) (*domain.Image, error) {
	return s.inner.Update(ctx, id, userID, params)
}
func (s *captureImageListSpy) Delete(ctx context.Context, id string, userID uuid.UUID) error {
	return s.inner.Delete(ctx, id, userID)
}

// --- UpdateImage ---

func TestUpdateImage_HappyPath(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{updateResult: image}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":"some note","artist_id":null,"character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", image.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got imageResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, image.ID.String(), got.ID)
}

func TestUpdateImage_ClearsTitle(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{updateResult: image}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":null,"character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", image.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, spy.lastUpdateParams.Title)
	require.Nil(t, spy.lastUpdateParams.Notes)
}

func TestUpdateImage_NotFound(t *testing.T) {
	spy := &spyImageUsecase{updateErr: gorm.ErrRecordNotFound}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":null,"character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateImage_CharacterNotOwned(t *testing.T) {
	spy := &spyImageUsecase{updateErr: usecase.ErrCharacterNotOwned}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":null,"character_ids":["` + uuid.New().String() + `"]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateImage_InvalidArtistID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":"not-a-uuid","character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var respBody struct {
		Errors []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respBody))
	require.Len(t, respBody.Errors, 1)
	require.Equal(t, "artist_id", respBody.Errors[0].Field)
	require.Equal(t, "artist_id must be a valid UUID", respBody.Errors[0].Message)
}

func TestUpdateImage_ArtistNotOwned(t *testing.T) {
	spy := &spyImageUsecase{updateErr: usecase.ErrArtistNotOwned}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":"` + uuid.New().String() + `","character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateImage_ArtistIDNull_ClearsArtist(t *testing.T) {
	image := makeImage()
	spy := &spyImageUsecase{updateResult: image}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":null,"character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", image.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, spy.lastUpdateParams.ArtistID)
}

func TestUpdateImage_MalformedJSON(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/images/%s", uuid.New()), strings.NewReader(`{bad`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateImage_InvalidUUID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.PUT("/images/:id", h.UpdateImage)

	body := `{"title":null,"notes":null,"artist_id":null,"character_ids":[]}`
	req := httptest.NewRequest(http.MethodPut, "/images/not-a-uuid", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- DeleteImage ---

func TestDeleteImage_HappyPath(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/images/:id", h.DeleteImage)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/images/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteImage_NotFound(t *testing.T) {
	spy := &spyImageUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/images/:id", h.DeleteImage)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/images/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteImage_InvalidUUID(t *testing.T) {
	spy := &spyImageUsecase{}
	h := newImageHandler(spy)

	e := setupEcho(testUserID)
	e.DELETE("/images/:id", h.DeleteImage)

	req := httptest.NewRequest(http.MethodDelete, "/images/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
