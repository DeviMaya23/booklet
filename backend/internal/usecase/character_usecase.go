package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// validateAndEnrichFolders calls Bookleaf to validate the given folder IDs for the user.
// It returns a pre-enriched []domain.CharacterFolder with FolderID and FolderName populated.
// IDs absent from the Bookleaf response are silently dropped and logged at INFO level.
// If Bookleaf returns any error, the entire operation should be aborted (the error is returned).
func (u *characterUsecase) validateAndEnrichFolders(ctx context.Context, characterID uuid.UUID, folderIDs []uuid.UUID, idpSubject string) ([]domain.CharacterFolder, error) {
	folderList, err := u.bookleafClient.GetPublicFolders(ctx, idpSubject)
	if err != nil {
		return nil, fmt.Errorf("validate folders via bookleaf: %w", err)
	}

	nameByID := make(map[string]string, len(folderList.FolderList))
	for _, f := range folderList.FolderList {
		nameByID[f.FolderID] = f.FolderName
	}

	logger := observability.LoggerFromContext(ctx, u.tel.Logger)
	var enriched []domain.CharacterFolder
	for _, id := range folderIDs {
		name, ok := nameByID[id.String()]
		if !ok {
			logger.Info("folder ID not found in Bookleaf response, dropping",
				zap.String("event", "character.folder.dropped"),
				zap.String("folder_id", id.String()),
			)
			continue
		}
		enriched = append(enriched, domain.CharacterFolder{
			CharacterID: characterID,
			FolderID:    id,
			FolderName:  name,
		})
	}
	return enriched, nil
}

var (
	ErrCharacterNotFound      = errors.New("character not found or does not belong to the user")
	ErrPendingUploadNotFound  = errors.New("pending avatar upload not found or does not belong to the user")
)

type CreateCharacterParams struct {
	Name       string
	Notes      *string
	IsPublic   bool
	FolderIDs  *[]uuid.UUID
	IDPSubject string
}

type AvatarUploadResult struct {
	ID        uuid.UUID
	UploadURL string
	ExpiresAt time.Time
}

type characterUsecase struct {
	characterRepo    CharacterRepository
	storage          StorageService
	avatarUploadRepo CharacterAvatarUploadRepository
	imageRepo        ImageRepository
	transactor       Transactor
	tel              *observability.Telemetry
	bookleafClient   BookleafClient
}

func NewCharacterUsecase(
	characterRepo CharacterRepository,
	storage StorageService,
	avatarUploadRepo CharacterAvatarUploadRepository,
	imageRepo ImageRepository,
	transactor Transactor,
	tel *observability.Telemetry,
	bookleafClient BookleafClient,
) *characterUsecase {
	return &characterUsecase{
		characterRepo:    characterRepo,
		storage:          storage,
		avatarUploadRepo: avatarUploadRepo,
		imageRepo:        imageRepo,
		transactor:       transactor,
		tel:              tel,
		bookleafClient:   bookleafClient,
	}
}

func (u *characterUsecase) Create(ctx context.Context, userID uuid.UUID, params CreateCharacterParams) (*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CreateCharacter")
	defer span.End()

	character := &domain.Character{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     params.Name,
		Notes:    params.Notes,
		IsPublic: params.IsPublic,
	}
	if params.FolderIDs != nil && len(*params.FolderIDs) > 0 {
		enriched, err := u.validateAndEnrichFolders(ctx, character.ID, *params.FolderIDs, params.IDPSubject)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		character.Folders = enriched
	}
	if err := u.characterRepo.Create(ctx, character); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return character, nil
}

