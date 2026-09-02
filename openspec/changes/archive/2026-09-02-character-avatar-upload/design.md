## Context

Characters have a `hero_image_r2_path` column that was reserved for a profile image but was never used or exposed. The image upload flow (`pending_uploads` + `POST /images` + `POST /images/:id/complete`) is image-specific: it carries image metadata (title, artist, notes, character associations) and triggers thumbnail generation on completion. A character avatar is a simpler case — a single file stored as-is, tied to one character, with no downstream processing.

The existing `StorageService` interface already exposes both `GeneratePresignedPutURL` and `DeleteObject`, so R2 integration is already abstracted.

## Goals / Non-Goals

**Goals:**
- Rename `hero_image_r2_path` → `avatar_r2_path` on `characters` via migration
- Two-phase upload flow for character avatars: init (presign) + complete (update character, delete old object)
- Synchronous deletion of the previous R2 object when a character already has an avatar
- Stale pending avatar upload cleanup integrated with the existing purge worker

**Non-Goals:**
- Thumbnail generation for avatars
- Serving/proxying the avatar (clients construct the R2 URL directly from the path)
- Generalizing the pending table across file types (deferred — see proposal)
- Updating the `HeroImageR2Path` field via the existing create/update character endpoints (those fields are removed entirely)

## Decisions

### Separate `pending_character_avatar_uploads` table
The existing `pending_uploads` table is image-specific: it carries `title`, `artist_id`, `notes`, `character_ids` which are meaningless for an avatar upload. A dedicated table keeps each flow self-contained and avoids adding a discriminator to an existing table. Future generalization cost is low (pending rows are ephemeral; blast radius is ~4 files).

Schema:
```
pending_character_avatar_uploads
  id           uuid         PK
  user_id      uuid         NOT NULL → users(id)
  character_id uuid         NOT NULL → characters(id) ON DELETE CASCADE
  r2_key       text         NOT NULL
  mime_type    text         NOT NULL
  created_at   timestamptz  NOT NULL
```

`character_id` is a FK to `characters`. If the character is deleted mid-upload, the pending row is cascaded away.

### Character ownership checked at init, not just complete
The init handler verifies that the character exists and belongs to the authenticated user before generating a presigned URL. This prevents generating presigned URLs for characters the caller does not own. Consistent with how `upload_usecase` validates artist ownership at init time.

### Synchronous R2 deletion of previous avatar on complete
When `avatar_r2_path` is already set on the character at complete time, the old R2 object is deleted before updating the column. This is a direct `storage.DeleteObject` call within the complete usecase — no River job. Rationale: R2 deletes are fast and cheap; a failure here orphans a file but does not break the user experience. Retry machinery is not worth the added complexity for this case.

### Complete updates character in-place; no new entity
Unlike the image flow which creates an `Image` row on complete, avatar complete calls a targeted `UpdateAvatarR2Path` on the character repository. This keeps the avatar as a property of the character, not a standalone entity.

### Stale cleanup via existing `PurgeExpiredUploads` worker pattern
A new `PurgeExpiredCharacterAvatarUploads` River job + worker is added alongside the existing `PurgeExpiredUploads` job, scheduled at the same 5-minute interval. `CleanupStaleAvatarUploads(ctx, threshold)` is added as a method on `characterUsecase`, mirroring the existing pattern in `upload_usecase.go`. This avoids coupling the two pending tables into a shared cleanup path.

### Avatar upload folded into `characterUsecase` and `CharacterHandler`
`InitAvatarUpload` and `CompleteAvatarUpload` are added as methods on the existing `characterUsecase` rather than a separate usecase. The avatar upload is not creating a new entity — it is updating an attribute on an already-existing character. Treating it as part of character management is the more accurate model. `characterUsecase` gains `StorageService` and `CharacterAvatarUploadRepository` as constructor dependencies; the two new handler methods are added to `CharacterHandler`. This keeps all character-related logic in one usecase and avoids a separate handler with duplicated wiring.

### R2 key format
`{userID}/files/{uuid}` — no file extension. Avatar files are served as-is; the mime type is stored in the pending row for the presign call but the key itself carries no extension. Consistent with the stated path format in the proposal.

## Risks / Trade-offs

- **Orphaned R2 objects on failed complete**: If the character update fails after a successful R2 delete, the old file is gone with no record. Mitigation: delete old object *after* a successful DB update (reverse the order — read old key, update DB, then delete old key).
- **Concurrent completes**: Two complete calls for the same pending upload ID race. The second will find the pending row deleted (by the first) and return 404. The presigned URL is single-use from R2's perspective anyway. Acceptable.
- **Stale pending rows**: Covered by the purge worker. The R2 objects at those keys are never confirmed and can be cleaned up by the same job that deletes the DB row.

## Migration Plan

1. Migration `000011`: Rename `hero_image_r2_path` → `avatar_r2_path` on `characters`
2. Migration `000012`: Create `pending_character_avatar_uploads` table
3. Deploy: domain, usecase, handler, routes updated atomically with migrations
4. Rollback: reverse migrations; restore old field names in code (column was unused so no data risk)
