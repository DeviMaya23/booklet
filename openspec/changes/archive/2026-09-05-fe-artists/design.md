## Context

The backend artist CRUD is fully implemented (`GET/POST/PATCH/DELETE /artists`). The route `/app/artists` is registered in `App.tsx` and `ArtistsPage.tsx` exists but returns null. Images and characters pages follow an established pattern: a page component owns search state, a query hook fetches the list, and a separate display component renders it. Artists follow the same page-level structure but differ in two meaningful ways: the display is a flat list (not a card grid), and the feature includes a create/edit modal with delete embedded in edit mode.

## Goals / Non-Goals

**Goals:**
- Implement a working artists list page consistent with the images/characters page structure
- Client-side search by name
- Create artist via modal
- Edit artist (name, link, blurb) via the same modal component
- Delete artist from within the edit modal, with a confirmation dialog
- Copy `artist_link` to clipboard via a link button; button is disabled when `artist_link` is null or empty

**Non-Goals:**
- Artist avatar/image support
- Server-side search (the `q` param exists on the backend but is unused here — client-side filtering is sufficient and matches the established pattern)
- Any change to the backend

## Decisions

### Decision: Row-based list, not a card grid

Images and characters use `ResourceCard` in a responsive grid because their primary content is visual (thumbnail/avatar). Artists are text-only entities. A flat row list with name on the left and action buttons on the right is a better fit and avoids forcing a card layout onto non-visual data.

Alternative considered: reuse `ResourceCard` with a text-only fallback (the `ImageOff` placeholder). Rejected — the card shape is semantically wrong for a list of named references, and the action surface (link + edit) doesn't map to the card's `(...)` pattern naturally.

### Decision: Single `ArtistFormModal` component with a mode prop

The create and edit flows share identical fields (name, link, blurb). A single modal component driven by an optional `artist` prop handles both: no `artist` → create mode, `artist` present → edit mode. The delete button renders only in edit mode.

The modal is rendered as a centered floating dialog using `@base-ui/react/dialog` — the same primitive already used by `sheet.tsx`. A new `ui/dialog.tsx` wrapper is added to provide styled `Dialog`/`DialogContent`/`DialogHeader`/`DialogFooter`/`DialogTitle` components, consistent with how `sheet.tsx` and `alert-dialog.tsx` wrap their respective primitives.

Alternative considered: two separate modal components. Rejected — it duplicates form logic with no meaningful separation of concerns.

### Decision: Delete lives inside the edit modal

The artist list row already has two inline action buttons (link, edit). Adding a third `(...)` dropdown with delete would make the row action-heavy for a text list. Placing delete at the bottom of the edit modal adds appropriate friction for a destructive action on reference data and keeps the row surface clean.

The delete button in the modal triggers an `AlertDialog` confirm before calling `DELETE /artists/:id`.

### Decision: Link button always visible, disabled when no link

A button that disappears based on data state can cause layout shift and makes the UI feel inconsistent as the list updates. The link button is always rendered; it is disabled (and non-interactive) when `artist_link` is null or empty string.

### Decision: Invalidate query on mutation success

Consistent with `useDeleteImage` and `useDeleteCharacter`: mutations call `queryClient.invalidateQueries` on success to refetch the list. No optimistic updates — the list is small and the round-trip is fast.

## Risks / Trade-offs

- **Name conflict on create/edit** → The backend returns 409. The modal must surface this as a field-level or inline error (not just a generic toast) so the user understands what went wrong. Toast-only error handling would be confusing here.
- **Clipboard API unavailability** → `navigator.clipboard.writeText` requires a secure context (HTTPS or localhost). In practice the app runs on HTTPS in production; no mitigation needed beyond the existing deployment assumption.

## File Structure

```
frontend/src/components/ui/
  dialog.tsx               — new: centered floating dialog wrapper over @base-ui/react/dialog

frontend/src/features/artists/
  api/
    useArtists.ts          — GET /artists query hook; exports Artist type and ARTISTS_QUERY_KEY
    useCreateArtist.ts     — POST /artists mutation hook
    useUpdateArtist.ts     — PATCH /artists/:id mutation hook
    useDeleteArtist.ts     — DELETE /artists/:id mutation hook
  components/
    ArtistsList.tsx        — row-based list; owns delete confirm dialog state
    ArtistFormModal.tsx    — create/edit modal; mode driven by presence of artist prop

frontend/src/pages/ArtistsPage.tsx  — search state, fetches list, opens modal
```
