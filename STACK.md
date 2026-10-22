# Stack

> Agent-agnostic. For critical constraints (what NOT to do), see `frontend/CLAUDE.md` and `backend/CLAUDE.md`.

## Backend
- Language: Go
- Framework: Echo
- ORM: GORM
- Database: PostgreSQL
- Auth: Kinde
- Storage: Cloudflare R2
- Logging: Zap
- Observability: OpenTelemetry

## Frontend
- Framework: React + Vite + TypeScript
- Styling: Tailwind CSS
- UI primitives: `@base-ui/react`

## UI Components

- We use shadcn/ui with Base UI as the primitive layer (this is shadcn's default, not a custom override).
- shadcn CLI-generated components, file structure, and styling conventions apply as normal.
- API differences from Radix: use `onClick` not `onSelect` on `ContextMenuItem`. If unsure, check @base-ui-components/react docs, not Radix docs.
- Do NOT import from @radix-ui/* packages.

## Library Usage

### Frontend
Before self-implementing any non-trivial logic — layout algorithms, data
transformation pipelines, anything with meaningful edge cases or more than
~20 lines — you must:

1. Flag that you are about to self-implement something
2. Note any well-known libraries that solve the same problem
3. State your reasoning for or against each option
4. Wait for a decision before proceeding

This rule exists because the React/JS ecosystem has strong conventions
that are not always obvious. Do not assume self-implementation is preferred.

### Backend (Go)
Flag non-trivial self-implementations and surface library alternatives,
but a decision is not required before proceeding — note the trade-off
inline and continue unless the implementation touches a shared interface
or introduces a new architectural pattern (see Decision Boundaries in CLAUDE.md).
