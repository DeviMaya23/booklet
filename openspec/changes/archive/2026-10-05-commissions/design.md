## Context

The app already has `artpieces`, `artists`, `characters`, and `files` with a well-established layered architecture: domain model → repository interface → repository implementation → usecase → handler. The `artpiece-files` relationship (`files.artpiece_id` nullable FK, attach/detach/replace endpoints) is the direct structural precedent for `commission-artpiece-attachment`. The `artpiece_characters` many2many join table is the direct precedent for `commission_characters`.

`domain.Character` currently uses GORM soft-delete (`gorm.DeletedAt`). This will be removed as a separate concern; the design assumes hard-delete for characters by the time this lands.

## Goals / Non-Goals

**Goals:**
- Introduce `commissions` as a first-class entity with full CRUD.
- Allow commissions to be linked to characters and to artpieces.
- Expose `commission_id` on existing artpiece responses without breaking them.
- Follow the exact same patterns already established for `artpiece-files` (conflict checks, bulk operations, transactional replace).

**Non-Goals:**
- Status transition enforcement (no state machine; any valid status value is always settable).
- `commission_characters` read/write endpoints (table exists, populated at create/update, not independently managed).
- Frontend changes.
- Limbo / pre-waitlist state.

## Decisions

### D1: `commission_id` lives on `artpieces` as a nullable FK (not a join table)

An artpiece belongs to at most one commission at a time. A join table would allow many, which is wrong for this domain. Nullable FK mirrors `files.artpiece_id` exactly — same constraint shape, same ON DELETE SET NULL behaviour.

**Alternative considered**: join table `commission_artpieces`. Rejected because it allows N commissions per artpiece, which is not the intended model.

### D2: Attach conflict policy — reject, not overwrite

If an artpiece being attached already has a `commission_id` pointing to a *different* commission, the request is rejected with `ErrArtpieceAlreadyAttached`. This is consistent with how `AttachFile` handles `ErrFileAlreadyAttached`.

The PUT full-replace endpoint follows the same check: any incoming artpiece ID that already belongs to a different commission is rejected. The caller is expected to detach first.

**Alternative considered**: silent overwrite. Rejected to stay consistent with the file precedent and to avoid surprising data loss.

### D3: Commission GET returns artpiece summaries (id + cover thumbnail), not full artpiece objects

The commission detail view needs to show which artpieces are linked and display their thumbnails. Returning full `artpieceResponse` objects (with all files, characters, etc.) would be over-fetching. A slim `{ id, thumbnail_url }` shape is sufficient and mirrors how the artpiece list returns `thumbnail_url` via the `cover_file_id → CoverFile → ThumbnailR2Path` chain.

The `Commission` domain model will preload `Artpieces` with their `CoverFile` (lean preload, not full artpiece load).

### D4: `status` enforced at app layer only, stored as `text`

No DB CHECK constraint. Valid values (`waitlist`, `wip`, `done`) are defined as Go constants in the domain or usecase package. A CHECK constraint can be added later without a data migration since existing values will already satisfy it.

### D5: Commission CRUD follows the same usecase constructor pattern

`CommissionUsecase` takes `CommissionRepository`, `ArtistRepository` (for artist ownership check on create/update), `CharacterRepository`, `ArtpieceRepository` (for artpiece ownership + commission_id updates), and `Transactor`. No new infrastructure dependencies.

### D6: `updated_at` on artpieces is touched by attachment operations

`UPDATE artpieces SET commission_id = $1` will bump GORM's `updated_at` on artpiece rows. Accepted — commission attachment is a meaningful change to an artpiece record.

## Risks / Trade-offs

- **Loose artist consistency between commission and artpiece**: A commission has `artist_id`, an artpiece has `artist_id`. They can diverge. No constraint enforced. → Accepted trade-off per proposal; a future migration can add a consistency check if needed.
- **No status transition guard**: Any status is directly settable. → Accepted; transition rules are out of scope and can be layered on later.
- **Cover thumbnail presigning in commission GET adds N presign calls**: One presign per attached artpiece's cover thumbnail. For large commissions this is a series of S3 presign calls. → Low risk in practice (commissions rarely have >10 artpieces), acceptable for now.

## Migration Plan

Three migrations, applied in order:

1. `CREATE TABLE commissions` + `CREATE TABLE commission_characters`
2. `ALTER TABLE artpieces ADD COLUMN commission_id UUID NULL REFERENCES commissions(id) ON DELETE SET NULL`

Down migrations:
1. `ALTER TABLE artpieces DROP COLUMN commission_id`
2. `DROP TABLE commission_characters; DROP TABLE commissions`

Both are non-destructive additions; rollback is safe as long as no data has been written to the new tables/column.
