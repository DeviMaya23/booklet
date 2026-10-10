package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/stretchr/testify/require"
)

type spyBookleafClientFolderUsecase struct {
	folderListResult  *bookleaf.FolderList
	folderListErr     error
	folderImagesResult *bookleaf.FolderImageList
	folderImagesErr    error
}

func (s *spyBookleafClientFolderUsecase) GetPublicFolders(_ context.Context, _ string) (*bookleaf.FolderList, error) {
	return s.folderListResult, s.folderListErr
}

func (s *spyBookleafClientFolderUsecase) GetFolderImages(_ context.Context, _, _ string) (*bookleaf.FolderImageList, error) {
	return s.folderImagesResult, s.folderImagesErr
}

func (s *spyBookleafClientFolderUsecase) DeleteAccount(_ context.Context, _ string) error {
	return nil
}

func TestGetFolderImages_SuccessForwardedFromClient(t *testing.T) {
	expected := &bookleaf.FolderImageList{
		Images: []bookleaf.FolderImage{
			{ImageID: "img-1", ThumbnailURL: "https://cdn.example.com/thumb1.jpg"},
		},
	}
	spy := &spyBookleafClientFolderUsecase{folderImagesResult: expected}
	uc := usecase.NewFolderUsecase(spy, observability.NewTelemetry(nil, nil, nil))

	result, err := uc.GetFolderImages(context.Background(), "user-1", "folder-1")

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestGetFolderImages_ClientErrorReturnedAsCaller(t *testing.T) {
	clientErr := errors.New("bookleaf unavailable")
	spy := &spyBookleafClientFolderUsecase{folderImagesErr: clientErr}
	uc := usecase.NewFolderUsecase(spy, observability.NewTelemetry(nil, nil, nil))

	_, err := uc.GetFolderImages(context.Background(), "user-1", "folder-1")

	require.ErrorIs(t, err, clientErr)
}
