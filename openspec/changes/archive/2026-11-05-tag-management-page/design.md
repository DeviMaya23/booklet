# Design

## Context

See `proposal.md` for motivation. Relevant current state:

- Backend tags already have `Create`, `List` (with `sort=recent`), `Update` (rename; 409 on name collision via the unique violation) and `Delete`. `artpiece_tags(artpiece_id, tag_id)` has a composite primary key and `ON DELETE CASCADE` on both foreign keys, so deleting a tag detaches it everywhere and the same artpiece cannot hold a tag twice.
- `tagRepository.List` returns `[]*domain.Tag` and orders by unqualified `name` / `last_used_at`. `TagUsecase` and `TagRepository` interfaces return `domain.Tag`. `GET /tags` is consumed by `TagPicker`, `ArtistCharacterFilter` and `SaveFilterModal` through `useTags`, none of which want counts.
- Artpieces have no soft-delete column, so a count of `artpiece_tags` rows per tag is the count of artpieces. A tag can only be assigned to an artpiece of the same user (`ErrTagNotOwned` on assignment), so no per-user join through `artpieces` is needed.
- Multi-step writes use `Transactor.InTransaction` from the usecase, with repositories reading the transaction via `dbFromContext`.
- Frontend: `features/tags/api` has `useTags` and `useCreateTag` (which invalidates `TAGS_QUERY_KEY`, a prefix of every tag query key). Feature-owned confirmation dialogs exist (`DeleteCommissionDialog`, `DeleteArtpieceDialog`, built on `components/ui/alert-dialog`). `ArtistsPage` is the model for a search-box list page. The sidebar items are a static array in `AppSidebar.tsx`.

## Goals / Non-Goals

**Goals:**
- A Tags page with search, sort, usage count, inline rename, delete and merge, per the mockup.
- Pay the cost of counting only for the caller that asks.
- Make merge atomic and impossible to trigger by accident through plain rename.

**Non-Goals:**
- No new sort value, and no sort-by-usage.
- No server-side search, pagination or bulk actions on tags.
- No undo for merge or delete (the UI says so).
- No change to `PATCH /tags/:id` semantics, `last_used_at` rules, or the tag picker.
- No change to the Bookleaf integration.

## Decisions

### D1: Usage count is opt-in via `with_usage`, served by a separate list path
`ListTagFilters` gains `WithUsage bool` (`query:"with_usage"`), bound with `c.Bind` + `c.Validate`, so a non-boolean value is a bind error (400). The handler branches: without the flag it calls the existing `List` and builds `{id, name}` exactly as today; with it, it calls a new `ListWithUsage` on the usecase and repository returning `[]*domain.TagWithUsage`.

Alternatives considered:
- *Change `List` to always return the new type*: would force the counting join on every caller, which the user explicitly wants to avoid.
- *Make `List` return a union / optional count on `domain.Tag`*: puts a non-column field on a pure table model (rejected by the user in favour of a new struct).
- *A separate endpoint (`/tags/usage`)*: duplicates sort handling and scoping. A flag on the same resource keeps one list contract.

### D2: `domain.TagWithUsage` carries the count
```
TagWithUsage { Tag (embedded); ArtpieceCount int }
```
It is not a gorm model with its own table; the repository scans a grouped query into it. `domain.Tag` is untouched.

The usage query is `tags LEFT JOIN artpiece_tags ON artpiece_tags.tag_id = tags.id`, filtered by `tags.user_id`, `GROUP BY tags.id`, selecting `COUNT(artpiece_tags.artpiece_id)` (counting the joined column, not `*`, so an unused tag yields 0 rather than 1). Ordering columns are qualified (`tags.name`, `tags.last_used_at`) so the same two orderings work with or without the join. The ordering strings are shared between the two repository methods rather than copied.

### D3: Response shape: `artpiece_count` present only when requested
The list response item for the usage path is a separate struct whose `artpiece_count` is a plain `int`; the non-usage item stays the existing `tagResponse`. Both serialise into the same JSON array shape; the handler picks one. This keeps the "absent unless requested" property by construction instead of via pointers and `omitempty`, and guarantees `0` is emitted for unused tags.

### D4: Merge is a dedicated endpoint, orchestrated in the usecase
`POST /tags/:id/merge` with `{target_id}`. The handler parses `:id` with `uuid.Parse` (400), binds and validates the body (`target_id` `required`; invalid JSON or non-UUID is a bind error → 400, absent → 422). The usecase:

1. rejects `source == target` with a new `ErrTagMergeSameTag` (→ 422);
2. inside `Transactor.InTransaction`:
   a. `GetByIDsAndUserID([source, target], user)`; fewer than two rows → `gorm.ErrRecordNotFound` (→ 404, the same mapping as update and delete, and it does not reveal whether the other user owns the id);
   b. new repository method `ReassignArtpieces(from, to)`: `INSERT INTO artpiece_tags (artpiece_id, tag_id) SELECT artpiece_id, :to FROM artpiece_tags WHERE tag_id = :from ON CONFLICT DO NOTHING`;
   c. existing `Delete(source)`: the cascade removes the source's remaining join rows.

