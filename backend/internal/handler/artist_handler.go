package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ArtistUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateArtistParams) (*domain.Artist, error)
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Artist, error)
	List(ctx context.Context, userID uuid.UUID, filters usecase.ListArtistFilters) ([]*domain.Artist, error)
	Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateArtistParams) (*domain.Artist, error)
	Delete(ctx context.Context, id string, userID uuid.UUID) error
}

type ArtistHandler struct {
	artistUsecase ArtistUsecase
	tel           *observability.Telemetry
}

func NewArtistHandler(artistUsecase ArtistUsecase, tel *observability.Telemetry) *ArtistHandler {
	return &ArtistHandler{artistUsecase: artistUsecase, tel: tel}
}

type createArtistRequest struct {
	Name       string  `json:"name" validate:"required"`
	Notes      *string `json:"notes"`
	ArtistLink *string `json:"artist_link" validate:"omitempty,url"`
}

type updateArtistRequest struct {
	Name       *string `json:"name" validate:"omitempty,min=1"`
	Notes      *string `json:"notes"`
	ArtistLink *string `json:"artist_link" validate:"omitempty,url"`
}

type artistResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Notes      *string `json:"notes"`
	ArtistLink *string `json:"artist_link"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func (h *ArtistHandler) CreateArtist(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateArtist")
	defer span.End()

	var req createArtistRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artist, err := h.artistUsecase.Create(ctx, userID, usecase.CreateArtistParams{
		Name:       req.Name,
		Notes:      req.Notes,
		ArtistLink: req.ArtistLink,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrArtistNameConflict) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create artist")
	}

	return c.JSON(http.StatusCreated, toArtistResponse(artist))
}

func (h *ArtistHandler) GetArtistByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetArtistByID")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artist id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artist, err := h.artistUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artist not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get artist")
	}

	return c.JSON(http.StatusOK, toArtistResponse(artist))
}

func (h *ArtistHandler) ListArtists(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListArtists")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var filters usecase.ListArtistFilters
	if err := c.Bind(&filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params")
	}

	artists, err := h.artistUsecase.List(ctx, userID, filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list artists")
	}

	responses := make([]artistResponse, len(artists))
	for i, artist := range artists {
		responses[i] = toArtistResponse(artist)
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *ArtistHandler) UpdateArtist(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateArtist")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artist id")
	}

	var req updateArtistRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artist, err := h.artistUsecase.Update(ctx, id, userID, usecase.UpdateArtistParams{
		Name:       req.Name,
		Notes:      req.Notes,
		ArtistLink: req.ArtistLink,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artist not found")
		}
		if errors.Is(err, usecase.ErrArtistNameConflict) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update artist")
	}

	return c.JSON(http.StatusOK, toArtistResponse(artist))
}

func (h *ArtistHandler) DeleteArtist(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteArtist")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artist id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err := h.artistUsecase.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artist not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete artist")
	}

	return c.NoContent(http.StatusNoContent)
}

func toArtistResponse(artist *domain.Artist) artistResponse {
	return artistResponse{
		ID:         artist.ID.String(),
		Name:       artist.Name,
		Notes:      artist.Notes,
		ArtistLink: artist.ArtistLink,
		CreatedAt:  artist.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  artist.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
