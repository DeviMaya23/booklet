package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/handler"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyTagUsecase struct {
	createResult *domain.Tag
	createErr    error
	lastCreate   usecase.CreateTagParams

	listResult  []*domain.Tag
	listErr     error
	lastFilters usecase.ListTagFilters

	updateResult *domain.Tag
	updateErr    error
	lastUpdate   usecase.UpdateTagParams

	deleteErr error

	usageResult []*domain.TagWithUsage
	usageErr    error

	mergeErr    error
	lastMergeID uuid.UUID
	lastMerge   usecase.MergeTagParams
	listCalled  bool
}

func (s *spyTagUsecase) ListWithUsage(_ context.Context, _ uuid.UUID, f usecase.ListTagFilters) ([]*domain.TagWithUsage, error) {
	s.lastFilters = f
	return s.usageResult, s.usageErr
}

func (s *spyTagUsecase) Merge(_ context.Context, id uuid.UUID, _ uuid.UUID, p usecase.MergeTagParams) error {
	s.lastMergeID = id
	s.lastMerge = p
	return s.mergeErr
}

func (s *spyTagUsecase) Create(_ context.Context, _ uuid.UUID, p usecase.CreateTagParams) (*domain.Tag, error) {
	s.lastCreate = p
	return s.createResult, s.createErr
}

func (s *spyTagUsecase) List(_ context.Context, _ uuid.UUID, f usecase.ListTagFilters) ([]*domain.Tag, error) {
	s.lastFilters = f
	s.listCalled = true
	return s.listResult, s.listErr
}

func (s *spyTagUsecase) Update(_ context.Context, _ uuid.UUID, _ uuid.UUID, p usecase.UpdateTagParams) (*domain.Tag, error) {
	s.lastUpdate = p
	return s.updateResult, s.updateErr
}

func (s *spyTagUsecase) Delete(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return s.deleteErr
}

func makeTag(name string) *domain.Tag {
	return &domain.Tag{ID: uuid.New(), UserID: testUserID, Name: name}
}

func newTagHandler(spy *spyTagUsecase) *handler.TagHandler {
	return handler.NewTagHandler(spy, observability.NewTelemetry(nil, nil, nil))
}

func doTagRequest(t *testing.T, h *handler.TagHandler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	e := setupEcho(testUserID)
	e.POST("/tags", h.CreateTag)
	e.GET("/tags", h.ListTags)
	e.PATCH("/tags/:id", h.UpdateTag)
	e.DELETE("/tags/:id", h.DeleteTag)
	e.POST("/tags/:id/merge", h.MergeTag)

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// --- CreateTag ---

func TestCreateTag_HappyPath(t *testing.T) {
	tag := makeTag("wip")
	spy := &spyTagUsecase{createResult: tag}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags", `{"name":"WIP"}`)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, map[string]any{"id": tag.ID.String(), "name": "wip"}, got)
	require.Equal(t, "WIP", spy.lastCreate.Name)
}

func TestCreateTag_MissingOrEmptyName(t *testing.T) {
	for _, body := range []string{`{}`, `{"name":""}`} {
		rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodPost, "/tags", body)
		require.Equal(t, http.StatusUnprocessableEntity, rec.Code, body)
	}
}

func TestCreateTag_BlankNameFromUsecase(t *testing.T) {
	spy := &spyTagUsecase{createErr: usecase.ErrTagNameInvalid}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags", `{"name":"   "}`)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateTag_Conflict(t *testing.T) {
	spy := &spyTagUsecase{createErr: usecase.ErrTagNameConflict}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags", `{"name":"wip"}`)

	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateTag_MalformedJSON(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodPost, "/tags", `{not json`)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTag_InternalError(t *testing.T) {
	spy := &spyTagUsecase{createErr: errors.New("boom")}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags", `{"name":"wip"}`)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- ListTags ---

func TestListTags_HappyPath(t *testing.T) {
	a, b := makeTag("icon"), makeTag("wip")
	spy := &spyTagUsecase{listResult: []*domain.Tag{a, b}}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags", "")

	require.Equal(t, http.StatusOK, rec.Code)
	var got []map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, []map[string]any{
		{"id": a.ID.String(), "name": "icon"},
		{"id": b.ID.String(), "name": "wip"},
	}, got)
	require.Nil(t, spy.lastFilters.Sort)
}

func TestListTags_EmptyIsArray(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodGet, "/tags", "")

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `[]`, rec.Body.String())
}

func TestListTags_SortReachesUsecase(t *testing.T) {
	spy := &spyTagUsecase{}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags?sort=recent", "")

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, spy.lastFilters.Sort)
	require.Equal(t, "recent", *spy.lastFilters.Sort)
}

func TestListTags_InvalidSort(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodGet, "/tags?sort=bogus", "")

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	require.Contains(t, rec.Body.String(), "sort")
}

func TestListTags_InternalError(t *testing.T) {
	spy := &spyTagUsecase{listErr: errors.New("boom")}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags", "")

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- UpdateTag ---

func TestUpdateTag_HappyPath(t *testing.T) {
	tag := makeTag("finished")
	spy := &spyTagUsecase{updateResult: tag}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodPatch, "/tags/"+tag.ID.String(), `{"name":"Finished"}`)

	require.Equal(t, http.StatusOK, rec.Code)
	var got map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, map[string]any{"id": tag.ID.String(), "name": "finished"}, got)
	require.Equal(t, "Finished", spy.lastUpdate.Name)
}