Only a new method on the repository interface is added; ownership and delete reuse existing methods. The target's `last_used_at` is never written. The handler responds 204.

Alternatives considered:
- *`PATCH` merges when the name collides (with a flag)*: overloads rename, contradicts the spec'd 409, and a stale client could destroy data by renaming.
- *One repository `Merge` method doing everything with its own transaction*: bypasses the `Transactor` convention and cannot be unit tested without a database; the usecase-orchestrated version keeps unit-testable logic (same-tag rule, ordering of calls) and leaves SQL to an integration-tested repository method.
- *`UPDATE artpiece_tags SET tag_id = :to WHERE tag_id = :from`*: violates the primary key for artpieces that have both tags.

### D5: Frontend collision detection and merge routing
Collision is decided client-side against the loaded list (the page always has the full list), by a small pure helper `findTagByName(tags, name, excludeId)` using trimmed, case-insensitive comparison, in `features/tags/lib/` with its own test. A submit that matches another tag opens the merge dialog; otherwise `PATCH` is called. A case-only change to the tag's own name (`wip` → `WIP`) excludes the tag itself and goes through `PATCH`, matching the backend's allowance.

A rename that races with a collision created elsewhere gets 409 from `PATCH`; the row stays in edit mode and shows an error (spec: "Rename failure"). No automatic fallback to merge.

### D6: `useTags` gets a `withUsage` option; the query key reflects it
`UseTagsOptions.withUsage?: boolean`. The key becomes `[...tagsQueryKey(sort), 'usage']` when set and is unchanged otherwise, so existing cache entries and tests are unaffected, and a no-usage list can never satisfy a usage request. All keys still start with `TAGS_QUERY_KEY`, so the existing prefix invalidation in mutations covers them. `Tag` gains an optional `artpiece_count?: number`; the Tags page treats a missing value as 0 only defensively, since it always requests usage.

Alternative considered: a separate `useTagsWithUsage` hook with its own exported type. Cleaner typing, but duplicates the fetch/key logic; the option keeps `useTags` the single reader of `GET /tags`.

### D7: Frontend structure
Following ARCHITECTURE.md and CONVENTIONS.md (feature-owned, matching existing shapes, tests written with each unit):

```
pages/TagsPage.tsx                       thin: search, sort, list, owns editingId + dialog targets
features/tags/
  api/useUpdateTag.ts                    PATCH, invalidates TAGS_QUERY_KEY
  api/useDeleteTag.ts                    DELETE, invalidates TAGS_QUERY_KEY
  api/useMergeTag.ts                     POST merge, invalidates TAGS_QUERY_KEY
  components/TagsList.tsx                rows + empty state
  components/TagRow.tsx                  one row: display <-> inline edit (draft state local)
  components/DeleteTagDialog.tsx         { item: Tag | null, onCancel, onConfirm }
  components/MergeTagDialog.tsx          { item: {source, target} | null, onCancel, onConfirm }
  lib/findTagByName.ts                   pure collision helper
```
Dialogs are two near-identical feature-owned `AlertDialog` wrappers matching `DeleteCommissionDialog`; no new shared component. Artpiece responses embed `tags`, so rename, merge and delete also invalidate `ARTPIECES_QUERY_KEY` (and the per-artpiece detail queries), the same way `useUpdateArtpiece` invalidates its related keys; otherwise an open gallery or detail page would keep showing a tag that no longer exists or has its old name.

### D8: Page-level orchestration
`TagsPage` holds `search`, `sort` (`'recent' | 'name'`, default `'recent'`), `editingId`, `deleteTarget`, `mergeTarget`. Sorting is the API's job: it requests `useTags({ sort: sort === 'recent' ? 'recent' : undefined, withUsage: true })`; search filters the result client-side (as `ArtistsPage` does). Changing sort changes the query key and fetches the other ordering.

## Risks / Trade-offs

- **Count query cost grows with `artpiece_tags` size** → only paid on the Tags page; grouped by indexed primary key / FK columns, and the library is single-user scale. Integration tests cover correctness, not performance.
- **Merge is irreversible** → confirmation dialog states consequences; the operation is a single transaction so a failure never leaves a half-merged state.
- **Client-side collision check can be stale (two tabs)** → worst cases are benign: a missed collision yields a 409 error on rename; a spurious collision (target deleted elsewhere) yields a 404 from merge, surfaced as an error and a list refresh.
- **Optional `artpiece_count` on the shared `Tag` type** is weaker than a distinct type → accepted to avoid a second fetch path (D6); only the Tags page reads it.
- **Qualifying ORDER BY columns** changes the existing `List` query text → covered by the existing tag repository integration tests, which must stay green.

## Migration Plan

No schema migration and no data backfill: `artpiece_tags` already has the needed primary key and cascades. Deploy backend and frontend together; the new query param and endpoint are additive, and old clients are unaffected. Rollback is a plain revert.

## Open Questions

None blocking. (The merge dialog's artpiece count wording is a copy detail resolved in implementation from the loaded `artpiece_count`.)
