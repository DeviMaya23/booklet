package usecase_test

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- fakes ---

type fakeCommissionRepository struct {
	commissions map[uuid.UUID]*domain.Commission
}

func newFakeCommissionRepository() *fakeCommissionRepository {
	return &fakeCommissionRepository{commissions: make(map[uuid.UUID]*domain.Commission)}
}

func (f *fakeCommissionRepository) Create(_ context.Context, c *domain.Commission) (*domain.Commission, error) {
	f.commissions[c.ID] = c
	return c, nil
}

func (f *fakeCommissionRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Commission, error) {
	c, ok := f.commissions[id]
	if !ok || c.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return c, nil
}

func (f *fakeCommissionRepository) List(_ context.Context, userID uuid.UUID) ([]*domain.Commission, error) {
	var result []*domain.Commission
	for _, c := range f.commissions {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (f *fakeCommissionRepository) Update(_ context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateCommissionParams) (*domain.Commission, error) {
	c, ok := f.commissions[id]
	if !ok || c.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	c.Status = params.Status
	c.ArtistID = params.ArtistID
	c.Price = params.Price
	c.Paid = params.Paid
	c.Notes = params.Notes
	return c, nil
}

func (f *fakeCommissionRepository) Patch(_ context.Context, id uuid.UUID, userID uuid.UUID, params usecase.PatchCommissionParams) (*domain.Commission, error) {
	c, ok := f.commissions[id]
	if !ok || c.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	if params.Status != nil {
		c.Status = *params.Status
	}
	if params.Paid != nil {
		c.Paid = *params.Paid
	}
	return c, nil
}

func (f *fakeCommissionRepository) Delete(_ context.Context, id uuid.UUID, userID uuid.UUID) error {
	c, ok := f.commissions[id]
	if !ok || c.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	delete(f.commissions, id)
	return nil
}

type fakeCommissionArtistRepository struct {
	artists            map[uuid.UUID]*domain.Artist
	lastUsedAtUpdated  []uuid.UUID
}

func newFakeCommissionArtistRepository() *fakeCommissionArtistRepository {
	return &fakeCommissionArtistRepository{artists: make(map[uuid.UUID]*domain.Artist)}
}

func (f *fakeCommissionArtistRepository) GetByID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error) {
	a, ok := f.artists[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

func (f *fakeCommissionArtistRepository) UpdateLastUsedAt(_ context.Context, id uuid.UUID, _ uuid.UUID) error {
	f.lastUsedAtUpdated = append(f.lastUsedAtUpdated, id)
	return nil
}

type fakeCommissionCharacterRepository struct {
	characters map[uuid.UUID]*domain.Character
}

func newFakeCommissionCharacterRepository() *fakeCommissionCharacterRepository {
	return &fakeCommissionCharacterRepository{characters: make(map[uuid.UUID]*domain.Character)}
}

func (f *fakeCommissionCharacterRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Character, error) {
	var result []domain.Character
	for _, id := range ids {
		if c, ok := f.characters[id]; ok && c.UserID == userID {
			result = append(result, *c)
		}
	}
	return result, nil
}

type fakeCommissionArtpieceRepository struct {
	artpieces map[uuid.UUID]*domain.Artpiece
}

func newFakeCommissionArtpieceRepository() *fakeCommissionArtpieceRepository {
	return &fakeCommissionArtpieceRepository{artpieces: make(map[uuid.UUID]*domain.Artpiece)}
}

func (f *fakeCommissionArtpieceRepository) GetByIDsAndUserID(_ context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.Artpiece, error) {
	var result []*domain.Artpiece
	for _, id := range ids {
		if a, ok := f.artpieces[id]; ok && a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (f *fakeCommissionArtpieceRepository) GetByIDAndUserID(_ context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	a, ok := f.artpieces[id]
	if !ok || a.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

func (f *fakeCommissionArtpieceRepository) BulkUpdateCommissionID(_ context.Context, artpieceIDs []uuid.UUID, commissionID *uuid.UUID) error {
	for _, id := range artpieceIDs {
		if a, ok := f.artpieces[id]; ok {
			a.CommissionID = commissionID
		}
	}
	return nil
}

func (f *fakeCommissionArtpieceRepository) GetArtpiecesForCommission(_ context.Context, commissionID uuid.UUID) ([]*domain.Artpiece, error) {
	var result []*domain.Artpiece
	for _, a := range f.artpieces {
		if a.CommissionID != nil && *a.CommissionID == commissionID {
			result = append(result, a)
		}
	}
	return result, nil
}

func newCommissionUsecase(
	commissionRepo *fakeCommissionRepository,
	artistRepo *fakeCommissionArtistRepository,
	charRepo *fakeCommissionCharacterRepository,
	artpieceRepo *fakeCommissionArtpieceRepository,
) *usecase.CommissionUsecase {
	return usecase.NewCommissionUsecase(commissionRepo, artistRepo, charRepo, artpieceRepo, &fakeTransactor{}, observability.NewTelemetry(nil, nil, nil))
}

// --- Create tests ---

func TestCreateCommission_Success(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	commission, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{
		Status: usecase.CommissionStatusWaitlist,
	})

	require.NoError(t, err)
	assert.Equal(t, userID, commission.UserID)
	assert.Equal(t, usecase.CommissionStatusWaitlist, commission.Status)
}

func TestCreateCommission_DefaultStatusWaitlist(t *testing.T) {
	userID := uuid.New()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	commission, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{})

	require.NoError(t, err)
	assert.Equal(t, usecase.CommissionStatusWaitlist, commission.Status)
}

func TestCreateCommission_InvalidStatus(t *testing.T) {
	userID := uuid.New()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{Status: "invalid"})

	require.ErrorIs(t, err, usecase.ErrInvalidCommissionStatus)
}

func TestCreateCommission_ArtistNotOwned(t *testing.T) {
	userID := uuid.New()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())
	artistID := uuid.New()

	_, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{
		Status:   usecase.CommissionStatusWaitlist,
		ArtistID: &artistID,
	})

	require.ErrorIs(t, err, usecase.ErrArtistNotOwned)
}

func TestCreateCommission_CharacterNotOwned(t *testing.T) {
	userID := uuid.New()
	charRepo := newFakeCommissionCharacterRepository()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), charRepo, newFakeCommissionArtpieceRepository())

	_, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{
		Status:       usecase.CommissionStatusWaitlist,
		CharacterIDs: []uuid.UUID{uuid.New()},
	})

	require.ErrorIs(t, err, usecase.ErrCharacterNotOwned)
}

