package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/stretchr/testify/require"
)

type spyFolderUsecase struct {
	listResult        *bookleaf.FolderList
	listErr           error
	imagesResult      *bookleaf.FolderImageList
	imagesErr         error
}

func (s *spyFolderUsecase) ListFolders(_ context.Context, _ string) (*bookleaf.FolderList, error) {
	return s.listResult, s.listErr
}

func (s *spyFolderUsecase) GetFolderImages(_ context.Context, _, _ string) (*bookleaf.FolderImageList, error) {
	return s.imagesResult, s.imagesErr
}

func TestListFolders_HappyPath(t *testing.T) {
	result := &bookleaf.FolderList{
		FolderList: []bookleaf.Folder{
			{FolderID: "folder-uuid-1", Token: "tok1", FolderName: "My Folder"},
		},
	}
	spy := &spyFolderUsecase{listResult: result}
	h := handler.NewFolderHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/folders", h.ListFolders)

	req := httptest.NewRequest(http.MethodGet, "/folders", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got bookleaf.FolderList
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.FolderList, 1)
	require.Equal(t, "folder-uuid-1", got.FolderList[0].FolderID)
	require.Equal(t, "My Folder", got.FolderList[0].FolderName)
}

func TestListFolders_UpstreamError(t *testing.T) {
	spy := &spyFolderUsecase{listErr: errors.New("bookleaf unavailable")}
	h := handler.NewFolderHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/folders", h.ListFolders)

	req := httptest.NewRequest(http.MethodGet, "/folders", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetFolderImages_HappyPath(t *testing.T) {
	folderID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	result := &bookleaf.FolderImageList{
		Images: []bookleaf.FolderImage{
			{ImageID: "img-1", ThumbnailURL: "https://cdn.example.com/thumb1.jpg"},
		},
	}
	spy := &spyFolderUsecase{imagesResult: result}
	h := handler.NewFolderHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/folders/:folderID/images", h.GetFolderImages)

	req := httptest.NewRequest(http.MethodGet, "/folders/"+folderID+"/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got bookleaf.FolderImageList
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.Images, 1)
	require.Equal(t, "img-1", got.Images[0].ImageID)
	require.Equal(t, "https://cdn.example.com/thumb1.jpg", got.Images[0].ThumbnailURL)
}

func TestGetFolderImages_EmptyFolder(t *testing.T) {
	folderID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	spy := &spyFolderUsecase{imagesResult: &bookleaf.FolderImageList{Images: []bookleaf.FolderImage{}}}
	h := handler.NewFolderHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/folders/:folderID/images", h.GetFolderImages)

	req := httptest.NewRequest(http.MethodGet, "/folders/"+folderID+"/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got bookleaf.FolderImageList
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Empty(t, got.Images)
}

func TestGetFolderImages_NotFound(t *testing.T) {
	folderID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	spy := &spyFolderUsecase{imagesErr: fmt.Errorf("wrapped: %w", bookleaf.ErrNotFound)}
	h := handler.NewFolderHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/folders/:folderID/images", h.GetFolderImages)

	req := httptest.NewRequest(http.MethodGet, "/folders/"+folderID+"/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetFolderImages_UpstreamError(t *testing.T) {
	folderID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	spy := &spyFolderUsecase{imagesErr: errors.New("bookleaf unavailable")}
	h := handler.NewFolderHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.GET("/folders/:folderID/images", h.GetFolderImages)

	req := httptest.NewRequest(http.MethodGet, "/folders/"+folderID+"/images", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
