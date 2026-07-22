## 1. Migration

- [x] 1.1 Create migration `000006_create_character_folders.up.sql` — `character_folders` table with `character_id` (uuid, FK → characters.id), `folder_id` (uuid), composite PK `(character_id, folder_id)`
- [x] 1.2 Create migration `000006_create_character_folders.down.sql` — drop `character_folders`

## 2. Config & Bookleaf Client

- [x] 2.1 Add `BookleafHost` and `BookleafInternalSecret` fields to a new `BookleafConfig` struct in `Config`; parse `BOOKLEAF_HOST` and `BOOKLEAF_INTERNAL_SECRET` as required env vars in `config.go`
- [x] 2.2 Create `internal/bookleaf/` package with a `Client` struct and `GetPublicFolders(ctx, userID string)` method that calls `GET <host>/internal/users/:user_id/public-folders` with the `X-Bookleaf-Internal-Secret` header set to the configured secret, and returns the response decoded into a typed struct (`FolderList` with `[]Folder{FolderID, Token, FolderName}`)
- [x] 2.3 Wire bookleaf client construction in `main.go` `initApp` using `cfg.Bookleaf`

## 3. Domain

- [x] 3.1 Add `domain.CharacterFolder` struct with `CharacterID uuid.UUID` and `FolderID uuid.UUID` (GORM model, no soft-delete)
- [x] 3.2 Add `Folders []CharacterFolder` association field to `domain.Character` with appropriate GORM tag (`has many`, `foreignKey:CharacterID`)

## 4. Repository

- [x] 4.1 Update `characterRepository.Create` — after inserting the character, if `FolderIDs` is non-nil and non-empty, insert `character_folders` rows in a DB transaction (begin transaction, insert character, bulk-insert folders, commit)
- [x] 4.2 Update `characterRepository.GetByID` — preload `Folders` association when fetching
- [x] 4.3 Update `characterRepository.List` — preload `Folders` association when fetching
- [x] 4.4 Update `characterRepository.Update` — when `FolderIDs != nil`, replace all folder rows for the character in a transaction (delete existing `character_folders` for the character, bulk-insert new rows, then run character field updates); when `FolderIDs == nil`, leave folders unchanged
- [x] 4.5 Update `characterRepository.Delete` — hard-delete all `character_folders` rows for the character in the same transaction as the soft-delete
- [x] 4.6 Update `CharacterRepository` interface in `usecase/character_repository.go` — add `FolderIDs *[]uuid.UUID` to `UpdateCharacterParams` and extend `CreateCharacterParams` (move params struct here or add field)

## 5. Usecase

- [x] 5.1 Add `FolderIDs *[]uuid.UUID` to `CreateCharacterParams` in `usecase/character_usecase.go` and pass it through to `characterRepo.Create`
- [x] 5.2 Add `FolderIDs *[]uuid.UUID` to `UpdateCharacterParams` in `usecase/character_repository.go` and pass it through to `characterRepo.Update`
- [x] 5.3 Create `usecase/folder_usecase.go` — define `BookleafClient` interface (method: `GetPublicFolders(ctx, userID string)`) and `folderUsecase` struct with a `ListFolders(ctx, userID string)` method that delegates to the client
- [x] 5.4 Create `usecase/folder_usecase_test.go` — no happy-path unit test (pure delegation); no error-path test (pure pass-through); file may be omitted or minimal per conventions

## 6. Handler

- [x] 6.1 Update `createCharacterRequest` — add `FolderIDs *[]uuid.UUID \`json:"folder_ids" validate:"omitempty,dive,uuid4"\``
- [x] 6.2 Update `updateCharacterRequest` — add `FolderIDs *[]uuid.UUID \`json:"folder_ids" validate:"omitempty,dive,uuid4"\``
- [x] 6.3 Update `characterResponse` — add `FolderIDs []string \`json:"folder_ids"\`` (emit empty array, never null)
- [x] 6.4 Update `toCharacterResponse` — populate `FolderIDs` from `character.Folders`
- [x] 6.5 Update `CreateCharacter` handler — pass `req.FolderIDs` into `CreateCharacterParams`
- [x] 6.6 Update `UpdateCharacter` handler — pass `req.FolderIDs` into `UpdateCharacterParams`
- [x] 6.7 Create `handler/folder_handler.go` — `FolderUsecase` interface with `ListFolders`, `FolderHandler` struct, `GET /folders` handler that extracts `userID` from JWT and returns the bookleaf response as-is
- [x] 6.8 Register `GET /folders` route in `main.go` `initApp`

## 7. Tests — Usecase Unit Tests

- [x] 7.1 Update `character_usecase_test.go` `TestCreate_AssemblesCharacter` — assert `folder_ids` is passed through to repo when provided
- [x] 7.2 Add `TestCreate_AssemblesCharacter_WithFolders` — verify `FolderIDs` is populated on the character passed to repo
- [x] 7.3 Update `fakes_test.go` `fakeCharacterRepository` — propagate `FolderIDs` from `CreateCharacterParams` and `UpdateCharacterParams` to the stored character

## 8. Tests — Handler Unit Tests

- [x] 8.1 Update `character_handler_test.go` — add folder_ids to create/update happy path assertions; assert `folder_ids: []` in response when no folders assigned
- [x] 8.2 Add handler test: `POST /characters` with invalid UUID in `folder_ids` returns 422
- [x] 8.3 Add handler test: `PATCH /characters/:id` with invalid UUID in `folder_ids` returns 422
- [x] 8.4 Create `handler/folder_handler_test.go` — test happy path (200 with proxied response body), test upstream error maps to 500

## 9. Tests — Repository Integration Tests

- [x] 9.1 Update `TestCharacterRepository_Create` — seed character with folder IDs, assert `character_folders` rows exist after creation
- [x] 9.2 Add `TestCharacterRepository_Create_WithFolders` — verify folder rows are inserted and preloaded on subsequent GetByID
- [x] 9.3 Update `TestCharacterRepository_GetByID` — assert `Folders` is populated (including empty case)
- [x] 9.4 Update `TestCharacterRepository_List` — assert `Folders` is populated on each returned character
- [x] 9.5 Add `TestCharacterRepository_Update_ReplaceFolders` — update with new folder set, assert old rows gone and new rows present
- [x] 9.6 Add `TestCharacterRepository_Update_ClearFolders` — update with empty slice, assert all folder rows removed
- [x] 9.7 Add `TestCharacterRepository_Update_NilFolders_NoOp` — update without folder_ids, assert folder rows unchanged
- [x] 9.8 Add `TestCharacterRepository_Delete_CascadesFolders` — delete character with folders, assert `character_folders` rows are hard-deleted

## 10. Bruno Collection

- [x] 10.1 Create Bruno file for `GET /folders`

## 11. Lint

- [x] 11.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
