## Context

The Inbox page (`FilesPage`) already has a floating "+" button wired to a no-op. The existing `ArtpieceFilesInput` component handles upload + thumbnail polling + blur-save for file rows inside the Create New Artpiece modal, but it's designed around a parent-controlled Save: fields are draft state, errors are surfaced by a global action, and the whole form can be disabled from the parent.

The new modal has a fundamentally different contract: uploads commit immediately, fields auto-save on blur, and field-save failures must be surfaced inline since there is no final Save action to catch them.

## Goals / Non-Goals

**Goals:**
- Wire the "+" button to open an "Add files to dump" modal
- Support multi-file upload (drag-and-drop + file picker) inside the modal
- Show each file as a row with thumbnail, editable name, editable notes, and a remove button (post-upload only)
- Auto-save name/notes on blur; surface field-save failures inline with a red border + error message
- Match existing upload failure handling: remove row + toast on upload error
- Thumbnail pending state: spinner in the thumbnail slot until `thumbnail_gen_state` resolves

**Non-Goals:**
- Sharing a `FileRow` sub-component between this modal and `ArtpieceFilesInput` — they differ enough in contract that a shared component would be prop-soup
- Any backend changes — all required endpoints already exist
- Cancelling an in-flight upload (X is not shown while uploading)

## Decisions

### Decision: Separate `AddFilesModal` component, no shared FileRow

The two upload UIs differ at the row level: `ArtpieceFilesInput` holds draft state committed by a parent Save, while `AddFilesModal` auto-saves on blur with inline error feedback. Sharing a row component would require conditional props to switch between these modes, making both harder to reason about.

The actual reuse that matters is already abstracted into hooks (`useInitFileUpload`, `useCompleteFileUpload`, `useUpdateFile`, `useFiles`). `AddFilesModal` uses those same hooks directly, mirroring the pattern used by `FileInboxGrid`.

**Considered**: Extending `ArtpieceFilesInput` with an `autoSave` flag and per-field error state. Rejected — the component was not designed for this and adding optional modes makes the existing Create New Artpiece flow harder to reason about.

### Decision: Per-field error state tracked in the row's local state

Each uploaded file row tracks `{ nameError: string | null, notesError: string | null }` as part of its `DraftFile` state slice. On blur, the component calls `useUpdateFile` and catches failure to set the error. On the next successful save of that field, the error is cleared.

This keeps error state co-located with the row that owns it, without lifting it to the modal level.

### Decision: No X button while uploading

An in-flight upload cannot be cancelled (the PUT to R2 is fire-and-forget). Showing an X that does nothing would be misleading. The uploading row renders as: spinner thumbnail + disabled name input + disabled notes input, no X.

### Decision: `useFiles` polling for thumbnail resolution

Same as `ArtpieceFilesInput`: the modal subscribes to the `useFiles` query (which polls), and a `useEffect` syncs `thumbnailUrl` and `thumbnailGenState` into local row state when the query updates. No new polling mechanism is introduced.

## Risks / Trade-offs

- **Field-save failure UX is new** — there is no existing pattern in the app for inline field-save errors. The proposal defines it clearly (red border + short message below the field), but it's worth verifying the error message copy feels right during implementation.
- **Row removal on upload failure** — if multiple files are uploading in parallel and one fails, its row disappears while others continue uploading. This is consistent with existing behavior but could be surprising if the failed file was the first row.
