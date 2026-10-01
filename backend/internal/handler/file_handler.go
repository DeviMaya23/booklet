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

type FileUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error)
	List(ctx context.Context, userID uuid.UUID, unassigned bool) ([]*domain.File, error)
	UpdateNotes(ctx context.Context, id uuid.UUID, userID uuid.UUID, notes *string) (*domain.File, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	BulkDelete(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error
}

type FileHandler struct {
	fileUsecase FileUsecase
	presigner   Presigner
	tel         *observability.Telemetry
}

func NewFileHandler(fileUsecase FileUsecase, presigner Presigner, tel *observability.Telemetry) *FileHandler {
	return &FileHandler{fileUsecase: fileUsecase, presigner: presigner, tel: tel}
}

type listFilesQuery struct {
	Unassigned bool `query:"unassigned"`
}

type updateFileRequest struct {
	Notes *string `json:"notes"`
}

type bulkDeleteFilesRequest struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

func (h *FileHandler) GetFile(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetFile")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	file, err := h.fileUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "file not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get file")
	}

	fileURL, _ := h.presignFileURL(ctx, file.FileR2Path)
	thumbnailURL, _ := h.presignThumbnailURL(ctx, file.ThumbnailR2Path)
	return c.JSON(http.StatusOK, toFileResponse(file, fileURL, thumbnailURL))
}

func (h *FileHandler) ListFiles(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListFiles")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var q listFilesQuery
	if err := c.Bind(&q); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params")
	}
	if err := c.Validate(&q); err != nil {
		return err
	}

	files, err := h.fileUsecase.List(ctx, userID, q.Unassigned)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list files")
	}

	responses := make([]fileResponse, len(files))
	for i, f := range files {
		fileURL, _ := h.presignFileURL(ctx, f.FileR2Path)
		thumbnailURL, _ := h.presignThumbnailURL(ctx, f.ThumbnailR2Path)
		responses[i] = toFileResponse(f, fileURL, thumbnailURL)
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *FileHandler) UpdateFile(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateFile")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req updateFileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	file, err := h.fileUsecase.UpdateNotes(ctx, id, userID, req.Notes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "file not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update file")
	}

	fileURL, _ := h.presignFileURL(ctx, file.FileR2Path)
	thumbnailURL, _ := h.presignThumbnailURL(ctx, file.ThumbnailR2Path)
	return c.JSON(http.StatusOK, toFileResponse(file, fileURL, thumbnailURL))
}

func (h *FileHandler) DeleteFile(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteFile")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.fileUsecase.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "file not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete file")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *FileHandler) BulkDeleteFiles(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.BulkDeleteFiles")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req bulkDeleteFilesRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, raw := range req.IDs {
		id, err := uuid.Parse(raw)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid file id: "+raw)
		}
		ids = append(ids, id)
	}

	if err := h.fileUsecase.BulkDelete(ctx, ids, userID); err != nil {
		if errors.Is(err, usecase.ErrFileOwnershipViolation) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "one or more file IDs do not belong to the user")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to bulk delete files")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *FileHandler) presignFileURL(ctx context.Context, r2Path string) (*string, error) {
	u, err := h.presigner.GeneratePresignedGetURL(ctx, r2Path, usecase.PresignGetTTL)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (h *FileHandler) presignThumbnailURL(ctx context.Context, r2Path *string) (*string, error) {
	if r2Path == nil {
		return nil, nil
	}
	u, err := h.presigner.GeneratePresignedGetURL(ctx, *r2Path, usecase.PresignGetTTL)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
