## 1. Migrations

- [x] 1.1 Create migration `000011`: rename `hero_image_r2_path` → `avatar_r2_path` on `characters` table
- [x] 1.2 Create migration `000012`: create `pending_character_avatar_uploads` table (`id`, `user_id`, `character_id` FK with ON DELETE CASCADE, `r2_key`, `mime_type`, `created_at`)

## 2. Domain & Usecase Interfaces

- [x] 2.1 Rename `Character.HeroImageR2Path` → `Character.AvatarR2Path` in `domain/character.go`; update GORM column tag
- [x] 2.2 Rename `CreateCharacterParams.HeroImageR2Path` → `AvatarR2Path` and `UpdateCharacterParams.HeroImageR2Path` → `AvatarR2Path` in `usecase/character_repository.go`
- [x] 2.3 Add `domain.PendingCharacterAvatarUpload` struct in `domain/` (`id`, `user_id`, `character_id`, `r2_key`, `mime_type`, `created_at`)
- [x] 2.4 Add to `usecase/character_repository.go`: `CharacterAvatarUploadRepository` interface (Create, GetByID, Delete, ListStale) and `UpdateAvatarR2Path` method to `CharacterRepository` interface

## 3. Character Avatar Upload Usecase

- [x] 3.1 Add `StorageService` and `CharacterAvatarUploadRepository` as dependencies to `characterUsecase` constructor in `usecase/character_usecase.go`
- [x] 3.2 Implement `InitAvatarUpload(ctx, userID, characterID, mimeType)` on `characterUsecase` — verify character ownership, generate R2 key (`{userID}/files/{uuid}`), presign PUT URL, create pending row, return `(id, upload_url, expires_at)`
- [x] 3.3 Implement `CompleteAvatarUpload(ctx, userID, characterID, uploadID)` on `characterUsecase` — fetch pending row (verify user ownership), snapshot old `avatar_r2_path`, update character and delete pending row in a transaction, then delete old R2 object if it existed
- [x] 3.4 Implement `CleanupStaleAvatarUploads(ctx, threshold)` on `characterUsecase` — list stale pending rows, delete R2 objects, delete DB rows (mirrors `upload_usecase.CleanupStaleUploads`)

## 4. Character Usecase Tests

- [x] 4.1 Test `InitAvatarUpload`: character not owned returns `ErrCharacterNotFound`; successful init returns result with upload URL and creates pending row
- [x] 4.2 Test `CompleteAvatarUpload`: pending upload not found returns error; successful complete with no previous avatar sets `avatar_r2_path`; successful complete with existing avatar sets new path and deletes old R2 object

## 5. Character Avatar Upload Repository

- [x] 5.1 Create `repository/character_avatar_repository.go` implementing `CharacterAvatarUploadRepository`: `Create`, `GetByID` (scoped by userID), `Delete`, `ListStale`
- [x] 5.2 Add `UpdateAvatarR2Path(ctx, characterID, userID, r2Key string) error` to `characterRepository`

## 6. Existing Character Code Updates

- [x] 6.1 Update `character_usecase.go`: remove `HeroImageR2Path` from `CreateCharacterParams` usage; `AvatarR2Path` is not accepted in create (set to nil — avatar is set only via upload flow)
- [x] 6.2 Update `character_repository.go` (`Update` method): rename `HeroImageR2Path` → `AvatarR2Path` in update map
- [x] 6.3 Update `handler/character_handler.go`: remove `hero_image_r2_path` from request structs; add `avatar_r2_path` (read-only) to `characterResponse` and `toCharacterResponse`; add `InitAvatarUpload` and `CompleteAvatarUpload` to the `CharacterUsecase` interface

## 7. Character Avatar Handler Methods

- [x] 7.1 Add `InitAvatarUpload` handler method to `CharacterHandler`: validate character UUID path param and mime_type body, call usecase, return 201 with `{ id, upload_url, expires_at }`
- [x] 7.2 Add `CompleteAvatarUpload` handler method to `CharacterHandler`: validate character UUID and upload UUID path params, call usecase, return 204; map not-found to 404

## 8. Handler Tests

- [x] 8.1 Test `InitAvatarUpload` handler: character not found returns 404; missing mime_type returns 422; unsupported mime_type returns 422; successful init returns 201 with correct body
- [x] 8.2 Test `CompleteAvatarUpload` handler: pending upload not found returns 404; successful complete returns 204

## 9. Stale Cleanup Worker

- [x] 9.1 Create `worker/purge_expired_character_avatar_uploads.go` with `PurgeExpiredCharacterAvatarUploadsArgs` and worker, mirroring the existing `purge_expired_uploads` worker
- [x] 9.2 Register worker and periodic job (5-minute interval) in `cmd/server/main.go` alongside the existing purge upload job

## 10. Routes & Wiring

- [x] 10.1 Instantiate `characterAvatarRepository` and pass it (along with `r2Storage`) into the updated `characterUsecase` constructor in `cmd/server/main.go`
- [x] 10.2 Register routes: `POST /characters/:id/avatar/init` and `POST /characters/:id/avatar/:uploadID/complete` on the existing `characterHandler`

## 11. Bruno Collection

- [x] 11.1 Create `collection/characters/init_avatar_upload.bru`
- [x] 11.2 Create `collection/characters/complete_avatar_upload.bru`
- [x] 11.3 Create `collection/characters/delete_avatar.bru`

## 12. Avatar Deletion

- [x] 12.1 Add `DeleteAvatar(ctx context.Context, userID uuid.UUID, characterID string) error` to `characterUsecase` — verify character ownership, if `avatar_r2_path` is null return nil (idempotent), set `avatar_r2_path = null` in DB, then synchronously delete R2 object
- [x] 12.2 Add `ClearAvatarR2Path(ctx context.Context, id string, userID uuid.UUID) (oldKey string, err error) error` to `CharacterRepository` interface and `characterRepository` — updates `avatar_r2_path` to null and returns the previous value so the usecase can delete the R2 object
- [x] 12.3 Add `DeleteAvatar` to the `CharacterUsecase` interface in `handler/character_handler.go`
- [x] 12.4 Add `DeleteAvatar` handler method to `CharacterHandler` — validate character UUID, call usecase, return 204; map not-found to 404
- [x] 12.5 Register route `DELETE /characters/:id/avatar` in `cmd/server/main.go`
- [x] 12.6 Test `DeleteAvatar` usecase: character not found returns error; no avatar returns nil without R2 call; existing avatar sets path to null and deletes R2 object
- [x] 12.7 Test `DeleteAvatar` handler: character not found returns 404; successful deletion returns 204

## 13. Lint

- [x] 13.1 Run `golangci-lint run ./...` and fix any issues
