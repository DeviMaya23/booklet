
# Project Information

## What is Booklet
Booklet is a single-user art library and commission tracker, built for the
owner's personal use.

- **Art library first**: files, artpieces, characters, artists, search.
  Commission management is a workflow layer on top that *produces* artpieces.
- **Artpiece** is the browsable unit: one distinct thing you'd show on its own.
  Has an artist, characters, and one or more files.
- **File** is a version/format of an artpiece (png, psd, clip, mp4...).
  `artpiece_id` is nullable: unassigned files live in the "Inbox".
- **Commission** is the transaction record (artist, price, paid, status:
  waitlist/wip/done). One commission → many artpieces.

## Bookleaf
Bookleaf is a **separate, live app** (own codebase and deployment): an image
moodboard/pinboard. Booklet shows Bookleaf moodboards on character pages by
linking to Bookleaf's *public folders*.

- Bookleaf's code and data are **not** in this repo. Ask me if you need specific information of its repo.
- Integration is API-only. Booklet stores folder IDs and validates them against
  Bookleaf's API.
- Both apps share one Kinde user pool but are separate applications.
- If a task seems to need a Bookleaf change, stop and flag it rather than
  working around it.

## Reference Files

Read these before the work type they cover — not all at once, only the relevant one:

| When you're about to… | Read |
|---|---|
| Orient to the project stack or pick a library | `STACK.md` |
| Place a new file, create a feature, or touch an existing structure | `ARCHITECTURE.md` |
| Write any frontend or backend code | `CONVENTIONS.md` |
| Write or plan any tests | `TESTING.md` |

`frontend/CLAUDE.md` and `backend/CLAUDE.md` are loaded automatically — they contain critical constraints for each side that override common defaults.


## Decision Boundaries

Before introducing a new pattern, new abstraction, new dependency, 
or new layer boundary not already present in the codebase, stop and 
propose it to me first. Do not proceed until confirmed.

## Mid-Implementation Adaptations

If, while implementing a task, you discover that something not covered by the design is needed — a structural workaround, a bridging type, a scope change, an unanticipated dependency — stop and inform me before proceeding. Describe what you found, why the adaptation is needed, and what you're proposing. Do not implement it silently.


## OpenSpec Proposals

- Before starting a new proposal, pull the latest main branch and checkout from there
- If current working branch has unstaged or uncommitted changes, confirm with me before proceeding
- Branch name format: `feat/<spec-name-here>`
- Generate each artifact during proposal step by step. Confirm with me before moving on to the next one.

### Testing in proposals

- Always plan for unit tests on the usecase and handler layers
- Always plan for integration tests on the repository layer — do not write unit tests for SQL repositories
- This applies to every repository touched by the change, including ones not modified but whose output the feature depends on (e.g. a repository that must preload a newly added association)
- Follow the unit testing rules in CONVENTIONS.md — they define what scenarios are worth writing and what test doubles to use
- Do not default to one success + one failure per function; write only the scenarios that have a reason to exist per the conventions

#### Assertion quality
- If a function returns a result, assert the result — not just the error
- Failure scenarios must assert the specific error type or message, not just that an error occurred

### Backend handler input processing

- Prefer `c.Bind` + `c.Validate` over manual parsing and validation for body and query params
- Use `uuid.Parse` manually only for path params, where `c.Bind` does not apply

### Others to keep in mind during proposals
- Every surface listed in the Impact section must have at least one corresponding spec scenario. The scenario must name the specific data path (e.g. "commission list returns Artist.Links preloaded and surfaces the primary link"), not just the surface. Impact without a scenario is documentation; Impact with a scenario is a contract.
- When a change modifies a shared contract (a function signature, API endpoint, hook, etc.), grep for every call site of that symbol across all layers (backend, frontend, extension) before finalizing Impact/Capabilities/tasks — do not rely on a named flow (e.g. "the upload flow") to be exhaustive, as parallel entry points (e.g. drag-and-drop vs. modal vs. batch upload) commonly funnel into the same shared function and are easy to miss.
- When creating tasks for a new endpoint, always include a bruno file creation.
- On any BE development, include a task to run golang-ci lint at the end, and fix whatever issue arises.
- On any FE development, include a task to run npm run build and npm run lint at the end, and fix whatever issue arises.
