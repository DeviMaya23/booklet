## Why

Artists often have multiple relevant links (commission storefront, social, portfolio) and the current single-link field forces users to pick one permanently. Separately, the artist lookup dropdown sorts alphabetically, which makes it hard to find the right artist when several share a similar name — users need to search by memory rather than recency.

## What Changes

- **BREAKING**: `artist_link` field removed from artist create/update request body and artist response; replaced by a `links[]` array (each item: `url`, `is_primary`)
- Artists can have multiple links; at most one is marked `is_primary`; links are managed inline in the artist form and saved with the artist (replace-all on update)
- Wherever a single artist link was previously surfaced (artist list, commission table, artpiece detail, dashboard), the primary link is used instead; if no primary is set, nothing shows
- New `last_used_at` column on `artists`; set to `now()` on artist creation and updated whenever an artist is assigned on a commission or artpiece (create or edit)
- Artist lookup dropdown (combobox) sorts by `last_used_at DESC, name ASC` instead of alphabetical
- Existing `artists.artist_link` column dropped; existing link values are not migrated — users re-enter links going forward
- Backfill: all existing artists get `last_used_at = now()` at migration time, breaking ties by name until real usage data accumulates

## Capabilities

### New Capabilities

- `artist-links`: Artists have a one-to-many `artist_links` table (url, is_primary); managed inline in the artist create/update API via a `links[]` array (replace-all); at most one link per artist may be `is_primary`
- `artist-recency-sort`: `artists.last_used_at` tracks when each artist was last used; updated on artist creation and on commission/artpiece create or update when an artist is assigned; artist lookup sorts by this field descending

### Modified Capabilities

- `artist-management`: **BREAKING** — `artist_link` field removed from create/update request and response; `links[]` array added to create/update request and to all artist responses (create, get, list, update)
- `web-artists`: Artist form modal gains a multi-link editor (add/remove links, mark one as primary); artist list's copy-link action uses the primary link

## Impact

**Backend**
- New DB table: `artist_links` (id, artist_id FK, url, is_primary, created_at, updated_at)
- New column: `artists.last_used_at timestamptz NULL`; backfilled to `now()`
- Dropped column: `artists.artist_link`
- `domain.Artist`: drop `ArtistLink`, add `LastUsedAt`, add `Links []ArtistLink` association
- `artistRepository`: drop `artist_link` from `Update()` map; add `UpdateLastUsedAt()`; add link replace via GORM association; change `List()` sort to `last_used_at DESC NULLS LAST, name ASC`
- `ArtistRepository` interface (usecase): drop `ArtistLink` from params; `CommissionArtistRepository` and `ArtpieceArtistRepository` interfaces both gain `UpdateLastUsedAt()`
- Commission and artpiece usecases: call `artistRepo.UpdateLastUsedAt()` after validating artist ownership on create/update
- `commission_repository.go`, `artpiece_repository.go`: `Preload("Artist")` → `Preload("Artist.Links")` so primary link is accessible in responses
- Handler responses for commissions/artpieces: `artist_link` field stays as `string | null` (server resolves to primary link URL); no change to downstream consumers

**Frontend**
- `Artist` type (`useArtists.ts`): drop `artist_link`, add `links: { id: string; url: string; is_primary: boolean }[]`
- `useCreateArtist.ts`, `useUpdateArtist.ts`: drop `artist_link` from payload types; add `links[]`
- `ArtistFormModal`: replace single link input with multi-link editor; frontend URL validation (`new URL()`) with `https://` auto-prepend on blur; inline error per invalid link
- `ArtistsList`: copy-link uses `links.find(l => l.is_primary)?.url`
- `ArtpieceDetailView`: derives artist link from `viewArtist.links.find(l => l.is_primary)?.url` (viewArtist sourced from useArtists cache, not artpiece preload — no preload change needed here)
- `CommissionFormModal`: constructs a minimal `Artist` object for combobox pre-population; `artist_link: commission.artist_link` becomes `links: []` (combobox only uses `id` and `name`)
- `CommissionsTable`, `DashboardPage`, `useCommissions.ts`, `useDashboard.ts`: no change — these consume `commission.artist_link: string | null` from the Commission response, which the server continues to resolve to the primary link URL

**Bruno collection**
- `create_artist.bru`, `update_artist.bru`: replace `artist_link` with `links[]`
- New bru files not needed (links are not a separate endpoint)

**Tests**
- `TestArtistRepository_List`: current alphabetical assertion must be updated to match new sort order
- `TestArtistRepository_Create`: drop `ArtistLink` assertion; add `Links` assertions
- New usecase tests: `last_used_at` set on create; `UpdateLastUsedAt` called on commission/artpiece assignment
