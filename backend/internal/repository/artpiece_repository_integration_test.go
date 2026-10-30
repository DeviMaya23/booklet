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
