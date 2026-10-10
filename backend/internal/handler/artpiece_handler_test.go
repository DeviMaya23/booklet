package handler_test

import (
	"context"
	"encoding/json"
	"io"
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

type spyArtpieceUsecase struct {
	createResult *domain.Artpiece
	createErr    error

	getByIDResult *domain.Artpiece
	getByIDErr    error

	listResult []*domain.Artpiece
	listErr    error

	updateResult *domain.Artpiece
	updateErr    error

	deleteErr error

	attachResult *domain.Artpiece
	attachErr    error

	detachResult *domain.Artpiece
	detachErr    error

	setCoverResult *domain.Artpiece
	setCoverErr    error
}

func (s *spyArtpieceUsecase) Create(_ context.Context, _ uuid.UUID, _ usecase.CreateArtpieceParams) (*domain.Artpiece, error) {
	return s.createResult, s.createErr
}
func (s *spyArtpieceUsecase) GetByID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Artpiece, error) {
	return s.getByIDResult, s.getByIDErr
}
func (s *spyArtpieceUsecase) List(_ context.Context, _ uuid.UUID, _ usecase.ListArtpieceFilters) ([]*domain.Artpiece, error) {
	return s.listResult, s.listErr
}
func (s *spyArtpieceUsecase) Update(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ usecase.UpdateArtpieceParams) (*domain.Artpiece, error) {
	return s.updateResult, s.updateErr
}
func (s *spyArtpieceUsecase) Delete(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ bool) error {
	return s.deleteErr
}
func (s *spyArtpieceUsecase) AttachFile(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ uuid.UUID) (*domain.Artpiece, error) {
	return s.attachResult, s.attachErr
}
func (s *spyArtpieceUsecase) DetachFile(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ uuid.UUID) (*domain.Artpiece, error) {
	return s.detachResult, s.detachErr
}
func (s *spyArtpieceUsecase) SetCover(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ uuid.UUID) (*domain.Artpiece, error) {
	return s.setCoverResult, s.setCoverErr
}
func (s *spyArtpieceUsecase) ReplaceFiles(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ []uuid.UUID) (*domain.Artpiece, error) {
	return nil, nil
}
func (s *spyArtpieceUsecase) DownloadFiles(_ context.Context, _ *domain.Artpiece, _ io.Writer) error {
	return nil
}

