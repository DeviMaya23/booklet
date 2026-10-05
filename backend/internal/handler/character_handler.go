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
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Character, error)
	List(ctx context.Context, userID uuid.UUID, filters usecase.ListCharacterFilters) ([]*domain.Character, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	InitAvatarUpload(ctx context.Context, userID uuid.UUID, characterID uuid.UUID, mimeType string) (*usecase.AvatarUploadResult, error)
	CompleteAvatarUpload(ctx context.Context, userID uuid.UUID, characterID uuid.UUID, uploadID uuid.UUID) error
	DeleteAvatar(ctx context.Context, userID uuid.UUID, characterID uuid.UUID) error
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
	Name      string       `json:"name" validate:"required"`
	Notes     *string      `json:"notes"`
	IsPublic  bool         `json:"is_public"`
	FolderIDs *[]uuid.UUID `json:"folder_ids"`
}

type updateCharacterRequest struct {
	Name      string      `json:"name" validate:"required,min=1"`
	Notes     *string     `json:"notes"`
	IsPublic  bool        `json:"is_public"`
	FolderIDs []uuid.UUID `json:"folder_ids"`
}

type folderResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type characterResponse struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	AvatarURL *string          `json:"avatar_url"`
	Notes     *string          `json:"notes"`
	IsPublic  bool             `json:"is_public"`
	Folders   []folderResponse `json:"folders"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
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

	idpSubject, ok := middleware.AuthenticatedIDPSubjectFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	character, err := h.characterUsecase.Create(ctx, userID, usecase.CreateCharacterParams{
		Name:       req.Name,
		Notes:      req.Notes,
		IsPublic:   req.IsPublic,
		FolderIDs:  req.FolderIDs,
		IDPSubject: idpSubject,
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
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

	var filters usecase.ListCharacterFilters
	if err := c.Bind(&filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params")
	}

	characters, err := h.characterUsecase.List(ctx, userID, filters)
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
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

	idpSubject, ok := middleware.AuthenticatedIDPSubjectFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	character, err := h.characterUsecase.Update(ctx, id, userID, usecase.UpdateCharacterParams{
		Name:       req.Name,
		Notes:      req.Notes,
		IsPublic:   req.IsPublic,
		FolderIDs:  req.FolderIDs,
		IDPSubject: idpSubject,
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid character id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.characterUsecase.Delete(ctx, id, userID)
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
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

	characterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
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


func toCharacterResponse(character *domain.Character, avatarURL *string) characterResponse {
	folders := make([]folderResponse, len(character.Folders))
	for i, f := range character.Folders {
		folders[i] = folderResponse{ID: f.FolderID.String(), Name: f.FolderName}
	}
	return characterResponse{
		ID:        character.ID.String(),
		Name:      character.Name,
		AvatarURL: avatarURL,
		Notes:     character.Notes,
		IsPublic:  character.IsPublic,
		Folders:   folders,
		CreatedAt: character.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: character.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
