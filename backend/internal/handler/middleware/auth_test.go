package middleware_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testIssuer   = "https://example.kinde.com"
	testAudience = "booklet"
	testKID      = "test-key-1"
)

type stubUserUsecase struct {
	user *domain.User
}

func (s *stubUserUsecase) GetOrProvision(_ context.Context, _ string) (*domain.User, error) {
	return s.user, nil
}

func buildTestStorage(t *testing.T) (*rsa.PrivateKey, jwkset.Storage) {
	t.Helper()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	store := jwkset.NewMemoryStorage()
	jwk, err := jwkset.NewJWKFromKey(privKey.Public(), jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{KID: testKID},
	})
	require.NoError(t, err)
	require.NoError(t, store.KeyWrite(context.Background(), jwk))

	return privKey, store
}

func signedToken(t *testing.T, privKey *rsa.PrivateKey, subject string) string {
	t.Helper()

	claims := jwt.RegisteredClaims{
		Issuer:    testIssuer,
		Audience:  jwt.ClaimStrings{testAudience},
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = testKID

	signed, err := token.SignedString(privKey)
	require.NoError(t, err)
	return signed
}

func buildMiddleware(t *testing.T, store jwkset.Storage, user *domain.User) echo.MiddlewareFunc {
	t.Helper()
	return middleware.NewAuthMiddlewareWithStorage(testIssuer, testAudience, store, &stubUserUsecase{user: user}, nil)
}

func TestAuth_PendingDeletionAccount_Returns401(t *testing.T) {
	privKey, store := buildTestStorage(t)
	user := &domain.User{ID: uuid.New(), IDPSubject: "sub-1", AccountState: domain.AccountStatePendingDeletion}
	mw := buildMiddleware(t, store, user)

	e := echo.New()
	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}, mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+signedToken(t, privKey, "sub-1"))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuth_PurgedAccount_Returns401(t *testing.T) {
	privKey, store := buildTestStorage(t)
	user := &domain.User{ID: uuid.New(), IDPSubject: "sub-2", AccountState: domain.AccountStatePurged}
	mw := buildMiddleware(t, store, user)

	e := echo.New()
	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}, mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+signedToken(t, privKey, "sub-2"))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuth_ActiveAccount_Passes(t *testing.T) {
	privKey, store := buildTestStorage(t)
	user := &domain.User{ID: uuid.New(), IDPSubject: "sub-3", AccountState: domain.AccountStateActive}
	mw := buildMiddleware(t, store, user)

	e := echo.New()
	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}, mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+signedToken(t, privKey, "sub-3"))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAuth_ActiveAccount_SetsUUIDAndIDPSubjectInContext(t *testing.T) {
	privKey, store := buildTestStorage(t)
	userID := uuid.New()
	user := &domain.User{ID: userID, IDPSubject: "sub-4", AccountState: domain.AccountStateActive}
	mw := buildMiddleware(t, store, user)

	var (
		gotUserID     uuid.UUID
		gotIDPSubject string
		userIDOK      bool
		idpSubjectOK  bool
	)

	e := echo.New()
	e.GET("/protected", func(c echo.Context) error {
		gotUserID, userIDOK = middleware.AuthenticatedUserIDFromContext(c)
		gotIDPSubject, idpSubjectOK = middleware.AuthenticatedIDPSubjectFromContext(c)
		return c.String(http.StatusOK, "ok")
	}, mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+signedToken(t, privKey, "sub-4"))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, userIDOK)
	assert.Equal(t, userID, gotUserID)
	assert.True(t, idpSubjectOK)
	assert.Equal(t, "sub-4", gotIDPSubject)
}
