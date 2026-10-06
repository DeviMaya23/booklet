package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/pkg/timeutil"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type CommissionUsecase struct {
	commissionRepo CommissionRepository
	artistRepo     CommissionArtistRepository
	characterRepo  CommissionCharacterRepository
	artpieceRepo   CommissionArtpieceRepository
	transactor     Transactor
	tel            *observability.Telemetry
}

func NewCommissionUsecase(
	commissionRepo CommissionRepository,
	artistRepo CommissionArtistRepository,
	characterRepo CommissionCharacterRepository,
	artpieceRepo CommissionArtpieceRepository,
	transactor Transactor,
	tel *observability.Telemetry,
) *CommissionUsecase {
	return &CommissionUsecase{
		commissionRepo: commissionRepo,
		artistRepo:     artistRepo,
		characterRepo:  characterRepo,
		artpieceRepo:   artpieceRepo,
		transactor:     transactor,
		tel:            tel,
	}
}

func (u *CommissionUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateCommissionParams) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateCommission")
	defer span.End()

	status := params.Status
	if status == "" {
		status = CommissionStatusWaitlist
	}
	if !isValidCommissionStatus(status) {
		return nil, ErrInvalidCommissionStatus
	}

	if params.ArtistID != nil {
		if _, err := u.artistRepo.GetByID(ctx, *params.ArtistID, userID); err != nil {
			return nil, ErrArtistNotOwned
		}
	}

	if len(params.CharacterIDs) > 0 {
		chars, err := u.characterRepo.GetByIDsAndUserID(ctx, params.CharacterIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if len(chars) != len(params.CharacterIDs) {
			return nil, ErrCharacterNotOwned
		}
	}

	artpieceIDs := dedupUUIDs(params.ArtpieceIDs)
	if len(artpieceIDs) > 0 {
		if err := u.validateArtpieceOwnershipAndConflict(ctx, artpieceIDs, userID, uuid.Nil); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	characters := make([]domain.Character, len(params.CharacterIDs))
	for i, cid := range params.CharacterIDs {
		characters[i] = domain.Character{ID: cid}
	}

	commission := &domain.Commission{
		ID:         uuid.New(),
		UserID:     userID,
		Title:      params.Title,
		ArtistID:   params.ArtistID,
		Status:     status,
		Price:      params.Price,
		Paid:       params.Paid,
		PaidDate:   timeutil.ParseDate(params.PaidDate),
		FinishDate: timeutil.ParseDate(params.FinishDate),
		Notes:      params.Notes,
		Characters: characters,
	}

	created, err := u.commissionRepo.Create(ctx, commission)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(artpieceIDs) > 0 {
		err = u.transactor.InTransaction(ctx, func(ctx context.Context) error {
			return u.artpieceRepo.BulkUpdateCommissionID(ctx, artpieceIDs, &created.ID)
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	return u.commissionRepo.GetByID(ctx, created.ID, userID)
}

func (u *CommissionUsecase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetCommissionByID")
	defer span.End()

	res, err := u.commissionRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *CommissionUsecase) List(ctx context.Context, userID uuid.UUID) ([]*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListCommissions")
	defer span.End()

	res, err := u.commissionRepo.List(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *CommissionUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateCommissionParams) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateCommission")
	defer span.End()

	if !isValidCommissionStatus(params.Status) {
		return nil, ErrInvalidCommissionStatus
	}

	if params.ArtistID != nil {
		if _, err := u.artistRepo.GetByID(ctx, *params.ArtistID, userID); err != nil {
			return nil, ErrArtistNotOwned
		}
	}

	if len(params.CharacterIDs) > 0 {
		chars, err := u.characterRepo.GetByIDsAndUserID(ctx, params.CharacterIDs, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if len(chars) != len(params.CharacterIDs) {
			return nil, ErrCharacterNotOwned
		}
	}

	res, err := u.commissionRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *CommissionUsecase) Patch(ctx context.Context, id uuid.UUID, userID uuid.UUID, params PatchCommissionParams) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.PatchCommission")
	defer span.End()

	if params.Status != nil && !isValidCommissionStatus(*params.Status) {
		return nil, ErrInvalidCommissionStatus
	}

	res, err := u.commissionRepo.Patch(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *CommissionUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteCommission")
	defer span.End()

	if err := u.commissionRepo.Delete(ctx, id, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (u *CommissionUsecase) AttachArtpieces(ctx context.Context, commissionID uuid.UUID, userID uuid.UUID, artpieceIDs []uuid.UUID) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.AttachArtpieces")
	defer span.End()

	artpieceIDs = dedupUUIDs(artpieceIDs)

	if _, err := u.commissionRepo.GetByID(ctx, commissionID, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if err := u.validateArtpieceOwnershipAndConflict(ctx, artpieceIDs, userID, commissionID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Filter out artpieces already on this commission (no-op for those).
	var toAttach []uuid.UUID
	artpieces, err := u.artpieceRepo.GetByIDsAndUserID(ctx, artpieceIDs, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	for _, a := range artpieces {
		if a.CommissionID == nil || *a.CommissionID != commissionID {
			toAttach = append(toAttach, a.ID)
		}
	}

	if len(toAttach) > 0 {
		if err := u.artpieceRepo.BulkUpdateCommissionID(ctx, toAttach, &commissionID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	return u.commissionRepo.GetByID(ctx, commissionID, userID)
}

func (u *CommissionUsecase) DetachArtpieces(ctx context.Context, commissionID uuid.UUID, userID uuid.UUID, artpieceIDs []uuid.UUID) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DetachArtpieces")
	defer span.End()

	artpieceIDs = dedupUUIDs(artpieceIDs)

	if _, err := u.commissionRepo.GetByID(ctx, commissionID, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Only detach artpieces that actually belong to this commission.
	artpieces, err := u.artpieceRepo.GetByIDsAndUserID(ctx, artpieceIDs, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var toDetach []uuid.UUID
	for _, a := range artpieces {
		if a.CommissionID != nil && *a.CommissionID == commissionID {
			toDetach = append(toDetach, a.ID)
		}
	}

	if len(toDetach) > 0 {
		if err := u.artpieceRepo.BulkUpdateCommissionID(ctx, toDetach, nil); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	return u.commissionRepo.GetByID(ctx, commissionID, userID)
}

func (u *CommissionUsecase) ReplaceArtpieces(ctx context.Context, commissionID uuid.UUID, userID uuid.UUID, artpieceIDs []uuid.UUID) (*domain.Commission, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ReplaceArtpieces")
	defer span.End()

	artpieceIDs = dedupUUIDs(artpieceIDs)

	if _, err := u.commissionRepo.GetByID(ctx, commissionID, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(artpieceIDs) > 0 {
		if err := u.validateArtpieceOwnershipAndConflict(ctx, artpieceIDs, userID, commissionID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	currentArtpieces, err := u.artpieceRepo.GetArtpiecesForCommission(ctx, commissionID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	newSet := make(map[uuid.UUID]struct{}, len(artpieceIDs))
	for _, id := range artpieceIDs {
		newSet[id] = struct{}{}
	}
	currentSet := make(map[uuid.UUID]struct{}, len(currentArtpieces))
	for _, a := range currentArtpieces {
		currentSet[a.ID] = struct{}{}
	}

	var toAttach, toDetach []uuid.UUID
	for id := range newSet {
		if _, exists := currentSet[id]; !exists {
			toAttach = append(toAttach, id)
		}
	}
	for id := range currentSet {
		if _, exists := newSet[id]; !exists {
			toDetach = append(toDetach, id)
		}
	}

	err = u.transactor.InTransaction(ctx, func(ctx context.Context) error {
		if len(toAttach) > 0 {
			if err := u.artpieceRepo.BulkUpdateCommissionID(ctx, toAttach, &commissionID); err != nil {
				return err
			}
		}
		if len(toDetach) > 0 {
			if err := u.artpieceRepo.BulkUpdateCommissionID(ctx, toDetach, nil); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return u.commissionRepo.GetByID(ctx, commissionID, userID)
}

// validateArtpieceOwnershipAndConflict checks that all artpiece IDs are owned by the user
// and that none are attached to a different commission. Pass uuid.Nil as currentCommissionID
// when there is no "current" commission to allow (e.g., on Create).
func (u *CommissionUsecase) validateArtpieceOwnershipAndConflict(ctx context.Context, artpieceIDs []uuid.UUID, userID uuid.UUID, currentCommissionID uuid.UUID) error {
	artpieces, err := u.artpieceRepo.GetByIDsAndUserID(ctx, artpieceIDs, userID)
	if err != nil {
		return err
	}
	if len(artpieces) != len(artpieceIDs) {
		return ErrArtpieceNotOwned
	}
	for _, a := range artpieces {
		if a.CommissionID != nil && *a.CommissionID != currentCommissionID {
			return ErrArtpieceAlreadyAttached
		}
	}
	return nil
}