func TestCreateCommission_ArtpieceNotOwned(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeCommissionArtpieceRepository()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{
		Status:      usecase.CommissionStatusWaitlist,
		ArtpieceIDs: []uuid.UUID{uuid.New()},
	})

	require.ErrorIs(t, err, usecase.ErrArtpieceNotOwned)
}

func TestCreateCommission_ArtpieceAlreadyAttached(t *testing.T) {
	userID := uuid.New()
	artpieceRepo := newFakeCommissionArtpieceRepository()
	otherCommissionID := uuid.New()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &otherCommissionID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{
		Status:      usecase.CommissionStatusWaitlist,
		ArtpieceIDs: []uuid.UUID{artpiece.ID},
	})

	require.ErrorIs(t, err, usecase.ErrArtpieceAlreadyAttached)
}

// --- Update tests ---

func TestUpdateCommission_Success(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	updated, err := uc.Update(context.Background(), commission.ID, userID, usecase.UpdateCommissionParams{
		Status: usecase.CommissionStatusWIP,
	})

	require.NoError(t, err)
	assert.Equal(t, usecase.CommissionStatusWIP, updated.Status)
}

func TestUpdateCommission_ArtistNotOwned(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artistID := uuid.New()
	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.Update(context.Background(), commission.ID, userID, usecase.UpdateCommissionParams{
		Status:   usecase.CommissionStatusWaitlist,
		ArtistID: &artistID,
	})

	require.ErrorIs(t, err, usecase.ErrArtistNotOwned)
}

