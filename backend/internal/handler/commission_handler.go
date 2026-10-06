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

type CommissionUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, params usecase.CreateCommissionParams) (*domain.Commission, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Commission, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Commission, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateCommissionParams) (*domain.Commission, error)
	Patch(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.PatchCommissionParams) (*domain.Commission, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	AttachArtpieces(ctx context.Context, commissionID uuid.UUID, userID uuid.UUID, artpieceIDs []uuid.UUID) (*domain.Commission, error)
	DetachArtpieces(ctx context.Context, commissionID uuid.UUID, userID uuid.UUID, artpieceIDs []uuid.UUID) (*domain.Commission, error)
	ReplaceArtpieces(ctx context.Context, commissionID uuid.UUID, userID uuid.UUID, artpieceIDs []uuid.UUID) (*domain.Commission, error)
}

type CommissionHandler struct {
	commissionUsecase CommissionUsecase
	presigner         Presigner
	tel               *observability.Telemetry
}

func NewCommissionHandler(commissionUsecase CommissionUsecase, presigner Presigner, tel *observability.Telemetry) *CommissionHandler {
	return &CommissionHandler{
		commissionUsecase: commissionUsecase,
		presigner:         presigner,
		tel:               tel,
	}
}

type createCommissionRequest struct {
	Title        *string     `json:"title"`
	ArtistID     *uuid.UUID  `json:"artist_id"`
	Status       string      `json:"status"`
	Price        *float64    `json:"price"`
	Paid         bool        `json:"paid"`
	PaidDate     *string     `json:"paid_date"`
	FinishDate   *string     `json:"finish_date"`
	Notes        *string     `json:"notes"`
	CharacterIDs []uuid.UUID `json:"character_ids"`
	ArtpieceIDs  []uuid.UUID `json:"artpiece_ids"`
}

type updateCommissionRequest struct {
	Title        *string     `json:"title"`
	ArtistID     *uuid.UUID  `json:"artist_id"`
	Status       string      `json:"status"`
	Price        *float64    `json:"price"`
	Paid         bool        `json:"paid"`
	PaidDate     *string     `json:"paid_date"`
	FinishDate   *string     `json:"finish_date"`
	Notes        *string     `json:"notes"`
	CharacterIDs []uuid.UUID `json:"character_ids"`
}

type artpieceIDsRequest struct {
	ArtpieceIDs []uuid.UUID `json:"artpiece_ids"`
}

type artpieceSummary struct {
	ID           string  `json:"id"`
	ThumbnailURL *string `json:"thumbnail_url"`
}

type commissionCharacterRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type commissionResponse struct {
	ID              string                   `json:"id"`
	Title           *string                  `json:"title"`
	ArtistID        *string                  `json:"artist_id"`
	ArtistName      *string                  `json:"artist_name"`
	ArtistLink      *string                  `json:"artist_link"`
	Status          string                   `json:"status"`
	Price           *float64                 `json:"price"`
	Paid            bool                     `json:"paid"`
	PaidDate        *string                  `json:"paid_date"`
	FinishDate      *string                  `json:"finish_date"`
	LastContactedAt *string                  `json:"last_contacted_at"`
	Notes           *string                  `json:"notes"`
	Characters      []commissionCharacterRef `json:"characters"`
	Artpieces       []artpieceSummary        `json:"artpieces,omitempty"`
	CreatedAt    string                   `json:"created_at"`
	UpdatedAt    string                   `json:"updated_at"`
}

func (h *CommissionHandler) CreateCommission(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.CreateCommission")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req createCommissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	commission, err := h.commissionUsecase.Create(ctx, userID, usecase.CreateCommissionParams{
		Title:        req.Title,
		ArtistID:     req.ArtistID,
		Status:       req.Status,
		Price:        req.Price,
		Paid:         req.Paid,
		PaidDate:     req.PaidDate,
		FinishDate:   req.FinishDate,
		Notes:        req.Notes,
		CharacterIDs: req.CharacterIDs,
		ArtpieceIDs:  req.ArtpieceIDs,
	})
	if err != nil {
		return h.mapCommissionError(err)
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *CommissionHandler) GetCommissionByID(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.GetCommissionByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	commission, err := h.commissionUsecase.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "commission not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get commission")
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusOK, resp)
}

