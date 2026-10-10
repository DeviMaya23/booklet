package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
)

type FolderUsecase interface {
	ListFolders(ctx context.Context, idpSubject string) (*bookleaf.FolderList, error)
	GetFolderImages(ctx context.Context, idpSubject, folderID string) (*bookleaf.FolderImageList, error)
}

type FolderHandler struct {
	folderUsecase FolderUsecase
	tel           *observability.Telemetry
}

func NewFolderHandler(folderUsecase FolderUsecase, tel *observability.Telemetry) *FolderHandler {
	return &FolderHandler{folderUsecase: folderUsecase, tel: tel}
}

func (h *FolderHandler) GetFolderImages(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetFolderImages")
	defer span.End()

	idpSubject, ok := middleware.AuthenticatedIDPSubjectFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	folderID, err := uuid.Parse(c.Param("folderID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid folder id")
	}

	result, err := h.folderUsecase.GetFolderImages(ctx, idpSubject, folderID.String())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, bookleaf.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "folder not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get folder images")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *FolderHandler) ListFolders(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListFolders")
	defer span.End()

	idpSubject, ok := middleware.AuthenticatedIDPSubjectFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	result, err := h.folderUsecase.ListFolders(ctx, idpSubject)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list folders")
	}

	return c.JSON(http.StatusOK, result)
}
