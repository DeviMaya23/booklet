package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const defaultHealthProbeTimeout = 3 * time.Second

type r2Storage interface {
	Ping(ctx context.Context) error
}

type bookleafClient interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db            *gorm.DB
	store         r2Storage
	bookleaf      bookleafClient
	dbProbe       func(ctx context.Context) error
	r2Probe       func(ctx context.Context) error
	bookleafProbe func(ctx context.Context) error
	probeTimeout  time.Duration
}

type healthResponse struct {
	Status   string `json:"status"`
	DB       string `json:"db"`
	R2       string `json:"r2"`
	Bookleaf string `json:"bookleaf"`
}

type livenessResponse struct {
	Status string `json:"status"`
}

func NewHealthHandler(db *gorm.DB, store r2Storage, bl bookleafClient) *HealthHandler {
	h := &HealthHandler{
		db:           db,
		store:        store,
		bookleaf:     bl,
		probeTimeout: defaultHealthProbeTimeout,
	}

	h.dbProbe = func(ctx context.Context) error {
		if h.db == nil {
			return errors.New("db is not configured")
		}
		return h.db.WithContext(ctx).Exec("SELECT 1").Error
	}

	h.r2Probe = func(ctx context.Context) error {
		if h.store == nil {
			return errors.New("r2 storage is not configured")
		}
		return h.store.Ping(ctx)
	}

	h.bookleafProbe = func(ctx context.Context) error {
		if h.bookleaf == nil {
			return errors.New("bookleaf is not configured")
		}
		return h.bookleaf.Ping(ctx)
	}

	return h
}

// GetHealth is the shallow liveness check: it never touches a dependency.
func (h *HealthHandler) GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, livenessResponse{Status: "ok"})
}

// GetDeepHealth probes the database, R2 and Bookleaf in parallel, each with its
// own timeout. It returns 503 when the database or R2 probe fails. Bookleaf is
// reported in the body but does not affect the overall status.
func (h *HealthHandler) GetDeepHealth(c echo.Context) error {
	var dbErr, r2Err, bookleafErr error

	var wg sync.WaitGroup
	run := func(probe func(context.Context) error, out *error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Probes run outside Echo's per-request recovery, so a panic here
			// would take down the whole process. Report it as a probe failure.
			defer func() {
				if r := recover(); r != nil {
					*out = fmt.Errorf("probe panicked: %v", r)
				}
			}()
			ctx, cancel := context.WithTimeout(c.Request().Context(), h.probeTimeout)
			defer cancel()
			*out = probe(ctx)
		}()
	}

	run(h.dbProbe, &dbErr)
	run(h.r2Probe, &r2Err)
	run(h.bookleafProbe, &bookleafErr)
	wg.Wait()

	res := healthResponse{
		Status:   "ok",
		DB:       "ok",
		R2:       "ok",
		Bookleaf: "ok",
	}
	code := http.StatusOK

	if dbErr != nil {
		res.DB = dbErr.Error()
		res.Status = "degraded"
		code = http.StatusServiceUnavailable
	}
	if r2Err != nil {
		res.R2 = r2Err.Error()
		res.Status = "degraded"
		code = http.StatusServiceUnavailable
	}
	if bookleafErr != nil {
		res.Bookleaf = bookleafErr.Error()
	}

	return c.JSON(code, res)
}
