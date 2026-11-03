package repository

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedArtpiece(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.Artpiece {
	t.Helper()
	a := &domain.Artpiece{
		ID:     uuid.New(),
		UserID: userID,
	}
	require.NoError(t, tx.Create(a).Error)
	return a
}

func TestArtpieceRepository_GetByID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpiece(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), artpiece.ID, user.ID)

	require.NoError(t, err)
	assert.Equal(t, artpiece.ID, got.ID)
	assert.Empty(t, got.Characters)
	assert.Nil(t, got.CoverFile)
}

func TestArtpieceRepository_GetByID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	artpiece := seedArtpiece(t, tx, user1.ID)

	_, err := repo.GetByID(context.Background(), artpiece.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestArtpieceRepository_List_FilterByArtist(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	artist := seedArtist(t, tx, user.ID, "Picasso")

	a1 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, ArtistID: &artist.ID}
	a2 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID}
	require.NoError(t, tx.Create(a1).Error)
	require.NoError(t, tx.Create(a2).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		ArtistIDs: []uuid.UUID{artist.ID},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, a1.ID, results[0].ID)
}

func TestArtpieceRepository_List_FilterByCharacter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	char := seedCharacter(t, tx, user.ID)

	a1 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, Characters: []domain.Character{{ID: char.ID}}}
	a2 := &domain.Artpiece{ID: uuid.New(), UserID: user.ID}
	require.NoError(t, tx.Create(a1).Error)
	require.NoError(t, tx.Create(a2).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		CharacterIDs: []uuid.UUID{char.ID},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, a1.ID, results[0].ID)
}

func TestArtpieceRepository_List_CoverFileFieldsAccessible(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")

	fileName := "cover.jpg"
	file := &domain.File{
		ID:         uuid.New(),
		UserID:     user.ID,
		FileR2Path: "users/x/files/cover.jpg",
		MimeType:   "image/jpeg",
		Name:       &fileName,
	}
	require.NoError(t, tx.Create(file).Error)

	artpiece := &domain.Artpiece{
		ID:          uuid.New(),
		UserID:      user.ID,
		CoverFileID: &file.ID,
	}
	require.NoError(t, tx.Create(artpiece).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{})

	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NotNil(t, results[0].CoverFile)
	assert.Equal(t, "users/x/files/cover.jpg", results[0].CoverFile.FileR2Path)
	assert.Equal(t, "image/jpeg", results[0].CoverFile.MimeType)
	assert.Equal(t, "cover.jpg", *results[0].CoverFile.Name)
}

func TestArtpieceRepository_Create_IsReferencePersisted(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	a := &domain.Artpiece{
		ID:          uuid.New(),
		UserID:      user.ID,
		IsReference: true,
	}
	created, err := repo.Create(context.Background(), a)

	require.NoError(t, err)
	assert.True(t, created.IsReference)
}

func TestArtpieceRepository_Update_IsReferenceFromTrueToFalse(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	a := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: true}
	require.NoError(t, tx.Create(a).Error)

	f := false
	updated, err := repo.Update(context.Background(), a.ID, user.ID, usecase.UpdateArtpieceParams{
		IsReference: &f,
	})

	require.NoError(t, err)
	assert.False(t, updated.IsReference)
}

func TestArtpieceRepository_Update_IsReferenceOmittedLeavesUnchanged(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	a := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: true}
	require.NoError(t, tx.Create(a).Error)

	updated, err := repo.Update(context.Background(), a.ID, user.ID, usecase.UpdateArtpieceParams{})

	require.NoError(t, err)
	assert.True(t, updated.IsReference, "omitting IsReference must leave the field unchanged")
}

func TestArtpieceRepository_List_NoFilterReturnsBothKinds(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	ref := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: true}
	normal := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: false}
	require.NoError(t, tx.Create(ref).Error)
	require.NoError(t, tx.Create(normal).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{})

	require.NoError(t, err)
	require.Len(t, results, 2)
}

func TestArtpieceRepository_List_IsReferenceTrue_ReturnsOnlyReferences(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	ref := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: true}
	normal := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: false}
	require.NoError(t, tx.Create(ref).Error)
	require.NoError(t, tx.Create(normal).Error)

	v := true
	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{IsReference: &v})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, ref.ID, results[0].ID)
}

