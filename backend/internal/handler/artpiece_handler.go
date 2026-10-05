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

type ArtpieceUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateArtpieceParams) (*domain.Artpiece, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
	List(ctx context.Context, userID uuid.UUID, filters usecase.ListArtpieceFilters) ([]*domain.Artpiece, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateArtpieceParams) (*domain.Artpiece, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	AttachFile(ctx context.Context, artpieceID uuid.UUID, fileID uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
	DetachFile(ctx context.Context, artpieceID uuid.UUID, fileID uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
	SetCover(ctx context.Context, artpieceID uuid.UUID, fileID uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
	ReplaceFiles(ctx context.Context, artpieceID uuid.UUID, userID uuid.UUID, fileIDs []uuid.UUID) (*domain.Artpiece, error)
}

type ArtpieceHandler struct {
	artpieceUsecase ArtpieceUsecase
	presigner       Presigner
	tel             *observability.Telemetry
}

func NewArtpieceHandler(artpieceUsecase ArtpieceUsecase, presigner Presigner, tel *observability.Telemetry) *ArtpieceHandler {
	return &ArtpieceHandler{artpieceUsecase: artpieceUsecase, presigner: presigner, tel: tel}
}

type createArtpieceRequest struct {
	Title        *string     `json:"title"`
	ArtistID     *uuid.UUID  `json:"artist_id"`
	Notes        *string     `json:"notes"`
	CharacterIDs []uuid.UUID `json:"character_ids"`
	FileIDs      []uuid.UUID `json:"file_ids"`
}

type replaceFilesRequest struct {
	FileIDs []uuid.UUID `json:"file_ids"`
}

type updateArtpieceRequest struct {
	Title        *string     `json:"title"`
	ArtistID     *uuid.UUID  `json:"artist_id"`
	Notes        *string     `json:"notes"`
	CharacterIDs []uuid.UUID `json:"character_ids"`
}

type setCoverRequest struct {
	FileID uuid.UUID `json:"file_id"`
}

type characterRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type fileRef struct {
	ID           string  `json:"id"`
	FileURL      string  `json:"file_url"`
	ThumbnailURL *string `json:"thumbnail_url"`
}

type artpieceResponse struct {
	ID           string         `json:"id"`
	Title        *string        `json:"title"`
	ArtistID     *string        `json:"artist_id"`
	ArtistName   *string        `json:"artist_name"`
	CoverFileID  *string        `json:"cover_file_id"`
	ThumbnailURL *string        `json:"thumbnail_url"`
	Notes        *string        `json:"notes"`
	Characters   []characterRef `json:"characters"`
	Files        []fileRef      `json:"files"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

func (h *ArtpieceHandler) CreateArtpiece(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateArtpiece")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req createArtpieceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	artpiece, err := h.artpieceUsecase.Create(ctx, userID, usecase.CreateArtpieceParams{
		Title:        req.Title,
		ArtistID:     req.ArtistID,
		Notes:        req.Notes,
		CharacterIDs: req.CharacterIDs,
		FileIDs:      req.FileIDs,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrArtistNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, usecase.ErrCharacterNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, usecase.ErrFileNotOwned) || errors.Is(err, usecase.ErrFileAlreadyAttached) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create artpiece")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	return c.JSON(http.StatusCreated, toArtpieceResponse(artpiece, thumbnailURL, nil))
}

func (h *ArtpieceHandler) GetArtpieceByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetArtpieceByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artpiece, err := h.artpieceUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get artpiece")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	fileRefs := make([]fileRef, 0, len(artpiece.Files))
	for _, f := range artpiece.Files {
		fileURL, err := h.presigner.GeneratePresignedGetURL(ctx, f.FileR2Path, usecase.PresignGetTTL)
		if err != nil {
			continue
		}
		var thumbURL *string
		if f.ThumbnailR2Path != nil {
			u, err := h.presigner.GeneratePresignedGetURL(ctx, *f.ThumbnailR2Path, usecase.PresignGetTTL)
			if err == nil {
				thumbURL = &u
			}
		}
		fileRefs = append(fileRefs, fileRef{
			ID:           f.ID.String(),
			FileURL:      fileURL,
			ThumbnailURL: thumbURL,
		})
	}
	return c.JSON(http.StatusOK, toArtpieceResponse(artpiece, thumbnailURL, fileRefs))
}

func (h *ArtpieceHandler) ListArtpieces(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListArtpieces")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var filters usecase.ListArtpieceFilters
	if err := c.Bind(&filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params")
	}
	if err := c.Validate(&filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params: character_ids and artist_ids must be valid UUIDs")
	}

	artpieces, err := h.artpieceUsecase.List(ctx, userID, filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list artpieces")
	}

	responses := make([]artpieceResponse, len(artpieces))
	for i, a := range artpieces {
		thumbnailURL, _ := h.presignCoverThumbnail(ctx, a)
		responses[i] = toArtpieceResponse(a, thumbnailURL, nil)
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *ArtpieceHandler) UpdateArtpiece(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateArtpiece")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}

	var req updateArtpieceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artpiece, err := h.artpieceUsecase.Update(ctx, id, userID, usecase.UpdateArtpieceParams{
		Title:        req.Title,
		ArtistID:     req.ArtistID,
		Notes:        req.Notes,
		CharacterIDs: req.CharacterIDs,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		if errors.Is(err, usecase.ErrArtistNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, usecase.ErrCharacterNotOwned) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update artpiece")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	return c.JSON(http.StatusOK, toArtpieceResponse(artpiece, thumbnailURL, nil))
}

func (h *ArtpieceHandler) DeleteArtpiece(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteArtpiece")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.artpieceUsecase.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete artpiece")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArtpieceHandler) AttachFile(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.AttachFile")
	defer span.End()

	artpieceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artpiece, err := h.artpieceUsecase.AttachFile(ctx, artpieceID, fileID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		if errors.Is(err, usecase.ErrFileNotOwned) || errors.Is(err, usecase.ErrFileAlreadyAttached) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to attach file")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	return c.JSON(http.StatusOK, toArtpieceResponse(artpiece, thumbnailURL, nil))
}

func (h *ArtpieceHandler) ReplaceFiles(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ReplaceFiles")
	defer span.End()

	artpieceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req replaceFilesRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	artpiece, err := h.artpieceUsecase.ReplaceFiles(ctx, artpieceID, userID, req.FileIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		if errors.Is(err, usecase.ErrFileNotOwned) || errors.Is(err, usecase.ErrFileAlreadyAttached) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to replace files")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	return c.JSON(http.StatusOK, toArtpieceResponse(artpiece, thumbnailURL, nil))
}

func (h *ArtpieceHandler) DetachFile(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DetachFile")
	defer span.End()

	artpieceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artpiece, err := h.artpieceUsecase.DetachFile(ctx, artpieceID, fileID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		if errors.Is(err, usecase.ErrFileNotOwned) || errors.Is(err, usecase.ErrFileNotInArtpiece) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to detach file")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	return c.JSON(http.StatusOK, toArtpieceResponse(artpiece, thumbnailURL, nil))
}

func (h *ArtpieceHandler) SetCover(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.SetCover")
	defer span.End()

	artpieceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid artpiece id")
	}

	var req setCoverRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	artpiece, err := h.artpieceUsecase.SetCover(ctx, artpieceID, req.FileID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "artpiece not found")
		}
		if errors.Is(err, usecase.ErrFileNotOwned) || errors.Is(err, usecase.ErrFileNotInArtpiece) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to set cover")
	}

	thumbnailURL, _ := h.presignCoverThumbnail(ctx, artpiece)
	return c.JSON(http.StatusOK, toArtpieceResponse(artpiece, thumbnailURL, nil))
}

func (h *ArtpieceHandler) presignCoverThumbnail(ctx context.Context, a *domain.Artpiece) (*string, error) {
	if a.CoverFile == nil || a.CoverFile.ThumbnailR2Path == nil {
		return nil, nil
	}
	u, err := h.presigner.GeneratePresignedGetURL(ctx, *a.CoverFile.ThumbnailR2Path, usecase.PresignGetTTL)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func toArtpieceResponse(a *domain.Artpiece, thumbnailURL *string, files []fileRef) artpieceResponse {
	chars := make([]characterRef, len(a.Characters))
	for i, c := range a.Characters {
		chars[i] = characterRef{ID: c.ID.String(), Name: c.Name}
	}

	var artistID, artistName *string
	if a.ArtistID != nil {
		s := a.ArtistID.String()
		artistID = &s
	}
	if a.Artist != nil {
		artistName = &a.Artist.Name
	}

	var coverFileID *string
	if a.CoverFileID != nil {
		s := a.CoverFileID.String()
		coverFileID = &s
	}

	return artpieceResponse{
		ID:           a.ID.String(),
		Title:        a.Title,
		ArtistID:     artistID,
		ArtistName:   artistName,
		CoverFileID:  coverFileID,
		ThumbnailURL: thumbnailURL,
		Notes:        a.Notes,
		Characters:   chars,
		Files:        files,
		CreatedAt:    a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
