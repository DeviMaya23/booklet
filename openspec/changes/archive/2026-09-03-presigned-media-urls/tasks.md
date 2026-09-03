## 1. Shared Presigner Interface

- [x] 1.1 Create `handler/presign.go` declaring the `Presigner` interface with `GeneratePresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error)`
- [x] 1.2 Add `PresignGetTTL = 1 * time.Hour` constant to `usecase/upload_usecase.go` alongside the existing `PresignTTL`

## 2. Character Handler — Presigned Avatar URL

- [x] 2.1 Add `presigner Presigner` field to `CharacterHandler` and update `NewCharacterHandler` to accept and store it
- [x] 2.2 Update `characterResponse` struct: replace `AvatarR2Path *string json:"avatar_r2_path"` with `AvatarURL *string json:"avatar_url"`
- [x] 2.3 Update `toCharacterResponse` to accept a presigned avatar URL string pointer instead of reading `AvatarR2Path` directly — signature becomes `toCharacterResponse(character *domain.Character, avatarURL *string) characterResponse`
- [x] 2.4 Update `GetCharacterByID`: if `character.AvatarR2Path != nil`, call `h.presigner.GeneratePresignedGetURL` to produce the URL; pass result to `toCharacterResponse`
- [x] 2.5 Update `ListCharacters`: presign avatar URL for each character in the loop; pass result to `toCharacterResponse`

## 3. Image Handler — Presigned Image and Thumbnail URLs

- [x] 3.1 Add `presigner Presigner` field to `ImageHandler` and update `NewImageHandler` to accept and store it
- [x] 3.2 Update `imageResponse` struct: remove `ImageR2Path`, rename `ThumbnailR2Path *string` → `ThumbnailURL *string json:"thumbnail_url"`; add `ImageURL *string json:"image_url,omitempty"` (omitted from list responses)
- [x] 3.3 Update `toImageResponse` to accept presigned URL arguments: `imageURL *string, thumbnailURL *string`; use these instead of raw R2 paths
- [x] 3.4 Update `GetImageByID`: presign `ImageR2Path` → `imageURL`; presign `ThumbnailR2Path` → `thumbnailURL` if non-nil; pass both to `toImageResponse`
- [x] 3.5 Update `ListImages`: presign `ThumbnailR2Path` → `thumbnailURL` for each image if non-nil; pass `nil` for `imageURL` (not included in list); pass to `toImageResponse`

## 4. Character Images — New Endpoint

- [x] 4.1 Add `ListByCharacterID(ctx context.Context, characterID uuid.UUID, userID uuid.UUID) ([]*domain.Image, error)` to the `ImageRepository` interface in `usecase/image_repository.go`
- [x] 4.2 Implement `ListByCharacterID` in `repository/image_repository.go`: join `images` with `image_characters` on `image_id`, filter by `character_id = ?` and `images.user_id = ?`, no preloads needed (only `id`, `title`, `thumbnail_r2_path` are used)
- [x] 4.3 Add `GetCharacterImages(ctx context.Context, characterID string, userID uuid.UUID) ([]*domain.Image, error)` to the `CharacterUsecase` interface in `handler/character_handler.go`
- [x] 4.4 Implement `GetCharacterImages` in `usecase/character_usecase.go`: validate UUID, call `imageRepo.ListByCharacterID`, return results. Add `imageRepo ImageRepository` dependency to `characterUsecase` struct and `NewCharacterUsecase`
- [x] 4.5 Add `characterImageResponse` struct to `handler/character_handler.go` with fields: `ImageID string`, `ImageName *string`, `ThumbnailURL *string`
- [x] 4.6 Implement `GetCharacterImages` handler method in `handler/character_handler.go`: validate UUID param, get userID from context, call usecase, presign thumbnail for each result, return 200 with array
- [x] 4.7 Register route `GET /characters/:id/images` in `main.go`
- [x] 4.8 Wire `imageRepository` into `characterUsecase` constructor call in `main.go`
- [x] 4.9 Create Bruno file `collection/characters/get_character_images.bru`

## 5. Wiring

- [x] 5.1 Update `NewCharacterHandler` call in `main.go` to pass `r2Storage` as `Presigner`
- [x] 5.2 Update `NewImageHandler` call in `main.go` to pass `r2Storage` as `Presigner`

## 6. Unit Tests

- [x] 6.1 Add `CharacterUsecase.GetCharacterImages` unit tests: success (returns images), success with empty result (character has no images or is unowned)
- [x] 6.2 Add `CharacterHandler.GetCharacterByID` unit test: avatar URL is presigned (not raw path) when avatar is set; avatar_url is null when no avatar
- [x] 6.3 Add `CharacterHandler.ListCharacters` unit test: avatar URLs are presigned for each character with an avatar
- [x] 6.4 Add `ImageHandler.GetImageByID` unit test: `image_url` and `thumbnail_url` are presigned; `image_r2_path` and `thumbnail_r2_path` absent from response
- [x] 6.5 Add `ImageHandler.ListImages` unit test: `thumbnail_url` is presigned when present; `image_url` is absent from list items
- [x] 6.6 Add `CharacterHandler.GetCharacterImages` unit test: returns presigned thumbnail URLs; returns empty array when no images

## 7. Lint

- [x] 7.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