func TestArtpieceRepository_List_IsReferenceFalse_ReturnsOnlyNonReferences(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)

	user := seedUser(t, tx, "user_1")
	ref := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: true}
	normal := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, IsReference: false}
	require.NoError(t, tx.Create(ref).Error)
	require.NoError(t, tx.Create(normal).Error)

	v := false
	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{IsReference: &v})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, normal.ID, results[0].ID)
}

func TestArtpieceRepository_Delete_DoesNotCascadeToFiles(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	artpieceRepo := NewArtpieceRepository(tx)
	fileRepo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpiece(t, tx, user.ID)

	f := &domain.File{
		ID:         uuid.New(),
		UserID:     user.ID,
		ArtpieceID: &artpiece.ID,
		FileR2Path: "users/x/files/abc.jpg",
		MimeType:   "image/jpeg",
	}
	created, err := fileRepo.Create(context.Background(), f)
	require.NoError(t, err)

	err = artpieceRepo.Delete(context.Background(), artpiece.ID, user.ID)
	require.NoError(t, err)

	// file should still exist with artpiece_id NULL
	got, err := fileRepo.GetByIDAndUserID(context.Background(), created.ID, user.ID)
	require.NoError(t, err)
	assert.Nil(t, got.ArtpieceID, "file should have artpiece_id set to NULL after artpiece deletion")
}

// --- tags ---

func tagIDsOf(tags []domain.Tag) []uuid.UUID {
	ids := make([]uuid.UUID, len(tags))
	for i, tg := range tags {
		ids[i] = tg.ID
	}
	return ids
}

func artpieceIDsOf(as []*domain.Artpiece) []uuid.UUID {
	ids := make([]uuid.UUID, len(as))
	for i, a := range as {
		ids[i] = a.ID
	}
	return ids
}

func seedArtpieceWithTags(t *testing.T, tx *gorm.DB, userID uuid.UUID, tags ...*domain.Tag) *domain.Artpiece {
	t.Helper()
	a := &domain.Artpiece{ID: uuid.New(), UserID: userID}
	require.NoError(t, tx.Create(a).Error)
	for _, tg := range tags {
		require.NoError(t, tx.Exec("INSERT INTO artpiece_tags (artpiece_id, tag_id) VALUES (?, ?)", a.ID, tg.ID).Error)
	}
	return a
}

func TestArtpieceRepository_Create_WithTags(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	wip := seedTag(t, tx, user.ID, "wip")
	sketch := seedTag(t, tx, user.ID, "sketch")

	got, err := repo.Create(context.Background(), &domain.Artpiece{
		ID:     uuid.New(),
		UserID: user.ID,
		Tags:   []domain.Tag{{ID: wip.ID}, {ID: sketch.ID}},
	})

	require.NoError(t, err)
	assert.ElementsMatch(t, []uuid.UUID{wip.ID, sketch.ID}, tagIDsOf(got.Tags))
	var joinRows int64
	require.NoError(t, tx.Table("artpiece_tags").Where("artpiece_id = ?", got.ID).Count(&joinRows).Error)
	assert.Equal(t, int64(2), joinRows)

	// The association upsert must not alter the existing tag rows.
	var after domain.Tag
	require.NoError(t, tx.First(&after, "id = ?", wip.ID).Error)
	assert.Equal(t, "wip", after.Name)
	assert.Equal(t, user.ID, after.UserID)
}

func TestArtpieceRepository_Update_ReplacesTagSet(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	old := seedTag(t, tx, user.ID, "old")
	kept := seedTag(t, tx, user.ID, "kept")
	added := seedTag(t, tx, user.ID, "added")
	a := seedArtpieceWithTags(t, tx, user.ID, old, kept)

	got, err := repo.Update(context.Background(), a.ID, user.ID, usecase.UpdateArtpieceParams{
		TagIDs: []uuid.UUID{kept.ID, added.ID},
	})

	require.NoError(t, err)
	assert.ElementsMatch(t, []uuid.UUID{kept.ID, added.ID}, tagIDsOf(got.Tags))
	var oldRows int64
	require.NoError(t, tx.Model(&domain.Tag{}).Where("id = ?", old.ID).Count(&oldRows).Error)
	assert.Equal(t, int64(1), oldRows, "replaced tag must still exist")
}

