package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/devi/booklet/internal/bookleaf"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type spyUserRepository struct {
	setPendingDeletionErr         error
	deleteAllUserDataKeys         []string
	deleteAllUserDataErr          error
	deleteExpiredTombstonesErr    error
	setPendingDeletionCalled      bool
	deleteExpiredTombstonesCalled bool
}

func (s *spyUserRepository) GetOrCreate(_ context.Context, idpSubject string) (*domain.User, error) {
	return &domain.User{ID: uuid.New(), IDPSubject: idpSubject}, nil
}

func (s *spyUserRepository) GetByIDPSubject(_ context.Context, _ string) (*domain.User, error) {
	return nil, errors.New("not found")
}

func (s *spyUserRepository) GetByID(_ context.Context, _ uuid.UUID) (*domain.User, error) {
	return nil, errors.New("not found")
}

func (s *spyUserRepository) SetPendingDeletion(_ context.Context, _ uuid.UUID) error {
	s.setPendingDeletionCalled = true
	return s.setPendingDeletionErr
}

func (s *spyUserRepository) DeleteAllUserData(_ context.Context, _ uuid.UUID) ([]string, error) {
	return s.deleteAllUserDataKeys, s.deleteAllUserDataErr
}

func (s *spyUserRepository) DeleteExpiredTombstones(_ context.Context) error {
	s.deleteExpiredTombstonesCalled = true
	return s.deleteExpiredTombstonesErr
}

type spyBookleafClient struct {
	deleteAccountErr error
}

func (s *spyBookleafClient) GetPublicFolders(_ context.Context, _ string) (*bookleaf.FolderList, error) {
	return nil, nil
}

func (s *spyBookleafClient) DeleteAccount(_ context.Context, _ string) error {
	return s.deleteAccountErr
}

func TestMarkPendingDeletion_BookleafSuccess(t *testing.T) {
	repo := &spyUserRepository{}
	bl := &spyBookleafClient{}
	uc := usecase.NewUserUsecase(repo, bl, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.MarkPendingDeletion(context.Background(), uuid.New(), "kp_user1")

	require.NoError(t, err)
	require.True(t, repo.setPendingDeletionCalled)
}

func TestMarkPendingDeletion_Bookleaf401_ReturnsConfigError(t *testing.T) {
	repo := &spyUserRepository{}
	bl := &spyBookleafClient{deleteAccountErr: bookleaf.ErrUnauthorized}
	uc := usecase.NewUserUsecase(repo, bl, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.MarkPendingDeletion(context.Background(), uuid.New(), "kp_user1")

	require.ErrorIs(t, err, usecase.ErrBookleafConfigError)
}

func TestCleanupExpiredTombstones_CallsRepo(t *testing.T) {
	repo := &spyUserRepository{}
	uc := usecase.NewUserUsecase(repo, &spyBookleafClient{}, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.CleanupExpiredTombstones(context.Background())

	require.NoError(t, err)
	require.True(t, repo.deleteExpiredTombstonesCalled)
}

func TestMarkPendingDeletion_BookleafNon2xx_ReturnsError(t *testing.T) {
	unexpectedErr := fmt.Errorf("%w: status 503", bookleaf.ErrUnexpectedStatus)
	repo := &spyUserRepository{}
	bl := &spyBookleafClient{deleteAccountErr: unexpectedErr}
	uc := usecase.NewUserUsecase(repo, bl, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	err := uc.MarkPendingDeletion(context.Background(), uuid.New(), "kp_user1")

	require.ErrorIs(t, err, bookleaf.ErrUnexpectedStatus)
	require.NotErrorIs(t, err, usecase.ErrBookleafConfigError)
}

func TestGetOrProvision_NewUser_CreatesWithIDPSubject(t *testing.T) {
	repo := &spyUserRepository{}
	uc := usecase.NewUserUsecase(repo, &spyBookleafClient{}, &spyTransactor{}, observability.NewTelemetry(nil, nil, nil))

	user, err := uc.GetOrProvision(context.Background(), "kp_newuser")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "kp_newuser", user.IDPSubject)
	assert.NotEqual(t, uuid.Nil, user.ID)
}
