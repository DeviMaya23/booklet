## 1. Query Config

- [x] 1.1 Add `retry: false` and `staleTime: 15 * 60 * 1000` to the `useQuery` call in `usePublicFolders.ts`

## 2. FolderPicker — Inline Error Affordance

- [x] 2.1 Add `isError?: boolean` and `onRetry?: () => void` props to `FolderPickerProps`
- [x] 2.2 Render inline message below `TokenInput` when `isError` is true: "Couldn't reach folder list." followed by an underlined "Retry" button that calls `onRetry`; use `text-muted-foreground` color

## 3. CharacterFormModal — Reconciliation + Wiring

- [x] 3.1 Add `renamedFolders` state (`useState<Map<string, string>>(new Map())`) and reset it in `resetForm`
- [x] 3.2 Update `selectedFolders` computation to apply name overrides: `.map((f) => ({ ...f, name: renamedFolders.get(f.id) ?? f.name }))` on the original folders slice
- [x] 3.3 Add reconciliation `useEffect` watching `publicFoldersQuery.data`: build a `Map` from the response, then (a) collect name overrides for drifted original folders into `renamedFolders`, and (b) add missing original folder IDs to `removedIds`
- [x] 3.4 Pass `isError={publicFoldersQuery.isError}` and `onRetry={() => publicFoldersQuery.refetch()}` to `FolderPicker`

## 4. Spec Update

- [x] 4.1 Archive the `folder-linking-enhancements` change (`openspec archive change folder-linking-enhancements`) to merge the delta into `openspec/specs/folder-linking-ui/spec.md`

## 5. Build & Lint

- [x] 5.1 Run `npm run build` from the `frontend/` directory and fix any errors
- [x] 5.2 Run `npm run lint` from the `frontend/` directory and fix any issues