func TestArtpieceRepository_Update_EmptyTagSetClears(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	tag := seedTag(t, tx, user.ID, "wip")
	a := seedArtpieceWithTags(t, tx, user.ID, tag)

	got, err := repo.Update(context.Background(), a.ID, user.ID, usecase.UpdateArtpieceParams{TagIDs: []uuid.UUID{}})

	require.NoError(t, err)
	assert.Empty(t, got.Tags)
	var tagRows int64
	require.NoError(t, tx.Model(&domain.Tag{}).Where("id = ?", tag.ID).Count(&tagRows).Error)
	assert.Equal(t, int64(1), tagRows)
}

func TestArtpieceRepository_Update_TagNotOwnedLeavesArtpieceUnchanged(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	other := seedUser(t, tx, "user_2")
	mine := seedTag(t, tx, user.ID, "mine")
	theirs := seedTag(t, tx, other.ID, "theirs")
	title := "original"
	a := &domain.Artpiece{ID: uuid.New(), UserID: user.ID, Title: &title}
	require.NoError(t, tx.Create(a).Error)
	require.NoError(t, tx.Exec("INSERT INTO artpiece_tags (artpiece_id, tag_id) VALUES (?, ?)", a.ID, mine.ID).Error)
	newTitle := "changed"

	got, err := repo.Update(context.Background(), a.ID, user.ID, usecase.UpdateArtpieceParams{
		Title:  &newTitle,
		TagIDs: []uuid.UUID{theirs.ID},
	})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, usecase.ErrTagNotOwned)
	after, err := repo.GetByID(context.Background(), a.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, after.Title)
	assert.Equal(t, "original", *after.Title)
	assert.Equal(t, []uuid.UUID{mine.ID}, tagIDsOf(after.Tags))
}

func TestArtpieceRepository_GetByID_PreloadsTags(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	wip := seedTag(t, tx, user.ID, "wip")
	tagged := seedArtpieceWithTags(t, tx, user.ID, wip)
	untagged := seedArtpiece(t, tx, user.ID)

	got, err := repo.GetByID(context.Background(), tagged.ID, user.ID)
	require.NoError(t, err)
	require.Len(t, got.Tags, 1)
	assert.Equal(t, wip.ID, got.Tags[0].ID)
	assert.Equal(t, "wip", got.Tags[0].Name)

	gotUntagged, err := repo.GetByID(context.Background(), untagged.ID, user.ID)
	require.NoError(t, err)
	assert.Empty(t, gotUntagged.Tags)
}

func TestArtpieceRepository_List_PreloadsTags(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	wip := seedTag(t, tx, user.ID, "wip")
	sketch := seedTag(t, tx, user.ID, "sketch")
	tagged := seedArtpieceWithTags(t, tx, user.ID, wip, sketch)
	untagged := seedArtpiece(t, tx, user.ID)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{})

	require.NoError(t, err)
	require.Len(t, results, 2)
	for _, r := range results {
		switch r.ID {
		case tagged.ID:
			assert.ElementsMatch(t, []uuid.UUID{wip.ID, sketch.ID}, tagIDsOf(r.Tags))
		case untagged.ID:
			assert.Empty(t, r.Tags)
		default:
			t.Fatalf("unexpected artpiece %s", r.ID)
		}
	}
}

func TestArtpieceRepository_List_FilterBySingleTag(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	a := seedTag(t, tx, user.ID, "a")
	b := seedTag(t, tx, user.ID, "b")
	withBoth := seedArtpieceWithTags(t, tx, user.ID, a, b)
	seedArtpieceWithTags(t, tx, user.ID, b)
	seedArtpiece(t, tx, user.ID)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{TagIDs: []uuid.UUID{a.ID}})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, withBoth.ID, results[0].ID)
	assert.ElementsMatch(t, []uuid.UUID{a.ID, b.ID}, tagIDsOf(results[0].Tags), "filtered result still shows all its tags")
}

func TestArtpieceRepository_List_TagMatchAnyReturnsUnionWithoutDuplicates(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	a := seedTag(t, tx, user.ID, "a")
	b := seedTag(t, tx, user.ID, "b")
	c := seedTag(t, tx, user.ID, "c")
	onlyA := seedArtpieceWithTags(t, tx, user.ID, a)
	onlyB := seedArtpieceWithTags(t, tx, user.ID, b)
	both := seedArtpieceWithTags(t, tx, user.ID, a, b)
	seedArtpieceWithTags(t, tx, user.ID, c)

	for _, match := range []string{"", "any"} {
		results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
			TagIDs:   []uuid.UUID{a.ID, b.ID},
			TagMatch: match,
		})

		require.NoError(t, err)
		assert.ElementsMatch(t, []uuid.UUID{onlyA.ID, onlyB.ID, both.ID}, artpieceIDsOf(results), "tag_match=%q", match)
	}
}

