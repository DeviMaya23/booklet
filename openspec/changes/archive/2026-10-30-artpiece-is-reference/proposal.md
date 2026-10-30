# Proposal

## Why

References (style sheets, anatomy guides, mood boards saved as files) are a distinct kind of artpiece that the user wants to separate from commissioned art in the gallery. Without a flag, the only way to separate them is client-side search, which doesn't scale and gives no affordance to the user.

## What Changes

- `artpieces` table gains `is_reference boolean NOT NULL DEFAULT false`. Existing rows become normal artpieces.
- Create and update endpoints accept and return `is_reference`.
- `GET /artpieces` gains an optional `is_reference` query param (tri-state: absent = all, `true` = references only, `false` = non-references only). The gallery does not use this param in v1; the capability is built for future use.
- The gallery filter popover gains a "References only" toggle (default off). Off shows only non-references; on shows only references. All other filters (artist, characters, match mode, sort, search) remain active in both states.
- The Filter button's active-filter summary includes a "references" indicator when the toggle is on.
- "Clear all" in the filter popover also resets the reference toggle.
- The artpiece create modal and artpiece detail edit mode gain a "This is a reference" checkbox (default off).
- Saved filter payload adds `is_reference` (boolean, absent → false). `filterToHref` and filter restoration from URL both carry the new field.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `artpiece-management`: Create and update accept/return `is_reference`; list gains optional `is_reference` filter param.
- `web-artpieces-gallery`: Filter popover gains reference toggle; active-filter summary reflects it; gallery shows only non-references by default.
- `web-artist-character-filter`: Component gains `isReference` / `onIsReferenceChange` props; "Clear all" resets the toggle; trigger-button summary includes "references" when active.
- `web-artpiece-create-modal`: Create form gains "This is a reference" checkbox.
- `web-artpiece-detail`: Edit mode gains "This is a reference" checkbox.
- `saved-filters`: `filter_payload` schema extends to include `is_reference` (boolean).

## Impact

- **DB**: one migration — `ALTER TABLE artpieces ADD COLUMN is_reference boolean NOT NULL DEFAULT false`.
- **Backend**: `domain.Artpiece`, `CreateArtpieceParams`, `UpdateArtpieceParams`, `ListArtpieceFilters`, `createArtpieceRequest`, `updateArtpieceRequest`, `artpieceResponse` — all gain `is_reference`. Repository `List` gains a conditional `WHERE is_reference = ?` clause. Repository `Update` map gains the key when the pointer is non-nil.
- **Frontend types**: `ArtpieceSummary` gains `is_reference: boolean`. `FilterPayload` / `parseFilterPayload` gain `is_reference`. `useCreateSavedFilter` / `useUpdateSavedFilter` param types updated.
- **`useArtpiecesFilter`**: gains `isReference` state (URL param `ref`), filter logic, and returns the value for the page.
- **`ArtistCharacterFilter`**: new optional props `isReference` / `onIsReferenceChange`; `buildSummaryLabel` updated; Clear all handler updated.
- **`ArtpiecesPage`**: wires `isReference` from `useArtpiecesFilter` into `ArtistCharacterFilter`; `SavedFilterSnapshot` and `SaveFilterModal.filterState` include `isReference`.
- **`parseFilterPayload`**: coerces absent/non-boolean `is_reference` to `false`.
- **`filterToHref`** (DashboardPage): appends `?ref=true` when `is_reference` is true.
- **Bruno file**: new/updated request files for create and update artpiece.
