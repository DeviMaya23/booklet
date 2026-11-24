package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	failIfCalled := func(name string) func(context.Context) error {
		return func(_ context.Context) error {
			t.Errorf("%s probe must not run for the shallow health check", name)
			return errors.New(name + " probed")
		}
	}

	tests := []struct {
		name string
		h    *HealthHandler
	}{
		{
			name: "process is up",
			h: &HealthHandler{
				dbProbe:       failIfCalled("db"),
				r2Probe:       failIfCalled("r2"),
				bookleafProbe: failIfCalled("bookleaf"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			require.NoError(t, tt.h.GetHealth(c))

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
		})
	}
}

func TestHealthHandler_GetDeepHealth(t *testing.T) {
	probeResult := func(err error) func(context.Context) error {
		return func(_ context.Context) error { return err }
	}

	tests := []struct {
		name         string
		dbErr        error
		r2Err        error
		bookleafErr  error
		wantCode     int
		wantStatus   string
		wantDB       string
		wantR2       string
		wantBookleaf string
	}{
		{
			name:         "all healthy",
			wantCode:     http.StatusOK,
			wantStatus:   "ok",
			wantDB:       "ok",
			wantR2:       "ok",
			wantBookleaf: "ok",
		},
		{
			name:         "db probe fails",
			dbErr:        errors.New("db timeout"),
			wantCode:     http.StatusServiceUnavailable,
			wantStatus:   "degraded",
			wantDB:       "db timeout",
			wantR2:       "ok",
			wantBookleaf: "ok",
		},
		{
			name:         "r2 probe fails",
			r2Err:        errors.New("r2 unavailable"),
			wantCode:     http.StatusServiceUnavailable,
			wantStatus:   "degraded",
			wantDB:       "ok",
			wantR2:       "r2 unavailable",
			wantBookleaf: "ok",
		},
		{
			name:         "only bookleaf probe fails",
			bookleafErr:  errors.New("connection refused"),
			wantCode:     http.StatusOK,
			wantStatus:   "ok",
			wantDB:       "ok",
			wantR2:       "ok",
			wantBookleaf: "connection refused",
		},
		{
			name:         "bookleaf and db probes fail",
			dbErr:        errors.New("db timeout"),
			bookleafErr:  errors.New("connection refused"),
			wantCode:     http.StatusServiceUnavailable,
			wantStatus:   "degraded",
			wantDB:       "db timeout",
			wantR2:       "ok",
			wantBookleaf: "connection refused",
		},
		{
			name:         "db and r2 probes fail",
			dbErr:        errors.New("db timeout"),
			r2Err:        errors.New("r2 unavailable"),
			wantCode:     http.StatusServiceUnavailable,
			wantStatus:   "degraded",
			wantDB:       "db timeout",
			wantR2:       "r2 unavailable",
			wantBookleaf: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/health/deep", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &HealthHandler{
				dbProbe:       probeResult(tt.dbErr),
				r2Probe:       probeResult(tt.r2Err),
				bookleafProbe: probeResult(tt.bookleafErr),
				probeTimeout:  time.Second,
			}

			require.NoError(t, h.GetDeepHealth(c))
			assert.Equal(t, tt.wantCode, rec.Code)

			var resp healthResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tt.wantStatus, resp.Status)
			assert.Equal(t, tt.wantDB, resp.DB)
			assert.Equal(t, tt.wantR2, resp.R2)
			assert.Equal(t, tt.wantBookleaf, resp.Bookleaf)
		})
	}
}

func TestHealthHandler_GetDeepHealth_SlowProbeDoesNotAffectOthers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/deep", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &HealthHandler{
		// Blocks until its own timeout fires.
		dbProbe: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
		// Return the state of their own contexts; a context shared with the slow
		// db probe would already be expired by the time these run.
		r2Probe:       func(ctx context.Context) error { return ctx.Err() },
		bookleafProbe: func(ctx context.Context) error { return ctx.Err() },
		probeTimeout:  20 * time.Millisecond,
	}

	require.NoError(t, h.GetDeepHealth(c))

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var resp healthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "degraded", resp.Status)
	assert.Equal(t, context.DeadlineExceeded.Error(), resp.DB)
	assert.Equal(t, "ok", resp.R2)
	assert.Equal(t, "ok", resp.Bookleaf)
}

func TestHealthHandler_GetDeepHealth_PanickingProbeIsReportedAsFailure(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health/deep", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &HealthHandler{
		dbProbe:       func(_ context.Context) error { return nil },
		r2Probe:       func(_ context.Context) error { panic("boom") },
		bookleafProbe: func(_ context.Context) error { return nil },
		probeTimeout:  time.Second,
	}

	require.NotPanics(t, func() { require.NoError(t, h.GetDeepHealth(c)) })

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var resp healthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "degraded", resp.Status)
	assert.Equal(t, "ok", resp.DB)
	assert.Equal(t, "probe panicked: boom", resp.R2)
	assert.Equal(t, "ok", resp.Bookleaf)
}
