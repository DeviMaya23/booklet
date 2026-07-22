package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type characterRepository struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *characterRepository {
	return &characterRepository{db: db}
}

func (r *characterRepository) Create(ctx context.Context, character *domain.Character) error {
	return dbFromContext(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Folders").Create(character).Error; err != nil {
			return err
		}
		if len(character.Folders) > 0 {
			return tx.Create(&character.Folders).Error
		}
		return nil
	})
}

func (r *characterRepository) GetByID(ctx context.Context, id, userID string) (*domain.Character, error) {
	var character domain.Character
	err := dbFromContext(ctx, r.db).
		Preload("Folders").
		Where("id = ? AND user_id = ?", id, userID).
		First(&character).Error
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}
	return &character, nil
}

func (r *characterRepository) List(ctx context.Context, userID string) ([]*domain.Character, error) {
	var characters []*domain.Character
	err := dbFromContext(ctx, r.db).
		Preload("Folders").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&characters).Error
	if err != nil {
		return nil, fmt.Errorf("list characters: %w", err)
	}
	return characters, nil
}

func (r *characterRepository) Update(ctx context.Context, id, userID string, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	updates := map[string]interface{}{}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.HeroImageR2Path != nil {
		updates["hero_image_r2_path"] = *params.HeroImageR2Path
	}
	if params.Biography != nil {
		updates["biography"] = *params.Biography
	}
	if params.IsPublic != nil {
		updates["is_public"] = *params.IsPublic
	}

	if params.FolderIDs != nil {
		return r.updateWithFolders(ctx, id, userID, updates, *params.FolderIDs)
	}

	if len(updates) > 0 {
		result := dbFromContext(ctx, r.db).
			Model(&domain.Character{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(updates)
		if result.Error != nil {
			return nil, fmt.Errorf("update character: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return r.GetByID(ctx, id, userID)
}

func (r *characterRepository) updateWithFolders(ctx context.Context, id, userID string, updates map[string]interface{}, folderIDs []uuid.UUID) (*domain.Character, error) {
	err := dbFromContext(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&domain.Character{}).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
			return fmt.Errorf("check character exists: %w", err)
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Where("character_id = ?", id).Delete(&domain.CharacterFolder{}).Error; err != nil {
			return fmt.Errorf("delete character folders: %w", err)
		}

		if len(folderIDs) > 0 {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				return fmt.Errorf("parse character id: %w", err)
			}
			folders := make([]domain.CharacterFolder, len(folderIDs))
			for i, fid := range folderIDs {
				folders[i] = domain.CharacterFolder{CharacterID: parsedID, FolderID: fid}
			}
			if err := tx.Create(&folders).Error; err != nil {
				return fmt.Errorf("insert character folders: %w", err)
			}
		}

		if len(updates) > 0 {
			result := tx.
				Model(&domain.Character{}).
				Where("id = ? AND user_id = ?", id, userID).
				Updates(updates)
			if result.Error != nil {
				return fmt.Errorf("update character: %w", result.Error)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id, userID)
}

func (r *characterRepository) GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID string) ([]domain.Character, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var characters []domain.Character
	err := dbFromContext(ctx, r.db).
		Where("id IN ? AND user_id = ?", ids, userID).
		Find(&characters).Error
	if err != nil {
		return nil, fmt.Errorf("get characters by ids: %w", err)
	}
	return characters, nil
}

func (r *characterRepository) Delete(ctx context.Context, id, userID string) error {
	return dbFromContext(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Character{})
		if result.Error != nil {
			return fmt.Errorf("delete character: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Where("character_id = ?", id).Delete(&domain.CharacterFolder{}).Error
	})
}
