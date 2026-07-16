## Context

The images module introduces a new top-level resource scoped to authenticated users. It mirrors the character module's layered architecture (domain → usecase → repository → handler) and telemetry patterns exactly. The only structural novelty is the `image_characters` join table, which requires GORM Many2Many associations — the first use of that pattern in this codebase.

Image creation (which requires R2 URL generation and thumbnail processing) is deferred to a later proposal. This change implements the remaining four operations: list, get, update, and delete.

## Goals / Non-Goals

**Goals:**
- Introduce the `images` table and `image_characters` join table
- Implement `GET /images`, `GET /images/:id`, `PATCH /images/:id`, `DELETE /images/:id`
- Embed associated character `id` and `name` in every image response
- Support full-replace of character associations on PATCH
- Enforce character ownership before writing associations

**Non-Goals:**
- Image creation endpoint (deferred — requires R2 and thumbnail flow)
- Pagination or filtering on list
- Public/private visibility on images
- R2 cleanup on delete (deferred to creation proposal)

## Decisions

### 1. GORM Many2Many for the join table

**Decision:** `domain.Image` carries a `Characters []Character` field tagged `gorm:"many2many:image_characters;"`. The repository uses `Preload("Characters")` on all reads and `Association("Characters").Replace(...)` for full-replace writes.

**Alternatives considered:**
- *Manual queries* — two raw SQL queries assembled in Go. Consistent with current single-table style but adds ~30 lines of plumbing with no benefit at this scale.

**Rationale:** The join table is exactly the use case GORM Many2Many is designed for. `Association.Replace` handles the transactional delete-then-insert atomically. The pattern is new to the codebase but isolated to the image repository — it introduces no cross-cutting change.

---

### 2. Hard delete for images

**Decision:** `domain.Image` has no `DeletedAt` field. `DELETE /images/:id` issues a hard delete.

**Alternatives considered:**
- *Soft delete* (consistent with character) — keeps history. Rejected: images are tied to R2 assets which are permanently deleted. Soft-deleted rows pointing at non-existent R2 paths have no operational value. The mismatch with character's soft delete is deliberate and acceptable.

---

### 3. Full-replace semantics for character associations

**Decision:** `PATCH /images/:id` accepts an optional `character_ids []string` field. When present (even as an empty array), the entire `image_characters` set for that image is replaced. When absent, associations are untouched.

**Rationale:** Simplest contract. The client always knows the full desired state. Add/remove individually would require separate endpoints and round-trips.

---

### 4. Character ownership validation before association replace

**Decision:** Before calling `Association.Replace`, the repository validates that every supplied character ID exists in the `characters` table scoped to `user_id`. Any mismatch returns an error that maps to 422 at the handler.

**Rationale:** Without this check, a user could link another user's character to their image. The character row would be unreachable via normal APIs but the foreign key would exist.

---

### 5. Nullable fields

**Decision:** Only `image_r2_path` is NOT NULL. All other text fields (`title`, `thumbnail_r2_path`, `artist_name`, `artist_link`, `notes`) are nullable and represented as `*string` in the domain struct and update params.

**Rationale:** Image creation is deferred. Keeping `image_r2_path` NOT NULL enforces that a valid path is always supplied at creation time, even though the creation endpoint isn't in this proposal.

## Risks / Trade-offs

- **GORM Many2Many is new to the codebase** → Contained to `image_repository.go`. If the pattern proves problematic, it can be replaced with manual queries without touching other layers.
- **Hard delete with no R2 cleanup** → Acceptable for now. The creation proposal will define the full delete lifecycle (R2 + DB). Until then, no images exist in production.
- **No pagination on list** → Acceptable for early stage. Images are user-scoped so the ceiling is per-user volume, not global.
