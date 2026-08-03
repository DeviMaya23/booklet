package handler

import (
	"context"
	"net/http"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
)

type FolderUsecase interface {
	ListFolders(ctx context.Context, userID string) (*bookleaf.FolderList, error)
}

type FolderHandler struct {
	folderUsecase FolderUsecase
	tel           *observability.Telemetry
}

func NewFolderHandler(folderUsecase FolderUsecase, tel *observability.Telemetry) *FolderHandler {
	return &FolderHandler{folderUsecase: folderUsecase, tel: tel}
}

func (h *FolderHandler) ListFolders(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListFolders")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	result, err := h.folderUsecase.ListFolders(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list folders")
	}

	return c.JSON(http.StatusOK, result)
}
