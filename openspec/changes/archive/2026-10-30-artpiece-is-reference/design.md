# Design

## Context

See proposal.md — Why.

Relevant existing state:
- `domain.Artpiece` is a GORM model. The repository `Update` method assembles a `map[string]interface{}` for the GORM `Updates()` call — zero values in a map are written, unlike struct-based updates.
- `ListArtpieceFilters` binds directly from Echo query params via `c.Bind`.
- `artpieceResponse` is a plain Go struct; fields with `omitempty` would suppress `false`, so the new field must not carry that tag.
- On the FE, the gallery fetches all artpieces in one unparameterized call and filters client-side. `useArtpiecesFilter` owns all filter state and syncs it to URL params. `ArtistCharacterFilter` is a fully controlled component — all value state lives in the caller.
- `parseFilterPayload` defensively coerces missing/invalid fields to their defaults (e.g. `match_mode` missing → `'all'`).
- `SavedFilterSnapshot` in `ArtpiecesPage` tracks the last-saved filter state to show the "Saved as …" confirmation.

## Goals / Non-Goals

**Goals:**
- Add `is_reference` to the data model and all read/write paths.
- Gallery toggle filters client-side like all other gallery filters.
- Saved filter round-trips the flag correctly (including absent → false coercion).
- Filter button summary and Clear all both account for the new toggle.

**Non-Goals:**
- Server-side gallery filtering via `?is_reference=` (BE capability is built, gallery does not use it in v1).
- Bulk mark-as-reference.
- Any reference-specific view, badge, or metadata beyond the toggle.

## Decisions

### D1: `*bool` in `UpdateArtpieceParams`, plain `bool` in `CreateArtpieceParams`

Update must distinguish "field omitted" from "field set to false". A `*bool` pointer solves this: nil means omitted (leave unchanged), non-nil means set. Create has no ambiguity — absent `is_reference` in the JSON body maps to `false`, which is the correct default. The handler struct `updateArtpieceRequest` also uses `*bool json:"is_reference"` for the same reason.

In the repository `Update` map: `if params.IsReference != nil { updates["is_reference"] = *params.IsReference }`. Using a map (not a GORM struct update) means `false` is written correctly — GORM only skips zero values in struct-based updates.

Alternative considered: plain `bool` on update with a separate `ClearIsReference bool` field. Rejected as unnecessarily complex when a pointer is idiomatic for optional fields in this codebase.

### D2: Client-side gallery filtering (fetch stays unparameterized)

The gallery already fetches all artpieces and filters client-side. Adding `is_reference` as a client-side filter keeps the architecture uniform, avoids cache-key complexity (`ARTPIECES_QUERY_KEY` stays a fixed tuple), and requires no changes to invalidation call sites. The BE query param is built for completeness but not called by the gallery.

### D3: `isReference` state lives in `useArtpiecesFilter`, URL param `ref`

All other gallery filter state (artist, characters, match mode) lives in `useArtpiecesFilter` via URL params. `isReference` belongs there too — it is a plain attribute filter, no different from the others. URL param chosen: `ref` (short, unambiguous in this domain). Absent → false (same defensive coercion pattern as `match_mode`).

### D4: `is_reference` in `FilterPayload`, coerced in `parseFilterPayload`

`parseFilterPayload` already defensively coerces missing/bad values. Adding `is_reference: p.is_reference === true` (strict equality against `true`) means any absent, null, or non-boolean value becomes `false`. This handles the forward-compat case where old saved filters have no `is_reference` field.

`filterToHref` in `DashboardPage` appends `ref=true` only when `payload.is_reference` is true (false is the default; no param needed).

### D5: `ArtistCharacterFilter` gains optional `isReference` / `onIsReferenceChange` props

Making the new props optional preserves the existing usage of `ArtistCharacterFilter` in the artpiece attach modal and any future caller that doesn't need the reference toggle. `buildSummaryLabel` and the Clear all handler both guard on prop presence.

## Risks / Trade-offs

- **GORM zero-value false**: If `is_reference` is ever added to the update map unconditionally as a plain `bool`, GORM would silently zero it to false on every update that doesn't intend to change it. Mitigation: the `*bool` pattern in params makes this structurally impossible — the map key is only added when the pointer is non-nil.
- **`omitempty` on response field**: If `omitempty` is accidentally added to `artpieceResponse.IsReference`, `false` will be omitted from the JSON. Mitigation: the field must be `json:"is_reference"` with no `omitempty`, and the response tests assert the false case explicitly.
- **Old saved filters**: Any saved filter created before this change has no `is_reference` in its payload. `parseFilterPayload`'s `=== true` coercion handles this correctly. No migration needed.

## Migration Plan

1. DB migration: `ALTER TABLE artpieces ADD COLUMN is_reference boolean NOT NULL DEFAULT false`. Non-destructive; all existing rows become `is_reference = false`.
2. Deploy backend with the new column and updated response. No breaking change — field is additive.
3. Deploy frontend. No coordination required; new field appears in responses before the frontend reads it.
4. Rollback: remove the FE changes; the column remains harmlessly on the table.
