package handler

import (
	"context"
	"net/http"

	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/usecase"
	"github.com/labstack/echo/v4"
)

type FolderUsecase interface {
	ListFolders(ctx context.Context, userID string) (*usecase.FolderList, error)
}

type FolderHandler struct {
	folderUsecase FolderUsecase
}

func NewFolderHandler(folderUsecase FolderUsecase) *FolderHandler {
	return &FolderHandler{folderUsecase: folderUsecase}
}

func (h *FolderHandler) ListFolders(c echo.Context) error {
	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	result, err := h.folderUsecase.ListFolders(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list folders")
	}

	return c.JSON(http.StatusOK, result)
}
