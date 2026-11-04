# Design

## Context

See proposal.md for motivation and scope. The closest existing code to copy is:

- **Verses** (`domain/verse.go`, `usecase/verse_*.go`, `repository/verse_repository.go`, `handler/verse_handler.go`) for a user-owned, unique-named entity with CRUD and a `409` on name conflict.
- **Characters** for `last_used_at` (`UpdateLastUsedAt(ctx, ids, userID)`, opt-in `sort=recent`) and for the artpiece association (`many2many:artpiece_characters`, `Association(...).Replace`, a narrow `ArtpieceCharacterRepository` interface, ownership checked in the usecase by comparing counts).
- **`ListArtpieceFilters`** (`usecase/artpiece_repository.go`), bound with `c.Bind` + `c.Validate` in `ListArtpieces`. Validation failures there already map to 400.

Constraints: no `AutoMigrate` (the migration and the domain struct change together), and migrations auto-apply on merge to `main`. The next migration number is `000038`.

## Goals / Non-Goals

**Goals:**
- A tag entity that is independent of artpieces, so tagging characters later needs no change to tags.
- Reuse the verse, character and artpiece-filter patterns as-is. No new layer or abstraction.

**Non-Goals:**
- Frontend, tagging characters, tag usage counts, tag colors or hierarchy.
- Cleaning up tags in `DeleteAllUserData` (separate proposal).

## Decisions

### 1. One join table per tagged entity (`artpiece_tags`), not a polymorphic `taggings` table

`tags` has no reference to what it tags. `artpiece_tags(artpiece_id, tag_id)` has a composite primary key and `ON DELETE CASCADE` on both foreign keys, the same shape as `artpiece_characters`. Characters can later get `character_tags` with a migration, a many2many field and a `tag_ids` filter, with no change to the tag table, API, or usecase.

*Alternative:* a single `taggings(tag_id, taggable_type, taggable_id)`. Rejected because `taggable_id` cannot carry a foreign key or cascade, orphan rows would need manual cleanup, and nothing else in the schema works this way.

### 2. Case-insensitivity by normalizing on write, with a `CHECK` backstop

The usecase applies `strings.TrimSpace` then `strings.ToLower` to the name on create and update. The table has `UNIQUE (user_id, name)` and `CHECK (name = lower(name))`. Filtering and conflict detection need no `LOWER()` calls or functional indexes. `CHECK` constraints already exist in earlier migrations (`000008`, `000009`).

*Alternative:* `CITEXT`. Rejected because it needs an extension and stores the original casing, while the requirement is to store lowercase.

A blank name after trimming returns `ErrTagNameInvalid` (HTTP 422). The request struct also has `validate:"required,min=1"`, which catches an empty or missing name before the usecase runs.

### 3. Tag repository surface

`TagRepository` lives in `usecase/tag_repository.go` (more than 2 methods):

```
Create(ctx, *domain.Tag) (*domain.Tag, error)          // unique violation -> ErrTagNameConflict
List(ctx, userID, ListTagFilters) ([]*domain.Tag, error)
Update(ctx, id, userID, UpdateTagParams) (*domain.Tag, error)   // 0 rows -> gorm.ErrRecordNotFound
Delete(ctx, id, userID) error                                    // 0 rows -> gorm.ErrRecordNotFound
GetByIDsAndUserID(ctx, ids, userID) ([]domain.Tag, error)
UpdateLastUsedAt(ctx, ids, userID) error
```

There is no `GetByID`. `Update` reads the row back with `First` inside the repository, and the artpiece flow validates ownership through `GetByIDsAndUserID`. `ListTagFilters` has `Sort *string` with `query:"sort" validate:"omitempty,oneof=recent"`, and `List` orders by `name ASC` or `last_used_at DESC NULLS LAST, name ASC`. The handler maps a validation failure to 422 with `validationErrResponse`, as the character list does.

### 4. Artpiece integration

