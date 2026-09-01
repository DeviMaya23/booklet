package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyUserUsecase struct {
	markPendingDeletionErr error
	purgeUserDataErr       error
}

func (s *spyUserUsecase) MarkPendingDeletion(_ context.Context, _ uuid.UUID, _ string) error {
	return s.markPendingDeletionErr
}

func (s *spyUserUsecase) PurgeUserData(_ context.Context, _ uuid.UUID) error {
	return s.purgeUserDataErr
}

// --- DeleteMe ---

func TestDeleteMe_Success(t *testing.T) {
	spy := &spyUserUsecase{}
	h := handler.NewUserHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/me", h.DeleteMe)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusAccepted, rec.Code)
}

func TestDeleteMe_ConfigError_Returns500(t *testing.T) {
	spy := &spyUserUsecase{markPendingDeletionErr: usecase.ErrBookleafConfigError}
	h := handler.NewUserHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/me", h.DeleteMe)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteMe_BookleafError_Returns502(t *testing.T) {
	spy := &spyUserUsecase{markPendingDeletionErr: errors.New("bookleaf unavailable")}
	h := handler.NewUserHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho(testUserID)
	e.DELETE("/me", h.DeleteMe)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadGateway, rec.Code)
}

// --- DeleteUserByID ---

func TestDeleteUserByID_Success(t *testing.T) {
	spy := &spyUserUsecase{}
	h := handler.NewUserHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := echo.New()
	e.DELETE("/internal/users/:id", h.DeleteUserByID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/internal/users/%s", testUserID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusAccepted, rec.Code)
}

func TestDeleteUserByID_NotFound_Returns404(t *testing.T) {
	spy := &spyUserUsecase{purgeUserDataErr: gorm.ErrRecordNotFound}
	h := handler.NewUserHandler(spy, observability.NewTelemetry(nil, nil, nil))

	e := echo.New()
	e.DELETE("/internal/users/:id", h.DeleteUserByID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/internal/users/%s", uuid.New()), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
