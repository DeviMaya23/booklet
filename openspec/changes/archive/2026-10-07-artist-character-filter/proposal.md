## Why

Three places in the UI independently implement "filter by artist and character" — the artpieces gallery, the Add Files to Artpiece modal, and the Commission form modal — each with slightly different behaviour and presentation. Consolidating them into one shared component eliminates drift risk, normalises behaviour, and makes a future change to either field (e.g. multi-artist support) a single-site update.

## What Changes

- Extract a new shared `ArtistCharacterFilter` component at `components/ArtistCharacterFilter.tsx`, replacing all three existing filter implementations.
- The component renders a trigger button (with active-filter summary label) that opens a floating panel containing: an Artist combobox (single-select), a Characters token input (multi-select), an optional All/Any match toggle (shown only when 2+ characters are selected and the caller opts in), and a "Clear all" link.
- Replace the current `ArtpiecesFilterPopover` in `ArtpiecesPage` with the new component (behaviour unchanged, match toggle retained).
- Replace the inline collapsible filter disclosure in `CommissionFormModal` with the new component (no match toggle; trigger button placed beside the artpiece search input).
- Replace the inline collapsible filter in `AddToExistingArtpieceModal` with the new component (no match toggle; trigger button beside the search input). Fix the character match logic, which currently hardcodes ALL — it should be ANY.
- Delete `ArtpiecesFilterPopover` once replaced.

## Capabilities

### New Capabilities

- `web-artist-character-filter`: Shared filter component — trigger button, floating panel, Artist combobox, Characters token input, optional All/Any toggle, Clear all.

### Modified Capabilities

- `web-artpieces-gallery`: Filter trigger moves to a button beside the search field; panel design updates to match the new component.
- `web-add-files-to-artpiece`: Inline filter disclosure replaced with button + panel; character match corrected from ALL to ANY.
- `web-commissions-list`: Commission form modal artpiece filter replaced with button + panel.

## Impact

- **New file**: `frontend/src/components/ArtistCharacterFilter.tsx`
- **Deleted**: `frontend/src/features/artpieces/components/ArtpiecesFilterPopover.tsx` (and its test)
- **Modified**: `frontend/src/pages/ArtpiecesPage.tsx`, `frontend/src/features/files/components/AddToExistingArtpieceModal.tsx`, `frontend/src/features/commissions/components/CommissionFormModal.tsx`
- **No backend changes**
- **No API changes**
