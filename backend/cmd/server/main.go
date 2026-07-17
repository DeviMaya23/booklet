package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httphandler "github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/config"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/repository"
	"github.com/devi/booklet/internal/storage"
	"github.com/devi/booklet/internal/usecase"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	otelgorm "gorm.io/plugin/opentelemetry/tracing"

	authmiddleware "github.com/devi/booklet/internal/handler/middleware"
)

type server struct {
	echo        *echo.Echo
	logger      *zap.Logger
	shutdownTel func(context.Context)
}

func newServer(ctx context.Context, cfg *config.Config, logger *zap.Logger) *server {
	e := initEcho(cfg)
	tel, shutdownTel := initTelemetry(ctx, cfg, e, logger)
	db := initDB(cfg, logger)
	initApp(ctx, cfg, db, tel, e, logger)
	return &server{
		echo:        e,
		logger:      logger,
		shutdownTel: shutdownTel,
	}
}

func (s *server) start(port string) error {
	return s.echo.Start(":" + port)
}

func (s *server) shutdown(ctx context.Context) {
	s.shutdownTel(ctx)
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Error("echo shutdown", zap.Error(err))
	}
	_ = s.logger.Sync()
}

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	logger, err := observability.NewLogger(cfg.Obs.LogFormat)
	if err != nil {
		panic(fmt.Errorf("init logger: %w", err))
	}

	srv := newServer(ctx, cfg, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		srv.logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.shutdown(shutdownCtx)
	}()

	if err := srv.start(cfg.Port); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server stopped", zap.Error(err))
	}
}

func initEcho(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Validator = httphandler.NewEchoValidator()
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderContentType,
			"X-Booklet-Bypass",
		},
		ExposeHeaders: []string{
			"X-Booklet-Maintenance",
		},
	}))
	return e
}

func initTelemetry(ctx context.Context, cfg *config.Config, e *echo.Echo, logger *zap.Logger) (*observability.Telemetry, func(context.Context)) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error("otel error", zap.Error(err))
	}))

	if !cfg.Obs.OTELEnabled {
		return observability.NewTelemetry(logger, nil, nil), func(context.Context) {}
	}

	tp, err := observability.NewTracerProvider(ctx, cfg.Obs.OTELExporter, cfg.Obs.SampleRatio)
	if err != nil {
		logger.Fatal("init tracer provider", zap.Error(err))
	}

	mp, metricsHandler, err := observability.NewMeterProvider(cfg.Obs.OTELMetricsExporter)
	if err != nil {
		logger.Fatal("init meter provider", zap.Error(err))
	}

	tel := observability.NewTelemetry(logger, otel.Tracer("booklet"), otel.Meter("booklet"))
	e.Use(observability.TraceMiddleware(otel.Tracer("booklet")))
	e.Use(observability.MetricsMiddleware(otel.Meter("booklet")))
	if metricsHandler != nil {
		e.GET("/metrics", echo.WrapHandler(metricsHandler))
	}

	shutdown := func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("tracer shutdown", zap.Error(err))
		}
		if err := mp.Shutdown(ctx); err != nil {
			logger.Error("meter shutdown", zap.Error(err))
		}
	}

	return tel, shutdown
}

func initDB(cfg *config.Config, logger *zap.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DB.URL), &gorm.Config{
		Logger: repository.NewZapGORMLogger(logger),
	})
	if err != nil {
		logger.Fatal("open database connection", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("get underlying sql.DB", zap.Error(err))
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if cfg.Obs.OTELEnabled {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			logger.Fatal("register otelgorm plugin", zap.Error(err))
		}
	}

	return db
}

func initApp(ctx context.Context, cfg *config.Config, db *gorm.DB, tel *observability.Telemetry, e *echo.Echo, logger *zap.Logger) {
	r2Storage := storage.NewR2Storage(cfg.R2, tel)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, tel)

	characterRepository := repository.NewCharacterRepository(db)
	characterUsecase := usecase.NewCharacterUsecase(characterRepository, tel)
	characterHandler := httphandler.NewCharacterHandler(characterUsecase, tel)

	imageRepository := repository.NewImageRepository(db)
	imageUsecase := usecase.NewImageUsecase(imageRepository, tel)
	imageHandler := httphandler.NewImageHandler(imageUsecase, tel)

	uploadRepository := repository.NewUploadRepository(db)
	uploadUsecase := usecase.NewUploadUsecase(uploadRepository, r2Storage, characterRepository, tel)
	uploadHandler := httphandler.NewUploadHandler(uploadUsecase, tel)

	authMiddleware, err := authmiddleware.NewAuthMiddleware(cfg.Kinde.IssuerURL, cfg.Kinde.Audience, userUsecase, logger)
	if err != nil {
		logger.Fatal("initialise auth middleware", zap.Error(err))
	}

	healthHandler := httphandler.NewHealthHandler(db, r2Storage)

	e.GET("/health", healthHandler.GetHealth)
	protected := e.Group("")
	// protected.Use(authmiddleware.NewMaintenanceMiddleware(cfg.Maintenance))
	protected.Use(authMiddleware)
	protected.POST("/characters", characterHandler.CreateCharacter)
	protected.GET("/characters", characterHandler.ListCharacters)
	protected.GET("/characters/:id", characterHandler.GetCharacterByID)
	protected.PATCH("/characters/:id", characterHandler.UpdateCharacter)
	protected.DELETE("/characters/:id", characterHandler.DeleteCharacter)

	protected.GET("/images", imageHandler.ListImages)
	protected.GET("/images/:id", imageHandler.GetImageByID)
	protected.PATCH("/images/:id", imageHandler.UpdateImage)
	protected.DELETE("/images/:id", imageHandler.DeleteImage)
	protected.POST("/images", uploadHandler.InitialUpload)
	protected.POST("/images/:id/complete", uploadHandler.CompleteUpload)
}
