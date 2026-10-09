package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type DashboardUsecase struct {
	repo DashboardRepository
	tel  *observability.Telemetry
}

func NewDashboardUsecase(repo DashboardRepository, tel *observability.Telemetry) *DashboardUsecase {
	return &DashboardUsecase{repo: repo, tel: tel}
}

type DashboardResult struct {
	RecentArtpieces      []*domain.Artpiece
	InProgress           []*domain.Commission
	CommissionsNoArtist  []DashboardHousekeepingItem
	ArtpiecesNoArtist    []DashboardHousekeepingItem
	DoneNoArtpieces      []DashboardHousekeepingItem
	ArtpiecesNoFiles     []DashboardHousekeepingItem
}

func (u *DashboardUsecase) GetDashboard(ctx context.Context, userID uuid.UUID) (*DashboardResult, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetDashboard")
	defer span.End()

	recentArtpieces, err := u.repo.GetRecentArtpieces(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	inProgress, err := u.repo.GetInProgressCommissions(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	commissionsNoArtist, err := u.repo.GetCommissionsNoArtist(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	artpiecesNoArtist, err := u.repo.GetArtpiecesNoArtist(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	doneNoArtpieces, err := u.repo.GetDoneCommissionsNoArtpieces(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	artpiecesNoFiles, err := u.repo.GetArtpiecesNoFiles(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &DashboardResult{
		RecentArtpieces:     recentArtpieces,
		InProgress:          inProgress,
		CommissionsNoArtist: commissionsNoArtist,
		ArtpiecesNoArtist:   artpiecesNoArtist,
		DoneNoArtpieces:     doneNoArtpieces,
		ArtpiecesNoFiles:    artpiecesNoFiles,
	}, nil
}
