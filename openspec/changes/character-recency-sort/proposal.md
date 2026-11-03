# Proposal

## Why

Artists in the lookup combobox surface most-recently-used first (`artist-recency-sort`), but characters in the pickers are ordered oldest-created-last with no regard for use. Picking characters for a new artpiece means scrolling past ones you haven't touched in months. Characters should behave like artists in the picker context.

## What Changes

- Add `characters.last_used_at timestamptz NULL`, backfilled from `created_at` so existing order is preserved at rollout.
- Bump `last_used_at` for every character in a non-empty `character_ids` on artpiece create/update and commission create/update, and set it when a character is created (mirrors the artist rule).
- `GET /characters` accepts an optional `sort` query param. `sort=recent` orders by `last_used_at DESC NULLS LAST, name ASC`. Absent, the order stays `created_at DESC`. An unrecognised value is rejected with 422. `last_used_at` is not exposed in the response.
- The three character pickers (`ArtpieceFormModal`, `ArtpieceDetailView`, `ArtistCharacterFilter`) request `sort=recent`. The Characters page and id-lookup consumers keep the default order.
- Artpiece and commission create/update mutations invalidate the characters query so picker order is fresh after a save.
- `ArtpieceUsecase.Update` validates character ownership itself (it currently relies on the repository), so the bump has a validated ID set to work from.

**Non-goal:** a shared `CharacterCombobox` wrapper replacing the three `TokenInput` + `useCharacters` call sites. The duplication is accepted for now.

## Capabilities

### New Capabilities
- `character-recency-sort`: the `last_used_at` tracking column on characters, when it is bumped, and the opt-in recency ordering of the character list and its picker consumers.

### Modified Capabilities
- `character-management`: "List characters" gains the optional `sort` query parameter; the data model gains the `last_used_at` column (not exposed in responses).
- `artpiece-management`: create/update bump the `last_used_at` of assigned characters; update validates character ownership in the usecase.
- `commission-management`: create/update bump the `last_used_at` of assigned characters.
- `web-artpiece-create-modal`: character suggestions are shown in recency order.
- `web-artpiece-detail`: character suggestions in edit mode are shown in recency order.
- `web-artist-character-filter`: character suggestions are shown in recency order.

## Impact

- **Backend:** new migration `000037`; `domain/character.go`; `usecase/character_repository.go` (`ListCharacterFilters.Sort`, `UpdateLastUsedAt`); `repository/character_repository.go`; `handler/character_handler.go` (validate `sort`); `usecase/character_usecase.go`, `artpiece_usecase.go`, `commission_usecase.go`; every test double implementing the character repository interface.
- **Frontend:** `useCharacters` (optional `sort`, query key `['characters', 'recent']`); `ArtpieceFormModal`, `ArtpieceDetailView`, `ArtistCharacterFilter`; `useCreateArtpiece`, `useUpdateArtpiece`, `useCreateCommission`, `useUpdateCommission`. `CommissionFormModal` has no character picker of its own and only embeds `ArtistCharacterFilter`.
- **Unchanged consumers:** `CharactersPage`, `ChangeThumbnailModal`, `useArtpiecesFilter` keep the default order.
- **Docs/tooling:** `DATABASE.md`, Bruno `collection/characters/list_characters.bru`.
- **Bookleaf:** no impact; the integration is folder IDs only.