type patchCommissionRequest struct {
	Status          *string `json:"status"          validate:"omitempty,oneof=waitlist wip done"`
	Paid            *bool   `json:"paid"`
	PaidDate        *string `json:"paid_date"        validate:"omitempty,datetime=2006-01-02"`
	LastContactedAt *string `json:"last_contacted_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

func (h *CommissionHandler) PatchCommission(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.PatchCommission")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req patchCommissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	commission, err := h.commissionUsecase.Patch(ctx, id, userID, usecase.PatchCommissionParams{
		Status:          req.Status,
		Paid:            req.Paid,
		PaidDate:        req.PaidDate,
		LastContactedAt: req.LastContactedAt,
	})
	if err != nil {
		return h.mapCommissionError(err)
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommissionHandler) ListCommissions(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ListCommissions")
	defer span.End()

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	commissions, err := h.commissionUsecase.List(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list commissions")
	}

	responses := make([]commissionResponse, 0, len(commissions))
	for _, comm := range commissions {
		resp, err := h.toCommissionResponse(ctx, comm)
		if err != nil {
			continue
		}
		responses = append(responses, resp)
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *CommissionHandler) UpdateCommission(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.UpdateCommission")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req updateCommissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrResponse(err))
	}

	commission, err := h.commissionUsecase.Update(ctx, id, userID, usecase.UpdateCommissionParams{
		Title:        req.Title,
		ArtistID:     req.ArtistID,
		Status:       req.Status,
		Price:        req.Price,
		Paid:         req.Paid,
		PaidDate:     req.PaidDate,
		FinishDate:   req.FinishDate,
		Notes:        req.Notes,
		CharacterIDs: req.CharacterIDs,
	})
	if err != nil {
		return h.mapCommissionError(err)
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommissionHandler) DeleteCommission(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DeleteCommission")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if err := h.commissionUsecase.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "commission not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete commission")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CommissionHandler) AttachArtpieces(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.AttachArtpieces")
	defer span.End()

	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req artpieceIDsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	commission, err := h.commissionUsecase.AttachArtpieces(ctx, commissionID, userID, req.ArtpieceIDs)
	if err != nil {
		return h.mapCommissionError(err)
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommissionHandler) DetachArtpieces(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.DetachArtpieces")
	defer span.End()

	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req artpieceIDsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	commission, err := h.commissionUsecase.DetachArtpieces(ctx, commissionID, userID, req.ArtpieceIDs)
	if err != nil {
		return h.mapCommissionError(err)
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommissionHandler) ReplaceArtpieces(c echo.Context) error {
	ctx, span := h.tel.Tracer.Start(c.Request().Context(), "handler.ReplaceArtpieces")
	defer span.End()

	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid commission id")
	}

	userID, ok := middleware.AuthenticatedUserIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req artpieceIDsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	commission, err := h.commissionUsecase.ReplaceArtpieces(ctx, commissionID, userID, req.ArtpieceIDs)
	if err != nil {
		return h.mapCommissionError(err)
	}

	resp, err := h.toCommissionResponse(ctx, commission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to build response")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommissionHandler) toCommissionResponse(ctx context.Context, c *domain.Commission) (commissionResponse, error) {
	chars := make([]commissionCharacterRef, len(c.Characters))
	for i, ch := range c.Characters {
		chars[i] = commissionCharacterRef{ID: ch.ID.String(), Name: ch.Name}
	}

	artpieces := make([]artpieceSummary, 0, len(c.Artpieces))
	for _, a := range c.Artpieces {
		summary := artpieceSummary{ID: a.ID.String()}
		if a.CoverFile != nil && a.CoverFile.ThumbnailR2Path != nil {
			u, err := h.presigner.GeneratePresignedGetURL(ctx, *a.CoverFile.ThumbnailR2Path, usecase.PresignGetTTL)
			if err == nil {
				summary.ThumbnailURL = &u
			}
		}
		artpieces = append(artpieces, summary)
	}

	var artistID, artistName, artistLink *string
	if c.ArtistID != nil {
		s := c.ArtistID.String()
		artistID = &s
	}
	if c.Artist != nil {
		artistName = &c.Artist.Name
		artistLink = c.Artist.ArtistLink
	}

	var lastContactedAt *string
	if c.LastContactedAt != nil {
		s := c.LastContactedAt.Format("2006-01-02T15:04:05Z07:00")
		lastContactedAt = &s
	}

	return commissionResponse{
		ID:              c.ID.String(),
		Title:           c.Title,
		ArtistID:        artistID,
		ArtistName:      artistName,
		ArtistLink:      artistLink,
		Status:          c.Status,
		Price:           c.Price,
		Paid:            c.Paid,
		PaidDate:        formatDate(c.PaidDate),
		FinishDate:      formatDate(c.FinishDate),
		LastContactedAt: lastContactedAt,
		Notes:           c.Notes,
		Characters:      chars,
		Artpieces:       artpieces,
		CreatedAt:       c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *CommissionHandler) mapCommissionError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "commission not found")
	}
	if errors.Is(err, usecase.ErrArtistNotOwned) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if errors.Is(err, usecase.ErrCharacterNotOwned) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if errors.Is(err, usecase.ErrArtpieceNotOwned) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if errors.Is(err, usecase.ErrArtpieceAlreadyAttached) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if errors.Is(err, usecase.ErrInvalidCommissionStatus) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
}

func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}
