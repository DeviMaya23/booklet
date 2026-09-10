## Why

The update endpoints for artists, characters, and images use HTTP PATCH, but PATCH provides no real value for this app. PATCH is justified when clients need to send partial state — for collaborative scenarios, large payloads, or clients that don't hold the full resource. None of those apply here: this is a single-user app where every edit form loads the full entity before opening. PUT with replace semantics is the correct fit: the client always has complete state, always sends it, and the server always replaces it. The PATCH choice has introduced unnecessary complexity (three-state field wrappers, conditional payload logic) with no benefit in return.

## What Changes

- **BREAKING**: `PATCH /artists/:id` → `PUT /artists/:id`. Request body always includes all user-editable fields; absent optional fields are treated as `null` (clear).
- **BREAKING**: `PATCH /characters/:id` → `PUT /characters/:id`. Same semantics.
- **BREAKING**: `PATCH /images/:id` → `PUT /images/:id`. `artist_id` drops the `Patch[T]` wrapper and becomes a plain nullable string.
- BE: `UpdateArtistParams`, `UpdateCharacterParams`, `UpdateImageParams` structs simplified — no more pointer-of-pointer or `Patch[T]` fields.
- BE: Repository update logic simplified — no `if param != nil` guards for nullable fields; always write what is given.
- FE: `useUpdateArtist` hook updated to `PUT` and `UpdateArtistInput` type made fully required.
- FE: `ArtistFormModal` payload logic simplified — always sends all fields, converts empty string to `null`.

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `artist-management`: update endpoint changes from PATCH to PUT with replace semantics
- `character-management`: update endpoint changes from PATCH to PUT with replace semantics
- `image-management`: update endpoint changes from PATCH to PUT; `artist_id` handling simplified
- `web-artists`: `useUpdateArtist` hook and `ArtistFormModal` payload logic updated

## Impact

- `backend/internal/handler/artist_handler.go` — route method, request struct
- `backend/internal/handler/character_handler.go` — route method, request struct
- `backend/internal/handler/image_handler.go` — route method, request struct, `Patch[T]` removal
- `backend/internal/usecase/artist_repository.go` — `UpdateArtistParams` struct
- `backend/internal/usecase/character_repository.go` — `UpdateCharacterParams` struct
- `backend/internal/usecase/image_repository.go` — `UpdateImageParams` struct
- `backend/internal/repository/artist_repository.go` — update logic
- `backend/internal/repository/character_repository.go` — update logic
- `backend/internal/repository/image_repository.go` — update logic, `Patch` removal
- `backend/internal/handler/artist_handler_test.go` — test cases
- `backend/internal/handler/character_handler_test.go` — test cases
- `backend/internal/handler/image_handler_test.go` — test cases
- `frontend/src/features/artists/api/useUpdateArtist.ts` — method + type
- `frontend/src/features/artists/components/ArtistFormModal.tsx` — payload construction
- Bruno collection files for artist, character, image update requests
