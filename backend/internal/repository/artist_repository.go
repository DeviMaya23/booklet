package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type artistRepository struct {
	db *gorm.DB
}

func NewArtistRepository(db *gorm.DB) *artistRepository {
	return &artistRepository{db: db}
}

func (r *artistRepository) Create(ctx context.Context, artist *domain.Artist) (*domain.Artist, error) {
	if err := dbFromContext(ctx, r.db).Create(artist).Error; err != nil {
		if isUniqueConstraintViolation(err) {
			return nil, usecase.ErrArtistNameConflict
		}
		return nil, fmt.Errorf("insert artist: %w", err)
	}
	return artist, nil
}

func (r *artistRepository) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Artist, error) {
	var artist domain.Artist
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&artist).Error
	if err != nil {
		return nil, fmt.Errorf("get artist: %w", err)
	}
	return &artist, nil
}

func (r *artistRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error) {
	var artist domain.Artist
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&artist).Error
	if err != nil {
		return nil, fmt.Errorf("get artist: %w", err)
	}
	return &artist, nil
}

func (r *artistRepository) List(ctx context.Context, userID uuid.UUID) ([]*domain.Artist, error) {
	var artists []*domain.Artist
	err := dbFromContext(ctx, r.db).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&artists).Error
	if err != nil {
		return nil, fmt.Errorf("list artists: %w", err)
	}
	return artists, nil
}

func (r *artistRepository) Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateArtistParams) (*domain.Artist, error) {
	updates := map[string]interface{}{}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.Notes != nil {
		updates["notes"] = *params.Notes
	}
	if params.ArtistLink != nil {
		updates["artist_link"] = *params.ArtistLink
	}

	if len(updates) > 0 {
		result := dbFromContext(ctx, r.db).
			Model(&domain.Artist{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(updates)
		if result.Error != nil {
			if isUniqueConstraintViolation(result.Error) {
				return nil, usecase.ErrArtistNameConflict
			}
			return nil, fmt.Errorf("update artist: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		var count int64
		if err := dbFromContext(ctx, r.db).Model(&domain.Artist{}).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("check artist exists: %w", err)
		}
		if count == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return r.GetByID(ctx, id, userID)
}

func (r *artistRepository) Delete(ctx context.Context, id string, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Artist{})
	if result.Error != nil {
		return fmt.Errorf("delete artist: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func isUniqueConstraintViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
