## Why

The commissions list page has a non-functional "+ New Commission" button, and there is no way to create, edit, or delete a commission from the UI. This change wires up the full create/edit/delete flow via a modal, closing the gap between the data model (which already supports everything) and what the user can actually do.

## What Changes

- The "+ New Commission" button on the commissions list page opens a new `CommissionFormModal` for creating a commission.
- Each commission row gains two action controls (brush icon and trash icon) at the end — brush opens the same modal in edit mode pre-filled with that commission's data, trash opens a delete confirmation dialog.
- **Commission form modal** — handles both create and edit, with fields: Title (optional), Artist (combobox + create artist shortcut), Notes, Status (chip dropdown), Price, Paid (checkbox), Paid Date (datepicker), Finish Date (datepicker), and a collapsible Artpieces section (thumbnail strip, title search, artist/character filter).
- **Artpieces section** — attach-only; search results exclude artpieces already attached to a different commission. Artpiece creation is not possible from inside the modal. Edit mode pre-populates the strip from a GET /commissions/:id fetch.
- **Delete flow** — confirmation dialog stating that linked artpieces are left intact (the DB already enforces `ON DELETE SET NULL`).
- **Status chip** — extracted as a shared controlled `<StatusChip>` component, replacing the inline-only implementation in `CommissionsTable`. The table continues to autosave on change; the modal drives local form state.
- `ArtpieceSummary` frontend type gains a `commission_id` field (the backend already sends it) to enable client-side filtering of already-attached artpieces.

## Capabilities

### New Capabilities

- `web-commission-form-modal`: Create/edit commission modal and delete confirmation dialog, including the artpieces attachment section, the shared StatusChip component, and all new FE mutation hooks.

### Modified Capabilities

- `web-commissions-list`: The "+ New Commission" button now opens the form modal (was a no-op). Each table row gains edit and delete action controls. The StatusChip refactor touches the table's inline status cell.

## Impact

**Frontend**
- New component: `CommissionFormModal` (create + edit, shared modal)
- New component: `DeleteCommissionDialog`
- New shared component: `StatusChip` (replaces inline `DropdownMenu` in `CommissionsTable`)
- New hooks: `useCreateCommission`, `useUpdateCommission`, `useCommission` (single fetch by ID), `useReplaceArtpieces`, `useDeleteCommission`
- Type change: `ArtpieceSummary` — add `commission_id: string | null`
- Modified: `CommissionsPage` — wire "+ New Commission" button
- Modified: `CommissionsTable` — add actions column, swap status cell to `<StatusChip>`

**Backend** — no changes required. All endpoints (`POST /commissions`, `PUT /commissions/:id`, `PUT /commissions/:id/artpieces`, `GET /commissions/:id`, `DELETE /commissions/:id`) are already implemented and correct.

**APIs used (existing)**
- `POST /commissions` — create with optional artpiece_ids
- `PUT /commissions/:id` — full field update
- `PUT /commissions/:id/artpieces` — full artpiece set replace (fired only if set changed)
- `GET /commissions/:id` — edit pre-populate (includes artpieces with thumbnail URLs)
- `DELETE /commissions/:id` — delete, artpieces are unlinked not deleted
