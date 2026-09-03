## Why

The list endpoints for images, artists, and characters return all records for a user with no filtering. The frontend needs freetext search and ID-based filtering on these endpoints before list views can be built.

## What Changes

- `GET /images` accepts optional query params: `q` (freetext title search), `character_ids` (multi-value UUID), `artist_ids` (multi-value UUID)
- `GET /artists` accepts optional query param: `q` (freetext name search)
- `GET /characters` accepts optional query param: `q` (freetext name search)
- Freetext search is case-insensitive substring (`ILIKE '%q%'`)
- ID filters use `IN` semantics — return images matching any of the given IDs
- All filters are optional; omitting them preserves current behavior (return all)
- Filter structs carry `query` tags and are defined in the usecase package; `c.Bind` is used in handlers — one struct per entity, no mapping layer

## Capabilities

### New Capabilities

- `list-filters`: Query-param filtering on the image, artist, and character list endpoints

### Modified Capabilities

- `image-management`: List endpoint gains `q`, `character_ids`, `artist_ids` filter params
- `artist-management`: List endpoint gains `q` filter param
- `character-management`: List endpoint gains `q` filter param

## Impact

- **Handler layer**: `ListImages`, `ListArtists`, `ListCharacters` — bind filter struct from query params, validate UUIDs, pass to usecase
- **Handler usecase interfaces**: `ImageUsecase`, `ArtistUsecase`, `CharacterUsecase` — `List` signatures updated
- **Usecase layer**: `List` signatures updated; usecases pass filter struct through to repo unchanged
- **Repository interfaces** (`usecase/*_repository.go`): `List` signatures updated
- **Repository layer**: SQL queries gain conditional `ILIKE` and `IN` clauses; `DISTINCT` added to image list when ID filters are present
- **Tests**: handler spies, usecase fakes, and repository integration tests updated across all three entities
- No new dependencies; no breaking changes (all params optional)
