## Context

The character folder picker (`CharacterFormModal` + `FolderPicker`) currently uses React Query's default retry behavior (3 retries, ~7s) before `isError` becomes true, and shows no indication during that window or after. The stored folder data (names, IDs) comes from the character list query cache and is never cross-checked against Bookleaf's live state on modal open.

Two concerns being addressed:
1. **UX gap**: disabled field with no explanation when Bookleaf is unreachable
2. **Data staleness**: chips can display stale names or deleted folders until next successful PUT

## Goals / Non-Goals

**Goals:**
- Inform the user immediately when the folder list can't be loaded, with a manual retry path
- Silently reconcile stale names and removed folder IDs when Bookleaf responds successfully
- Cache the folder list client-side to avoid redundant fetches across modal open/close cycles

**Non-Goals:**
- Auto-retry or background polling on failure
- Notifying the user when a folder is silently removed during reconciliation
- Fetching a fresh character record on modal open (reconciliation works against the existing prop)

## Decisions

### retry: false on usePublicFolders

React Query's default 3-retry behavior means users wait ~7 seconds of silent disabled state before an error surfaces. With `retry: false`, failure is immediate — the inline message appears at the same time the picker disables, which is the right moment. A manual retry button covers the recovery path.

Alternative considered: keep auto-retry and add a "Loading folders…" message during the retry window. Rejected — adds a transient loading sub-state to manage; the retry window is purely overhead for users who need to act anyway.

### Reconciliation via useEffect on publicFoldersQuery.data

The cross-check runs in a `useEffect` watching `publicFoldersQuery.data`. It re-runs whenever `data` changes — which covers both initial success and successful manual retry. No `hasReconciled` guard: re-running reconciliation on a retry response is correct behavior (fresh Bookleaf data may differ from the first attempt's data, if any).

Name overrides for original folders are stored in a separate `renamedFolders: Map<string, string>` state rather than mutating `originalFolders`. This keeps the diff logic clean — `originalFolders` stays as the reference point for which folders were added/removed, while `renamedFolders` is a display-only override layer applied at render time.

Missing IDs are handled by adding the folder's ID to `removedIds` — the same state used for manual chip removals. This means the diff submitted on save is naturally correct without any special-case logic.

### isError + onRetry props on FolderPicker

`FolderPicker` receives `isError` and `onRetry` from `CharacterFormModal`. The inline message + retry button render conditionally on `isError`. The retry button calls `publicFoldersQuery.refetch()` passed down as `onRetry`.

Alternative considered: handle the inline message inside `CharacterFormModal` directly, outside `FolderPicker`. Rejected — the message is logically part of the picker's degraded-state presentation; collocating it in `FolderPicker` keeps the layout and error state together.

### 15-minute staleTime

Folder list changes in Bookleaf are infrequent. A 15-minute cache means repeated modal open/close cycles within a session skip redundant network requests. A stale cache hit is treated as a successful response for reconciliation — accepted tradeoff given the low change frequency.

## Risks / Trade-offs

- **Reconciliation on stale cache**: If the cached folder list is up to 15 minutes old, a renamed or deleted folder in Bookleaf may not be caught until cache expiry or manual retry. Accepted — the same staleness tradeoff applies to the picker's add list, which was already established as staleness-tolerant.
- **No undo for silently removed folders**: A folder missing from Bookleaf's list is removed from the selection without user awareness. The user cannot re-add it from this UI (it's not in the picker). This is intentional per spec; the only recovery is Bookleaf-side (restore the folder's public visibility).
- **Effect dependency on originalFolders**: The reconciliation effect depends on `originalFolders`, which is derived from the `character` prop. This is stable for the modal's lifetime but would need revisiting if the modal ever refreshes the character mid-session.
