## 1. Database Migration

- [x] 1.1 Write migration `000029` up: add `last_used_at timestamptz NULL` to `artists`, create `artist_links` table (id, artist_id FK CASCADE, url, is_primary, created_at, updated_at), drop `artist_link` column from `artists`, backfill `UPDATE artists SET last_used_at = NOW()`
- [x] 1.2 Write migration `000029` down: re-add `artist_link text NULL` to `artists`, drop `artist_links` table, drop `last_used_at` from `artists`

## 2. Domain & Repository

- [x] 2.1 Add `domain.ArtistLink` struct (ID, ArtistID, URL, IsPrimary, CreatedAt, UpdatedAt) with GORM tags
- [x] 2.2 Update `domain.Artist`: drop `ArtistLink *string`, add `LastUsedAt *time.Time`, add `Links []ArtistLink` GORM association (`foreignKey:ArtistID`)
- [x] 2.3 Update `artistRepository.Update()`: drop `artist_link` from the updates map; add `db.Model(&artist).Association("Links").Replace(links)` to replace all links
- [x] 2.4 Update `artistRepository.Create()`: pass links through `Association("Links").Replace(links)` after insert (or include in the initial Create if GORM handles association on create)
- [x] 2.5 Add `artistRepository.UpdateLastUsedAt(ctx, id, userID uuid.UUID) error` method
- [x] 2.6 Update `artistRepository.List()`: change `Order("name ASC")` to `Order("last_used_at DESC NULLS LAST, name ASC")`
- [x] 2.7 Update `commission_repository.go`: change `Preload("Artist")` to `Preload("Artist.Links")` in `GetByID` and `List`
- [x] 2.8 Update `artpiece_repository.go`: add `Preload("Artist.Links")` wherever `Preload("Artist")` is used

## 3. Usecase Layer

- [x] 3.1 Update `usecase/artist_repository.go`: drop `ArtistLink` from `CreateArtistParams` and `UpdateArtistParams`; add `Links []ArtistLinkInput` (url string, is_primary bool) to both; add `UpdateLastUsedAt(ctx context.Context, id, userID uuid.UUID) error` to `ArtistRepository` interface
- [x] 3.2 Update `usecase/commission_repository.go`: add `UpdateLastUsedAt(ctx context.Context, id, userID uuid.UUID) error` to `CommissionArtistRepository` interface
- [x] 3.3 Update `usecase/artpiece_repository.go`: add `UpdateLastUsedAt(ctx context.Context, id, userID uuid.UUID) error` to `ArtpieceArtistRepository` interface
- [x] 3.4 Update `artist_usecase.go Create()`: set `LastUsedAt` to `now()` on the new artist domain object before passing to repo; pass `params.Links` through to repo
- [x] 3.5 Update `commission_usecase.go Create()` and `Update()`: after successfully validating artist ownership, call `u.artistRepo.UpdateLastUsedAt(ctx, *params.ArtistID, userID)` when `ArtistID` is non-nil
- [x] 3.6 Update `artpiece_usecase.go Create()` and `Update()`: same pattern — call `u.artistRepo.UpdateLastUsedAt(ctx, *params.ArtistID, userID)` when `ArtistID` is non-nil

## 4. Handler Layer

- [x] 4.1 Update `artist_handler.go`: add `artistLinkInput` (url, is_primary) and `artistLinkResponse` (id, url, is_primary) types; update `createArtistRequest` and `updateArtistRequest` to replace `ArtistLink *string` with `Links []artistLinkInput validate:"dive"`; update `artistResponse` to replace `ArtistLink *string` with `Links []artistLinkResponse`; update `toArtistResponse()` to map `artist.Links`; add `url` case to `tagMessage()` in `validator.go` returning a useful message
- [x] 4.2 Update `commission_handler.go` `toCommissionResponse()`: replace `artistLink = c.Artist.ArtistLink` with a scan of `c.Artist.Links` for the entry with `IsPrimary == true`; use its URL as `artistLink`, or `nil` if none found
- [x] 4.3 Update Bruno collection: `artists/create_artist.bru` and `artists/update_artist.bru` — replace `artist_link` field with a `links` array example

## 5. Frontend

- [x] 5.1 Update `Artist` type in `features/artists/api/useArtists.ts`: drop `artist_link: string | null`; add `links: { id: string; url: string; is_primary: boolean }[]`
- [x] 5.2 Update `features/artists/api/useCreateArtist.ts`: drop `artist_link` from request payload type; add `links?: { url: string; is_primary: boolean }[]`
- [x] 5.3 Update `features/artists/api/useUpdateArtist.ts`: drop `artist_link` from request payload type; add `links: { url: string; is_primary: boolean }[]`
- [x] 5.4 Rewrite `features/artists/components/ArtistFormModal.tsx`: replace single link input with multi-link editor — state as `{ url: string; is_primary: boolean }[]` and parallel `(string | null)[]` for per-link errors; implement add/remove, star toggle, first-link-is-primary default, remove-primary-promotes-first logic, `https://` prepend on blur, `new URL()` validation on blur with inline error, blank-link filtering on submit
- [x] 5.5 Update `features/artists/components/ArtistsList.tsx`: derive copy-link value from `artist.links.find(l => l.is_primary)?.url ?? null`; update disabled/title logic accordingly
- [x] 5.6 Update `features/artpieces/components/ArtpieceDetailView.tsx`: replace `viewArtist.artist_link` with `viewArtist.links.find(l => l.is_primary)?.url ?? null` (two references around line 275)
- [x] 5.7 Update `features/commissions/components/CommissionFormModal.tsx`: in the artist pre-population block (~line 121), replace `artist_link: commission.artist_link` with `links: []`

## 6. Tests

- [x] 6.1 Update `repository/artist_repository_integration_test.go`: fix `TestArtistRepository_List` — current assertion expects alphabetical order; update to reflect `last_used_at DESC NULLS LAST, name ASC` sort; update `TestArtistRepository_Create` to remove `ArtistLink` assertion; add integration test for link replace behaviour on `Update`
- [x] 6.2 Update `usecase/artist_usecase_test.go`: remove `ArtistLink` from `fakeArtistRepository.Update()` and params structs; add test `TestCreateArtist_SetsLastUsedAt` verifying `last_used_at` is set on the created artist
- [x] 6.3 Update `usecase/commission_usecase_test.go`: add `UpdateLastUsedAt` method to `fakeCommissionArtistRepository`; add test verifying `UpdateLastUsedAt` is called when `ArtistID` is set on commission create; add test verifying it is called on commission update with non-nil `ArtistID`
- [x] 6.4 Update `usecase/artpiece_usecase_test.go`: add `UpdateLastUsedAt` method to the fake artist repo used there; add test verifying `UpdateLastUsedAt` is called when `ArtistID` is set on artpiece create and update

## 7. Quality

- [x] 7.1 Run `golangci-lint run ./...` from `backend/` and fix all reported issues
- [x] 7.2 Run `npm run build` from `frontend/` and fix any type errors or build failures
- [x] 7.3 Run `npm run lint` from `frontend/` and fix any lint issues