- `domain.Artpiece` gets `Tags []Tag` with `gorm:"many2many:artpiece_tags;"`.
- `usecase.ArtpieceTagRepository` is a narrow interface (`GetByIDsAndUserID`, `UpdateLastUsedAt`) beside `ArtpieceCharacterRepository`, and `NewArtpieceUsecase` takes it as a new parameter.
- On create and update with a non-empty `TagIDs`, the usecase de-duplicates the IDs with the existing `dedupUUIDs`. Without that, a repeated ID would fail the count comparison and also violate the join table's primary key. It then loads the tags, returns `ErrTagNotOwned` when the count differs, and calls `UpdateLastUsedAt`. This mirrors the character block, including bumping before the artpiece write.
- `artpieceRepository.Create` sets `Tags` on the model, as it does for `Characters`. `Update` validates tag ownership at the top, before any write, so a bad tag leaves the artpiece untouched, then calls `Association("Tags").Replace(...)`. Update always replaces the set, and an empty or absent `tag_ids` clears it, the same as `character_ids`.
- `GetByID` and `List` add `Preload("Tags")`. Other handler paths that return an artpiece (attach, detach, set cover, replace files) go through `GetByID`, so they pick up tags without further changes.
- `ErrTagNotOwned` maps to 422 in the create and update handlers.

### 5. Tag filter on `GET /artpieces`

`ListArtpieceFilters` gains `TagIDs []uuid.UUID` (`query:"tag_ids"`) and `TagMatch string` (`query:"tag_match" validate:"omitempty,oneof=any all"`). The existing handler already turns validation errors into 400. The repository uses a subquery rather than extending the character join:

```
any:  artpieces.id IN (SELECT artpiece_id FROM artpiece_tags WHERE tag_id IN ?)
all:  artpieces.id IN (SELECT artpiece_id FROM artpiece_tags WHERE tag_id IN ?
                       GROUP BY artpiece_id HAVING COUNT(DISTINCT tag_id) = ?)   -- ? = number of distinct requested tag IDs
```

A subquery does not multiply rows, so it needs no `DISTINCT` and does not interact with the `Distinct()` and join already used for `character_ids`. The artpiece query is already scoped by `artpieces.user_id`, so another user's tag IDs cannot leak artpieces. The `HAVING` count compares against the de-duplicated requested IDs.

The error text in the handler (`character_ids and artist_ids must be valid UUIDs`) is updated to mention `tag_ids` and `tag_match`.

### 6. Recency semantics

`tags.last_used_at timestamptz NULL`. The usecase stamps it with the current time on create, and the artpiece usecase bumps it through `UpdateLastUsedAt`. Rename does not touch it. It is never serialized.

## Risks / Trade-offs

- **Pre-write bump** → The artpiece usecase bumps `last_used_at` before the artpiece write, as it does for characters, so a failed write can still bump a tag. Accepted for consistency, and the effect is a slightly wrong sort order.
- **`CHECK (name = lower(name))` vs Go's `ToLower`** → The two could disagree on exotic Unicode and the insert would fail with a 500 instead of a 409/422. Tag names are expected to be short and mostly ASCII, so the risk is accepted, and a repository integration test covers a non-ASCII name.
- **`user_id` foreign key and account deletion** → `tags.user_id` references `users(id)` without a cascade, like `verses`. `DeleteAllUserData` does not clean up verses or artists today, so tags join that existing gap. The shelved account-deletion proposal must include tags. `artpiece_tags` rows cascade with artpieces.
- **Constructor change** → `NewArtpieceUsecase` gets one more parameter, so `cmd/server/main.go` and the usecase test constructor must be updated together.
- **Association upsert on create** → GORM upserts the associated `Tag` rows (`ON CONFLICT DO NOTHING`) when creating an artpiece, as it does for characters. The artpiece repository integration test for create-with-tags verifies the join rows are written and the existing tag rows are untouched.

## Migration Plan

`000038_create_tags.up.sql` creates `tags` (with the unique constraint, `CHECK`, and foreign key to `users`), then `artpiece_tags` (composite primary key, cascading foreign keys, and an index on `tag_id` for the filter subquery and tag-delete cascade). `000038_create_tags.down.sql` drops `artpiece_tags` then `tags`. The migration is additive and creates no data, so rollback is the down migration. Deploying the new code before the migration would break artpiece reads, but migrations apply on merge to `main`, ahead of the deploy.
