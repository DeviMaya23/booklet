## 1. Artist — Backend

- [x] 1.1 Update `UpdateArtistParams` in `usecase/artist_repository.go`: `Notes` and `ArtistLink` remain `*string` (semantics unchanged at param level; nil now means "clear")
- [x] 1.2 Update `updateArtistRequest` in `handler/artist_handler.go`: change method registration to `PUT`; `Notes` and `ArtistLink` are still `*string` in the struct (JSON null decodes to nil, which now means clear)
- [x] 1.3 Update `artistRepository.Update` in `repository/artist_repository.go`: remove `if params.Notes != nil` / `if params.ArtistLink != nil` guards — always include both fields in the updates map (use `gorm.Expr("NULL")` or direct nil assignment for null values)
- [x] 1.4 Update artist handler tests in `handler/artist_handler_test.go`: change test requests to `PUT`, add scenarios for clearing `notes` and `artist_link` via null
- [x] 1.5 Update artist Bruno file: change method from PATCH to PUT, add null-field examples

## 2. Character — Backend

- [x] 2.1 Update `UpdateCharacterParams` in `usecase/character_repository.go`: `Biography` remains `*string`; mark `Name` and `IsPublic` as always-present (non-pointer) — assess if handler can enforce presence via validate tags without param struct change
- [x] 2.2 Update `updateCharacterRequest` in `handler/character_handler.go`: change method to `PUT`; `Name` becomes required (validate:"required,min=1"), `IsPublic` becomes required `bool` (non-pointer), `Biography` stays `*string` (null = clear), `FolderIDs` becomes required `[]string` (always sent)
- [x] 2.3 Update `characterRepository.Update` in `repository/character_repository.go`: remove `if params.Biography != nil` guard — always write biography; `FolderIDs` path unchanged (replace-all already works)
- [x] 2.4 Update character handler tests in `handler/character_handler_test.go`: change test requests to `PUT`, add scenarios for clearing `biography`, always-present `is_public` and `folder_ids`
- [x] 2.5 Update character Bruno file: change method from PATCH to PUT

## 3. Image — Backend

- [x] 3.1 Update `UpdateImageParams` in `usecase/image_repository.go`: `ArtistID` changes from `**uuid.UUID` to `*uuid.UUID` (null = clear); `Title` and `Notes` remain `*string` (null = clear); `CharacterIDs` changes from `*[]string` to `[]string` (always present)
- [x] 3.2 Update `updateImageRequest` in `handler/image_handler.go`: change method to `PUT`; replace `ArtistID Patch[string]` with `ArtistID *string`; `Title` and `Notes` stay `*string`; `CharacterIDs` becomes `[]string` (required, non-pointer); remove `Patch[T]` conversion logic from `UpdateImage` handler
- [x] 3.3 Update `imageRepository.Update` in `repository/image_repository.go`: remove `if params.Title != nil` / `if params.Notes != nil` guards — always write both; simplify `ArtistID` handling from double-pointer to single-pointer; `CharacterIDs` replace-all logic unchanged
- [x] 3.4 Remove `Patch[T]` usage from `image_handler.go` only; leave `handler/patch.go` and its validator registration in place
- [x] 3.5 Update image handler tests in `handler/image_handler_test.go`: change test requests to `PUT`, update `artist_id` field format, add null-clear scenarios for `title` and `notes`
- [x] 3.6 Update image Bruno file: change method from PATCH to PUT

## 4. Artist — Frontend

- [x] 4.1 Update `useUpdateArtist.ts`: change `method` from `'PATCH'` to `'PUT'`; update `UpdateArtistInput` type — `name` required, `notes` and `artist_link` typed as `string | null` (not optional)
- [x] 4.2 Update `ArtistFormModal.tsx` `handleSubmit`: always include all three fields in the payload; convert empty string to `null` for `artist_link` and `notes` (`artistLink.trim() || null`, `notes.trim() || null`)

## 5. Final Checks

- [x] 5.1 Run `npm run build` and `npm run lint` in `frontend/`; fix any issues
- [x] 5.2 Run `golangci-lint run ./...` in `backend/`; fix any issues
