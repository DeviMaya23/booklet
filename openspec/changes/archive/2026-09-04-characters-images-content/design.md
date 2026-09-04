## Context

The app shell is in place with routes for `/app/characters` and `/app/images` rendering empty stubs. Backend list and delete endpoints for both resources exist and are stable. The stack is React + Vite + TypeScript, Tailwind, shadcn/ui (Base UI primitives), `@tanstack/react-query` v5, `lucide-react`, and `sonner`.

## Goals / Non-Goals

**Goals:**
- Render browsable, searchable grids of character and image cards.
- Share a single `ResourceCard` component across both pages.
- Support per-card delete with a confirmation dialog.
- Client-side search filtering on the fetched list.

**Non-Goals:**
- Create/upload flows (+ New button is a no-op).
- Detail view / resource modal (card click is a no-op).
- Pagination or infinite scroll.
- Artists page.
- Server-side search.

## Decisions

### 1. One shared `ResourceCard`, two feature-specific grids

Characters and images have near-identical card shapes (nullable image, string label). A single `ResourceCard` in `src/components/` accepts `imageUrl?: string | null` and `label: string` and handles both. Each feature owns its own grid component (`CharactersGrid`, `ImagesGrid`) that maps API data to card props — the label mapping (`title ?? "Untitled"` for images, `name` for characters) lives at the grid level, not inside the card.

Alternative considered: separate `CharacterCard` and `ImageCard` — rejected as premature duplication given identical structure.

### 2. Client-side search

The full list is fetched once; the search bar filters in JS on `name` (characters) or `title` (images). Null titles are treated as empty string for filtering purposes — they do not match any query, including the literal string `"Untitled"`. The "Untitled" display label is cosmetic only.

Alternative considered: server-side `?q=` — deferred. The backend supports it, but client-side is simpler today and sufficient for expected collection sizes. Switching later requires only changing the query key and adding a debounced param.

### 3. Feature folder structure

Each resource gets a feature folder following the existing `features/auth` pattern:

```
src/features/
  characters/
    api/
      useCharacters.ts       ← React Query list hook
      useDeleteCharacter.ts  ← React Query mutation hook
    components/
      CharactersGrid.tsx
  images/
    api/
      useImages.ts
      useDeleteImage.ts
    components/
      ImagesGrid.tsx
src/components/
  ResourceCard.tsx           ← shared card
```

### 4. Delete flow

Clicking Delete in the `(...)` menu opens a shadcn `<AlertDialog>` (destructive confirmation). On confirm, the mutation fires `DELETE /characters/:id` or `DELETE /images/:id`, invalidates the list query, and shows a `sonner` toast on success or error. No optimistic removal — the list refetches after the mutation settles to avoid a flash of stale state.

A `<Dialog>` component does not exist in `src/components/ui/` yet; it will be added via the shadcn CLI (`npx shadcn add alert-dialog`).

### 5. Auth token propagation

Both hooks call `apiFetch`, which takes a `getToken` callback. Hooks receive `getToken` from `useKindeAuth()` at the call site (inside the page or grid component), consistent with how the rest of the app accesses Kinde.

### 6. No-image placeholder

When `imageUrl` is null or absent, `ResourceCard` renders a muted-background square with a centered `lucide-react` `ImageOff` icon. No external asset fetch required.

## Risks / Trade-offs

- **Presigned URL expiry**: Thumbnail/avatar URLs are short-lived. If a user leaves the tab open, cards may show broken images on return. Mitigation: React Query's `staleTime` defaults mean the list refetches on window focus, regenerating fresh URLs. No extra work needed today.
- **Client-side search scales poorly**: Filtering in JS on a large list will lag. Mitigation: acceptable for now; switch to server-side `?q=` when collections grow.
- **"Untitled" not findable**: Images with null titles are invisible to search. This is intentional per product decision, but could surprise users. No mitigation planned.
