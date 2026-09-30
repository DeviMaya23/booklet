## 1. Usecase interfaces

- [x] 1.1 Update `usecase/artist_repository.go`: change `GetByID`, `Update`, `Delete` signatures from `id string` to `id uuid.UUID`
- [x] 1.2 Update `usecase/character_repository.go`: change `GetByID`, `Update`, `Delete`, `UpdateAvatarR2Path`, `ClearAvatarR2Path` signatures from `id string` to `id uuid.UUID`
- [x] 1.3 Update `usecase/artpiece_repository.go`: change `ListArtpieceFilters.CharacterIDs` and `ArtistIDs` from `[]string` with `validate` tag to `[]uuid.UUID`
- [x] 1.4 Update `usecase/image_repository.go`: change `ListImageFilters.CharacterIDs` and `ArtistIDs` from `[]string` with `validate` tag to `[]uuid.UUID`; change `UpdateImageParams.CharacterIDs` from `[]string` to `[]uuid.UUID`

## 2. Usecase implementations

- [x] 2.1 Update `artist_usecase.go`: change `GetByID`, `Update`, `Delete` to take `id uuid.UUID`; remove any internal `uuid.Parse` calls on the ID
- [x] 2.2 Update `character_usecase.go`: change `GetByID`, `Update`, `Delete`, `InitAvatarUpload`, `CompleteAvatarUpload`, `DeleteAvatar` to take `id uuid.UUID`; remove any internal `uuid.Parse` calls on the ID

## 3. Repository implementations

- [x] 3.1 Update `artist_repository.go`: change `GetByID`, `Update`, `Delete` to take `id uuid.UUID`; merge `GetByIDAndUserID` into `GetByID` (remove `GetByIDAndUserID`, update call sites in `artpiece_usecase.go` to use `GetByID`)
- [x] 3.2 Update `character_repository.go`: change `GetByID`, `Update`, `Delete`, `UpdateAvatarR2Path`, `ClearAvatarR2Path`, `updateWithFolders` to take `id uuid.UUID`
- [x] 3.3 Update `image_repository.go` `List`: remove `::text` SQL casts on `CharacterIDs` and `ArtistIDs` filter conditions; update `Update` to remove per-element `uuid.Parse` loop over `CharacterIDs`

## 4. Handler layer — artpiece, upload, file_upload

- [x] 4.1 Update `artpiece_handler.go` request structs: change `createArtpieceRequest.ArtistID` to `*uuid.UUID`, `createArtpieceRequest.CharacterIDs` to `[]uuid.UUID`, `updateArtpieceRequest.ArtistID` to `*uuid.UUID`, `updateArtpieceRequest.CharacterIDs` to `[]uuid.UUID`, `setCoverRequest.FileID` to `uuid.UUID`; remove all `validate:"uuid4"` tags; remove post-bind `uuid.Parse` calls in `CreateArtpiece`, `UpdateArtpiece`, `SetCover`
- [x] 4.2 Update `upload_handler.go` request struct: change `ArtistID` to `*uuid.UUID`, `CharacterIDs` to `[]uuid.UUID`; remove `validate:"uuid4"` tags; remove post-bind parse loop and `uuid.Parse` call
- [x] 4.3 Update `file_upload_handler.go` request struct: change `ArtpieceID` to `*uuid.UUID`; remove `validate:"omitempty,uuid4"` tag; remove post-bind `uuid.Parse` call

## 5. Handler layer — artist

- [x] 5.1 Update `artist_handler.go` path param handling: replace two-line guard pattern (`c.Param` + discard parse) with single capturing `uuid.Parse` in `GetArtistByID`, `UpdateArtist`, `DeleteArtist`
- [x] 5.2 Update `ArtistUsecase` interface in `artist_handler.go`: change `id string` to `id uuid.UUID` on `GetByID`, `Update`, `Delete`

## 6. Handler layer — character

- [x] 6.1 Update `character_handler.go` path param handling: replace two-line guard pattern with single capturing `uuid.Parse` in `GetCharacterByID`, `UpdateCharacter`, `DeleteCharacter`, `InitAvatarUpload`, `CompleteAvatarUpload`, `DeleteAvatar`
- [x] 6.2 Update `CharacterUsecase` interface in `character_handler.go`: change `characterID string` to `characterID uuid.UUID` on `GetByID`, `Update`, `Delete`, `InitAvatarUpload`, `CompleteAvatarUpload`, `DeleteAvatar`
- [x] 6.3 Update `createCharacterRequest.FolderIDs` from `*[]string` to `*[]uuid.UUID`; update `updateCharacterRequest.FolderIDs` from `[]string` to `[]uuid.UUID`; remove `validate:"omitempty,dive,uuid4"` / `validate:"dive,uuid4"` tags
- [x] 6.4 Delete `parseFolderIDs` and `parseFolderIDSlice` helper functions; update `CreateCharacter` and `UpdateCharacter` to pass `req.FolderIDs` directly to usecase params

## 7. Tests

- [x] 7.1 Update artpiece handler tests: update any stub calls affected by request struct type changes
- [x] 7.2 Update artist handler tests: update stub interface to use `uuid.UUID` IDs
- [x] 7.3 Update character handler tests: update stub interface to use `uuid.UUID` IDs; update any `FolderIDs` string slices in test inputs
- [x] 7.4 Update `character_repository_integration_test.go`: update call sites for changed signatures

## 8. Lint

- [x] 8.1 Run `golangci-lint run ./...` from the backend directory and fix any issues
