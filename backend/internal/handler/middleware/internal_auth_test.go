package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const internalSecretHeader = "X-Booklet-Internal-Secret"

func newInternalAuthEcho(secretByPath map[string]string) (*echo.Echo, map[string]*int) {
	e := echo.New()
	calls := make(map[string]*int, len(secretByPath))
	for path, secret := range secretByPath {
		n := 0
		calls[path] = &n
		counter := &n
		e.GET(path, func(c echo.Context) error {
			*counter++
			return c.String(http.StatusOK, "reached")
		}, middleware.NewInternalAuthMiddleware(secret))
	}
	return e, calls
}

func TestInternalAuthMiddleware(t *testing.T) {
	const secret = "correct-secret"

	tests := []struct {
		name       string
		header     *string
		wantStatus int
		wantBody   string
		wantCalled int
	}{
		{
			name:       "missing header",
			header:     nil,
			wantStatus: http.StatusUnauthorized,
			wantCalled: 0,
		},
		{
			name:       "empty header",
			header:     strPtr(""),
			wantStatus: http.StatusUnauthorized,
			wantCalled: 0,
		},
		{
			name:       "wrong value",
			header:     strPtr("wrong-secret"),
			wantStatus: http.StatusUnauthorized,
			wantCalled: 0,
		},
		{
			name:       "header has the secret as a prefix",
			header:     strPtr(secret + "-extra"),
			wantStatus: http.StatusUnauthorized,
			wantCalled: 0,
		},
		{
			name:       "header is a prefix of the secret",
			header:     strPtr("correct-"),
			wantStatus: http.StatusUnauthorized,
			wantCalled: 0,
		},
		{
			name:       "correct value",
			header:     strPtr(secret),
			wantStatus: http.StatusOK,
			wantBody:   "reached",
			wantCalled: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, calls := newInternalAuthEcho(map[string]string{"/guarded": secret})

			req := httptest.NewRequest(http.MethodGet, "/guarded", nil)
			if tt.header != nil {
				req.Header.Set(internalSecretHeader, *tt.header)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, rec.Body.String())
			}
			assert.Equal(t, tt.wantCalled, *calls["/guarded"])
		})
	}
}

func TestInternalAuthMiddleware_SecretsAreIsolatedPerInstance(t *testing.T) {
	const (
		internalAPISecret = "internal-api-secret"
		healthSecret      = "health-check-secret"
	)

	e, calls := newInternalAuthEcho(map[string]string{
		"/internal": internalAPISecret,
		"/health":   healthSecret,
	})

	tests := []struct {
		name       string
		path       string
		secret     string
		wantStatus int
		wantCalled int
	}{
		{"internal secret on internal route", "/internal", internalAPISecret, http.StatusOK, 1},
		{"health secret on health route", "/health", healthSecret, http.StatusOK, 1},
		{"health secret on internal route", "/internal", healthSecret, http.StatusUnauthorized, 0},
		{"internal secret on health route", "/health", internalAPISecret, http.StatusUnauthorized, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*calls[tt.path] = 0

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set(internalSecretHeader, tt.secret)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantCalled, *calls[tt.path])
		})
	}
}

func strPtr(s string) *string { return &s }
