# Proposal

## Why

Artpieces can only be grouped by artist and character today. There is no free-form label for things like "wip", "sketch", or "icon", so the gallery cannot be narrowed by theme or state. Tags add that, and the backend comes first so the frontend can build on a stable API.

## What Changes

- New `tags` table: `id`, `user_id`, `name`, `last_used_at`, `created_at`, `updated_at`, unique on `(user_id, name)`. Names are trimmed and stored lowercase, so tags are case-insensitive. A `CHECK (name = lower(name))` backs this up at the database level.
- New `artpiece_tags` join table (`artpiece_id`, `tag_id`, both `ON DELETE CASCADE`), mirroring `artpiece_characters`. The `tags` table itself has no knowledge of artpieces, so tagging another entity later (e.g. characters) only needs a new join table and filter, not a change to tags.
- New tag endpoints: `POST /tags`, `GET /tags`, `PATCH /tags/:id`, `DELETE /tags/:id`. There is no get-by-id endpoint.
- `GET /tags` accepts an optional `sort=recent` (order by `last_used_at DESC NULLS LAST, name ASC`). The default order is `name ASC`. Any other `sort` value is rejected with 422.
- `last_used_at` follows the artist/character recency pattern: set on create, bumped for every tag in `tag_ids` when an artpiece is created or updated with a non-empty `tag_ids`, not bumped on rename, and not exposed in responses.
- Artpiece create and update accept `tag_ids` (update replaces the whole set, like `character_ids`). A tag the user does not own is rejected with 422.
- Artpiece responses include `tags` (id and name). The list and get-by-id queries preload them.
- `GET /artpieces` accepts `tag_ids` (repeated) and `tag_match=any|all`, defaulting to `any`. `all` returns only artpieces carrying every given tag. An invalid `tag_match` or non-UUID `tag_ids` returns 400, as the existing artpiece list validation does.
- Deleting a tag removes it from all artpieces through the cascade and leaves the artpieces untouched.
- Bruno requests for the four tag endpoints, plus updates to the artpiece create, update and list requests.
- `DATABASE.md` is updated with the new tables.

Out of scope: any frontend work, tagging characters or other entities, a per-tag artpiece count, and account-deletion cleanup of tags (`DeleteAllUserData` is being handled in a later proposal).

## Capabilities

### New Capabilities
- `tag-management`: CRUD for user-owned, case-insensitive tags, the `last_used_at` recency tracking, and the opt-in recency sort on the tag list.

### Modified Capabilities
- `artpiece-management`: artpieces can be assigned tags on create and update, return their tags in the list and get-by-id responses, and the list endpoint can be filtered by tags.

## Impact

- **Database**: migration `000038` creates `tags` and `artpiece_tags`, and its `.down.sql` drops both. It auto-applies on merge to `main`.
- **Domain**: new `domain.Tag`; `domain.Artpiece` gains `Tags []Tag` as `many2many:artpiece_tags`.
- **Usecase**: new `TagRepository`, `tagUsecase` and tag sentinel errors. `ArtpieceUsecase` gains a narrow tag-repository dependency for ownership checks and `UpdateLastUsedAt`, plus `TagIDs` on `CreateArtpieceParams`, `UpdateArtpieceParams` and `ListArtpieceFilters`. Its constructor signature changes, so every call site (`cmd/server/main.go` and the usecase tests) must be updated.
- **Repository**: new `tagRepository`. `artpieceRepository` changes in `Create`, `Update`, `GetByID` and `List` (preload plus the filter join, with `tag_match=all` as a `GROUP BY ... HAVING COUNT(DISTINCT tag_id)`).
- **Handler**: new `TagHandler`; `ArtpieceHandler` request and response types, the list-filter binding, and error mapping change (`ErrTagNotOwned` returns 422).
- **Wiring and docs**: route registration in `cmd/server/main.go`, new `collection/tags/` Bruno files, and `DATABASE.md`.
- **Not affected**: commission responses (the `GetArtpiecesForCommission` query does not preload characters either), and saved-filter payloads, which the backend treats as opaque JSONB.
- **Dependencies**: none added.
