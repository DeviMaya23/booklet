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

type VerseUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateVerseParams) (*domain.Verse, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Verse, error)
	List(ctx context.Context, userID uuid.UUID, filters usecase.ListVerseFilters) ([]*domain.Verse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateVerseParams) (*domain.Verse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type VerseHandler struct {
	verseUsecase VerseUsecase
	tel          *observability.Telemetry
}

func NewVerseHandler(verseUsecase VerseUsecase, tel *observability.Telemetry) *VerseHandler {
	return &VerseHandler{verseUsecase: verseUsecase, tel: tel}
}

type createVerseRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type updateVerseRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type verseResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func (h *VerseHandler) CreateVerse(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateVerse")
	defer span.End()

	var req createVerseRequest
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

	verse, err := h.verseUsecase.Create(ctx, userID, usecase.CreateVerseParams{Name: req.Name})
	if err != nil {
		if errors.Is(err, usecase.ErrVerseNameConflict) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create verse")
	}

	return c.JSON(http.StatusCreated, toVerseResponse(verse))
}

func (h *VerseHandler) GetVerseByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetVerseByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid verse id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	verse, err := h.verseUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "verse not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get verse")
	}

	return c.JSON(http.StatusOK, toVerseResponse(verse))
}

func (h *VerseHandler) ListVerses(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListVerses")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	verses, err := h.verseUsecase.List(ctx, userID, usecase.ListVerseFilters{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list verses")
	}

	responses := make([]verseResponse, len(verses))
	for i, v := range verses {
		responses[i] = toVerseResponse(v)
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *VerseHandler) UpdateVerse(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateVerse")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid verse id")
	}

	var req updateVerseRequest
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

	verse, err := h.verseUsecase.Update(ctx, id, userID, usecase.UpdateVerseParams{Name: req.Name})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "verse not found")
		}
		if errors.Is(err, usecase.ErrVerseNameConflict) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update verse")
	}

	return c.JSON(http.StatusOK, toVerseResponse(verse))
}

func (h *VerseHandler) DeleteVerse(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteVerse")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid verse id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.verseUsecase.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "verse not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete verse")
	}

	return c.NoContent(http.StatusNoContent)
}

func toVerseResponse(verse *domain.Verse) verseResponse {
	var lastUsedAt *string
	if verse.LastUsedAt != nil {
		s := verse.LastUsedAt.Format("2006-01-02T15:04:05Z07:00")
		lastUsedAt = &s
	}
	return verseResponse{
		ID:         verse.ID.String(),
		Name:       verse.Name,
		LastUsedAt: lastUsedAt,
		CreatedAt:  verse.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  verse.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