func TestArtpieceRepository_List_TagMatchAllRequiresEveryTag(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	a := seedTag(t, tx, user.ID, "a")
	b := seedTag(t, tx, user.ID, "b")
	c := seedTag(t, tx, user.ID, "c")
	seedArtpieceWithTags(t, tx, user.ID, a)
	both := seedArtpieceWithTags(t, tx, user.ID, a, b)
	bothPlusExtra := seedArtpieceWithTags(t, tx, user.ID, a, b, c)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		TagIDs:   []uuid.UUID{a.ID, b.ID},
		TagMatch: "all",
	})

	require.NoError(t, err)
	assert.ElementsMatch(t, []uuid.UUID{both.ID, bothPlusExtra.ID}, artpieceIDsOf(results))
}

func TestArtpieceRepository_List_TagMatchAllWithRepeatedIDs(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	a := seedTag(t, tx, user.ID, "a")
	tagged := seedArtpieceWithTags(t, tx, user.ID, a)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		TagIDs:   []uuid.UUID{a.ID, a.ID},
		TagMatch: "all",
	})

	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{tagged.ID}, artpieceIDsOf(results))
}

func TestArtpieceRepository_List_TagFilterCombinesWithArtistAndCharacter(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	tag := seedTag(t, tx, user.ID, "wip")
	artist := seedArtist(t, tx, user.ID, "Picasso")
	char := seedCharacter(t, tx, user.ID)

	match := seedArtpieceWithTags(t, tx, user.ID, tag)
	require.NoError(t, tx.Model(match).Update("artist_id", artist.ID).Error)
	require.NoError(t, tx.Exec("INSERT INTO artpiece_characters (artpiece_id, character_id) VALUES (?, ?)", match.ID, char.ID).Error)

	wrongArtist := seedArtpieceWithTags(t, tx, user.ID, tag)
	require.NoError(t, tx.Exec("INSERT INTO artpiece_characters (artpiece_id, character_id) VALUES (?, ?)", wrongArtist.ID, char.ID).Error)

	noCharacter := seedArtpieceWithTags(t, tx, user.ID, tag)
	require.NoError(t, tx.Model(noCharacter).Update("artist_id", artist.ID).Error)

	results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
		TagIDs:       []uuid.UUID{tag.ID},
		ArtistIDs:    []uuid.UUID{artist.ID},
		CharacterIDs: []uuid.UUID{char.ID},
	})

	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{match.ID}, artpieceIDsOf(results))
}

func TestArtpieceRepository_List_OtherUsersTagIDReturnsNothing(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	user := seedUser(t, tx, "user_1")
	other := seedUser(t, tx, "user_2")
	theirTag := seedTag(t, tx, other.ID, "theirs")
	seedArtpieceWithTags(t, tx, other.ID, theirTag)
	seedArtpiece(t, tx, user.ID)

	for _, match := range []string{"any", "all"} {
		results, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{
			TagIDs:   []uuid.UUID{theirTag.ID},
			TagMatch: match,
		})

		require.NoError(t, err)
		assert.Empty(t, results, "tag_match=%s", match)
	}
}

func TestArtpieceRepository_List_DeletedTagDisappearsFromResults(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewArtpieceRepository(tx)
	tagRepo := NewTagRepository(tx)
	user := seedUser(t, tx, "user_1")
	keep := seedTag(t, tx, user.ID, "keep")
	doomed := seedTag(t, tx, user.ID, "doomed")
	a := seedArtpieceWithTags(t, tx, user.ID, keep, doomed)

	require.NoError(t, tagRepo.Delete(context.Background(), doomed.ID, user.ID))

	byFilter, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{TagIDs: []uuid.UUID{doomed.ID}})
	require.NoError(t, err)
	assert.Empty(t, byFilter)

	all, err := repo.List(context.Background(), user.ID, usecase.ListArtpieceFilters{})
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, a.ID, all[0].ID)
	assert.Equal(t, []uuid.UUID{keep.ID}, tagIDsOf(all[0].Tags))
}
