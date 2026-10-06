## Context

Three surfaces independently implement "filter by artist and character":

1. **ArtpiecesPage** — `ArtpiecesFilterPopover`, a hand-rolled absolutely-positioned panel behind a "Filter" button. Includes an All/Any match toggle for characters, caller owns all state via `useArtpiecesFilter`.
2. **AddToExistingArtpieceModal** — an inline collapsible section (always visible on desktop, toggle on mobile). No match toggle. Incorrectly hardcodes character match as ALL — should be ANY.
3. **CommissionFormModal** — an inline "Filters ▷" disclosure inside the artpieces box. No match toggle. Correctly uses ANY.

All three use `ArtistCombobox` + `TokenInput` but differ in container, disclosure pattern, and match semantics.

## Goals / Non-Goals

**Goals:**
- One shared `ArtistCharacterFilter` component used by all three surfaces.
- Consistent panel design matching the new mock (Artist combobox, Characters token input, optional All/Any toggle at 2+ characters, Clear all link).
- Button trigger with active-filter summary label in all three locations, placed beside the search input.
- Fix `AddToExistingArtpieceModal` character match from ALL to ANY.
- Delete `ArtpiecesFilterPopover` once fully replaced.

**Non-Goals:**
- Multi-artist selection (data model constrains each piece/commission to one artist; deferred).
- Backend filter endpoints (all filtering remains client-side).
- Changing filter logic in `useArtpiecesFilter` (the hook is correct; only the UI component changes).

## Decisions

### Component location: `components/ArtistCharacterFilter.tsx`

Used by three distinct features (`artpieces`, `files`, `commissions`), which meets the three-feature threshold for promotion to the shared `components/` layer alongside `TokenInput.tsx`. Staying in `features/artpieces/` would require cross-feature imports which violates the directory convention.

### Controlled component — caller owns all state

`ArtistCharacterFilter` accepts `artist / onArtistChange`, `characters / onCharactersChange`, and optionally `characterMatch / onCharacterMatchChange`. It owns only the open/closed state of the popover panel internally. This mirrors how `ArtistCombobox` and `TokenInput` already work, and keeps the filter logic in the callers (where the filtered data lives).

Callers that don't want the match toggle simply omit `characterMatch` and `onCharacterMatchChange`. The toggle is not rendered, and the caller handles its own match semantics (defaulting to ANY in both modal cases).

### Popover mechanism: Base UI `Popover` primitive

The current `ArtpiecesFilterPopover` uses a hand-rolled `mousedown` listener for outside-click dismissal. The project already uses Base UI `Popover` / `PopoverTrigger` / `PopoverContent` for the datepicker fields in `CommissionFormModal`. Using the same primitive gives correct keyboard behaviour (Escape to close, focus management) for free and keeps the implementation consistent.

Note: Base UI's `PopoverTrigger` does not support `asChild` — the trigger must be styled directly via `className` rather than wrapping a `Button` component. This is already established practice in this codebase.

### Match toggle visibility: 2+ characters (not 1+)

The current `ArtpiecesFilterPopover` shows the All/Any toggle as soon as one character is selected. The new design shows it only at 2+, which is semantically correct — any/all is meaningless with a single character. Callers using the toggle (`ArtpiecesPage`) see this change in behaviour; modal callers are unaffected (no toggle).

### Summary label on trigger button

When filters are active, the trigger shows a compact summary instead of "Filter": `1 artist`, `2 characters`, `1 artist · 2 characters`. When `characterMatch` is in use, append the mode: `2 characters, any`. This matches the current `ArtpiecesFilterPopover` label logic and is preserved in the new component.

### `ArtpiecesFilterPopover` is deleted

The new component fully supersedes it. Its test (`ArtpiecesFilterPopover.test.tsx`) is also deleted. `ArtpiecesPage` imports `ArtistCharacterFilter` directly.

## Risks / Trade-offs

- **Base UI `PopoverTrigger` className constraint** — styling the trigger button directly (not via `asChild`) means the trigger's visual appearance must be built with raw className rather than the `Button` component. The result looks identical but bypasses `Button`'s variant system. This is a known limitation already accepted elsewhere in the codebase.
- **`AddToExistingArtpieceModal` behaviour change** — switching from inline-always-visible-on-desktop to a popover-behind-a-button removes the desktop affordance of filters always being visible. This is a net improvement in consistency but is a visible layout change.
- **Character match fix in `AddToExistingArtpieceModal`** — correcting ALL → ANY changes which artpieces appear in results when multiple characters are selected. This is a bug fix, not a regression.
