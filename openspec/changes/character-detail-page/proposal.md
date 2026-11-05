# Proposal

## Why

Viewing or editing a character today means opening `CharacterFormModal`, the same dialog used to create one. A character is more than its form fields: it has an avatar, a verse, attached Bookleaf folders, and the artpieces and references it appears in. None of that is visible without leaving the character. Folder management is cramped inside a dialog that widens to fit the "Choose folders" column. A dedicated detail page, modeled on the artpiece detail page, gives characters room to be browsed as well as edited.

## What Changes

- Add a character detail page at `/app/characters/:id`, rendered by its own `CharacterDetailPage` (not hosted inside `CharactersPage`).
- **View mode**: back link ("← Characters"), avatar tile, name, verse chip, Edit button, notes, a Folders row (Bookleaf share links when available), and an Artpieces / References tab strip with counts.
- **Edit mode** (in place, same page): name input, Cancel/Save, an amber "Unsaved changes: …" summary line, verse combobox, notes, avatar tile with Change / remove, folder chips, "Choose folders" button, and a "Delete character" action. Leaving with unsaved changes is guarded by `useBlocker` and a discard dialog, as in the artpiece detail.
- **Choose folders drawer**: the folder picker opens in a right-side `Sheet` on the detail page. The shared `ChooseFoldersPanel` gains an "Updated N min ago" label and a "View all in Bookleaf" link, and folder previews show 24 images per page instead of 15.
- **Tabs**: add a Base UI Tabs wrapper (`components/ui/tabs.tsx`). The Artpieces tab lists the character's non-reference artpieces; References lists the character's `is_reference` artpieces. Items never appear in both. Both are derived client-side from `useArtpieces()` and reuse `ArtpiecesGrid` and `FileViewer`.
- Add a `useCharacter(id)` hook for the existing `GET /characters/:id`.
- `CharactersPage`: clicking a character card navigates to the detail page instead of opening the modal.
- `CharacterFormModal`: **BREAKING (internal)** — the edit branch is removed and the modal becomes create-only. "New Character" keeps the modal, including its inline folder picker.
- The avatar pick / drop / paste / crop state and the init → R2 PUT → complete upload pipeline are factored out of the modal so the modal (create) and the page (edit) share one implementation.
- No backend or Bookleaf changes.

## Capabilities

### New Capabilities
- `web-character-detail`: the character detail page: routing and navigation, view mode, edit mode (fields, avatar, unsaved-changes guard, save, delete), folder chips and the Choose folders drawer, Artpieces / References tabs.

### Modified Capabilities
- `web-characters`: card click navigates to the detail page; `CharacterFormModal` becomes create-only (edit-mode requirements and the in-modal delete are removed); the avatar picker's edit-mode "remove" behavior moves to the detail page.
- `folder-linking-ui`: the picker also appears as a drawer on the detail page; folder preview page size 15 → 24; the panel shows "Updated N min ago" and a "View all in Bookleaf" link; reconciliation, diff and degraded-state behavior also apply on the detail page.
- `web-character-avatar-crop`: clipboard paste and source-size validation also apply on the detail page in edit mode, not only when the modal is open.

## Impact

- **Frontend only.** New: `pages/CharacterDetailPage.tsx`, `features/characters/api/useCharacter.ts`, detail-page components and hooks under `features/characters/`, `components/ui/tabs.tsx`. Modified: `router.tsx` (new route), `pages/CharactersPage.tsx`, `CharacterFormModal.tsx`, `ChooseFoldersPanel.tsx`, `FolderPreview.tsx`.
- **Sidebar**: `AppSidebar` already uses `useMatch(..., { end: false })`, so Characters stays highlighted on `/app/characters/:id` with no code change; a scenario locks that in.
- **Data**: `GET /characters/:id` (exists, no FE hook yet), `GET /artpieces` (already returns `characters[]` and `is_reference`), `GET /folders` and `GET /folders/:id/images` (unchanged).
- **Call sites checked**: `CharacterFormModal` is used only by `CharactersPage`; `ChooseFoldersPanel`, `FolderList` and `FolderPreview` are used only inside the modal (the page adds a second consumer of the first and third).
- **Dependencies**: none new; `@base-ui/react` is already installed and provides Tabs.
- **Backend / Bookleaf**: none.
