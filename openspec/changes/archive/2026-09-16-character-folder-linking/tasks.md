## 1. Database Migration

- [x] 1.1 Write migration to add `folder_name text not null default ''` to `character_folders` table

## 2. Domain & Usecase Layer

- [x] 2.1 Add `FolderName string` to `domain.CharacterFolder` struct
- [x] 2.2 Add `IDPSubject string` to `CreateCharacterParams` and `UpdateCharacterParams`
- [x] 2.3 Inject `BookleafClient` into `characterUsecase` constructor and update `NewCharacterUsecase` signature in `main.go`
- [x] 2.4 Add Bookleaf validation helper in `character_usecase.go`: call `GetPublicFolders`, build ID→name map, filter incoming IDs, log INFO for drops, return error on Bookleaf failure
- [x] 2.5 Call validation helper in `Create`: skip if no folder IDs, otherwise enrich `character.Folders` with names before calling repo
- [x] 2.6 Call validation helper in `Update`: skip if folder IDs list is empty, otherwise pass enriched `[]domain.CharacterFolder` to repo

## 3. Repository Layer

- [x] 3.1 Update `updateWithFolders` to accept `[]domain.CharacterFolder` (pre-enriched with names) instead of `[]uuid.UUID`; insert rows with `folder_name` populated
- [x] 3.2 Update `Update` method to pass the enriched folder slice through to `updateWithFolders`
- [x] 3.3 Update `usecase.UpdateCharacterParams` usage in repository to reflect the new folder param shape

## 4. Handler Layer

- [x] 4.1 In `CreateCharacter` handler: extract `idpSubject` via `middleware.AuthenticatedIDPSubjectFromContext`; return 401 if missing; pass into `CreateCharacterParams.IDPSubject`
- [x] 4.2 In `UpdateCharacter` handler: extract `idpSubject` the same way; pass into `UpdateCharacterParams.IDPSubject`
- [x] 4.3 Update `CharacterResponse` struct: replace `FolderIDs []string` with `Folders []FolderResponse` where `FolderResponse` has `ID` and `Name` string fields
- [x] 4.4 Update `toCharacterResponse` helper to map `character.Folders` to `[]FolderResponse{ID: f.FolderID.String(), Name: f.FolderName}`

## 5. Bruno API File

- [x] 5.1 Update the existing character PUT Bruno request to include a sample `folder_ids` array
- [x] 5.2 Update any Bruno character response examples to reflect `folders: [{id, name}]` shape

## 6. Backend Unit Tests

- [x] 6.1 Usecase `Create` — test: Bookleaf returns folders, all IDs valid → character created with enriched folder names
- [x] 6.2 Usecase `Create` — test: some folder IDs absent from Bookleaf → absent IDs dropped, valid ones persisted
- [x] 6.3 Usecase `Create` — test: Bookleaf call fails → error returned, no character created
- [x] 6.4 Usecase `Create` — test: empty folder_ids → Bookleaf not called, character created with no folders
- [x] 6.5 Usecase `Update` — test: Bookleaf returns folders, all IDs valid → character updated with enriched folder names
- [x] 6.6 Usecase `Update` — test: some folder IDs absent from Bookleaf → absent IDs dropped
- [x] 6.7 Usecase `Update` — test: Bookleaf call fails → error returned, character unchanged
- [x] 6.8 Usecase `Update` — test: empty folder_ids → Bookleaf not called, all assignments removed
- [x] 6.9 Handler `CreateCharacter` — test: missing idpSubject in context → 401
- [x] 6.10 Handler `UpdateCharacter` — test: missing idpSubject in context → 401
- [x] 6.11 Handler `CreateCharacter` — test: Bookleaf failure propagates → 500
- [x] 6.12 Handler `UpdateCharacter` — test: Bookleaf failure propagates → 500

## 7. Backend Lint

- [x] 7.1 Run `golangci-lint run ./...` and fix any issues

## 8. Frontend — Type & API Updates

- [x] 8.1 Update `Character` type in `useCharacters.ts`: replace `folder_ids: string[]` with `folders: { id: string; name: string }[]`
- [x] 8.2 Update `useUpdateCharacter` to accept `folder_ids: string[]` in its params and include it in the request body (remove the hardcoded `folder_ids: []`)
- [x] 8.3 Update `useCreateCharacter` to accept `folder_ids?: string[]` in its params and include it in the request body if provided
- [x] 8.4 Add `usePublicFolders` hook: `GET /folders`, map `folder_list` response to `{ id: string; name: string }[]`

## 9. Frontend — FolderPicker Component

- [x] 9.1 Create `FolderPicker` component accepting: `selected: { id: string; name: string }[]`, `available: { id: string; name: string }[]`, `onAdd: (folder) => void`, `onRemove: (id: string) => void`, `disabled: boolean`
- [x] 9.2 Render selected folders as chips with a remove button each
- [x] 9.3 Render a search input that filters `available` folders by name (case-insensitive); selected folders excluded from the list
- [x] 9.4 Clicking a result in the dropdown calls `onAdd` and closes/clears the search
- [x] 9.5 Render a Folders section label with a tooltip (question mark icon) explaining the list comes from Bookleaf

## 10. Frontend — CharacterFormModal Integration

- [x] 10.1 Add folder state to the modal: `addedFolders`, `removedIds` tracked against the original `character.folders` list
- [x] 10.2 Initialize selected folders from `character.folders` on modal open; reset on close
- [x] 10.3 Call `usePublicFolders` in the modal; disable the `FolderPicker` add input if the query fails or is loading
- [x] 10.4 Wire `FolderPicker` into the modal form with `selected`, `available` (public folders minus already-selected), `onAdd`, `onRemove`, and `disabled` props
- [x] 10.5 On submit, compute final folder IDs: `(original + added) - removed`; pass to `useCreateCharacter` / `useUpdateCharacter`
- [x] 10.6 Update `resetForm` to clear folder diff state

## 11. Frontend Build & Lint

- [x] 11.1 Run `npm run build` and fix any type errors
- [x] 11.2 Run `npm run lint` and fix any issues
