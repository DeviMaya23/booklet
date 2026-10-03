package repository

import (
	"context"
	"testing"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedFile(t *testing.T, tx *gorm.DB, userID uuid.UUID, mimeType string) *domain.File {
	t.Helper()
	f := &domain.File{
		ID:         uuid.New(),
		UserID:     userID,
		FileR2Path: "users/" + userID.String() + "/files/test",
		MimeType:   mimeType,
	}
	require.NoError(t, tx.Create(f).Error)
	return f
}

func seedArtpieceForFile(t *testing.T, tx *gorm.DB, userID uuid.UUID) *domain.Artpiece {
	t.Helper()
	a := &domain.Artpiece{
		ID:     uuid.New(),
		UserID: userID,
	}
	require.NoError(t, tx.Create(a).Error)
	return a
}

func TestFileRepository_Create(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := &domain.File{
		ID:         uuid.New(),
		UserID:     user.ID,
		FileR2Path: "users/" + user.ID.String() + "/files/abc.jpg",
		MimeType:   "image/jpeg",
	}

	created, err := repo.Create(context.Background(), f)

	require.NoError(t, err)
	assert.Equal(t, f.ID, created.ID)
	assert.Equal(t, "image/jpeg", created.MimeType)
}

func TestFileRepository_GetByIDAndUserID_WrongUser(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	f := seedFile(t, tx, user1.ID, "image/jpeg")

	_, err := repo.GetByIDAndUserID(context.Background(), f.ID, user2.ID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileRepository_UpdateThumbnailPath(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := seedFile(t, tx, user.ID, "image/jpeg")

	err := repo.UpdateThumbnailPath(context.Background(), f.ID, "users/x/thumbnails/abc.jpg")
	require.NoError(t, err)

	got, err := repo.GetByIDAndUserID(context.Background(), f.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got.ThumbnailR2Path)
	assert.Equal(t, "users/x/thumbnails/abc.jpg", *got.ThumbnailR2Path)
}

func TestFileRepository_UpdateThumbnailGenState(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := seedFile(t, tx, user.ID, "image/jpeg")

	err := repo.UpdateThumbnailGenState(context.Background(), f.ID, "done")
	require.NoError(t, err)

	got, err := repo.GetByIDAndUserID(context.Background(), f.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got.ThumbnailGenState)
	assert.Equal(t, "done", *got.ThumbnailGenState)
}

func TestFileRepository_GetFilesForArtpiece(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpieceForFile(t, tx, user.ID)

	f1 := seedFile(t, tx, user.ID, "image/jpeg")
	f2 := seedFile(t, tx, user.ID, "video/mp4")
	artpieceID := artpiece.ID
	require.NoError(t, tx.Model(f1).Update("artpiece_id", artpieceID).Error)
	require.NoError(t, tx.Model(f2).Update("artpiece_id", artpieceID).Error)

	// file belonging to another artpiece — should not appear
	other := seedArtpieceForFile(t, tx, user.ID)
	f3 := seedFile(t, tx, user.ID, "image/png")
	require.NoError(t, tx.Model(f3).Update("artpiece_id", other.ID).Error)

	files, err := repo.GetFilesForArtpiece(context.Background(), artpieceID)
	require.NoError(t, err)
	assert.Len(t, files, 2)
}

func TestFileRepository_GetByIDsAndUserID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user1 := seedUser(t, tx, "user_1")
	user2 := seedUser(t, tx, "user_2")
	f1 := seedFile(t, tx, user1.ID, "image/jpeg")
	f2 := seedFile(t, tx, user1.ID, "video/mp4")
	_ = seedFile(t, tx, user2.ID, "image/png") // belongs to another user

	files, err := repo.GetByIDsAndUserID(context.Background(), []uuid.UUID{f1.ID, f2.ID}, user1.ID)
	require.NoError(t, err)
	assert.Len(t, files, 2)
}

func TestFileRepository_GetByIDsAndUserID_Empty(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	files, err := repo.GetByIDsAndUserID(context.Background(), nil, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, files)
}

func TestFileRepository_List_AllFiles(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpieceForFile(t, tx, user.ID)
	f1 := seedFile(t, tx, user.ID, "image/jpeg")
	f2 := seedFile(t, tx, user.ID, "video/mp4")
	require.NoError(t, tx.Model(f2).Update("artpiece_id", artpiece.ID).Error)

	files, err := repo.List(context.Background(), user.ID, false)
	require.NoError(t, err)
	ids := make([]uuid.UUID, len(files))
	for i, f := range files {
		ids[i] = f.ID
	}
	assert.Contains(t, ids, f1.ID)
	assert.Contains(t, ids, f2.ID)
}

func TestFileRepository_List_Unassigned(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpieceForFile(t, tx, user.ID)
	unassigned := seedFile(t, tx, user.ID, "image/jpeg")
	assigned := seedFile(t, tx, user.ID, "video/mp4")
	require.NoError(t, tx.Model(assigned).Update("artpiece_id", artpiece.ID).Error)

	files, err := repo.List(context.Background(), user.ID, true)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, unassigned.ID, files[0].ID)
}

func TestFileRepository_Update(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := seedFile(t, tx, user.ID, "image/jpeg")
	name := "my sketch.png"
	notes := "updated notes"

	err := repo.Update(context.Background(), f.ID, user.ID, &name, &notes)
	require.NoError(t, err)

	got, err := repo.GetByIDAndUserID(context.Background(), f.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got.Name)
	assert.Equal(t, name, *got.Name)
	require.NotNil(t, got.Notes)
	assert.Equal(t, notes, *got.Notes)
}

func TestFileRepository_Update_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	notes := "notes"

	err := repo.Update(context.Background(), uuid.New(), user.ID, nil, &notes)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileRepository_Delete(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	f := seedFile(t, tx, user.ID, "image/jpeg")

	err := repo.Delete(context.Background(), f.ID, user.ID)
	require.NoError(t, err)

	_, err = repo.GetByIDAndUserID(context.Background(), f.ID, user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileRepository_Delete_NotFound(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")

	err := repo.Delete(context.Background(), uuid.New(), user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFileRepository_BulkUpdateArtpieceID(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpieceForFile(t, tx, user.ID)
	f1 := seedFile(t, tx, user.ID, "image/jpeg")
	f2 := seedFile(t, tx, user.ID, "video/mp4")

	err := repo.BulkUpdateArtpieceID(context.Background(), []uuid.UUID{f1.ID, f2.ID}, &artpiece.ID)
	require.NoError(t, err)

	got1, err := repo.GetByIDAndUserID(context.Background(), f1.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got1.ArtpieceID)
	assert.Equal(t, artpiece.ID, *got1.ArtpieceID)

	got2, err := repo.GetByIDAndUserID(context.Background(), f2.ID, user.ID)
	require.NoError(t, err)
	require.NotNil(t, got2.ArtpieceID)
	assert.Equal(t, artpiece.ID, *got2.ArtpieceID)
}

func TestFileRepository_BulkUpdateArtpieceID_Detach(t *testing.T) {
	tx := testutil.NewTestTx(t, testDB)
	repo := NewFileRepository(tx)

	user := seedUser(t, tx, "user_1")
	artpiece := seedArtpieceForFile(t, tx, user.ID)
	f := seedFile(t, tx, user.ID, "image/jpeg")
	require.NoError(t, tx.Model(f).Update("artpiece_id", artpiece.ID).Error)

	err := repo.BulkUpdateArtpieceID(context.Background(), []uuid.UUID{f.ID}, nil)
	require.NoError(t, err)

	got, err := repo.GetByIDAndUserID(context.Background(), f.ID, user.ID)
	require.NoError(t, err)
	assert.Nil(t, got.ArtpieceID)
}
