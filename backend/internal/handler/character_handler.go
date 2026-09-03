package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler/middleware"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CharacterUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateCharacterParams) (*domain.Character, error)
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Character, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Character, error)
	Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id string, userID uuid.UUID) error
	InitAvatarUpload(ctx context.Context, userID uuid.UUID, characterID string, mimeType string) (*usecase.AvatarUploadResult, error)
	CompleteAvatarUpload(ctx context.Context, userID uuid.UUID, characterID string, uploadID uuid.UUID) error
	DeleteAvatar(ctx context.Context, userID uuid.UUID, characterID string) error
	GetCharacterImages(ctx context.Context, characterID uuid.UUID, userID uuid.UUID) ([]*domain.Image, error)
}

type CharacterHandler struct {
	characterUsecase CharacterUsecase
	presigner        Presigner
	tel              *observability.Telemetry
}

func NewCharacterHandler(characterUsecase CharacterUsecase, presigner Presigner, tel *observability.Telemetry) *CharacterHandler {
	return &CharacterHandler{characterUsecase: characterUsecase, presigner: presigner, tel: tel}
}

type createCharacterRequest struct {
	Name      string    `json:"name" validate:"required"`
	Biography *string   `json:"biography"`
	IsPublic  bool      `json:"is_public"`
	FolderIDs *[]string `json:"folder_ids" validate:"omitempty,dive,uuid4"`
}

type updateCharacterRequest struct {
	Name      *string   `json:"name" validate:"omitempty,min=1"`
	Biography *string   `json:"biography"`
	IsPublic  *bool     `json:"is_public"`
	FolderIDs *[]string `json:"folder_ids" validate:"omitempty,dive,uuid4"`
}

type characterResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	AvatarURL *string  `json:"avatar_url"`
	Biography *string  `json:"biography"`
	IsPublic  bool     `json:"is_public"`
	FolderIDs []string `json:"folder_ids"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type characterImageResponse struct {
	ImageID      string  `json:"image_id"`
	ImageName    *string `json:"image_name"`
	ThumbnailURL *string `json:"thumbnail_url"`
}

type initAvatarUploadRequest struct {
	MimeType string `json:"mime_type" validate:"required,oneof=image/jpeg image/png"`
}

type initAvatarUploadResponse struct {
	ID        string `json:"id"`
	UploadURL string `json:"upload_url"`
	ExpiresAt string `json:"expires_at"`
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
		Name:      req.Name,
		Biography: req.Biography,
		IsPublic:  req.IsPublic,
		FolderIDs: parseFolderIDs(req.FolderIDs),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create character")
	}

	// presigning is a local crypto op; failure means context cancellation, not a broken character
	avatarURL, _ := h.presignAvatarURL(ctx, character.AvatarR2Path)
	return c.JSON(http.StatusCreated, toCharacterResponse(character, avatarURL))
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

	avatarURL, _ := h.presignAvatarURL(ctx, character.AvatarR2Path)
	return c.JSON(http.StatusOK, toCharacterResponse(character, avatarURL))
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
		avatarURL, _ := h.presignAvatarURL(ctx, character.AvatarR2Path)
		responses[i] = toCharacterResponse(character, avatarURL)
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
		Name:      req.Name,
		Biography: req.Biography,
		IsPublic:  req.IsPublic,
		FolderIDs: parseFolderIDs(req.FolderIDs),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update character")
	}

	avatarURL, _ := h.presignAvatarURL(ctx, character.AvatarR2Path)
	return c.JSON(http.StatusOK, toCharacterResponse(character, avatarURL))
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

func (h *CharacterHandler) InitAvatarUpload(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.InitAvatarUpload")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	var req initAvatarUploadRequest
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

	result, err := h.characterUsecase.InitAvatarUpload(ctx, userID, id, req.MimeType)
	if err != nil {
		if errors.Is(err, usecase.ErrCharacterNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initiate avatar upload")
	}

	return c.JSON(http.StatusCreated, initAvatarUploadResponse{
		ID:        result.ID.String(),
		UploadURL: result.UploadURL,
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *CharacterHandler) CompleteAvatarUpload(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CompleteAvatarUpload")
	defer span.End()

	characterID := c.Param("id")
	if _, err := uuid.Parse(characterID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	uploadID, err := uuid.Parse(c.Param("uploadID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid upload id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.characterUsecase.CompleteAvatarUpload(ctx, userID, characterID, uploadID); err != nil {
		if errors.Is(err, usecase.ErrPendingUploadNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "pending upload not found")
		}
		if errors.Is(err, usecase.ErrCharacterNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to complete avatar upload")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CharacterHandler) DeleteAvatar(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteAvatar")
	defer span.End()

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.characterUsecase.DeleteAvatar(ctx, userID, id); err != nil {
		if errors.Is(err, usecase.ErrCharacterNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "character not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete avatar")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CharacterHandler) GetCharacterImages(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetCharacterImages")
	defer span.End()

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	images, err := h.characterUsecase.GetCharacterImages(ctx, characterID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get character images")
	}

	responses := make([]characterImageResponse, len(images))
	for i, img := range images {
		var thumbnailURL *string
		if img.ThumbnailR2Path != nil {
			u, presignErr := h.presigner.GeneratePresignedGetURL(ctx, *img.ThumbnailR2Path, usecase.PresignGetTTL)
			if presignErr == nil {
				thumbnailURL = &u
			}
		}
		responses[i] = characterImageResponse{
			ImageID:      img.ID.String(),
			ImageName:    img.Title,
			ThumbnailURL: thumbnailURL,
		}
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *CharacterHandler) presignAvatarURL(ctx context.Context, r2Path *string) (*string, error) {
	if r2Path == nil {
		return nil, nil
	}
	u, err := h.presigner.GeneratePresignedGetURL(ctx, *r2Path, usecase.PresignGetTTL)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func parseFolderIDs(strs *[]string) *[]uuid.UUID {
	if strs == nil {
		return nil
	}
	seen := make(map[uuid.UUID]struct{}, len(*strs))
	ids := make([]uuid.UUID, 0, len(*strs))
	for _, s := range *strs {
		id := uuid.MustParse(s)
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return &ids
}

func toCharacterResponse(character *domain.Character, avatarURL *string) characterResponse {
	folderIDs := make([]string, len(character.Folders))
	for i, f := range character.Folders {
		folderIDs[i] = f.FolderID.String()
	}
	return characterResponse{
		ID:        character.ID.String(),
		Name:      character.Name,
		AvatarURL: avatarURL,
		Biography: character.Biography,
		IsPublic:  character.IsPublic,
		FolderIDs: folderIDs,
		CreatedAt: character.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: character.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
