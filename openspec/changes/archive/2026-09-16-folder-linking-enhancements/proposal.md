## Why

The character folder picker silently disables when Bookleaf is unreachable, leaving users with no indication of why the field is inactive or how to recover. Separately, folder names and IDs persisted in Booklet can drift from Bookleaf's current state without any reconciliation on modal open — stale names render, and deleted folders stay visible until the next successful save.

## What Changes

- On successful GetPublicFolders load, silently reconcile against the character's stored folders: swap any drifted names, remove any IDs no longer present in Bookleaf (treated as manual removal; BE handles DB cleanup on next PUT)
- On GetPublicFolders failure, show an inline message below the disabled picker — "Couldn't reach folder list — Retry" — with a manual retry button; no auto-retry
- Configure the public folders query with `retry: false` (immediate failure, no silent auto-retry) and a 15-minute client-side cache

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `folder-linking-ui`: degraded-state requirement changes (silent disable → inline message + manual retry); adds folder name and ID reconciliation behavior on successful fetch

## Impact

- `frontend/src/features/characters/api/usePublicFolders.ts` — query config
- `frontend/src/features/characters/components/CharacterFormModal.tsx` — reconciliation logic, new state, updated FolderPicker props
- `frontend/src/features/characters/components/FolderPicker.tsx` — new `isError`/`onRetry` props, inline error message
- `openspec/specs/folder-linking-ui/spec.md` — updated requirements and scenarios
