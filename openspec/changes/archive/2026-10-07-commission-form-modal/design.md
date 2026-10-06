## Context

All backend endpoints for commission create, update, replace-artpieces, get-by-id, and delete are already implemented. This is a purely frontend change. The commissions list page exists with a non-functional "+ New Commission" button. The table has no row-level actions.

Key codebase facts that informed the decisions below:
- The list query (`GET /commissions`) does **not** include artpieces — only `GET /commissions/:id` does (with cover thumbnail presigned URLs).
- The backend sends `commission_id` on artpiece list items, but the frontend `ArtpieceSummary` type doesn't include it.
- The status chip is currently inlined in `CommissionsTable` as a `DropdownMenu` that fires `patchStatus` directly — it is not a controlled component.
- `ArtpiecesFilterPopover` is self-contained with its own open/close state; the caller owns artist + characters filter state.
- `ArtistCombobox` + `ArtistFormModal` are already used together in `ArtpieceFormModal` — reuse exactly that pattern.

## Goals / Non-Goals

**Goals:**
- Commission create via "+ New Commission" button → `CommissionFormModal`
- Commission edit via brush icon per row → same modal pre-filled
- Commission delete via trash icon per row → `DeleteCommissionDialog`
- Artpiece attachment section inside the modal (attach existing, remove, search/filter)
- Shared `StatusChip` controlled component, used by both the modal form and the existing table cell

**Non-Goals:**
- "Add Complete Commission" flow (artpiece + commission created together)
- Inbox flow changes for attaching artpieces to commissions
- Any backend changes
- Character match all/any toggle in the artpiece search filter (fixed to `'any'`)

## Decisions

### D1: Edit pre-populate via GET /commissions/:id, not list cache

The list cache doesn't include artpieces. When the user clicks edit, the modal triggers `useCommission(id)` (a new `useQuery` for `GET /commissions/:id`) and shows a loading skeleton in the artpiece strip until it resolves. The list cache is not modified.

**Alternative considered:** Add artpiece preload to the list query. Rejected — it adds N presigned URL generation calls to every list load, including when the user never edits.

### D2: StatusChip as a shared controlled component

Extract `<StatusChip value onChange hasError />` into `features/commissions/components/StatusChip.tsx`. The `CommissionsTable` cell becomes `<StatusChip value={commission.status} onChange={(s) => patchStatus(commission.id, s)} hasError={hasCellError(commission.id, 'status')} />`. The modal uses the same component with local form state as `onChange` target.

`STATUS_OPTIONS` and `STATUS_LABELS` move into the new file and are re-exported for any consumer that needs them.

**Alternative considered:** Inline a separate DropdownMenu in the modal. Rejected — status values will likely grow; centralising avoids two places to update.

### D3: Artpiece set change detection

On Save in edit mode, the modal compares the initial strip (from the GET /commissions/:id response) against the current strip by set equality. If identical, `PUT /commissions/:id/artpieces` is not called. If different, it fires — same logic as the artpiece detail page's file-PUT pattern.

Initial strip is captured from the GET response when the modal opens and is held in a `useRef` alongside the mutable strip state (tracked in `useState`).

### D4: commission_id added to ArtpieceSummary type

`ArtpieceSummary` in `useArtpieces.ts` gains `commission_id: string | null`. In the artpiece search section of the modal, results are filtered client-side: exclude items where `commission_id !== null && commission_id !== currentCommissionId` (for create mode, `currentCommissionId` is `null` so any artpiece with a commission_id is excluded).

No backend changes needed — the field is already sent.

### D5: Single modal for create and edit

`CommissionFormModal` receives a `mode: 'create' | 'edit'` prop and an optional `commission: Commission` for edit mode. Title, submit label, and data source differ by mode; field layout is identical.

### D6: Artpieces section — collapsible, state persists

The artpiece section is collapsible (collapsed by default) using a local `open` boolean state. Collapse/expand only toggles visibility; the underlying state (strip, search query, filter selections) is never reset on collapse. `useArtpieces()` is always fetched when the modal mounts.

### D7: Delete confirmation dialog

`DeleteCommissionDialog` is a standard `AlertDialog` (using the existing shadcn component). Body copy: *"This commission will be permanently deleted. Any artpieces linked to it will be kept intact."* Confirm fires `DELETE /commissions/:id`, then invalidates the commissions query. No additional backend call is needed after delete.

### D8: ArtpiecesFilterPopover character match mode

The filter popover requires a `characterMatch` prop. The commission modal passes a fixed `'any'` and a no-op `onCharacterMatchChange`. The all/any toggle is hidden from the user in this context — if the popover renders it conditionally based on `selectedCharacters.length > 0`, the no-op is safe (the toggle section only appears when characters are selected, and clicking it just calls the no-op).

## Risks / Trade-offs

**Presigned URL TTL on artpiece thumbnails in the strip** → If a user holds the edit modal open for an extended time (>1 hour), thumbnails in the strip may expire. Mitigation: acceptable for now — the URL TTL is 1 hour and modal sessions are short. A full fix would require re-fetching on display but is out of scope.

**commission_id type addition to ArtpieceSummary** → Any existing code that spreads or destructures `ArtpieceSummary` will need the compiler to confirm it's still valid. TypeScript will catch any breakage at build time — no runtime risk.

**Single modal for create and edit** → If the two modes diverge significantly in future (e.g., edit grows an audit trail section), the modal may need to be split. Current layout is identical so a single component is the right starting point.
