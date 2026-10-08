## Why

Users need a way to save and recall filter combinations (artist + characters + match mode) so they don't have to re-select them each time. This is a standalone CRUD feature with no dependencies on other in-progress work.

## What Changes

- New `saved_filters` table storing per-user filter snapshots with an optional thumbnail path and an opaque JSON payload
- Five new REST endpoints: create, list, get, patch, delete saved filters
- No changes to existing endpoints or data models

## Capabilities

### New Capabilities

- `saved-filters`: Full CRUD for user-owned saved filter records — create, list, get, update (name / thumbnail / payload), delete

### Modified Capabilities

_(none)_

## Impact

- **Backend**: new migration, new domain type, repository, usecase, handler, route registrations
- **API**: five new authenticated endpoints under `/saved_filters`
- **Frontend**: none (backend-only proposal; FE work is a separate change)
- **Dependencies**: no new packages — `json.RawMessage` from stdlib handles the jsonb field
