package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/worker"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

var ErrBookleafConfigError = errors.New("bookleaf: config error (unauthorized)")

type UserRepository interface {
	GetOrCreate(ctx context.Context, idpSubject string) (*domain.User, error)
	GetByIDPSubject(ctx context.Context, idpSubject string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	SetPendingDeletion(ctx context.Context, id uuid.UUID) error
	DeleteAllUserData(ctx context.Context, userID uuid.UUID) ([]string, error)
	DeleteExpiredTombstones(ctx context.Context) error
}

type userUsecase struct {
	userRepo       UserRepository
	bookleafClient BookleafClient
	transactor     Transactor
	jobInserter    JobInserter
	tel            *observability.Telemetry
}

func NewUserUsecase(userRepo UserRepository, bookleafClient BookleafClient, transactor Transactor, jobInserter JobInserter, tel *observability.Telemetry) *userUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		bookleafClient: bookleafClient,
		transactor:     transactor,
		jobInserter:    jobInserter,
		tel:            tel,
	}
}

func (u *userUsecase) GetOrProvision(ctx context.Context, kindeID string) (*domain.User, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetOrProvision")
	defer span.End()

	user, err := u.userRepo.GetByIDPSubject(ctx, kindeID)
	if err == nil {
		return user, nil
	}

	createdUser, err := u.userRepo.GetOrCreate(ctx, kindeID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	observability.LoggerFromContext(ctx, u.tel.Logger).Info(
		"user persisted",
		zap.String("event", "user.created"),
		zap.String("user_id", createdUser.ID.String()),
		zap.String("idp_subject", createdUser.IDPSubject),
	)
	return createdUser, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetByID")
	defer span.End()

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) MarkPendingDeletion(ctx context.Context, userID uuid.UUID, idpSubject string) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.MarkPendingDeletion")
	defer span.End()

	err := u.transactor.InTransaction(ctx, func(ctx context.Context) error {
		if err := u.userRepo.SetPendingDeletion(ctx, userID); err != nil {
			return err
		}
		return u.bookleafClient.DeleteAccount(ctx, idpSubject)
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, bookleaf.ErrUnauthorized) {
			err = fmt.Errorf("%w: %w", ErrBookleafConfigError, err)
		}
		return err
	}
	return nil
}

func (u *userUsecase) CleanupExpiredTombstones(ctx context.Context) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CleanupExpiredTombstones")
	defer span.End()

	if err := u.userRepo.DeleteExpiredTombstones(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (u *userUsecase) PurgeUserData(ctx context.Context, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.PurgeUserData")
	defer span.End()

	var keys []string
	err := u.transactor.InTransaction(ctx, func(ctx context.Context) error {
		var err error
		keys, err = u.userRepo.DeleteAllUserData(ctx, userID)
		return err
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	observability.LoggerFromContext(ctx, u.tel.Logger).Info(
		"user data purged",
		zap.String("event", "user.data.purged"),
		zap.String("user_id", userID.String()),
		zap.Int("images_to_be_deleted", len(keys)),
	)
	if len(keys) > 0 {
		if _, err := u.jobInserter.Insert(ctx, worker.PurgeUserStorageArgs{R2Keys: keys}, nil); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}
	return nil
}
