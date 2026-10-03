## Context

The Inbox page (`FilesPage.tsx`) and `FileInboxGrid.tsx` have two "New Artpiece" entry points that are no-op stubs: a dropdown in the selection toolbar and a right-click context menu item. The backend already supports creating an artpiece with `file_ids` in a single call (`POST /artpieces`), and file metadata updates exist via `PUT /files/:id`. No backend work is needed.

The artist dropdown and character chip multi-select patterns are established in `ImageFormModal.tsx` and are reused as-is. The file upload pattern (init → PUT → complete) is established in `FileInboxGrid.tsx` and is reused inside the modal.

## Goals / Non-Goals

**Goals:**
- Deliver a working "Create New Artpiece" modal reachable from both entry points
- Files sub-component (`ArtpieceFilesInput`) is self-contained and reusable for the future "Add to Existing Artpiece" flow
- Upload-within-modal follows the same flow as the inbox grid upload; thumbnail polling reuses the existing mechanism
- File name/notes edits save on blur via `PUT /files/:id`

**Non-Goals:**
- "Add to Existing Artpiece" flow (deferred)
- Per-tile "View Detail" action (existing no-op, not touched)
- Any backend changes

## Decisions

### D1: Two components — `ArtpieceFormModal` and `ArtpieceFilesInput`

The files region is split into its own component rather than inlined into the artpiece modal. This keeps the modal manageable in size (the files region has its own upload state, polling state, and per-row state) and makes the component available for the future "Add to Existing Artpiece" flow without restructuring.

Both files live in `frontend/src/features/files/components/` — the files feature owns both since the sub-component is fundamentally about file management.

**Alternative considered:** Single component. Rejected because `ImageFormModal.tsx` is already ~500 lines for a simpler form; adding upload state, per-row state, and thumbnail polling inline would make the modal hard to follow.

### D2: Draft file state shape

Files in the modal exist in two states: uploading (in-flight, no server ID yet) or uploaded (ID known, thumbnail may still be pending). A discriminated union makes state transitions explicit and eliminates null-checking:

```ts
type DraftFile =
  | { status: 'uploading'; clientId: string; fileName: string }
  | { status: 'uploaded'; clientId: string; id: string; name: string; notes: string; thumbnailUrl: string | null }
```

`ArtpieceFilesInput` manages this list internally and exposes `onFileIdsChange: (ids: string[]) => void` to the parent — the parent only needs the final ID list for the save call.

**Alternative considered:** Lifting the full draft state to the parent modal. Rejected because the parent only needs IDs at save time; keeping the draft list internal to `ArtpieceFilesInput` isolates upload complexity and supports future reuse with a clean interface.

### D3: Thumbnail polling via `useFiles` query

After a file is uploaded within the modal, it exists in the backend as an unassigned file, so it appears in the `useFiles` (`/files?unassigned=true`) query which already polls every 2s when any file lacks a `thumbnail_url`. `ArtpieceFilesInput` reads `useFiles` data and derives thumbnail readiness for its draft files by matching IDs — no separate polling mechanism is introduced.

**Alternative considered:** A separate per-file poll query inside the component. Rejected because `useFiles` is already running and doing the same work; a second poller is redundant.

### D4: File name/notes save on blur

When the user edits a file's name or notes field and blurs, `PUT /files/:id` fires immediately. This is consistent with how files are treated elsewhere — as independent entities from the artpiece they may eventually belong to. Orphan case (user edits then cancels) is intentionally accepted: the file's new name/notes persist in the dump, which is harmless.

The Save button is disabled while any upload is in-flight (`status: 'uploading'` entries exist). This prevents a race where Save is clicked before a blur fires on a name field that's still being edited while an upload completes — the user must let uploads finish before Save becomes available.

**Alternative considered:** Batch all name/notes updates at Save time alongside the artpiece creation. Rejected because it requires tracking dirty state per file and complicates the save flow with N+1 calls before the artpiece creation; save-on-blur is simpler and consistent with the "files are independent entities" model already in use.

### D5: Pre-population from selection

When opened from a file selection, the modal receives `initialFileIds: string[]`. It reads the current `useFiles` data (already fetched) to populate name/notes for those IDs. If a file's thumbnail is not yet ready (`thumbnail_url === null`), its row shows a spinner and the name/notes fields are disabled — same as for freshly uploaded files.

Files from the initial selection enter the draft list as `status: 'uploaded'` directly (no upload flow needed — they already exist on the server).

### D6: Upload affordance visual

The drop zone at the top of the scrollable files region uses a dashed border (`border-dashed border-border rounded-lg`) rather than the filled `bg-muted` style used in the image/character upload slots. The filled style communicates "this is the file slot" (single replace); the dashed style communicates "drop zone / add more" (additive). This is the first dashed-border upload affordance in the codebase; the pattern is intentionally distinct from the single-file slots.

## Risks / Trade-offs

- **Blur-save orphan files**: If a user edits a file name and cancels the modal, that name change is already persisted. Acceptable — dump file metadata is mutable and the user can correct it. The alternative (rollback on cancel) introduces a revert call that adds more complexity than the risk warrants.

- **`useFiles` polling dependency**: `ArtpieceFilesInput` reads `useFiles` to get thumbnail readiness. If the user navigates away from the inbox page and the query is unmounted, thumbnails won't update. In practice the modal is only reachable from the inbox page, so the query is always mounted while the modal is open. This assumption should be revisited if the modal is later added to other pages.

- **`status: 'uploading'` files on modal close**: If the user somehow closes the modal (e.g. Escape) while uploads are in flight, orphaned in-progress uploads may complete and add files to the dump. The modal should disable the close action while any upload is in-flight, or at minimum let in-flight uploads finish and add to the dump silently (no artpiece created). This is a UX edge case worth locking down in the spec.
