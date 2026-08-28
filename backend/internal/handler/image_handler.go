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

type ImageUsecase interface {
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Image, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Image, error)
	Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateImageParams) (*domain.Image, error)
	Delete(ctx context.Context, id string, userID uuid.UUID) error
}

type ImageHandler struct {
	imageUsecase ImageUsecase
	tel          *observability.Telemetry
}

func NewImageHandler(imageUsecase ImageUsecase, tel *observability.Telemetry) *ImageHandler {
	return &ImageHandler{imageUsecase: imageUsecase, tel: tel}
}

type updateImageRequest struct {
	Title           *string   `json:"title"`
	ThumbnailR2Path *string   `json:"thumbnail_r2_path"`
	ArtistName      *string   `json:"artist_name"`
	ArtistLink      *string   `json:"artist_link"`
	Notes           *string   `json:"notes"`
	CharacterIDs    *[]string `json:"character_ids"`
}

type characterRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type imageResponse struct {
	ID              string         `json:"id"`
	ImageR2Path     string         `json:"image_r2_path"`
	MimeType        string         `json:"mime_type"`
	Title           *string        `json:"title"`
	ThumbnailR2Path *string        `json:"thumbnail_r2_path"`
	ArtistName      *string        `json:"artist_name"`
	ArtistLink      *string        `json:"artist_link"`
	Notes           *string        `json:"notes"`
	Characters      []characterRef `json:"characters"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}

func (h *ImageHandler) GetImageByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetImageByID")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid image id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	image, err := h.imageUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "image not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get image")
	}

	return c.JSON(http.StatusOK, toImageResponse(image))
}

func (h *ImageHandler) ListImages(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListImages")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	images, err := h.imageUsecase.List(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list images")
	}

	responses := make([]imageResponse, len(images))
	for i, image := range images {
		responses[i] = toImageResponse(image)
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *ImageHandler) UpdateImage(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateImage")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid image id")
	}

	var req updateImageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	image, err := h.imageUsecase.Update(ctx, id, userID, usecase.UpdateImageParams{
		Title:           req.Title,
		ThumbnailR2Path: req.ThumbnailR2Path,
		ArtistName:      req.ArtistName,
		ArtistLink:      req.ArtistLink,
		Notes:           req.Notes,
		CharacterIDs:    req.CharacterIDs,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "image not found")
		}
		if errors.Is(err, usecase.ErrCharacterNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update image")
	}

	return c.JSON(http.StatusOK, toImageResponse(image))
}

func (h *ImageHandler) DeleteImage(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteImage")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid image id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err := h.imageUsecase.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "image not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete image")
	}

	return c.NoContent(http.StatusNoContent)
}

func toImageResponse(image *domain.Image) imageResponse {
	chars := make([]characterRef, len(image.Characters))
	for i, c := range image.Characters {
		chars[i] = characterRef{ID: c.ID.String(), Name: c.Name}
	}
	return imageResponse{
		ID:              image.ID.String(),
		ImageR2Path:     image.ImageR2Path,
		MimeType:        image.MimeType,
		Title:           image.Title,
		ThumbnailR2Path: image.ThumbnailR2Path,
		ArtistName:      image.ArtistName,
		ArtistLink:      image.ArtistLink,
		Notes:           image.Notes,
		Characters:      chars,
		CreatedAt:       image.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       image.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
