## 1. Scaffold & shared components

- [x] 1.1 Create feature directories: `src/features/characters/api/`, `src/features/characters/components/`, `src/features/images/api/`, `src/features/images/components/`
- [x] 1.2 Add shadcn `alert-dialog` component via `npx shadcn add alert-dialog`
- [x] 1.3 Add shadcn `dropdown-menu` component via `npx shadcn add dropdown-menu`
- [x] 1.4 Create `src/components/ResourceCard.tsx`: accepts `imageUrl?: string | null`, `label: string`, `onDeleteClick: () => void`; renders square aspect-ratio card with image or `ImageOff` placeholder, bottom-left label overlay, top-right `(...)` dropdown menu containing Delete

## 2. Characters feature

- [x] 2.1 Create `src/features/characters/api/useCharacters.ts`: React Query hook calling `GET /characters` via `apiFetch`
- [x] 2.2 Create `src/features/characters/api/useDeleteCharacter.ts`: React Query mutation calling `DELETE /characters/:id`, invalidates the characters list query on success
- [x] 2.3 Create `src/features/characters/components/CharactersGrid.tsx`: renders `ResourceCard` per character (passes `avatar_url` as `imageUrl`, `name` as `label`); wires delete mutation and confirmation dialog
- [x] 2.4 Implement `CharactersPage.tsx`: search bar state, filters characters by `name` (case-insensitive), `+ New` no-op button, renders `CharactersGrid`

## 3. Images feature

- [x] 3.1 Create `src/features/images/api/useImages.ts`: React Query hook calling `GET /images` via `apiFetch`
- [x] 3.2 Create `src/features/images/api/useDeleteImage.ts`: React Query mutation calling `DELETE /images/:id`, invalidates the images list query on success
- [x] 3.3 Create `src/features/images/components/ImagesGrid.tsx`: renders `ResourceCard` per image (passes `thumbnail_url` as `imageUrl`, `title ?? "Untitled"` as `label`); wires delete mutation and confirmation dialog; client-side filter excludes null-`title` images from any search query
- [x] 3.4 Implement `ImagesPage.tsx`: search bar state, filters images by `title` (case-insensitive, null titles excluded from all queries), `+ New` no-op button, renders `ImagesGrid`

## 4. Quality

- [x] 4.1 Run `npm run build` in `frontend/` and fix any TypeScript or build errors
- [x] 4.2 Run `npm run lint` in `frontend/` and fix any lint errors
