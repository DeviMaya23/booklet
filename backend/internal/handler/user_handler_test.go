package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/devi/booklet/internal/worker"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyUserUsecase struct {
	markPendingDeletionErr  error
	purgeUserDataKeys       []string
	purgeUserDataErr        error
}

func (s *spyUserUsecase) MarkPendingDeletion(_ context.Context, _ string) error {
	return s.markPendingDeletionErr
}

func (s *spyUserUsecase) PurgeUserData(_ context.Context, _ string) ([]string, error) {
	return s.purgeUserDataKeys, s.purgeUserDataErr
}

type spyJobInserter struct {
	lastArgs river.JobArgs
	insertErr error
}

func (s *spyJobInserter) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	s.lastArgs = args
	return nil, s.insertErr
}

// --- DeleteMe ---

func TestDeleteMe_Success(t *testing.T) {
	spy := &spyUserUsecase{}
	h := handler.NewUserHandler(spy, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.DELETE("/me", h.DeleteMe)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusAccepted, rec.Code)
}

func TestDeleteMe_ConfigError_Returns500(t *testing.T) {
	spy := &spyUserUsecase{markPendingDeletionErr: usecase.ErrBookleafConfigError}
	h := handler.NewUserHandler(spy, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.DELETE("/me", h.DeleteMe)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteMe_BookleafError_Returns502(t *testing.T) {
	spy := &spyUserUsecase{markPendingDeletionErr: errors.New("bookleaf unavailable")}
	h := handler.NewUserHandler(spy, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	e := setupEcho("user-1")
	e.DELETE("/me", h.DeleteMe)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadGateway, rec.Code)
}

// --- DeleteUserByID ---

func TestDeleteUserByID_Success(t *testing.T) {
	keys := []string{"users/u1/images/img.jpg", "users/u1/images/thumb.jpg"}
	spy := &spyUserUsecase{purgeUserDataKeys: keys}
	jobSpy := &spyJobInserter{}
	h := handler.NewUserHandler(spy, jobSpy, observability.NewTelemetry(nil, nil, nil))

	e := echo.New()
	e.DELETE("/internal/users/:id", h.DeleteUserByID)

	req := httptest.NewRequest(http.MethodDelete, "/internal/users/user-1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusAccepted, rec.Code)
	require.NotNil(t, jobSpy.lastArgs)
	got, ok := jobSpy.lastArgs.(worker.PurgeUserStorageArgs)
	require.True(t, ok)
	require.Equal(t, keys, got.R2Keys)
}

func TestDeleteUserByID_NotFound_Returns404(t *testing.T) {
	spy := &spyUserUsecase{purgeUserDataErr: gorm.ErrRecordNotFound}
	h := handler.NewUserHandler(spy, &spyJobInserter{}, observability.NewTelemetry(nil, nil, nil))

	e := echo.New()
	e.DELETE("/internal/users/:id", h.DeleteUserByID)

	req := httptest.NewRequest(http.MethodDelete, "/internal/users/unknown-user", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
