## Context

The codebase has two coexisting patterns for UUID handling. The artpiece/upload/file-upload stack (newer) uses `uuid.UUID` end-to-end: path params are parsed inline with `uuid.Parse`, body fields use `uuid.UUID` types, and usecases/repositories take `uuid.UUID`. The artist/character stack (older) validates with `validate:"uuid4"` tags and then re-parses at every boundary, flowing the raw `string` all the way down to the repository. This creates unnecessary parse calls, type-validation tags whose only job is to guard a subsequent `uuid.Parse`, and inconsistent signatures across the codebase.

Echo's `DefaultBinder` handles `uuid.UUID` natively via `encoding.TextUnmarshaler` for path params (`param` tag), query params (`query` tag), and JSON bodies. This means switching body fields and filter slices to `uuid.UUID` requires no custom binder work.

## Goals / Non-Goals

**Goals:**
- Eliminate all `validate:"uuid4"` tags and post-bind `uuid.Parse` calls from artist and character handler request structs
- Promote `id string` to `id uuid.UUID` across artist and character usecase interfaces and repository implementations
- Merge `artist_repository.GetByIDAndUserID` into `GetByID` (signatures become identical post-migration)
- Drop `::text` SQL casts in `image_repository.List` filter queries (incidental fix, same pass)
- Collapse `parseFolderIDs` / `parseFolderIDSlice` helpers in `character_handler.go` (become identity functions once `FolderIDs` fields are `[]uuid.UUID`)

**Non-Goals:**
- Image handler/usecase/repository stack — scheduled for deletion with the artpiece+files migration
- Any changes to error response format or HTTP status semantics beyond the uuid4→bind-error shift noted in the proposal

## Decisions

### Use inline `uuid.Parse` for path params, not struct binding

Path params on GET/DELETE endpoints have no request body. Putting a single `ID uuid.UUID \`param:"id"\`` into a struct just for `c.Bind` adds ceremony with no gain. The artpiece handler already establishes the clean precedent:

```go
id, err := uuid.Parse(c.Param("id"))
if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid X id")
}
```

For PUT/PATCH handlers that have both a path param and a body, the path param is still handled inline (same pattern) and the body struct switches to `uuid.UUID` typed fields.

### Accept HTTP 400 (no field info) for invalid UUID body fields

Currently: `validate:"uuid4"` tag → `c.Validate` → HTTP 422 with `{"errors": [{"field": "artist_id", "message": "..."}]}`.

After: `*uuid.UUID` field → `c.Bind` fails → HTTP 400 `{"message": "invalid request body"}`.

This is acceptable because UUIDs in request bodies are always app-issued IDs constructed by the frontend, not values typed by users. A malformed UUID in the body is a client-side bug, not a user input error. The 400 is sufficient signal.

### Merge `GetByIDAndUserID` into `GetByID` on artist repository

`GetByIDAndUserID(ctx, id uuid.UUID, userID uuid.UUID)` exists solely because the artpiece usecase needed to validate artist ownership and required `uuid.UUID`, while `GetByID(ctx, id string, userID uuid.UUID)` still took a string. After migration they are semantically and signature-identical. The `GetByIDAndUserID` call sites in `artpiece_usecase.go` are updated to use `GetByID`.

## Risks / Trade-offs

- **422 → 400 status shift for UUID body fields** → Acceptable per proposal rationale. No frontend code branches on this distinction for UUID-specific errors.
- **`uuid.MustParse` in `parseFolderIDs` / `parseFolderIDSlice`** → These currently panic on malformed input because the strings were pre-validated by the `validate:"uuid4"` tag. After the change, `FolderIDs` is `[]uuid.UUID` so these helpers are deleted entirely — no parse call needed, no panic risk.
- **`image_repository` `::text` cast removal** → Low risk; GORM's pgx driver passes `uuid.UUID` values as the native Postgres UUID type, which compares directly against UUID columns without casting.

## Migration Plan

No data migration required. This is a pure Go type change with no schema or API shape changes (request/response field names are unchanged). Deploy is a standard rollout.

Rollback: revert the PR; no state is affected.