func TestUpdateCommission_CharacterNotOwned(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.Update(context.Background(), commission.ID, userID, usecase.UpdateCommissionParams{
		Status:       usecase.CommissionStatusWaitlist,
		CharacterIDs: []uuid.UUID{uuid.New()},
	})

	require.ErrorIs(t, err, usecase.ErrCharacterNotOwned)
}

func TestUpdateCommission_NotFound(t *testing.T) {
	userID := uuid.New()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.Update(context.Background(), uuid.New(), userID, usecase.UpdateCommissionParams{
		Status: usecase.CommissionStatusWaitlist,
	})

	require.Error(t, err)
}

// --- AttachArtpieces tests ---

func TestAttachArtpieces_Success(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.AttachArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{artpiece.ID})

	require.NoError(t, err)
	assert.Equal(t, &commission.ID, artpiece.CommissionID)
}

func TestAttachArtpieces_ArtpieceNotOwned(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.AttachArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{uuid.New()})

	require.ErrorIs(t, err, usecase.ErrArtpieceNotOwned)
}

func TestAttachArtpieces_ArtpieceAlreadyAttachedToDifferentCommission(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	otherCommissionID := uuid.New()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &otherCommissionID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.AttachArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{artpiece.ID})

	require.ErrorIs(t, err, usecase.ErrArtpieceAlreadyAttached)
}

func TestAttachArtpieces_AlreadyOnThisCommissionIsNoop(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &commission.ID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.AttachArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{artpiece.ID})

	require.NoError(t, err)
	assert.Equal(t, &commission.ID, artpiece.CommissionID)
}

// --- DetachArtpieces tests ---

func TestDetachArtpieces_Success(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &commission.ID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.DetachArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{artpiece.ID})

	require.NoError(t, err)
	assert.Nil(t, artpiece.CommissionID)
}

func TestDetachArtpieces_NotAttachedIsNoop(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.DetachArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{artpiece.ID})

	require.NoError(t, err)
	assert.Nil(t, artpiece.CommissionID)
}

// --- ReplaceArtpieces tests ---

func TestReplaceArtpieces_SuccessFullReplace(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	oldArtpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &commission.ID}
	newArtpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID}
	artpieceRepo.artpieces[oldArtpiece.ID] = oldArtpiece
	artpieceRepo.artpieces[newArtpiece.ID] = newArtpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.ReplaceArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{newArtpiece.ID})

	require.NoError(t, err)
	assert.Nil(t, oldArtpiece.CommissionID)
	assert.Equal(t, &commission.ID, newArtpiece.CommissionID)
}

func TestReplaceArtpieces_EmptySetClearsAll(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	artpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &commission.ID}
	artpieceRepo.artpieces[artpiece.ID] = artpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.ReplaceArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{})

	require.NoError(t, err)
	assert.Nil(t, artpiece.CommissionID)
}

func TestReplaceArtpieces_ArtpieceNotOwnedRejectsWithoutChanges(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	currentArtpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &commission.ID}
	artpieceRepo.artpieces[currentArtpiece.ID] = currentArtpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.ReplaceArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{uuid.New()})

	require.ErrorIs(t, err, usecase.ErrArtpieceNotOwned)
	assert.Equal(t, &commission.ID, currentArtpiece.CommissionID)
}

