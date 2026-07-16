package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CharacterUsecase interface {
	Create(ctx context.Context, userID string, params usecase.CreateCharacterParams) (*domain.Character, error)
	GetByID(ctx context.Context, id, userID string) (*domain.Character, error)
	List(ctx context.Context, userID string) ([]*domain.Character, error)
	Update(ctx context.Context, id, userID string, params usecase.UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id, userID string) error
}

type CharacterHandler struct {
	characterUsecase CharacterUsecase
	tel              *observability.Telemetry
}

func NewCharacterHandler(characterUsecase CharacterUsecase, tel *observability.Telemetry) *CharacterHandler {
	return &CharacterHandler{characterUsecase: characterUsecase, tel: tel}
}

type createCharacterRequest struct {
	Name            string  `json:"name" validate:"required"`
	HeroImageR2Path *string `json:"hero_image_r2_path"`
	Biography       *string `json:"biography"`
	IsPublic        bool    `json:"is_public"`
}

type updateCharacterRequest struct {
	Name            *string `json:"name" validate:"omitempty,min=1"`
	HeroImageR2Path *string `json:"hero_image_r2_path"`
	Biography       *string `json:"biography"`
	IsPublic        *bool   `json:"is_public"`
}

type characterResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	HeroImageR2Path *string `json:"hero_image_r2_path"`
	Biography       *string `json:"biography"`
	IsPublic        bool    `json:"is_public"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func (h *CharacterHandler) CreateCharacter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateCharacter")
	defer span.End()

	var req createCharacterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	character, err := h.characterUsecase.Create(ctx, userID, usecase.CreateCharacterParams{
		Name:            req.Name,
		HeroImageR2Path: req.HeroImageR2Path,
		Biography:       req.Biography,
		IsPublic:        req.IsPublic,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create character")
	}

	return c.JSON(http.StatusCreated, toCharacterResponse(character))
}

func (h *CharacterHandler) GetCharacterByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetCharacterByID")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	character, err := h.characterUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get character")
	}

	return c.JSON(http.StatusOK, toCharacterResponse(character))
}

func (h *CharacterHandler) ListCharacters(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListCharacters")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	characters, err := h.characterUsecase.List(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list characters")
	}

	responses := make([]characterResponse, len(characters))
	for i, character := range characters {
		responses[i] = toCharacterResponse(character)
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *CharacterHandler) UpdateCharacter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateCharacter")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	var req updateCharacterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	character, err := h.characterUsecase.Update(ctx, id, userID, usecase.UpdateCharacterParams{
		Name:            req.Name,
		HeroImageR2Path: req.HeroImageR2Path,
		Biography:       req.Biography,
		IsPublic:        req.IsPublic,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update character")
	}

	return c.JSON(http.StatusOK, toCharacterResponse(character))
}

func (h *CharacterHandler) DeleteCharacter(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteCharacter")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err := h.characterUsecase.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete character")
	}

	return c.NoContent(http.StatusNoContent)
}

func toCharacterResponse(character *domain.Character) characterResponse {
	return characterResponse{
		ID:              character.ID.String(),
		Name:            character.Name,
		HeroImageR2Path: character.HeroImageR2Path,
		Biography:       character.Biography,
		IsPublic:        character.IsPublic,
		CreatedAt:       character.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       character.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
