package repository

import (
	"context"
	"testing"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedTag(t *testing.T, tx *gorm.DB, userID uuid.UUID, name string) *domain.Tag {
	t.Helper()
	tag := &domain.Tag{ID: uuid.New(), UserID: userID, Name: name}
	require.NoError(t, tx.Create(tag).Error)
	return tag
}

func seedTagWithLastUsed(t *testing.T, tx *gorm.DB, userID uuid.UUID, name string, lastUsed time.Time) *domain.Tag {
	t.Helper()
	tag := &domain.Tag{ID: uuid.New(), UserID: userID, Name: name, LastUsedAt: &lastUsed}
	require.NoError(t, tx.Create(tag).Error)
	return tag
}

func tagNames(tags []*domain.Tag) []string {
	names := make([]string, len(tags))
	for i, tg := range tags {
		names[i] = tg.Name
	}
	return names
}

func strPtr(s string) *string { return &s }

func TestTagRepository_CreateAndList(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")

	created, err := repo.Create(context.Background(), &domain.Tag{ID: uuid.New(), UserID: user.ID, Name: "wip"})
	require.NoError(t, err)
	assert.Equal(t, "wip", created.Name)

	got, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, created.ID, got[0].ID)
	assert.Equal(t, "wip", got[0].Name)
}

func TestTagRepository_List_DefaultOrderByName(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	now := time.Now()
	seedTagWithLastUsed(t, tx, user.ID, "wip", now)
	seedTagWithLastUsed(t, tx, user.ID, "sketch", now.Add(-time.Hour))
	seedTag(t, tx, user.ID, "icon")

	got, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{})

	require.NoError(t, err)
	assert.Equal(t, []string{"icon", "sketch", "wip"}, tagNames(got))
}

func TestTagRepository_List_SortRecent(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	now := time.Now().Truncate(time.Second)
	seedTagWithLastUsed(t, tx, user.ID, "old", now.Add(-time.Hour))
	seedTagWithLastUsed(t, tx, user.ID, "newest", now)
	seedTagWithLastUsed(t, tx, user.ID, "tie-b", now.Add(-time.Minute))
	seedTagWithLastUsed(t, tx, user.ID, "tie-a", now.Add(-time.Minute))
	seedTag(t, tx, user.ID, "never-b")
	seedTag(t, tx, user.ID, "never-a")

	got, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{Sort: strPtr("recent")})

	require.NoError(t, err)
	assert.Equal(t, []string{"newest", "tie-a", "tie-b", "old", "never-a", "never-b"}, tagNames(got))
}

func TestTagRepository_List_OnlyOwnTags(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedTag(t, tx, user1.ID, "mine")
	seedTag(t, tx, user2.ID, "theirs")

	got, err := repo.List(context.Background(), user1.ID, usecase.ListTagFilters{})

	require.NoError(t, err)
	assert.Equal(t, []string{"mine"}, tagNames(got))
}

func TestTagRepository_Create_DuplicateNameSameUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	seedTag(t, tx, user.ID, "wip")

	err := tx.Transaction(func(sp *gorm.DB) error {
		_, err := NewTagRepository(sp).Create(context.Background(), &domain.Tag{ID: uuid.New(), UserID: user.ID, Name: "wip"})
		return err
	})

	require.ErrorIs(t, err, usecase.ErrTagNameConflict)
	got, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestTagRepository_Create_SameNameDifferentUsers(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	seedTag(t, tx, user1.ID, "wip")

	got, err := repo.Create(context.Background(), &domain.Tag{ID: uuid.New(), UserID: user2.ID, Name: "wip"})

	require.NoError(t, err)
	assert.Equal(t, "wip", got.Name)
}

func TestTagRepository_Create_RejectsNonLowercaseName(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	user := seedUser(t, tx, "user_1")

	err := tx.Transaction(func(sp *gorm.DB) error {
		_, err := NewTagRepository(sp).Create(context.Background(), &domain.Tag{ID: uuid.New(), UserID: user.ID, Name: "WIP"})
		return err
	})

	require.Error(t, err)
	assert.NotErrorIs(t, err, usecase.ErrTagNameConflict)
	assert.ErrorContains(t, err, "tags_name_check")
}

func TestTagRepository_Create_NonASCIILowercaseRoundTrips(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")

	_, err := repo.Create(context.Background(), &domain.Tag{ID: uuid.New(), UserID: user.ID, Name: "éclair 日本語"})
	require.NoError(t, err)

	got, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	assert.Equal(t, []string{"éclair 日本語"}, tagNames(got))
}

func TestTagRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	tag := seedTag(t, tx, user.ID, "wip")

	got, err := repo.Update(context.Background(), tag.ID, user.ID, usecase.UpdateTagParams{Name: "finished"})

	require.NoError(t, err)
	assert.Equal(t, tag.ID, got.ID)
	assert.Equal(t, "finished", got.Name)
}

func TestTagRepository_Update_DoesNotTouchLastUsedAt(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	lastUsed := time.Now().Add(-48 * time.Hour).Truncate(time.Second)
	tag := seedTagWithLastUsed(t, tx, user.ID, "wip", lastUsed)

	got, err := repo.Update(context.Background(), tag.ID, user.ID, usecase.UpdateTagParams{Name: "finished"})

	require.NoError(t, err)
	require.NotNil(t, got.LastUsedAt)
	assert.True(t, lastUsed.Equal(*got.LastUsedAt))
}

