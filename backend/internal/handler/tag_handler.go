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

type TagUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateTagParams) (*domain.Tag, error)
	List(ctx context.Context, userID uuid.UUID, filters usecase.ListTagFilters) ([]*domain.Tag, error)
	ListWithUsage(ctx context.Context, userID uuid.UUID, filters usecase.ListTagFilters) ([]*domain.TagWithUsage, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateTagParams) (*domain.Tag, error)
	Merge(ctx context.Context, sourceID uuid.UUID, userID uuid.UUID, params usecase.MergeTagParams) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type TagHandler struct {
	tagUsecase TagUsecase
	tel        *observability.Telemetry
}

func NewTagHandler(tagUsecase TagUsecase, tel *observability.Telemetry) *TagHandler {
	return &TagHandler{tagUsecase: tagUsecase, tel: tel}
}

type createTagRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type updateTagRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type mergeTagRequest struct {
	TargetID uuid.UUID `json:"target_id" validate:"required"`
}

type tagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type tagWithUsageResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ArtpieceCount int    `json:"artpiece_count"`
}

func (h *TagHandler) CreateTag(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateTag")
	defer span.End()

	var req createTagRequest
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

	tag, err := h.tagUsecase.Create(ctx, userID, usecase.CreateTagParams{Name: req.Name})
	if err != nil {
		if errors.Is(err, usecase.ErrTagNameInvalid) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, usecase.ErrTagNameConflict) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create tag")
	}

	return c.JSON(http.StatusCreated, toTagResponse(tag))
}

func (h *TagHandler) ListTags(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListTags")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var filters usecase.ListTagFilters
	if err := c.Bind(&filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params")
	}
	if err := c.Validate(&filters); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	if filters.WithUsage {
		tags, err := h.tagUsecase.ListWithUsage(ctx, userID, filters)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to list tags")
		}
		responses := make([]tagWithUsageResponse, len(tags))
		for i, t := range tags {
			responses[i] = tagWithUsageResponse{ID: t.ID.String(), Name: t.Name, ArtpieceCount: t.ArtpieceCount}
		}
		return c.JSON(http.StatusOK, responses)
	}

	tags, err := h.tagUsecase.List(ctx, userID, filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list tags")
	}

	responses := make([]tagResponse, len(tags))
	for i, t := range tags {
		responses[i] = toTagResponse(t)
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *TagHandler) UpdateTag(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateTag")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag id")
	}

	var req updateTagRequest
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

	tag, err := h.tagUsecase.Update(ctx, id, userID, usecase.UpdateTagParams{Name: req.Name})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "tag not found")
		}
		if errors.Is(err, usecase.ErrTagNameInvalid) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, usecase.ErrTagNameConflict) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update tag")
	}

	return c.JSON(http.StatusOK, toTagResponse(tag))
}

func (h *TagHandler) MergeTag(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.MergeTag")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag id")
	}

	var req mergeTagRequest
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

	if err := h.tagUsecase.Merge(ctx, id, userID, usecase.MergeTagParams{TargetID: req.TargetID}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "tag not found")
		}
		if errors.Is(err, usecase.ErrTagMergeSameTag) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to merge tags")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TagHandler) DeleteTag(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteTag")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.tagUsecase.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "tag not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete tag")
	}

	return c.NoContent(http.StatusNoContent)
}

func toTagResponse(tag *domain.Tag) tagResponse {
	return tagResponse{ID: tag.ID.String(), Name: tag.Name}
}
