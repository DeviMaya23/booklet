package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/usecase"
	"github.com/stretchr/testify/require"
)

type spyFolderUsecase struct {
	listResult *usecase.FolderList
	listErr    error
}

func (s *spyFolderUsecase) ListFolders(_ context.Context, _ string) (*usecase.FolderList, error) {
	return s.listResult, s.listErr
}

func TestListFolders_HappyPath(t *testing.T) {
	result := &usecase.FolderList{
		FolderList: []usecase.Folder{
			{FolderID: "folder-uuid-1", Token: "tok1", FolderName: "My Folder"},
		},
	}
	spy := &spyFolderUsecase{listResult: result}
	h := handler.NewFolderHandler(spy)

	e := setupEcho("user-1")
	e.GET("/folders", h.ListFolders)

	req := httptest.NewRequest(http.MethodGet, "/folders", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got usecase.FolderList
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got.FolderList, 1)
	require.Equal(t, "folder-uuid-1", got.FolderList[0].FolderID)
	require.Equal(t, "My Folder", got.FolderList[0].FolderName)
}

func TestListFolders_UpstreamError(t *testing.T) {
	spy := &spyFolderUsecase{listErr: errors.New("bookleaf unavailable")}
	h := handler.NewFolderHandler(spy)

	e := setupEcho("user-1")
	e.GET("/folders", h.ListFolders)

	req := httptest.NewRequest(http.MethodGet, "/folders", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
