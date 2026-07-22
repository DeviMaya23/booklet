package usecase_test

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeCharacterRepository struct {
	characters map[uuid.UUID]*domain.Character

	lastCreated *domain.Character
	lastUpdated usecase.UpdateCharacterParams
}

func newFakeCharacterRepository() *fakeCharacterRepository {
	return &fakeCharacterRepository{
		characters: make(map[uuid.UUID]*domain.Character),
	}
}

func (f *fakeCharacterRepository) Create(_ context.Context, character *domain.Character) error {
	f.characters[character.ID] = character
	f.lastCreated = character
	return nil
}

func (f *fakeCharacterRepository) GetByID(_ context.Context, id, _ string) (*domain.Character, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}
	c, ok := f.characters[parsed]
	if !ok {
		return nil, fmt.Errorf("get character: %w", gorm.ErrRecordNotFound)
	}
	return c, nil
}

func (f *fakeCharacterRepository) List(_ context.Context, userID string) ([]*domain.Character, error) {
	var result []*domain.Character
	for _, c := range f.characters {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (f *fakeCharacterRepository) Update(_ context.Context, id, _ string, params usecase.UpdateCharacterParams) (*domain.Character, error) {
	f.lastUpdated = params
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("update character: %w", err)
	}
	c, ok := f.characters[parsed]
	if !ok {
		return nil, fmt.Errorf("update character: %w", gorm.ErrRecordNotFound)
	}
	if params.Name != nil {
		c.Name = *params.Name
	}
	if params.HeroImageR2Path != nil {
		c.HeroImageR2Path = params.HeroImageR2Path
	}
	if params.Biography != nil {
		c.Biography = params.Biography
	}
	if params.IsPublic != nil {
		c.IsPublic = *params.IsPublic
	}
	if params.FolderIDs != nil {
		folders := make([]domain.CharacterFolder, len(*params.FolderIDs))
		for i, id := range *params.FolderIDs {
			folders[i] = domain.CharacterFolder{CharacterID: parsed, FolderID: id}
		}
		c.Folders = folders
	}
	return c, nil
}

func (f *fakeCharacterRepository) Delete(_ context.Context, id, _ string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("delete character: %w", err)
	}
	if _, ok := f.characters[parsed]; !ok {
		return fmt.Errorf("delete character: %w", gorm.ErrRecordNotFound)
	}
	delete(f.characters, parsed)
	return nil
}