func TestTagRepository_Update_ConflictLeavesBothUnchanged(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	seedTag(t, tx, user.ID, "wip")
	other := seedTag(t, tx, user.ID, "sketch")

	err := tx.Transaction(func(sp *gorm.DB) error {
		_, err := NewTagRepository(sp).Update(context.Background(), other.ID, user.ID, usecase.UpdateTagParams{Name: "wip"})
		return err
	})

	require.ErrorIs(t, err, usecase.ErrTagNameConflict)
	got, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	assert.Equal(t, []string{"sketch", "wip"}, tagNames(got))
}

func TestTagRepository_Update_SelfRenameInDifferentCase(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	tag := seedTag(t, tx, user.ID, "wip")

	// The usecase lowercases before calling the repository, so "WIP" arrives as "wip".
	got, err := repo.Update(context.Background(), tag.ID, user.ID, usecase.UpdateTagParams{Name: "wip"})

	require.NoError(t, err)
	assert.Equal(t, "wip", got.Name)
}

func TestTagRepository_Update_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	owner := seedUser(t, tx, "user_1")
	other := seedUser(t, tx, "user_2")
	tag := seedTag(t, tx, owner.ID, "wip")

	got, err := repo.Update(context.Background(), tag.ID, other.ID, usecase.UpdateTagParams{Name: "hijacked"})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	list, err := repo.List(context.Background(), owner.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	assert.Equal(t, []string{"wip"}, tagNames(list))
}

func TestTagRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	tag := seedTag(t, tx, user.ID, "wip")

	require.NoError(t, repo.Delete(context.Background(), tag.ID, user.ID))

	list, err := repo.List(context.Background(), user.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestTagRepository_Delete_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	owner := seedUser(t, tx, "user_1")
	other := seedUser(t, tx, "user_2")
	tag := seedTag(t, tx, owner.ID, "wip")

	err := repo.Delete(context.Background(), tag.ID, other.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	list, err := repo.List(context.Background(), owner.ID, usecase.ListTagFilters{})
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestTagRepository_Delete_CascadesJoinRowsAndKeepsArtpiece(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	tag := seedTag(t, tx, user.ID, "wip")
	artpiece := seedArtpiece(t, tx, user.ID)
	require.NoError(t, tx.Exec("INSERT INTO artpiece_tags (artpiece_id, tag_id) VALUES (?, ?)", artpiece.ID, tag.ID).Error)

	require.NoError(t, repo.Delete(context.Background(), tag.ID, user.ID))

	var joinRows int64
	require.NoError(t, tx.Table("artpiece_tags").Where("artpiece_id = ?", artpiece.ID).Count(&joinRows).Error)
	assert.Zero(t, joinRows)
	var artpieceRows int64
	require.NoError(t, tx.Model(&domain.Artpiece{}).Where("id = ?", artpiece.ID).Count(&artpieceRows).Error)
	assert.Equal(t, int64(1), artpieceRows)
}

func TestTagRepository_GetByIDsAndUserID_ExcludesOtherUsersTags(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	mine := seedTag(t, tx, user1.ID, "mine")
	theirs := seedTag(t, tx, user2.ID, "theirs")

	got, err := repo.GetByIDsAndUserID(context.Background(), []uuid.UUID{mine.ID, theirs.ID}, user1.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, mine.ID, got[0].ID)
}

func TestTagRepository_UpdateLastUsedAt_BumpsOnlyGivenUsersTags(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	old := time.Now().Add(-48 * time.Hour).Truncate(time.Second)
	mine := seedTagWithLastUsed(t, tx, user1.ID, "mine", old)
	untouched := seedTagWithLastUsed(t, tx, user1.ID, "untouched", old)
	theirs := seedTagWithLastUsed(t, tx, user2.ID, "theirs", old)

	require.NoError(t, repo.UpdateLastUsedAt(context.Background(), []uuid.UUID{mine.ID, theirs.ID}, user1.ID))

	var after []domain.Tag
	require.NoError(t, tx.Find(&after).Error)
	byID := map[uuid.UUID]domain.Tag{}
	for _, tg := range after {
		byID[tg.ID] = tg
	}
	assert.True(t, byID[mine.ID].LastUsedAt.After(old))
	assert.True(t, byID[untouched.ID].LastUsedAt.Equal(old))
	assert.True(t, byID[theirs.ID].LastUsedAt.Equal(old))
}

func TestTagRepository_UpdateLastUsedAt_EmptyIDsIsNoOp(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	old := time.Now().Add(-48 * time.Hour).Truncate(time.Second)
	tag := seedTagWithLastUsed(t, tx, user.ID, "wip", old)

	require.NoError(t, repo.UpdateLastUsedAt(context.Background(), nil, user.ID))

	var after domain.Tag
	require.NoError(t, tx.First(&after, "id = ?", tag.ID).Error)
	assert.True(t, after.LastUsedAt.Equal(old))
}
