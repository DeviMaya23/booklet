package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devi/booklet/internal/bookleaf"
	httphandler "github.com/devi/booklet/internal/handler"
	authmiddleware "github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/config"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/repository"
	"github.com/devi/booklet/internal/storage"
	"github.com/devi/booklet/internal/usecase"
	"github.com/devi/booklet/internal/worker"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	rivertype "github.com/riverqueue/river/rivertype"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	otelgorm "gorm.io/plugin/opentelemetry/tracing"
)

type server struct {
	echo        *echo.Echo
	logger      *zap.Logger
	shutdownTel func(context.Context)
	riverPool   *pgxpool.Pool
	riverClient *river.Client[pgx.Tx]
}

func newServer(ctx context.Context, cfg *config.Config, logger *zap.Logger) *server {
	e := initEcho(cfg)
	tel, shutdownTel := initTelemetry(ctx, cfg, e, logger)
	db := initDB(cfg, logger)

	riverPool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Fatal("open river pgxpool", zap.Error(err))
	}

	riverClient := initApp(ctx, cfg, db, riverPool, tel, e, logger)
	if err := riverClient.Start(ctx); err != nil {
		logger.Fatal("start river client", zap.Error(err))
	}

	return &server{
		echo:        e,
		logger:      logger,
		shutdownTel: shutdownTel,
		riverPool:   riverPool,
		riverClient: riverClient,
	}
}

func (s *server) start(port string) error {
	return s.echo.Start(":" + port)
}

func (s *server) shutdown(ctx context.Context) {
	if err := s.riverClient.Stop(ctx); err != nil {
		s.logger.Error("river client stop", zap.Error(err))
	}
	s.riverPool.Close()
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

func initRiverClient(ctx context.Context, pool *pgxpool.Pool, workers *river.Workers, periodicJobs []*river.PeriodicJob, logger *zap.Logger) (*river.Client[pgx.Tx], error) {
	driver := riverpgxv5.New(pool)

	migrator, err := rivermigrate.New(driver, nil)
	if err != nil {
		return nil, fmt.Errorf("create river migrator: %w", err)
	}
	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		return nil, fmt.Errorf("run river migrations: %w", err)
	}
	logger.Info("river migrations applied")

	client, err := river.NewClient(driver, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 2},
		},
		Workers:      workers,
		PeriodicJobs: periodicJobs,
	})
	if err != nil {
		return nil, fmt.Errorf("create river client: %w", err)
	}

	return client, nil
}

// riverEnqueuer wraps river.Client[pgx.Tx] to satisfy usecase.JobInserter.
// Its client field is set after river.NewClient to break the
// uploadUsecase ↔ riverClient init cycle.
type riverEnqueuer struct {
	client *river.Client[pgx.Tx]
}

func (e *riverEnqueuer) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return e.client.Insert(ctx, args, opts)
}

