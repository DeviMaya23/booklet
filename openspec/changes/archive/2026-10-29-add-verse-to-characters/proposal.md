# Proposal

## Why

Characters need a way to group them into narrative contexts (universes, settings, canons). Without this, large character libraries have no way to scope browsing or filtering to a single story-world. Verse is a first-class entity — not just a tag — because it has recency tracking and its own CRUD surface, making it available as a picker in the future.

## What Changes

- New `verses` table: `id`, `user_id`, `name`, `last_used_at`, `created_at`, `updated_at`. Unique per user by name.
- New `verse_id` nullable FK column on `characters`, `ON DELETE SET NULL`.
- New CRUD endpoints: `GET /verses`, `POST /verses`, `PATCH /verses/:id`, `DELETE /verses/:id`.
- `GET /characters` gains an optional `verse_id` query param filter.
- `POST /characters` and `PUT /characters/:id` accept an optional `verse_id` field. When non-null, the character usecase validates verse ownership and bumps `verses.last_used_at`.
- Character API responses include `verse_id` (UUID string, nullable).
- `characterUsecase` gains a `VerseRepository` dependency for ownership validation and `last_used_at` bumping.

All FE integration (verse list UI, picker widget, character form wiring) is explicitly out of scope for this change.

## Capabilities

### New Capabilities

- `verse-management`: Full CRUD for verses — create, list, get-by-id, update name, delete. Owns the `verses` schema, the ownership rules, `last_used_at` tracking behavior, unique-name constraint, and `ON DELETE SET NULL` cascade.

### Modified Capabilities

- `character-management`: Adds `verse_id` as an optional assignable field on character create and update. Adds verse ownership validation in the usecase. Adds `verse_id` to the character response. Adds `verse_id` query filter on `GET /characters`.

## Impact

- **New migration**: `000033_create_verses` — creates `verses` table.
- **New migration**: `000034_add_verse_id_to_characters` — adds `verse_id` FK column to `characters`.
- **Backend**: new `domain.Verse`, `verse_repository`, `verse_usecase`, `verse_handler` files. `domain.Character` gains `VerseID *uuid.UUID`. `characterUsecase` constructor gains `verseRepo`. `UpdateCharacterParams` and `CreateCharacterParams` gain `VerseID *uuid.UUID`.
- **Bruno collection**: new `collection/verses/` folder with CRUD requests.
- **No FE changes** in this pass.
