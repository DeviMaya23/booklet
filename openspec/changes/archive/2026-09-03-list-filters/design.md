## Context

The three list endpoints (`GET /images`, `GET /artists`, `GET /characters`) accept no query parameters today and return all records scoped to the authenticated user. Each flows through a uniform handler → usecase → repository stack with `List(ctx, userID)` at every layer.

The frontend needs freetext search and ID-based filtering before building list views. All filters are optional — the no-param case must continue to behave identically to today.

## Goals / Non-Goals

**Goals:**
- Add `q` (case-insensitive substring) to the artist, character, and image list endpoints
- Add `character_ids` and `artist_ids` (multi-value UUID IN filters) to the image list endpoint
- Thread a typed filter struct through all three layers without a mapping step

**Non-Goals:**
- Pagination
- Sorting control (ordering stays as-is)
- Full-text search indexing
- Filtering on any field beyond what is specified above

## Decisions

### One filter struct per entity, defined in the usecase package, with `query` tags

Each entity gets one filter struct (e.g. `ListImageFilters`) defined in `usecase/*_repository.go` alongside the existing params structs (`UpdateImageParams`, etc.). The struct carries `query` struct tags so `c.Bind` works directly in the handler — no separate handler-side DTO, no mapping step.

The `query` tag is declarative metadata that brings no Echo imports into the usecase package. The codebase already accepts this tradeoff with GORM tags on domain types.

Alternative considered: separate handler-side query struct + usecase filter struct with a mapping step in the handler. Rejected — the usecase is a pure passthrough here, so the mapping adds ceremony with no benefit.

### Handler binds and validates; usecase passes through; repo applies conditions

- Handler: `c.Bind(&filters)` populates the struct from query params. UUID fields (`character_ids`, `artist_ids`) arrive as `[]string` via bind then are parsed to `[]uuid.UUID` inline; a malformed UUID returns 400.
- Usecase: receives the filter struct, forwards it to the repo unchanged. No logic added.
- Repository: appends WHERE conditions conditionally — only when the filter field is non-nil/non-empty.

### Repository builds WHERE clauses conditionally with GORM chaining

```
q != nil      → .Where("LOWER(title) LIKE ?", "%" + lower + "%")
character_ids → .Joins(...).Where("image_characters.character_id IN ?", ids)
artist_ids    → .Where("artist_id IN ?", ids)
```

When `character_ids` is non-empty, `.Distinct()` is added to avoid duplicate rows when an image is tagged with multiple matching characters.

Alternative for ID filters: subquery. Rejected — the existing `ListByCharacterID` uses a JOIN; staying consistent is simpler.

### Multi-value query params via repeated keys

`?character_ids=a&character_ids=b` — Echo's `c.Bind` populates a `[]string` slice field automatically with repeated keys. No custom parsing needed.

Alternative: comma-separated single param. Rejected — repeated keys are standard REST convention and what frontend query libraries produce.

## Risks / Trade-offs

- **ILIKE performance** → Acceptable at current scale. A `pg_trgm` index can be added later without API changes.
- **Duplicate rows from JOIN on character_ids** → Mitigated by `.Distinct()` in the repo when the filter is active.
- **`query` tags in usecase package** → Minor SoC leak; accepted, consistent with GORM tags on domain types.
- **Signature change ripples to all test doubles** → Handler spies and usecase fakes implementing `List` must be updated across all three entities. Mechanical but touches multiple files.
