package handler

import (
	"context"
	"net/http"

	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DashboardUsecase interface {
	GetDashboard(ctx context.Context, userID uuid.UUID) (*usecase.DashboardResult, error)
}

type DashboardHandler struct {
	dashboardUsecase DashboardUsecase
	presigner        Presigner
	tel              *observability.Telemetry
}

func NewDashboardHandler(dashboardUsecase DashboardUsecase, presigner Presigner, tel *observability.Telemetry) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: dashboardUsecase,
		presigner:        presigner,
		tel:              tel,
	}
}

type dashboardHousekeepingItemResponse struct {
	ID    string  `json:"id"`
	Title *string `json:"title"`
}

type dashboardHousekeepingResponse struct {
	CommissionsNoArtist []dashboardHousekeepingItemResponse `json:"commissions_no_artist"`
	ArtpiecesNoArtist   []dashboardHousekeepingItemResponse `json:"artpieces_no_artist"`
	DoneNoArtpieces     []dashboardHousekeepingItemResponse `json:"done_no_artpieces"`
	ArtpiecesNoFiles    []dashboardHousekeepingItemResponse `json:"artpieces_no_files"`
}

type dashboardRecentArtpieceResponse struct {
	ID           string  `json:"id"`
	Title        *string `json:"title"`
	ArtistName   *string `json:"artist_name"`
	ThumbnailURL *string `json:"thumbnail_url"`
}

type dashboardInProgressResponse struct {
	ID              string  `json:"id"`
	Title           *string `json:"title"`
	ArtistName      *string `json:"artist_name"`
	ArtistLink      *string `json:"artist_link"`
	Status          string  `json:"status"`
	Paid            bool    `json:"paid"`
	LastContactedAt *string `json:"last_contacted_at"`
	CreatedAt       string  `json:"created_at"`
}

type dashboardResponse struct {
	RecentArtpieces []dashboardRecentArtpieceResponse `json:"recent_artpieces"`
	InProgress      []dashboardInProgressResponse     `json:"in_progress"`
	Housekeeping    dashboardHousekeepingResponse     `json:"housekeeping"`
}

func (h *DashboardHandler) GetDashboard(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetDashboard")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	result, err := h.dashboardUsecase.GetDashboard(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load dashboard")
	}

	recentArtpieces := make([]dashboardRecentArtpieceResponse, 0, len(result.RecentArtpieces))
	for _, a := range result.RecentArtpieces {
		item := dashboardRecentArtpieceResponse{
			ID:    a.ID.String(),
			Title: a.Title,
		}
		if a.Artist != nil {
			item.ArtistName = &a.Artist.Name
		}
		if a.CoverFile != nil && a.CoverFile.ThumbnailR2Path != nil {
			u, err := h.presigner.GenerateDeterministicPresignedGetURL(ctx, *a.CoverFile.ThumbnailR2Path)
			if err == nil {
				item.ThumbnailURL = &u
			}
		}
		recentArtpieces = append(recentArtpieces, item)
	}

	inProgress := make([]dashboardInProgressResponse, 0, len(result.InProgress))
	for _, c := range result.InProgress {
		item := dashboardInProgressResponse{
			ID:        c.ID.String(),
			Title:     c.Title,
			Status:    c.Status,
			Paid:      c.Paid,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if c.Artist != nil {
			item.ArtistName = &c.Artist.Name
			for _, l := range c.Artist.Links {
				if l.IsPrimary {
					u := l.URL
					item.ArtistLink = &u
					break
				}
			}
		}
		if c.LastContactedAt != nil {
			s := c.LastContactedAt.Format("2006-01-02T15:04:05Z07:00")
			item.LastContactedAt = &s
		}
		inProgress = append(inProgress, item)
	}

	toItems := func(src []usecase.DashboardHousekeepingItem) []dashboardHousekeepingItemResponse {
		out := make([]dashboardHousekeepingItemResponse, len(src))
		for i, it := range src {
			out[i] = dashboardHousekeepingItemResponse{ID: it.ID.String(), Title: it.Title}
		}
		return out
	}

	return c.JSON(http.StatusOK, dashboardResponse{
		RecentArtpieces: recentArtpieces,
		InProgress:      inProgress,
		Housekeeping: dashboardHousekeepingResponse{
			CommissionsNoArtist: toItems(result.CommissionsNoArtist),
			ArtpiecesNoArtist:   toItems(result.ArtpiecesNoArtist),
			DoneNoArtpieces:     toItems(result.DoneNoArtpieces),
			ArtpiecesNoFiles:    toItems(result.ArtpiecesNoFiles),
		},
	})
}
