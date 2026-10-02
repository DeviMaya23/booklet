package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FileUploadUsecase interface {
	InitiateUpload(ctx context.Context, params usecase.InitiateFileUploadParams) (*usecase.InitiateFileUploadResult, error)
	CompleteUpload(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error)
}

type FileUploadHandler struct {
	fileUploadUsecase FileUploadUsecase
	presigner         Presigner
	tel               *observability.Telemetry
}

func NewFileUploadHandler(fileUploadUsecase FileUploadUsecase, presigner Presigner, tel *observability.Telemetry) *FileUploadHandler {
	return &FileUploadHandler{fileUploadUsecase: fileUploadUsecase, presigner: presigner, tel: tel}
}

type initiateFileUploadRequest struct {
	MimeType   string     `json:"mime_type" validate:"required"`
	ArtpieceID *uuid.UUID `json:"artpiece_id"`
	Name       *string    `json:"name"`
	Notes      *string    `json:"notes"`
}

type initiateFileUploadResponse struct {
	ID        string `json:"id"`
	UploadURL string `json:"upload_url"`
	ExpiresAt string `json:"expires_at"`
}

type fileResponse struct {
	ID           string  `json:"id"`
	FileURL      *string `json:"file_url,omitempty"`
	MimeType     string  `json:"mime_type"`
	ThumbnailURL *string `json:"thumbnail_url"`
	ArtpieceID   *string `json:"artpiece_id"`
	Name         *string `json:"name"`
	Notes        *string `json:"notes"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func (h *FileUploadHandler) InitiateUpload(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.InitiateFileUpload")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req initiateFileUploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	result, err := h.fileUploadUsecase.InitiateUpload(ctx, usecase.InitiateFileUploadParams{
		UserID:     userID,
		MimeType:   req.MimeType,
		ArtpieceID: req.ArtpieceID,
		Name:       req.Name,
		Notes:      req.Notes,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrArtpieceNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "artpiece_id does not exist or does not belong to the user")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initiate upload")
	}

	return c.JSON(http.StatusCreated, initiateFileUploadResponse{
		ID:        result.ID.String(),
		UploadURL: result.UploadURL,
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *FileUploadHandler) CompleteUpload(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CompleteFileUpload")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file upload id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	file, err := h.fileUploadUsecase.CompleteUpload(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "pending upload not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to complete upload")
	}

	fileURL, _ := h.presignFileURL(ctx, file.FileR2Path)
	thumbnailURL, _ := h.presignThumbnailURL(ctx, file.ThumbnailR2Path)
	return c.JSON(http.StatusCreated, toFileResponse(file, fileURL, thumbnailURL))
}

func (h *FileUploadHandler) presignFileURL(ctx context.Context, r2Path string) (*string, error) {
	u, err := h.presigner.GeneratePresignedGetURL(ctx, r2Path, usecase.PresignGetTTL)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (h *FileUploadHandler) presignThumbnailURL(ctx context.Context, r2Path *string) (*string, error) {
	if r2Path == nil {
		return nil, nil
	}
	u, err := h.presigner.GeneratePresignedGetURL(ctx, *r2Path, usecase.PresignGetTTL)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func toFileResponse(f *domain.File, fileURL *string, thumbnailURL *string) fileResponse {
	var artpieceID *string
	if f.ArtpieceID != nil {
		s := f.ArtpieceID.String()
		artpieceID = &s
	}
	return fileResponse{
		ID:           f.ID.String(),
		FileURL:      fileURL,
		MimeType:     f.MimeType,
		ThumbnailURL: thumbnailURL,
		ArtpieceID:   artpieceID,
		Name:         f.Name,
		Notes:        f.Notes,
		CreatedAt:    f.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    f.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
