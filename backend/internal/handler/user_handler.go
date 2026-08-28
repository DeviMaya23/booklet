package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/devi/booklet/internal/worker"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserUsecase interface {
	MarkPendingDeletion(ctx context.Context, userID uuid.UUID, idpSubject string) error
	PurgeUserData(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type JobInserter interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

type UserHandler struct {
	userUsecase UserUsecase
	jobInserter JobInserter
	tel         *observability.Telemetry
}

func NewUserHandler(userUsecase UserUsecase, jobInserter JobInserter, tel *observability.Telemetry) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		jobInserter: jobInserter,
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

	keys, err := h.userUsecase.PurgeUserData(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to purge user data")
	}

	if len(keys) > 0 {
		_, err = h.jobInserter.Insert(ctx, worker.PurgeUserStorageArgs{R2Keys: keys}, nil)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to enqueue storage cleanup")
		}
		observability.LoggerFromContext(ctx, h.tel.Logger).Info(
			"storage cleanup enqueued",
			zap.String("event", "user.storage.cleanup.enqueued"),
			zap.String("user_id", id.String()),
			zap.Int("images_to_be_deleted", len(keys)),
		)
	}

	return c.NoContent(http.StatusAccepted)
}
