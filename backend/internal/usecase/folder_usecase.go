package usecase

import (
	"context"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/platform/observability"
	"go.opentelemetry.io/otel/codes"
)

type BookleafClient interface {
	GetPublicFolders(ctx context.Context, userID string) (*bookleaf.FolderList, error)
	DeleteAccount(ctx context.Context, kindeUserID string) error
}

type folderUsecase struct {
	client BookleafClient
	tel    *observability.Telemetry
}

func NewFolderUsecase(client BookleafClient, tel *observability.Telemetry) *folderUsecase {
	return &folderUsecase{client: client, tel: tel}
}

func (u *folderUsecase) ListFolders(ctx context.Context, userID string) (*bookleaf.FolderList, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListFolders")
	defer span.End()

	result, err := u.client.GetPublicFolders(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return result, nil
}