func makeArtpiece() *domain.Artpiece {
	return &domain.Artpiece{
		ID:         uuid.New(),
		UserID:     testUserID,
		Characters: []domain.Character{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// --- Create ---

func TestCreateArtpiece_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{createResult: a}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artpieces", h.CreateArtpiece)

	req := httptest.NewRequest(http.MethodPost, "/artpieces", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, a.ID.String(), got["id"])
}

func TestCreateArtpiece_ArtistNotOwned(t *testing.T) {
	spy := &spyArtpieceUsecase{createErr: usecase.ErrArtistNotOwned}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artpieces", h.CreateArtpiece)

	body := `{"artist_id":"00000000-0000-0000-0000-000000000099"}`
	req := httptest.NewRequest(http.MethodPost, "/artpieces", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateArtpiece_CharacterNotOwned(t *testing.T) {
	spy := &spyArtpieceUsecase{createErr: usecase.ErrCharacterNotOwned}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artpieces", h.CreateArtpiece)

	body := `{"character_ids":["00000000-0000-0000-0000-000000000099"]}`
	req := httptest.NewRequest(http.MethodPost, "/artpieces", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// --- GetByID ---

func TestGetArtpieceByID_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{getByIDResult: a}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artpieces/:id", h.GetArtpieceByID)

	req := httptest.NewRequest(http.MethodGet, "/artpieces/"+a.ID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, a.ID.String(), got["id"])
}

func TestGetArtpieceByID_NotFound(t *testing.T) {
	spy := &spyArtpieceUsecase{getByIDErr: gorm.ErrRecordNotFound}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artpieces/:id", h.GetArtpieceByID)

	req := httptest.NewRequest(http.MethodGet, "/artpieces/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetArtpieceByID_InvalidUUID(t *testing.T) {
	spy := &spyArtpieceUsecase{}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artpieces/:id", h.GetArtpieceByID)

	req := httptest.NewRequest(http.MethodGet, "/artpieces/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- List ---

func TestListArtpieces_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{listResult: []*domain.Artpiece{a}}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artpieces", h.ListArtpieces)

	req := httptest.NewRequest(http.MethodGet, "/artpieces", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
}

func TestListArtpieces_CoverFieldsPopulatedWhenCoverExists(t *testing.T) {
	coverFileID := uuid.New()
	fileName := "artwork.jpg"
	coverFile := &domain.File{
		ID:         coverFileID,
		FileR2Path: "users/x/files/artwork.jpg",
		MimeType:   "image/jpeg",
		Name:       &fileName,
	}
	a := makeArtpiece()
	a.CoverFileID = &coverFileID
	a.CoverFile = coverFile

	presigner := &spyPresigner{presignedURL: "https://cdn.example.com/artwork.jpg?sig=xyz"}
	spy := &spyArtpieceUsecase{listResult: []*domain.Artpiece{a}}
	h := handler.NewArtpieceHandler(spy, presigner, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artpieces", h.ListArtpieces)

	req := httptest.NewRequest(http.MethodGet, "/artpieces", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	require.Equal(t, "https://cdn.example.com/artwork.jpg?sig=xyz", got[0]["cover_file_url"])
	require.Equal(t, "image/jpeg", got[0]["cover_file_mime_type"])
	require.Equal(t, "artwork.jpg", got[0]["cover_file_name"])
}

func TestListArtpieces_CoverFieldsNullWhenNoCover(t *testing.T) {
	a := makeArtpiece()
	// CoverFile and CoverFileID are nil

	spy := &spyArtpieceUsecase{listResult: []*domain.Artpiece{a}}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/artpieces", h.ListArtpieces)

	req := httptest.NewRequest(http.MethodGet, "/artpieces", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	require.Nil(t, got[0]["cover_file_url"])
	require.Nil(t, got[0]["cover_file_mime_type"])
	require.Nil(t, got[0]["cover_file_name"])
}

// --- Update ---

func TestUpdateArtpiece_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{updateResult: a}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PUT("/artpieces/:id", h.UpdateArtpiece)

	req := httptest.NewRequest(http.MethodPut, "/artpieces/"+a.ID.String(), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateArtpiece_NotFound(t *testing.T) {
	spy := &spyArtpieceUsecase{updateErr: gorm.ErrRecordNotFound}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PUT("/artpieces/:id", h.UpdateArtpiece)

	req := httptest.NewRequest(http.MethodPut, "/artpieces/"+uuid.New().String(), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateArtpiece_ArtistNotOwned(t *testing.T) {
	spy := &spyArtpieceUsecase{updateErr: usecase.ErrArtistNotOwned}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PUT("/artpieces/:id", h.UpdateArtpiece)

	body := `{"artist_id":"00000000-0000-0000-0000-000000000099"}`
	req := httptest.NewRequest(http.MethodPut, "/artpieces/"+uuid.New().String(), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// --- Delete ---

func TestDeleteArtpiece_HappyPath(t *testing.T) {
	spy := &spyArtpieceUsecase{}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artpieces/:id", h.DeleteArtpiece)

	req := httptest.NewRequest(http.MethodDelete, "/artpieces/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteArtpiece_NotFound(t *testing.T) {
	spy := &spyArtpieceUsecase{deleteErr: gorm.ErrRecordNotFound}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artpieces/:id", h.DeleteArtpiece)

	req := httptest.NewRequest(http.MethodDelete, "/artpieces/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteArtpiece_InternalError(t *testing.T) {
	spy := &spyArtpieceUsecase{deleteErr: gorm.ErrInvalidDB}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artpieces/:id", h.DeleteArtpiece)

	req := httptest.NewRequest(http.MethodDelete, "/artpieces/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- AttachFile ---

func TestAttachFile_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{attachResult: a}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artpieces/:id/files/:file_id", h.AttachFile)

	req := httptest.NewRequest(http.MethodPost, "/artpieces/"+a.ID.String()+"/files/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAttachFile_FileNotOwned(t *testing.T) {
	spy := &spyArtpieceUsecase{attachErr: usecase.ErrFileNotOwned}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artpieces/:id/files/:file_id", h.AttachFile)

	req := httptest.NewRequest(http.MethodPost, "/artpieces/"+uuid.New().String()+"/files/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestAttachFile_ArtpieceNotFound(t *testing.T) {
	spy := &spyArtpieceUsecase{attachErr: gorm.ErrRecordNotFound}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.POST("/artpieces/:id/files/:file_id", h.AttachFile)

	req := httptest.NewRequest(http.MethodPost, "/artpieces/"+uuid.New().String()+"/files/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

// --- DetachFile ---

func TestDetachFile_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{detachResult: a}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artpieces/:id/files/:file_id", h.DetachFile)

	req := httptest.NewRequest(http.MethodDelete, "/artpieces/"+a.ID.String()+"/files/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestDetachFile_FileNotInArtpiece(t *testing.T) {
	spy := &spyArtpieceUsecase{detachErr: usecase.ErrFileNotInArtpiece}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/artpieces/:id/files/:file_id", h.DetachFile)

	req := httptest.NewRequest(http.MethodDelete, "/artpieces/"+uuid.New().String()+"/files/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// --- SetCover ---

func TestSetCover_HappyPath(t *testing.T) {
	a := makeArtpiece()
	spy := &spyArtpieceUsecase{setCoverResult: a}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PUT("/artpieces/:id/cover", h.SetCover)

	fileID := uuid.New().String()
	body := `{"file_id":"` + fileID + `"}`
	req := httptest.NewRequest(http.MethodPut, "/artpieces/"+a.ID.String()+"/cover", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestSetCover_FileNotInArtpiece(t *testing.T) {
	spy := &spyArtpieceUsecase{setCoverErr: usecase.ErrFileNotInArtpiece}
	h := handler.NewArtpieceHandler(spy, &spyPresigner{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.PUT("/artpieces/:id/cover", h.SetCover)

	fileID := uuid.New().String()
	body := `{"file_id":"` + fileID + `"}`
	req := httptest.NewRequest(http.MethodPut, "/artpieces/"+uuid.New().String()+"/cover", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