func TestUpdateTag_ErrorMapping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"not found", gorm.ErrRecordNotFound, http.StatusNotFound},
		{"conflict", usecase.ErrTagNameConflict, http.StatusConflict},
		{"blank name", usecase.ErrTagNameInvalid, http.StatusUnprocessableEntity},
		{"internal", errors.New("boom"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyTagUsecase{updateErr: tt.err}

			rec := doTagRequest(t, newTagHandler(spy), http.MethodPatch, "/tags/"+uuid.NewString(), `{"name":"x"}`)

			require.Equal(t, tt.want, rec.Code)
		})
	}
}

func TestUpdateTag_EmptyName(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodPatch, "/tags/"+uuid.NewString(), `{"name":""}`)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateTag_MalformedJSON(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodPatch, "/tags/"+uuid.NewString(), `{not json`)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateTag_MalformedID(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodPatch, "/tags/not-a-uuid", `{"name":"x"}`)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- DeleteTag ---

func TestDeleteTag_HappyPath(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodDelete, "/tags/"+uuid.NewString(), "")

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
}

func TestDeleteTag_NotFound(t *testing.T) {
	spy := &spyTagUsecase{deleteErr: gorm.ErrRecordNotFound}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodDelete, "/tags/"+uuid.NewString(), "")

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteTag_MalformedID(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodDelete, "/tags/not-a-uuid", "")

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteTag_InternalError(t *testing.T) {
	spy := &spyTagUsecase{deleteErr: errors.New("boom")}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodDelete, "/tags/"+uuid.NewString(), "")

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- ListTags with usage ---

func TestListTags_WithUsageIncludesCount(t *testing.T) {
	a, b := makeTag("icon"), makeTag("wip")
	spy := &spyTagUsecase{usageResult: []*domain.TagWithUsage{
		{Tag: *a, ArtpieceCount: 3},
		{Tag: *b, ArtpieceCount: 0},
	}}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags?with_usage=true", "")

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `[
		{"id":"`+a.ID.String()+`","name":"icon","artpiece_count":3},
		{"id":"`+b.ID.String()+`","name":"wip","artpiece_count":0}
	]`, rec.Body.String())
	require.False(t, spy.listCalled)
}

func TestListTags_WithUsageAndSort(t *testing.T) {
	spy := &spyTagUsecase{}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags?with_usage=true&sort=recent", "")

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, spy.lastFilters.WithUsage)
	require.NotNil(t, spy.lastFilters.Sort)
	require.Equal(t, "recent", *spy.lastFilters.Sort)
	require.JSONEq(t, `[]`, rec.Body.String())
}

func TestListTags_WithUsageFalseIsDefault(t *testing.T) {
	a := makeTag("icon")
	spy := &spyTagUsecase{listResult: []*domain.Tag{a}}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags?with_usage=false", "")

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `[{"id":"`+a.ID.String()+`","name":"icon"}]`, rec.Body.String())
	require.True(t, spy.listCalled)
}

func TestListTags_WithUsageInvalidValue(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodGet, "/tags?with_usage=maybe", "")

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListTags_WithUsageInternalError(t *testing.T) {
	spy := &spyTagUsecase{usageErr: errors.New("boom")}

	rec := doTagRequest(t, newTagHandler(spy), http.MethodGet, "/tags?with_usage=true", "")

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- MergeTag ---

func TestMergeTag_HappyPath(t *testing.T) {
	spy := &spyTagUsecase{}
	source, target := uuid.New(), uuid.New()

	rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags/"+source.String()+"/merge", `{"target_id":"`+target.String()+`"}`)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, source, spy.lastMergeID)
	require.Equal(t, target, spy.lastMerge.TargetID)
}

func TestMergeTag_ErrorMapping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
		body string
	}{
		{"not found", gorm.ErrRecordNotFound, http.StatusNotFound, "tag not found"},
		{"same tag", usecase.ErrTagMergeSameTag, http.StatusUnprocessableEntity, usecase.ErrTagMergeSameTag.Error()},
		{"internal", errors.New("boom"), http.StatusInternalServerError, "failed to merge tags"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyTagUsecase{mergeErr: tt.err}

			rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags/"+uuid.NewString()+"/merge", `{"target_id":"`+uuid.NewString()+`"}`)

			require.Equal(t, tt.want, rec.Code)
			require.Contains(t, rec.Body.String(), tt.body)
		})
	}
}

func TestMergeTag_MalformedID(t *testing.T) {
	rec := doTagRequest(t, newTagHandler(&spyTagUsecase{}), http.MethodPost, "/tags/not-a-uuid/merge", `{"target_id":"`+uuid.NewString()+`"}`)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid tag id")
}

func TestMergeTag_BadBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		want int
	}{
		{"invalid json", `{not json`, http.StatusBadRequest},
		{"non-uuid target", `{"target_id":"nope"}`, http.StatusBadRequest},
		{"missing target", `{}`, http.StatusUnprocessableEntity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyTagUsecase{}

			rec := doTagRequest(t, newTagHandler(spy), http.MethodPost, "/tags/"+uuid.NewString()+"/merge", tt.body)

			require.Equal(t, tt.want, rec.Code)
			require.Equal(t, uuid.Nil, spy.lastMergeID)
		})
	}
}
