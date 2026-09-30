package repository

import (
	"context"
	"fmt"
	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type artpieceRepository struct {
	db *gorm.DB
}

func NewArtpieceRepository(db *gorm.DB) *artpieceRepository {
	return &artpieceRepository{db: db}
}

func (r *artpieceRepository) Create(ctx context.Context, a *domain.Artpiece) (*domain.Artpiece, error) {
	if err := dbFromContext(ctx, r.db).Create(a).Error; err != nil {
		return nil, fmt.Errorf("insert artpiece: %w", err)
	}
	return r.GetByID(ctx, a.ID, a.UserID)
}

func (r *artpieceRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	var a domain.Artpiece
	err := dbFromContext(ctx, r.db).
		Preload("Artist").
		Preload("CoverFile").
		Preload("Characters").
		Where("id = ? AND user_id = ?", id, userID).
		First(&a).Error
	if err != nil {
		return nil, fmt.Errorf("get artpiece: %w", err)
	}
	return &a, nil
}

func (r *artpieceRepository) List(ctx context.Context, userID uuid.UUID, filters usecase.ListArtpieceFilters) ([]*domain.Artpiece, error) {
	var artpieces []*domain.Artpiece
	q := dbFromContext(ctx, r.db).
		Preload("Artist").
		Preload("CoverFile").
		Preload("Characters").
		Where("artpieces.user_id = ?", userID)

	if len(filters.CharacterIDs) > 0 {
		q = q.Distinct().
			Joins("JOIN artpiece_characters ON artpiece_characters.artpiece_id = artpieces.id").
			Where("artpiece_characters.character_id::text IN ?", filters.CharacterIDs)
	}
	if len(filters.ArtistIDs) > 0 {
		q = q.Where("artpieces.artist_id::text IN ?", filters.ArtistIDs)
	}

	err := q.Order("artpieces.created_at DESC").Find(&artpieces).Error
	if err != nil {
		return nil, fmt.Errorf("list artpieces: %w", err)
	}
	return artpieces, nil
}

func (r *artpieceRepository) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateArtpieceParams) (*domain.Artpiece, error) {
	var a domain.Artpiece
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&a).Error
	if err != nil {
		return nil, fmt.Errorf("get artpiece: %w", err)
	}

	var title, notes interface{}
	if params.Title != nil {
		title = *params.Title
	}
	if params.Notes != nil {
		notes = *params.Notes
	}
	updates := map[string]interface{}{
		"title": title,
		"notes": notes,
	}

	if params.ArtistID == nil {
		updates["artist_id"] = nil
	} else {
		var count int64
		dbFromContext(ctx, r.db).Model(&domain.Artist{}).
			Where("id = ? AND user_id = ?", *params.ArtistID, userID).
			Count(&count)
		if count == 0 {
			return nil, usecase.ErrArtistNotOwned
		}
		updates["artist_id"] = *params.ArtistID
	}

	if err := dbFromContext(ctx, r.db).Model(&a).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update artpiece: %w", err)
	}

	charIDs := params.CharacterIDs
	if len(charIDs) > 0 {
		var count int64
		dbFromContext(ctx, r.db).Model(&domain.Character{}).
			Where("id IN ? AND user_id = ?", charIDs, userID).
			Count(&count)
		if count != int64(len(charIDs)) {
			return nil, usecase.ErrCharacterNotOwned
		}
	}

	characters := make([]domain.Character, len(charIDs))
	for i, cid := range charIDs {
		characters[i] = domain.Character{ID: cid}
	}
	if err := dbFromContext(ctx, r.db).Model(&a).Association("Characters").Replace(characters); err != nil {
		return nil, fmt.Errorf("replace characters: %w", err)
	}

	return r.GetByID(ctx, id, userID)
}

func (r *artpieceRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Artpiece{})
	if result.Error != nil {
		return fmt.Errorf("delete artpiece: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *artpieceRepository) UpdateCover(ctx context.Context, artpieceID uuid.UUID, coverFileID *uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Model(&domain.Artpiece{}).
		Where("id = ?", artpieceID).
		Update("cover_file_id", coverFileID)
	if result.Error != nil {
		return fmt.Errorf("update cover: %w", result.Error)
	}
	return nil
}

func (r *artpieceRepository) GetFilesForArtpiece(ctx context.Context, artpieceID uuid.UUID) ([]*domain.File, error) {
	var files []*domain.File
	err := dbFromContext(ctx, r.db).
		Where("artpiece_id = ?", artpieceID).
		Order("created_at ASC").
		Find(&files).Error
	if err != nil {
		return nil, fmt.Errorf("get files for artpiece: %w", err)
	}
	return files, nil
}

func (r *artpieceRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error) {
	var a domain.Artpiece
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&a).Error
	if err != nil {
		return nil, fmt.Errorf("get artpiece: %w", err)
	}
	return &a, nil
}

