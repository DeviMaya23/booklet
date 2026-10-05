## Why

Commissions are a core tracking need — users want to log what they ordered, from whom, for how much, and link the resulting artpieces once delivered. There is no model for this today; it is tracked externally or not at all.

## What Changes

- New `commissions` table: tracks per-user commission records with artist, status, price, payment, and notes fields.
- New `commission_characters` join table: links commissions to characters (many-to-many), for tagging which characters a commission was ordered for.
- New column `artpieces.commission_id`: nullable FK linking an artpiece back to the commission it was delivered for.
- New Commission CRUD endpoints: create (optionally attaching artpieces in the same call), get by ID, list, update, delete.
- New commission↔artpiece attachment endpoints: POST (bulk attach), DELETE (bulk detach), PUT (full replace).
- Commission GET response includes per-artpiece summary: id + cover thumbnail URL.
- Artpiece GET/list responses gain a `commission_id` field.

## Capabilities

### New Capabilities

- `commission-management`: Commission CRUD. Covers the `commissions` table, `commission_characters` join table, character ownership validation on create/update, and status as an app-layer enum (`waitlist`, `wip`, `done`).
- `commission-artpiece-attachment`: Managing the commission↔artpiece link. Covers attach (POST, bulk), detach (DELETE, bulk), full replace (PUT), and the ownership + conflict validation rules. Also covers the artpiece summary shape returned from commission GET.

### Modified Capabilities

- `artpiece-management`: Artpieces gain a nullable `commission_id` column. The artpiece response (GET and list) must expose `commission_id`.

## Impact

- **Migrations**: new `commissions` table, new `commission_characters` table, `ALTER TABLE artpieces ADD COLUMN commission_id`.
- **Domain**: new `Commission` model; `domain.Artpiece` gains `CommissionID *uuid.UUID`.
- **Backend layers**: new commission repository, usecase, and handler; artpiece response updated.
- **Routes**: commission CRUD + attachment sub-routes registered in `main.go`.
- **Bruno**: new request files for all commission endpoints.
- **Tests**: usecase unit tests for commission-management and commission-artpiece-attachment; artpiece usecase tests updated for any changed behavior.
- **No frontend changes** in this proposal.
