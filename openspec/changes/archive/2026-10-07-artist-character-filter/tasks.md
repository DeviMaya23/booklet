## 1. ArtistCharacterFilter Component

- [x] 1.1 Create `frontend/src/components/ArtistCharacterFilter.tsx` — controlled component with props `artist`, `onArtistChange`, `characters`, `onCharactersChange`, optional `characterMatch`, `onCharacterMatchChange`; owns open/closed panel state internally
- [x] 1.2 Implement trigger button — shows "Filter" when no filters active; shows summary label ("1 artist", "N characters", "1 artist · N characters, any/all") when filters are active
- [x] 1.3 Implement floating panel using Base UI `Popover` / `PopoverTrigger` / `PopoverContent` — opens on trigger click, closes on outside click or Escape
- [x] 1.4 Implement Artist field — `ArtistCombobox` single-select inside the panel, calls `onArtistChange` on select/clear
- [x] 1.5 Implement Characters field — `TokenInput` multi-select inside the panel, calls `onCharactersChange` on add/remove; load character suggestions via `useCharacters()`
- [x] 1.6 Implement All/Any toggle — render only when `characterMatch` prop is present AND 2+ characters are selected; clicking calls `onCharacterMatchChange`
- [x] 1.7 Implement Clear all — link at bottom of panel; calls `onArtistChange(null)` and `onCharactersChange([])` together
- [x] 1.8 Write `frontend/src/components/ArtistCharacterFilter.test.tsx` — cover: summary label with no filters, summary label with artist only, summary label with characters + match mode, toggle hidden at <2 characters, toggle visible at 2+ characters, toggle hidden when prop omitted, Clear all fires both handlers

## 2. Replace ArtpiecesFilterPopover in Gallery

- [x] 2.1 In `frontend/src/pages/ArtpiecesPage.tsx`, replace `<ArtpiecesFilterPopover>` with `<ArtistCharacterFilter>` — pass `selectedArtist`, `onArtistChange`, `selectedCharacters`, `onCharactersChange`, `characterMatch`, `onCharacterMatchChange` from `useArtpiecesFilter`
- [x] 2.2 Delete `frontend/src/features/artpieces/components/ArtpiecesFilterPopover.tsx` and `ArtpiecesFilterPopover.test.tsx`

## 3. Replace Inline Filter in AddToExistingArtpieceModal

- [x] 3.1 In `frontend/src/features/files/components/AddToExistingArtpieceModal.tsx`, replace the inline collapsible filter section with `<ArtistCharacterFilter>` placed beside the artpiece search input; remove `filtersOpen` state and the collapsible toggle button
- [x] 3.2 Fix character match logic from ALL (`every`) to ANY (`some`) in `filteredArtpieces`

## 4. Replace Inline Filter Disclosure in CommissionFormModal

- [x] 4.1 In `frontend/src/features/commissions/components/CommissionFormModal.tsx`, replace the "Filters ▷" inline disclosure with `<ArtistCharacterFilter>` placed beside the artpiece search input; remove `filtersOpen` state and the disclosure button

## 5. Build and Lint

- [x] 5.1 Run `npm run build` in `frontend/` and fix any TypeScript or build errors
- [x] 5.2 Run `npm run lint` in `frontend/` and fix any lint errors
