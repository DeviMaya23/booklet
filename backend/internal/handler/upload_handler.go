package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UploadUsecase interface {
	InitialUpload(ctx context.Context, params usecase.InitialUploadParams) (*usecase.InitialUploadResult, error)
	CompleteUpload(ctx context.Context, pendingID uuid.UUID, userID uuid.UUID) error
}

type UploadHandler struct {
	uploadUsecase UploadUsecase
	tel           *observability.Telemetry
}

func NewUploadHandler(uploadUsecase UploadUsecase, tel *observability.Telemetry) *UploadHandler {
	return &UploadHandler{uploadUsecase: uploadUsecase, tel: tel}
}

type initialUploadRequest struct {
	MimeType     string      `json:"mime_type" validate:"required,oneof=image/jpeg image/png"`
	Title        *string     `json:"title"`
	ArtistID     *uuid.UUID  `json:"artist_id"`
	Notes        *string     `json:"notes"`
	CharacterIDs []uuid.UUID `json:"character_ids"`
}

type initialUploadResponse struct {
	ID        string `json:"id"`
	UploadURL string `json:"upload_url"`
	ExpiresAt string `json:"expires_at"`
}

func (h *UploadHandler) InitialUpload(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.InitialUpload")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req initialUploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	result, err := h.uploadUsecase.InitialUpload(ctx, usecase.InitialUploadParams{
		UserID:       userID,
		MimeType:     req.MimeType,
		Title:        req.Title,
		ArtistID:     req.ArtistID,
		Notes:        req.Notes,
		CharacterIDs: req.CharacterIDs,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrArtistNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "artist_id does not exist or does not belong to the user")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initiate upload")
	}

	return c.JSON(http.StatusCreated, initialUploadResponse{
		ID:        result.ID.String(),
		UploadURL: result.UploadURL,
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *UploadHandler) CompleteUpload(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CompleteUpload")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid pending upload id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.uploadUsecase.CompleteUpload(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "pending upload not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to complete upload")
	}

	return c.NoContent(http.StatusCreated)
}