func TestReplaceArtpieces_ArtpieceBelongingToOtherCommissionRejects(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	artpieceRepo := newFakeCommissionArtpieceRepository()
	otherCommissionID := uuid.New()
	conflictingArtpiece := &domain.Artpiece{ID: uuid.New(), UserID: userID, CommissionID: &otherCommissionID}
	artpieceRepo.artpieces[conflictingArtpiece.ID] = conflictingArtpiece

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), artpieceRepo)

	_, err := uc.ReplaceArtpieces(context.Background(), commission.ID, userID, []uuid.UUID{conflictingArtpiece.ID})

	require.ErrorIs(t, err, usecase.ErrArtpieceAlreadyAttached)
	assert.Equal(t, &otherCommissionID, conflictingArtpiece.CommissionID)
}

// --- Patch tests ---

func TestPatchCommission_InvalidStatus(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())
	status := "invalid"

	_, err := uc.Patch(context.Background(), commission.ID, userID, usecase.PatchCommissionParams{
		Status: &status,
	})

	require.ErrorIs(t, err, usecase.ErrInvalidCommissionStatus)
}

func TestPatchCommission_NotFound(t *testing.T) {
	userID := uuid.New()
	uc := newCommissionUsecase(newFakeCommissionRepository(), newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())
	status := usecase.CommissionStatusWIP

	_, err := uc.Patch(context.Background(), uuid.New(), userID, usecase.PatchCommissionParams{
		Status: &status,
	})

	require.Error(t, err)
}

func TestPatchCommission_PatchStatus(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())
	status := usecase.CommissionStatusWIP

	updated, err := uc.Patch(context.Background(), commission.ID, userID, usecase.PatchCommissionParams{
		Status: &status,
	})

	require.NoError(t, err)
	assert.Equal(t, usecase.CommissionStatusWIP, updated.Status)
}

func TestPatchCommission_PatchPaidTrue(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWIP, Paid: false}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())
	paid := true

	updated, err := uc.Patch(context.Background(), commission.ID, userID, usecase.PatchCommissionParams{
		Paid: &paid,
	})

	require.NoError(t, err)
	assert.True(t, updated.Paid)
}

func TestCreateCommission_UpdatesArtistLastUsedAt(t *testing.T) {
	userID := uuid.New()
	artistRepo := newFakeCommissionArtistRepository()
	artistID := uuid.New()
	artistRepo.artists[artistID] = &domain.Artist{ID: artistID, UserID: userID}

	uc := newCommissionUsecase(newFakeCommissionRepository(), artistRepo, newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.Create(context.Background(), userID, usecase.CreateCommissionParams{
		Status:   usecase.CommissionStatusWaitlist,
		ArtistID: &artistID,
	})

	require.NoError(t, err)
	require.Contains(t, artistRepo.lastUsedAtUpdated, artistID)
}

func TestUpdateCommission_UpdatesArtistLastUsedAt(t *testing.T) {
	userID := uuid.New()
	artistRepo := newFakeCommissionArtistRepository()
	artistID := uuid.New()
	artistRepo.artists[artistID] = &domain.Artist{ID: artistID, UserID: userID}

	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWaitlist}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, artistRepo, newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())

	_, err := uc.Update(context.Background(), commission.ID, userID, usecase.UpdateCommissionParams{
		Status:   usecase.CommissionStatusWaitlist,
		ArtistID: &artistID,
	})

	require.NoError(t, err)
	require.Contains(t, artistRepo.lastUsedAtUpdated, artistID)
}

func TestPatchCommission_PatchPaidFalse(t *testing.T) {
	userID := uuid.New()
	commissionRepo := newFakeCommissionRepository()
	commission := &domain.Commission{ID: uuid.New(), UserID: userID, Status: usecase.CommissionStatusWIP, Paid: true}
	commissionRepo.commissions[commission.ID] = commission

	uc := newCommissionUsecase(commissionRepo, newFakeCommissionArtistRepository(), newFakeCommissionCharacterRepository(), newFakeCommissionArtpieceRepository())
	paid := false

	updated, err := uc.Patch(context.Background(), commission.ID, userID, usecase.PatchCommissionParams{
		Paid: &paid,
	})

	require.NoError(t, err)
	assert.False(t, updated.Paid)
}
