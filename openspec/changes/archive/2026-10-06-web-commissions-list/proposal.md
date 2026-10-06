## Why

The app has no UI for commissions. The backend CRUD already exists, but the commissions data is inaccessible from the frontend. This adds the list screen — the primary surface for scanning and managing commission status, payment state, and artist contact history.

## What Changes

- New migration `000027` adds `last_contacted_at TIMESTAMPTZ` to the `commissions` table
- Commission API response gains `artist_link` (from the joined `Artist`) and `last_contacted_at`
- New `PATCH /commissions/:id` endpoint for partial updates of inline-editable fields: `status`, `paid`, `paid_date`, `last_contacted_at`
- New frontend page at `/app/commissions` with a table layout, toggle between "Waitlist / In Progress" and "Finished" views
- Inline autosave on Status (dropdown), Paid (checkbox), Paid Date (datepicker), and Last Contact stamp action
- New shadcn components installed: `table`, `checkbox`, `calendar`, `popover`

## Capabilities

### New Capabilities

- `web-commissions-list`: Commissions list page — table layout, toggle view states, inline-editable cells with autosave, Last Contact stamp, Time Taken computed column

### Modified Capabilities

- `commission-management`: New `PATCH /commissions/:id` endpoint; `last_contacted_at` field added to domain, response, and DB schema; `artist_link` added to response

## Impact

- **Backend**: new migration, domain model update, new PATCH handler + usecase method, updated response struct
- **Frontend**: new page + route + sidebar entry, new `features/commissions/` folder, new shadcn components
- **Spec**: `commission-management` spec updated to document PATCH endpoint and new fields
