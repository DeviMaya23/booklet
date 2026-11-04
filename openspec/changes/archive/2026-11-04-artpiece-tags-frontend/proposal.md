# Proposal

## Why

The tags backend is merged (tag CRUD, `tag_ids` on artpiece create/update, `tags` on artpiece responses, `last_used_at` recency), but the frontend cannot assign, display or filter by tags. Until it can, tags are unusable. Wiring them in now also forces a decision on how saved filters treat ids that no longer exist, and that decision needs to cover characters too, because the thumbnail modal currently mishandles deleted characters.

## What Changes

- **Tag picker**: a new `TagPicker` chip input (chips + recency-ordered dropdown) that also supports inline tag creation. Typing a name that matches no existing tag offers a "Create 'name'" row; confirming it calls `POST /tags` and selects the result. A name matching an existing tag (ignoring case) selects the existing tag instead of creating. The picker requests tags with `sort=recent`.
- **Artpiece create** (`ArtpieceFormModal`, opened from the artpieces page and from the inbox "New artpiece" flows): a Tags field using `TagPicker`; `tag_ids` is sent on `POST /artpieces`.
- **Artpiece edit** (`ArtpieceDetailView` edit mode): a Tags field using `TagPicker`; `tag_ids` is sent on `PUT /artpieces/:id` (full replace); tag changes count toward the unsaved-changes caption and route blocker. View mode shows tags as chips.
- **Artpiece data**: `ArtpieceSummary` and `ArtpieceDetail` gain `tags: {id, name}[]`.
- **Gallery tag filter**: a Tags section in the filter popover with its own Any/All toggle (default `all`, visible at 2+ tags), selecting existing tags only (no creation). Filtering is client-side over the fetched list, like characters. State is persisted in the URL (`tags`, `tag_match`) and `localStorage`, shown in the trigger summary label, cleared by "Clear all", and AND-combined with the other filters. The popover shows no per-tag counts in this change.
- **Saved filters**: `filter_payload` gains `tag_ids` (string[]) and `tag_match_mode` (`all`|`any`); both are written by the Save Filter modal, included in its summary chip, read by `parseFilterPayload`, and mapped to URL params by `filterToHref`. Absent fields read as `[]` / `all`, so existing saved filters keep working. No backend change: the payload is stored opaquely.
- **Deleted ids are a no-op**: a tag id in the URL, `localStorage` or a saved filter that no longer resolves against the current tag list is dropped from the active filter, not an error and not a "match nothing" filter. The same rule is stated for character ids.
- **Fix (characters)**: `ChangeThumbnailModal` currently passes raw saved-filter character ids with no resolution, so a deleted character in `all` mode empties the candidate grid. It resolves character ids against the current character list, as it already does for the verse.
- Mutations that assign tags invalidate the tags query so picker order reflects `last_used_at`, as characters do.

Out of scope: tag management (rename/delete UI), per-tag counts, a tag filter on `ArtistCharacterFilter` callers other than the gallery, and any backend change.

## Capabilities

### New Capabilities
- `web-tag-picker`: the `TagPicker` component and its data hooks: recency-ordered suggestions, chip selection and removal, and inline creation (including existing-name handling and creation failure).

### Modified Capabilities
- `web-artpiece-create-modal`: adds the Tags field and `tag_ids` on save.
- `web-artpiece-detail`: adds tag chips in view mode, the Tags field in edit mode, `tag_ids` on save, and tags in the dirty check.
- `web-artpieces-gallery`: adds the tag filter (client-side, any/all, URL and `localStorage` persistence, summary, save-to-dashboard eligibility) and the rule that deleted tag and character ids resolve to "not selected".
- `web-artist-character-filter`: adds optional tag props (token input, any/all toggle, summary label, "Clear all"), absent for existing callers.
- `web-save-filter-modal`: summary chip and saved payload include tags.
- `web-dashboard-saved-filters`: click navigation maps tags to URL params; the Change thumbnail modal applies the tag filter and resolves tag and character ids against current data.
- `saved-filters`: `filter_payload` supports `tag_ids` and `tag_match_mode` with defaults when absent (same pattern as `verse_id`).

## Impact

- **New**: `features/tags/` (`api/useTags`, `api/useCreateTag`, `components/TagPicker`).
- **`components/TokenInput.tsx`**: additive optional props only (async create, create row, immediate dropdown); existing callers (`ArtistCharacterFilter`, `ArtpieceFormModal`, `ArtpieceDetailView`, `FolderPicker`) must be unchanged.
- **Create/edit surfaces**: `ArtpieceFormModal` (used by `ArtpiecesPage` and `FilesPage`), `ArtpieceDetailView`. `useCreateArtpiece`, `useUpdateArtpiece`.
- **Types**: `useArtpieces.ts`, `useArtpiece.ts`.
- **Filter**: `useArtpiecesFilter` (`StoredFilter`, `FilterParams`, `filterAndSortArtpieces`, URL commit and restore), `ArtistCharacterFilter` (gallery is the only caller that passes tag props; `CommissionFormModal` and `AddToExistingArtpieceModal` are unchanged), `ArtpiecesPage` (`SavedFilterSnapshot`, `snapshotMatchesCurrent`, change handlers, `hasActiveFilter`).
- **Saved filter shape** (all copies): `parseFilterPayload`, `useCreateSavedFilter` (params and `SavedFilter` types), `SaveFilterModal`, `filterToHref`, `ChangeThumbnailModal` (filter and chip label), `DashboardPage` (consumes via `filterToHref`).
- **APIs consumed**: `GET /tags?sort=recent`, `POST /tags`, `tag_ids` on `POST /artpieces` and `PUT /artpieces/:id`, `tags` on artpiece list and detail responses. The backend `tag_ids`/`tag_match` list query params are not used (filtering stays client-side).
- **Tests**: component tests for the picker, both forms, the filter, the save modal and the thumbnail modal; hook and lib tests for `useArtpiecesFilter`, `parseFilterPayload`, `filterToHref`.
- **No** backend, migration or Bruno changes.
