## Why

Characters can be associated with Bookleaf folders via `character_folders`, but there is currently no UI to manage those links and no backend validation that the linked folder IDs actually exist and are public. This change closes both gaps: it adds Bookleaf-validated folder linking to the character create/update flow and surfaces a folder picker in the character modal.

## What Changes

- **BREAKING** `character-management`: Character API responses replace the `folder_ids: string[]` field with `folders: [{id, name}][]` — callers must update to the new shape
- `character-folders`: `character_folders` table gains a `folder_name` column; name is written by BE on every validated save, never accepted from FE input
- `character-management`: Character create and update now validate incoming folder IDs against Bookleaf's `GetPublicFolders` on every request — unknown IDs are silently dropped, Bookleaf unavailability fails the entire operation
- `character-management`: `characterUsecase` gains a `BookleafClient` dependency for folder validation
- `folder-linking-ui` (new): Character modal gains a multi-select folder picker backed by `GET /folders`; existing folder names render from persisted data so the modal works even when Bookleaf is unreachable at open time

## Capabilities

### New Capabilities
- `folder-linking-ui`: Multi-select folder picker in the character create/edit modal — shows chips for selected folders, searches available public folders from Bookleaf, gracefully degrades when Bookleaf is unreachable

### Modified Capabilities
- `character-folders`: Adds `folder_name` column to `character_folders` and introduces Bookleaf validation on every folder write — changes the data model and write semantics
- `character-management`: Response shape change (`folder_ids → folders`); create/update now fail when Bookleaf is unreachable and silently drop IDs absent from Bookleaf's public list

## Impact

- **Backend**
  - New migration: add `folder_name text not null` to `character_folders`
  - `domain.CharacterFolder`: add `FolderName string`
  - `CreateCharacterParams` / `UpdateCharacterParams`: add `IDPSubject string`
  - `characterUsecase`: inject `BookleafClient`; validate + enrich folder IDs on create and update
  - Character repository: `updateWithFolders` accepts `[]domain.CharacterFolder` (pre-enriched by usecase) instead of `[]uuid.UUID`
  - Character handler: extract `idpSubject` from JWT context on create and update; 401 if missing
  - Character response: `folder_ids []string` → `folders []FolderResponse{ID, Name}`
- **Frontend**
  - `Character` type: `folder_ids: string[]` → `folders: {id: string; name: string}[]`
  - New `usePublicFolders` hook against existing `GET /folders`
  - New `FolderPicker` component (chips + search dropdown, no add option)
  - `CharacterFormModal`: folder state management (diff tracking), picker integration, tooltip explaining Bookleaf source
  - `useCreateCharacter` / `useUpdateCharacter`: pass `folder_ids` on submit
