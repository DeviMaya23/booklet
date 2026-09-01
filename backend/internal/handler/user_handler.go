package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

type UserUsecase interface {
	MarkPendingDeletion(ctx context.Context, userID uuid.UUID, idpSubject string) error
	PurgeUserData(ctx context.Context, userID uuid.UUID) error
}

type UserHandler struct {
	userUsecase UserUsecase
	tel         *observability.Telemetry
}

func NewUserHandler(userUsecase UserUsecase, tel *observability.Telemetry) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		tel:         tel,
	}
}

func (h *UserHandler) DeleteMe(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteMe")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	idpSubject, ok := middleware.AuthenticatedIDPSubjectFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.userUsecase.MarkPendingDeletion(ctx, userID, idpSubject); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, usecase.ErrBookleafConfigError) {
			return echo.NewHTTPError(http.StatusInternalServerError, "account deletion configuration error")
		}
		return echo.NewHTTPError(http.StatusBadGateway, "failed to coordinate account deletion")
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *UserHandler) DeleteUserByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteUserByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	if err := h.userUsecase.PurgeUserData(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to purge user data")
	}

	return c.NoContent(http.StatusAccepted)
}
