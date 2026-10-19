package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SavedFilterUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateSavedFilterParams) (*domain.SavedFilter, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.SavedFilter, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.SavedFilter, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateSavedFilterParams) (*domain.SavedFilter, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type SavedFilterHandler struct {
	uc        SavedFilterUsecase
	presigner Presigner
	tel       *observability.Telemetry
}

func NewSavedFilterHandler(uc SavedFilterUsecase, presigner Presigner, tel *observability.Telemetry) *SavedFilterHandler {
	return &SavedFilterHandler{uc: uc, presigner: presigner, tel: tel}
}

type createSavedFilterRequest struct {
	Name                string          `json:"name"               validate:"required"`
	ThumbnailArtpieceID *string         `json:"thumbnail_artpiece_id"`
	FilterPayload       json.RawMessage `json:"filter_payload"     validate:"required"`
}

type patchSavedFilterRequest struct {
	Name            usecase.Patch[string] `json:"name"`
	ThumbnailR2Path usecase.Patch[string] `json:"thumbnail_r2_path"`
	FilterPayload   *json.RawMessage      `json:"filter_payload"`
}

type savedFilterResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	ThumbnailURL *string         `json:"thumbnail_url"`
	FilterPayload json.RawMessage `json:"filter_payload"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

func (h *SavedFilterHandler) CreateSavedFilter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateSavedFilter")
	defer span.End()

	var req createSavedFilterRequest
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

	var thumbnailArtpieceID *uuid.UUID
	if req.ThumbnailArtpieceID != nil {
		parsed, err := uuid.Parse(*req.ThumbnailArtpieceID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid thumbnail_artpiece_id")
		}
		thumbnailArtpieceID = &parsed
	}

	sf, err := h.uc.Create(ctx, userID, usecase.CreateSavedFilterParams{
		Name:                req.Name,
		ThumbnailArtpieceID: thumbnailArtpieceID,
		FilterPayload:       req.FilterPayload,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicateSavedFilterName) {
			return c.JSON(http.StatusUnprocessableEntity, validationErrBody{
				Errors: []validationFieldError{{Field: "name", Message: "a saved filter with this name already exists"}},
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create saved filter")
	}

	resp, err := h.toSavedFilterResponse(ctx, sf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate response")
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *SavedFilterHandler) ListSavedFilters(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListSavedFilters")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	filters, err := h.uc.List(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list saved filters")
	}

	responses := make([]savedFilterResponse, 0, len(filters))
	for _, sf := range filters {
		resp, err := h.toSavedFilterResponse(ctx, sf)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate response")
		}
		responses = append(responses, resp)
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *SavedFilterHandler) GetSavedFilter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetSavedFilter")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid saved filter id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	sf, err := h.uc.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrSavedFilterNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "saved filter not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get saved filter")
	}

	resp, err := h.toSavedFilterResponse(ctx, sf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *SavedFilterHandler) PatchSavedFilter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.PatchSavedFilter")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid saved filter id")
	}

	var req patchSavedFilterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, validationErrResponse(err))
	}
	if req.Name.Set && req.Name.Value != nil && *req.Name.Value == "" {
		return c.JSON(http.StatusUnprocessableEntity, validationErrBody{Errors: []validationFieldError{{Field: "name", Message: "name must not be empty"}}})
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	sf, err := h.uc.Update(ctx, id, userID, usecase.UpdateSavedFilterParams{
		Name:            req.Name,
		ThumbnailR2Path: req.ThumbnailR2Path,
		FilterPayload:   req.FilterPayload,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrSavedFilterNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "saved filter not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update saved filter")
	}

	resp, err := h.toSavedFilterResponse(ctx, sf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *SavedFilterHandler) DeleteSavedFilter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteSavedFilter")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid saved filter id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.uc.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrSavedFilterNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "saved filter not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete saved filter")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *SavedFilterHandler) toSavedFilterResponse(ctx context.Context, sf *domain.SavedFilter) (savedFilterResponse, error) {
	var thumbnailURL *string
	if sf.ThumbnailR2Path != nil {
		u, err := h.presigner.GenerateDeterministicPresignedGetURL(ctx, *sf.ThumbnailR2Path)
		if err != nil {
			return savedFilterResponse{}, err
		}
		thumbnailURL = &u
	}
	return savedFilterResponse{
		ID:            sf.ID.String(),
		Name:          sf.Name,
		ThumbnailURL:  thumbnailURL,
		FilterPayload: sf.FilterPayload,
		CreatedAt:     sf.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     sf.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
