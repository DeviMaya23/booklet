# Design

## Context

`ArtpieceDetailView` edit mode renders `<ArtistCombobox value={editArtist} onChange={setEditArtist} />` bare, inside a `grid grid-cols-2` row next to Characters. `ArtistCombobox` has no add affordance. The "+" button is defined per consumer: `ArtpieceFormModal` and `CommissionFormModal` each hold `artistModalOpen` state, wrap the combobox in a `flex gap-2` row with an outline icon `Button`, and render `<ArtistFormModal onCreated={...} />`. See proposal.md for motivation.

Relevant existing behavior:
- `editArtist` is seeded once in `handleEnterEdit` from the cached artists list, so a refetch of artists does not reset it.
- The unsaved-changes caption already compares `editArtist?.id` with the original artist id, so a newly created artist is reported as "artist changed" without extra code.
- `useCreateArtist` invalidates `ARTISTS_QUERY_KEY`, so the new artist shows up in the combobox list.

## Goals / Non-Goals

**Goals:**
- Same "+" pathway in edit mode as in the other two forms, with the same look and behavior.

**Non-Goals:**
- Extracting a shared `ArtistPicker` or adding an `onAddNew` prop to `ArtistCombobox`. Considered and deferred as a separate refactor.
- Changing `ArtistFormModal`, `ArtistCombobox`, or any backend code.

## Decisions

**Copy the existing local pattern into `ArtpieceDetailView`.**
Add `artistModalOpen` state; wrap the combobox in `<div className="flex gap-2">` with a `flex-1` child and an outline icon `Button` (`type="button"`, `<Plus />`, `disabled={isSaving}`); render `<ArtistFormModal open onOpenChange onCreated={(artist) => setEditArtist(artist)} />` alongside the other dialogs.
- *Why:* smallest change, no new abstraction (CLAUDE.md requires approval for those), and it keeps edit mode identical to the other forms.
- *Alternative, shared wrapper component:* removes the duplication across three call sites but touches two working modals and their tests for no user-visible gain. Better as its own change.
- *Alternative, `onAddNew` prop on `ArtistCombobox`:* parents still own the modal and state, so it saves little.

**Selecting the created artist only updates local edit state.**
`onCreated` calls `setEditArtist`; nothing is persisted to the artpiece until Save. This matches the other forms and the spec's "pending change" scenario.

**Layout stays in the existing two-column row.**
The artist column gets the flex row inside it, so the Characters column and grid are untouched.

## Risks / Trade-offs

- [Third copy of the same boilerplate] → Accepted; a shared picker is a candidate follow-up refactor.
- [Narrow artist column leaves less room for the combobox next to the button] → The button is icon-sized; verify visually at the layout's usual widths during implementation.
- [`ArtistFormModal` is a dialog rendered while the detail view's unsaved-change navigation blocker is active] → Opening it is not a navigation, so the blocker is unaffected; confirm in a test that the edit form values survive opening and closing the modal.
