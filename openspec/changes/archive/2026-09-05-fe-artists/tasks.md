## 1. API Hooks

- [x] 1.1 Create `features/artists/api/useArtists.ts` — `GET /artists` query hook exporting `Artist` type and `ARTISTS_QUERY_KEY`
- [x] 1.2 Create `features/artists/api/useCreateArtist.ts` — `POST /artists` mutation hook; invalidates `ARTISTS_QUERY_KEY` on success
- [x] 1.3 Create `features/artists/api/useUpdateArtist.ts` — `PATCH /artists/:id` mutation hook; invalidates `ARTISTS_QUERY_KEY` on success
- [x] 1.4 Create `features/artists/api/useDeleteArtist.ts` — `DELETE /artists/:id` mutation hook; invalidates `ARTISTS_QUERY_KEY` on success

## 2. Artist Form Modal

- [x] 2.1 Create `features/artists/components/ArtistFormModal.tsx` — modal with Name (required), Link (optional), Blurb (optional textarea), Save and Cancel buttons; mode driven by presence of `artist` prop
- [x] 2.2 Wire create mode: submit calls `useCreateArtist`, closes modal and shows success toast on success, shows error toast on 409, modal stays open on any error
- [x] 2.3 Wire edit mode: form pre-populated from `artist` prop; submit calls `useUpdateArtist` with same success/error behaviour as create
- [x] 2.4 Add delete button (edit mode only): clicking opens an `AlertDialog` confirm; confirming calls `useDeleteArtist`, closes modal, shows success toast
- [x] 2.5 Block submission client-side when name field is empty

## 3. Artists List

- [x] 3.1 Create `features/artists/components/ArtistsList.tsx` — renders a flat list of rows; each row has a bullet, artist name, copy-link button, and edit button
- [x] 3.2 Copy-link button: calls `navigator.clipboard.writeText(artist_link)` and shows success toast; disabled when `artist_link` is null or empty

## 4. Artists Page

- [x] 4.1 Implement `pages/ArtistsPage.tsx` — owns search `useState`, fetches via `useArtists`, filters client-side by name (case-insensitive), renders search bar, `+ New` button, and `ArtistsList`
- [x] 4.2 Wire `+ New` button to open `ArtistFormModal` in create mode
- [x] 4.3 Wire edit button on each row to open `ArtistFormModal` in edit mode with the selected artist

## 5. Quality

- [x] 5.1 Run `npm run build` and fix any type or build errors
- [x] 5.2 Run `npm run lint` and fix any lint errors
