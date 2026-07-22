package usecase

import "context"

type Folder struct {
	FolderID   string `json:"folder_id"`
	Token      string `json:"token"`
	FolderName string `json:"folder_name"`
}

type FolderList struct {
	FolderList []Folder `json:"folder_list"`
}

type BookleafClient interface {
	GetPublicFolders(ctx context.Context, userID string) (*FolderList, error)
}

type folderUsecase struct {
	client BookleafClient
}

func NewFolderUsecase(client BookleafClient) *folderUsecase {
	return &folderUsecase{client: client}
}

func (u *folderUsecase) ListFolders(ctx context.Context, userID string) (*FolderList, error) {
	return u.client.GetPublicFolders(ctx, userID)
}
