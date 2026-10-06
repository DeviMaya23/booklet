## Context

The Artpieces gallery page (`/app/artpieces`) currently allows browsing, creating, and deleting artpieces. There is no way to view full artpiece details (notes, artist, characters, full file list) or edit them after creation. All backend endpoints needed for this feature already exist: `GET /artpieces/:id`, `PUT /artpieces/:id`, `PUT /artpieces/:id/files`, `PUT /artpieces/:id/cover`, `DELETE /artpieces/:id/files/:file_id`, and `DELETE /artpieces/:id?delete_files=true`.

## Goals / Non-Goals

**Goals:**
- Add an inline artpiece detail view accessible by double-clicking a gallery card.
- View mode: display all artpiece metadata and a thumbnail grid of attached files with a cover badge.
- Edit mode: editable fields using existing input components; drag-drop upload strip; per-thumbnail `...` menu for "Set as cover" and "Remove from artpiece."
- Shared delete dialog (used by both gallery card and detail page) with "Also delete attached files" checkbox.

**Non-Goals:**
- Routing change — no new route added; detail view is in-page state.
- Commission linking/display — commission_id is not shown or editable here.
- File name/notes editing from this page — that lives in the file inbox tile edit overlay.
- Adding inbox files directly from this page (that path is the "Add to Existing Artpiece" flow).

## Decisions

### Decision: State toggle over nested route
The detail view is rendered via `selectedArtpieceId` state in `ArtpiecesPage`, not a new route. The back chevron is a close button (`setSelectedArtpieceId(null)`), not browser navigation.

**Alternatives considered:** `/app/artpieces/:id` nested route — would give free browser-history support but adds routing complexity and deviates from how every other entity (artists, characters, files) works in this app. A future navigation use-case (e.g., jumping from a commission) can layer routing on top when it actually materialises.

### Decision: File uploads attach immediately — no staging
When the user drops or picks files in edit mode, the upload pipeline runs (init → S3 PUT → complete → `POST /artpieces/:id/files/:file_id`) immediately on completion, exactly as every other upload flow in the app. Cancel cannot undo uploads; uploaded files remain attached and appear in the grid on the next open.

**Alternatives considered:** Stage uploads and only commit on Save — adds a "pending upload" state inconsistent with the rest of the app; the extra complexity has no user-visible benefit given the existing app convention.

### Decision: Working set model — track locally-removed IDs only
Edit mode maintains a `locallyRemovedIds: Set<string>` and a `pendingCoverFileId: string | null`. The displayed file grid = artpiece's current files (from `useArtpiece` cache, updated live as uploads complete) minus `locallyRemovedIds`. Cancel clears both local states and resets field inputs. Save calls `PUT /artpieces/:id/files` with `(currentFileIds - locallyRemovedIds)` only if `locallyRemovedIds` is non-empty (uploads are already committed immediately and do not require a Save-time PUT); calls `PUT /artpieces/:id/cover` only if `pendingCoverFileId` changed.

No entry-snapshot is needed because uploads update the server immediately. The only reversible local changes are removes and cover selection.

### Decision: New `ArtpieceDetailFileGrid` instead of reusing `ArtpieceFilesInput`
`ArtpieceFilesInput` is built for the create flow: a scrollable list with editable name/notes per file, accumulating IDs for a future submission. The detail file grid is a thumbnail grid with a cover badge and per-tile `...` menu — a fundamentally different shape. The upload machinery (init → S3 PUT → complete) can be extracted as a hook or inline, but the component itself must be new.

### Decision: Shared `DeleteArtpieceDialog` with optional file count
A single `DeleteArtpieceDialog` component (props: `artpieceId`, `open`, `onOpenChange`, `onSuccess`, optional `fileCount`) replaces the inline AlertDialog in `ArtpiecesGrid` and is also used by the detail page. When `fileCount` is provided, the checkbox label reads "Also delete X attached files"; when absent, it reads "Also delete attached files." `useDeleteArtpiece` is updated to accept `{ id: string, deleteFiles: boolean }`.

## Risks / Trade-offs

- **Presigned URLs expire** — `GET /artpieces/:id` returns short-lived presigned URLs for file thumbnails. A user who stays on the detail page for an extended period will see stale/broken images. Mitigation: consistent with existing behaviour across the app; acceptable for now.
- **Upload-then-cancel leaves orphaned attachments** — A user who uploads files and then hits Cancel has attached files they may not have intended to keep. Mitigation: this is the stated, intentional design; users can remove files in a subsequent edit. The behaviour is identical to every other upload flow.
- **`useArtpiece` cache coherence during edit** — live cache updates (from React Query background refetches) could cause the displayed file list to shift while the user is in the middle of an edit. Mitigation: `locallyRemovedIds` is applied on top of whatever the cache returns, so a background refetch won't restore a file the user just removed locally.
