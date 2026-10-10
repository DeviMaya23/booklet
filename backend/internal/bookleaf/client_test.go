package bookleaf_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/stretchr/testify/require"
)

func newTestClient(server *httptest.Server) *bookleaf.Client {
	return bookleaf.NewClient(server.URL, "test-secret")
}

func TestGetFolderImages_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "test-secret", r.Header.Get("X-Bookleaf-Internal-Secret"))
		require.Equal(t, "/internal/users/user-1/folders/folder-1/images", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(bookleaf.FolderImageList{
			Images: []bookleaf.FolderImage{
				{ImageID: "img-1", ThumbnailURL: "https://cdn.example.com/thumb1.jpg"},
				{ImageID: "img-2", ThumbnailURL: "https://cdn.example.com/thumb2.jpg"},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)
	result, err := client.GetFolderImages(context.Background(), "user-1", "folder-1")

	require.NoError(t, err)
	require.Len(t, result.Images, 2)
	require.Equal(t, "img-1", result.Images[0].ImageID)
	require.Equal(t, "https://cdn.example.com/thumb1.jpg", result.Images[0].ThumbnailURL)
}

func TestGetFolderImages_NonOKStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{name: "not found", statusCode: http.StatusNotFound, wantErr: bookleaf.ErrNotFound},
		{name: "internal server error", statusCode: http.StatusInternalServerError, wantErr: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))
			defer srv.Close()

			client := newTestClient(srv)
			_, err := client.GetFolderImages(context.Background(), "user-1", "folder-1")

			require.Error(t, err)
			if tc.wantErr != nil {
				require.True(t, errors.Is(err, tc.wantErr), "expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
