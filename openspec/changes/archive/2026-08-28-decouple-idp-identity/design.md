## Context

Every table with a `user_id` column uses the Kinde-issued subject (e.g. `kp_abc123`) as a TEXT foreign key. This makes the database schema structurally dependent on the IdP's identifier format. The auth middleware currently extracts `claims.Subject` from the JWT and puts it directly into the Echo context as the user's identity, making the Kinde subject the single identity token flowing through every layer of the application.

Two external systems consume the Kinde subject today:
- **Kinde itself** (via bookleaf's `DeleteAccount` endpoint)
- **Bookleaf's internal API** (`GetPublicFolders`, `DeleteAccount`) — bookleaf has made the same migration but the two apps use Kinde subject as their cross-system bridge rather than syncing UUIDs, which is acceptable at current scale.

There is no live production data. Drop-and-recreate is the migration strategy.

## Goals / Non-Goals

**Goals:**
- Replace `users.id` (TEXT Kinde subject) with an app-generated UUID
- Demote the Kinde subject to `users.idp_subject` — a lookup column used only at the IdP/cross-system boundary
- Change `user_id` FK columns on `characters`, `images`, `pending_uploads` from TEXT to UUID
- Preserve the existing behavior of all endpoints; no API contract changes

**Non-Goals:**
- Syncing UUIDs between booklet and bookleaf — cross-system calls continue using `idp_subject`
- Changing the Kinde JWT validation flow
- Any frontend changes (the frontend is empty; no impact)

## Decisions

### D1 — UUID generation is app-side, not DB-side

`users.id` uses `uuid.New()` in Go at `GetOrCreate` time, not a Postgres `gen_random_uuid()` default. This matches the existing pattern for `characters`, `images`, and `pending_uploads`, keeps generation visible in application code, and avoids a GORM-vs-DB-default mismatch.

**Alternative considered**: Postgres `gen_random_uuid()` default. Rejected because GORM's `Clauses(clause.OnConflict{DoNothing: true})` pattern reads back the row after insert — the UUID would be assigned regardless — but generates it in Go keeps the value in scope for logging without a second query.

### D2 — Both app UUID and idp_subject stored in Echo context after auth

The middleware already has both `user.ID` and `claims.Subject` in scope when provisioning completes. Both are written to context:

```
"authenticatedUserID"    → user.ID (uuid.UUID)   — used by all internal routing
"authenticatedIDPSubject" → claims.Subject (string) — used only at bookleaf call sites
```

Handlers that call bookleaf-facing usecases (`DeleteMe`, `ListFolders`) extract `idpSubject` from context and pass it to their usecase.

**Alternative considered**: Usecases fetch the user by UUID to retrieve `idp_subject` when needed. Rejected because it adds a DB round-trip per bookleaf call and gives the folder usecase an unneeded UserRepository dependency.

### D3 — `MarkPendingDeletion` takes both `userID uuid.UUID` and `idpSubject string`

The usecase must write to the DB (UUID) and call bookleaf (Kinde subject). Both values are available at the handler level from context. Passing both as arguments is the simplest path.

**Alternative considered**: Usecase fetches user inside `MarkPendingDeletion` to get `idp_subject`. Same rejection as D2 — extra round-trip, and the value is already in the caller's hands.

### D4 — `AuthenticatedUserIDFromContext` return type changes from `string` to `uuid.UUID`

The return type change causes compile errors at every call site, which is intentional — the compiler audits the blast radius for us. All ~10 handler call sites must be updated explicitly.

### D5 — Migration is drop-and-recreate

No live data. The simplest approach is a new migration that drops the four affected tables and recreates them with the correct schema. Previous migrations remain as historical record but the effective schema starts fresh from migration 9.

**Alternative considered**: Alter-column migrations with backfill. Rejected — unnecessary complexity when there's no data to preserve.

## Risks / Trade-offs

- **Cross-system coupling persists at IdP level** — bookleaf calls still use Kinde subject. An IdP migration still requires updating `idp_subject` values and coordinating with bookleaf. This is a known, accepted trade-off at current scale.
- **R2 key path shape changes for new uploads** — keys were `users/kp_abc123/images/...`, will become `users/<uuid>/images/...`. Existing objects in R2 are unaffected (paths stored in DB rows); only new uploads use the new shape.
- **`AuthenticatedUserIDFromContext` return type is a breaking change** — compile errors at all call sites are the safety net. No runtime risk if the build passes.

## Migration Plan

1. Write migration `000009` that drops `pending_uploads`, `image_characters`, `images`, `character_folders`, `characters`, then `users`, and recreates them with the new schema
2. No rollback script needed (no live data; down migration is informational only)
3. Deploy: `make migrate` as usual