func (u *characterUsecase) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetCharacterByID")
	defer span.End()

	res, err := u.characterRepo.GetByID(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *characterUsecase) List(ctx context.Context, userID uuid.UUID, filters ListCharacterFilters) ([]*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.ListCharacters")
	defer span.End()

	res, err := u.characterRepo.List(ctx, userID, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *characterUsecase) Update(ctx context.Context, id string, userID uuid.UUID, params UpdateCharacterParams) (*domain.Character, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.UpdateCharacter")
	defer span.End()

	if len(params.FolderIDs) > 0 {
		charID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("parse character id: %w", err)
		}
		enriched, err := u.validateAndEnrichFolders(ctx, charID, params.FolderIDs, params.IDPSubject)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		params.Folders = enriched
	}
	// If FolderIDs is empty, params.Folders remains nil/empty → repo will clear all assignments

	res, err := u.characterRepo.Update(ctx, id, userID, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return res, nil
}

func (u *characterUsecase) Delete(ctx context.Context, id string, userID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteCharacter")
	defer span.End()

	err := u.characterRepo.Delete(ctx, id, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (u *characterUsecase) InitAvatarUpload(ctx context.Context, userID uuid.UUID, characterID string, mimeType string) (*AvatarUploadResult, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.InitAvatarUpload")
	defer span.End()

	character, err := u.characterRepo.GetByID(ctx, characterID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCharacterNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	id := uuid.New()
	r2Key := fmt.Sprintf("%s/files/%s", character.UserID.String(), id.String())
	expiresAt := time.Now().Add(PresignTTL)

	observability.LoggerFromContext(ctx, u.tel.Logger).Info("avatar upload initiated",
		zap.String("event", "r2.avatar.upload.started"),
		zap.String("pending_id", id.String()),
		zap.String("character_id", characterID),
		zap.String("user_id", userID.String()),
		zap.String("r2_key", r2Key),
	)

	uploadURL, err := u.storage.GeneratePresignedPutURL(ctx, r2Key, mimeType, PresignTTL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	charID, err := uuid.Parse(characterID)
	if err != nil {
		return nil, fmt.Errorf("parse character id: %w", err)
	}

	pending := &domain.PendingCharacterAvatarUpload{
		ID:          id,
		UserID:      userID,
		CharacterID: charID,
		R2Key:       r2Key,
		MimeType:    mimeType,
	}
	if _, err := u.avatarUploadRepo.Create(ctx, pending); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &AvatarUploadResult{ID: id, UploadURL: uploadURL, ExpiresAt: expiresAt}, nil
}

func (u *characterUsecase) CompleteAvatarUpload(ctx context.Context, userID uuid.UUID, characterID string, uploadID uuid.UUID) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.CompleteAvatarUpload")
	defer span.End()

	pending, err := u.avatarUploadRepo.GetByID(ctx, uploadID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPendingUploadNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if pending.CharacterID.String() != characterID {
		return ErrPendingUploadNotFound
	}

	character, err := u.characterRepo.GetByID(ctx, characterID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCharacterNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	oldKey := character.AvatarR2Path

	if err := u.transactor.InTransaction(ctx, func(txCtx context.Context) error {
		if err := u.characterRepo.UpdateAvatarR2Path(txCtx, characterID, userID, pending.R2Key); err != nil {
			return err
		}
		return u.avatarUploadRepo.Delete(txCtx, uploadID)
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if oldKey != nil {
		if err := u.storage.DeleteObject(ctx, *oldKey); err != nil {
			observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to delete old avatar from storage",
				zap.String("r2_key", *oldKey),
				zap.String("character_id", characterID),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (u *characterUsecase) DeleteAvatar(ctx context.Context, userID uuid.UUID, characterID string) error {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.DeleteAvatar")
	defer span.End()

	oldKey, err := u.characterRepo.ClearAvatarR2Path(ctx, characterID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCharacterNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if oldKey == "" {
		return nil
	}

	if err := u.storage.DeleteObject(ctx, oldKey); err != nil {
		observability.LoggerFromContext(ctx, u.tel.Logger).Error("failed to delete avatar from storage",
			zap.String("r2_key", oldKey),
			zap.String("character_id", characterID),
			zap.Error(err),
		)
	}

	return nil
}

func (u *characterUsecase) GetCharacterImages(ctx context.Context, characterID uuid.UUID, userID uuid.UUID) ([]*domain.Image, error) {
	ctx, span := u.tel.Tracer.Start(ctx, "usecase.GetCharacterImages")
	defer span.End()

	images, err := u.imageRepo.ListByCharacterID(ctx, characterID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return images, nil
}

func (u *characterUsecase) CleanupStaleAvatarUploads(ctx context.Context, threshold time.Duration) error {
	cutoff := time.Now().Add(-threshold)
	records, err := u.avatarUploadRepo.ListStale(ctx, cutoff)
	if err != nil {
		return err
	}

	deleted := 0
	for _, p := range records {
		if err := u.storage.DeleteObject(ctx, p.R2Key); err != nil {
			return err
		}
		if err := u.avatarUploadRepo.Delete(ctx, p.ID); err != nil {
			return err
		}
		deleted++
	}

	observability.LoggerFromContext(ctx, u.tel.Logger).Info("stale avatar uploads purged",
		zap.String("event", "r2.avatar.upload.purged"),
		zap.Int("count", deleted),
	)
	return nil
}