func initApp(ctx context.Context, cfg *config.Config, db *gorm.DB, riverPool *pgxpool.Pool, tel *observability.Telemetry, e *echo.Echo, logger *zap.Logger) *river.Client[pgx.Tx] {
	r2Storage := storage.NewR2Storage(cfg.R2, tel)

	// Deferred enqueuer: client field is set after river.NewClient to break the
	// uploadUsecase ↔ riverClient init cycle.
	enqueuer := &riverEnqueuer{}

	bookleafClient := bookleaf.NewClient(cfg.Bookleaf.Host, cfg.Bookleaf.InternalSecret)
	folderUsecase := usecase.NewFolderUsecase(bookleafClient, tel)
	folderHandler := httphandler.NewFolderHandler(folderUsecase, tel)

	transactor := repository.NewGormTransactor(db)
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, bookleafClient, transactor, enqueuer, tel)

	characterRepository := repository.NewCharacterRepository(db)
	characterAvatarRepository := repository.NewCharacterAvatarRepository(db)

	artistRepository := repository.NewArtistRepository(db)
	artistUsecase := usecase.NewArtistUsecase(artistRepository, tel)
	artistHandler := httphandler.NewArtistHandler(artistUsecase, tel)

	imageRepository := repository.NewImageRepository(db)
	imageUsecase := usecase.NewImageUsecase(imageRepository, tel)
	imageHandler := httphandler.NewImageHandler(imageUsecase, r2Storage, tel)

	characterUsecase := usecase.NewCharacterUsecase(characterRepository, r2Storage, characterAvatarRepository, imageRepository, transactor, tel)
	characterHandler := httphandler.NewCharacterHandler(characterUsecase, r2Storage, tel)

	uploadRepository := repository.NewUploadRepository(db)
	uploadUsecase := usecase.NewUploadUsecase(uploadRepository, r2Storage, characterRepository, artistRepository, imageRepository, transactor, enqueuer, tel)
	uploadHandler := httphandler.NewUploadHandler(uploadUsecase, tel)

	authMiddleware, err := authmiddleware.NewAuthMiddleware(cfg.Kinde.IssuerURL, cfg.Kinde.Audience, userUsecase, logger)
	if err != nil {
		logger.Fatal("initialise auth middleware", zap.Error(err))
	}

	healthHandler := httphandler.NewHealthHandler(db, r2Storage)

	workers := river.NewWorkers()
	river.AddWorker(workers, worker.NewPurgeExpiredUploadsWorker(uploadUsecase, usecase.PresignTTL))
	river.AddWorker(workers, worker.NewPurgeExpiredCharacterAvatarUploadsWorker(characterUsecase, usecase.PresignTTL))
	river.AddWorker(workers, worker.NewPurgeUserStorageWorker(r2Storage, logger))
	river.AddWorker(workers, worker.NewPurgeTombstonesWorker(userUsecase))
	river.AddWorker(workers, worker.NewGenerateThumbnailWorker(imageRepository, r2Storage))

	periodicJobs := []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return worker.PurgeExpiredUploadsArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return worker.PurgeExpiredCharacterAvatarUploadsArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return worker.PurgeTombstonesArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: false},
		),
	}

	riverClient, err := initRiverClient(ctx, riverPool, workers, periodicJobs, logger)
	if err != nil {
		logger.Fatal("init river client", zap.Error(err))
	}
	enqueuer.client = riverClient

	userHandler := httphandler.NewUserHandler(userUsecase, tel)

	e.GET("/health", healthHandler.GetHealth)
	protected := e.Group("")
	// protected.Use(authmiddleware.NewMaintenanceMiddleware(cfg.Maintenance))
	protected.Use(authMiddleware)
	protected.POST("/characters", characterHandler.CreateCharacter)
	protected.GET("/characters", characterHandler.ListCharacters)
	protected.GET("/characters/:id", characterHandler.GetCharacterByID)
	protected.PATCH("/characters/:id", characterHandler.UpdateCharacter)
	protected.DELETE("/characters/:id", characterHandler.DeleteCharacter)
	protected.POST("/characters/:id/avatar/init", characterHandler.InitAvatarUpload)
	protected.POST("/characters/:id/avatar/:uploadID/complete", characterHandler.CompleteAvatarUpload)
	protected.DELETE("/characters/:id/avatar", characterHandler.DeleteAvatar)
	protected.GET("/characters/:id/images", characterHandler.GetCharacterImages)

	protected.GET("/folders", folderHandler.ListFolders)

	protected.POST("/artists", artistHandler.CreateArtist)
	protected.GET("/artists", artistHandler.ListArtists)
	protected.GET("/artists/:id", artistHandler.GetArtistByID)
	protected.PATCH("/artists/:id", artistHandler.UpdateArtist)
	protected.DELETE("/artists/:id", artistHandler.DeleteArtist)

	protected.GET("/images", imageHandler.ListImages)
	protected.GET("/images/:id", imageHandler.GetImageByID)
	protected.PATCH("/images/:id", imageHandler.UpdateImage)
	protected.DELETE("/images/:id", imageHandler.DeleteImage)
	protected.POST("/images", uploadHandler.InitialUpload)
	protected.POST("/images/:id/complete", uploadHandler.CompleteUpload)

	protected.DELETE("/me", userHandler.DeleteMe)

	internal := e.Group("")
	internal.Use(authmiddleware.NewInternalAuthMiddleware(cfg.BookletInternalSecret))
	internal.DELETE("/internal/users/:id", userHandler.DeleteUserByID)

	return riverClient
}
