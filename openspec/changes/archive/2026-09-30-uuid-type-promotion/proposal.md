## Why

UUID path params and body fields flow as `string` through the artist and character handler/usecase/repository stack, forcing repeated `uuid.Parse` calls at every layer boundary and `validate:"uuid4"` tags that only exist to guard a manual parse. Artpiece, upload, and file-upload handlers have the same issue in their request body structs despite having correct path-param handling. This change eliminates all remaining `string`-typed UUID fields across the full BE handler surface.

## What Changes

- Path param validation in artist and character handlers: two-line guard (`c.Param` + `uuid.Parse` discard) collapses to a single capturing `uuid.Parse`
- Request body UUID fields (`ArtistID *string`, `CharacterIDs []string`, `FolderIDs *[]string / []string`) change to native UUID types — `c.Bind` handles parsing via `encoding.TextUnmarshaler`; `validate:"uuid4"` tags and post-bind parse calls are removed
- `ListArtpieceFilters` and `ListImageFilters` query-param slices (`CharacterIDs []string`, `ArtistIDs []string`) change to `[]uuid.UUID`; `validate:"omitempty,dive,uuid4"` tags removed
- Usecase interfaces for artist and character change `id string` → `id uuid.UUID` on all methods that take a record ID
- Repository implementations for artist and character updated to match; `artist_repository.GetByIDAndUserID` (already `uuid.UUID`) merges into `GetByID` since signatures become identical
- `image_repository.List` drops `::text` SQL casts on UUID filter columns
- **BREAKING** (internal): invalid UUID in a request body now returns HTTP 400 (bind failure) instead of HTTP 422 (validator failure). No field name in the error. Acceptable because UUIDs are app-managed and never entered by users directly.

## Capabilities

### New Capabilities

_None — this is a consistency/cleanup change with no new features._

### Modified Capabilities

- `request-validation`: UUID body fields no longer go through `go-playground/validator`; they fail at `c.Bind` with HTTP 400 instead of HTTP 422 with a field-level error

## Impact

**Backend — handler layer**
- `artist_handler.go`: 3 handler functions, handler interface
- `character_handler.go`: 6 handler functions, handler interface, 2 request structs, `parseFolderIDs` + `parseFolderIDSlice` helpers deleted
- `artpiece_handler.go`: `createArtpieceRequest` and `updateArtpieceRequest` (`ArtistID *string`, `CharacterIDs []string`); `setCoverRequest` (`FileID string`); post-bind parse calls in Create, Update, SetCover
- `upload_handler.go`: `ArtistID *string`, `CharacterIDs []string`; post-bind parse calls
- `file_upload_handler.go`: `ArtpieceID *string`; post-bind parse call

**Backend — usecase layer**
- `usecase/artist_repository.go` (interface): `GetByID`, `Update`, `Delete`
- `usecase/character_repository.go` (interface): `GetByID`, `Update`, `Delete`, `UpdateAvatarR2Path`, `ClearAvatarR2Path`
- `usecase/artpiece_repository.go`: `ListArtpieceFilters` field types
- `usecase/image_repository.go`: `ListImageFilters` field types, `UpdateImageParams.CharacterIDs`
- `artist_usecase.go` impl: 3 functions
- `character_usecase.go` impl: 6 functions

**Backend — repository layer**
- `artist_repository.go`: 3 functions; `GetByIDAndUserID` merged into `GetByID`
- `character_repository.go`: 5 functions + `updateWithFolders`
- `image_repository.go`: `List` drops `::text` SQL casts; `Update` drops per-element `uuid.Parse` loop

**Tests**
- Handler tests: artist, character — stub interfaces updated
- Integration tests: `character_repository_integration_test.go`
