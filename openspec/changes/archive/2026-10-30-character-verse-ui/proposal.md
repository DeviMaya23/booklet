# Proposal

## Why

Characters have a `verse_id` field and the verse API is fully implemented in the backend, but there is no FE to manage verses or assign them to characters. The character list page is a flat unsorted grid — verses as an organisational layer are invisible to the user.

## What Changes

- **BE:** Add `sort_order INT` column to `verses` (migration), update `LIST /verses` to order by `sort_order ASC`, and add a new `POST /verses/reorder` endpoint that accepts an ordered array of verse IDs and resets sort positions in a single transaction.
- **BE:** On `POST /verses`, set `sort_order` to `MAX(sort_order) + 1` for the user so new verses append to the bottom.
- **FE:** New `features/verses/` feature directory with API hooks (`useVerses`, `useCreateVerse`, `useUpdateVerse`, `useDeleteVerse`, `useReorderVerses`).
- **FE:** New "Verses" button on the `CharactersPage` header (left of "New character") that opens a `VersesModal`.
- **FE:** `VersesModal` — a management overlay that lists verses in drag-to-reorder order. Each row shows verse name, character count, an inline rename action (Enter to save / Esc to cancel, duplicate name shown as inline error), and a delete action with inline confirmation showing how many characters will move to Unsorted. Includes "+ New verse" button and an empty state.
- **FE:** `CharactersPage` restructured from a flat grid to a grouped layout: one section per verse (in `sort_order` order), then an "Unsorted" section at the bottom. Verses with no characters are hidden. Search filters characters within groups; empty groups after filtering are also hidden.
- **FE:** `VerseCombobox` added to `CharacterFormModal` above the Notes field, mirroring the `ArtistCombobox` pattern (searchable, clearable). Updates `isDirty` and `resetForm`.
- **FE dependency:** `@dnd-kit/core` and `@dnd-kit/sortable` added for drag-to-reorder in `VersesModal`.

## Capabilities

### New Capabilities

- `verse-management`: Verse entity lifecycle — CRUD, per-user name uniqueness, explicit sort order with a reorder endpoint, and `last_used_at` tracking when a verse is assigned to a character.

### Modified Capabilities

_(none — `character-management` already specifies `verse_id` on character create/update/response)_

## Impact

- **Backend:** New migration `000035_add_sort_order_to_verses`; new route `POST /verses/reorder`; `verse_usecase.Create` updated to compute sort_order; `verse_repository.List` order changed from `last_used_at DESC` to `sort_order ASC`.
- **Frontend:** `CharactersPage` and `CharactersGrid` refactored; `CharacterFormModal` updated; new `features/verses/` directory.
- **FE type:** `Character` interface in `useCharacters.ts` needs `verse_id: string | null` added.
- **New FE dependency:** `@dnd-kit/core`, `@dnd-kit/sortable`.
- **Bruno:** New file for `POST /verses/reorder`.
