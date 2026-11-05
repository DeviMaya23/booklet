# Proposal

## Why

Tags can be created, attached and picked, and the backend already supports rename and delete, but there is no place in the app to manage them. Typos, near-duplicates (`sketch` / `sketches`) and abandoned tags accumulate with no way to clean them up short of calling the API by hand. A dedicated Tags page under MANAGE closes that gap.

## What Changes

- Add a **Tags** page at `/app/tags`, reached from a new "Tags" item in the sidebar's MANAGE group, placed after Artists.
- The page lists the user's tags with a client-side search box and a sort dropdown ("Recently used" / "Name (A–Z)"), each row showing the tag name and "N artpieces →" (linking to the gallery filtered by that tag), an "Unused" badge when the count is 0, and edit and delete actions. An empty state shows when the user has no tags.
- **Inline rename**: the row turns into an input (Enter saves, Esc cancels). When the new name matches another existing tag (case-insensitive), the app shows a merge confirmation instead of saving.
- **Delete** with a confirmation dialog that states how many artpieces will lose the tag.
- **Merge**: new `POST /tags/:id/merge` with `{target_id}`. In one transaction it moves the source tag's artpieces to the target (artpieces that already have both keep one), then deletes the source. The target's `last_used_at` is left unchanged.
- `GET /tags` gains an opt-in `with_usage=true` query param. When set, each tag also contains `artpiece_count`. Without it the response is unchanged, so the existing tag pickers and filters pay no extra query cost. A new `domain.TagWithUsage` type carries the count.
- Sorting reuses the existing `GET /tags?sort=recent`; no new sort value is added.

Not changed: `PATCH /tags/:id` stays a plain rename and still returns 409 on a name collision, so a rename can never destroy data. Bookleaf is not involved.

## Capabilities

### New Capabilities
- `web-tag-management`: the Tags page: list, search, sort, usage display, inline rename, delete and merge confirmations, empty state, and link to the filtered gallery.

### Modified Capabilities
- `tag-management`: `GET /tags` gains the opt-in `with_usage` param and `artpiece_count` field (the current requirement says each tag contains only `id` and `name`); adds a merge-tags requirement.
- `web-nav-sidebar`: the MANAGE group now contains Commissions, Artists and Tags, and the icon assignments gain Tags.

## Impact

- **Backend**: `domain/tag.go` (new `TagWithUsage`), `usecase/tag_usecase.go`, `usecase/tag_repository.go` (`ListTagFilters.WithUsage`, merge method and errors), `handler/tag_handler.go` (merge handler, usage-aware list response), `repository/tag_repository.go` (usage count join, transactional merge), route registration in `cmd/server/main.go`. Merge runs inside the existing `Transactor`.
- **Bruno**: new `collection/tags/merge_tag.bru`; the existing `list_tags.bru` gains a with-usage request.
- **Frontend**: `components/AppSidebar.tsx`, `router.tsx`, new `pages/TagsPage.tsx`, new tag-row, delete and merge components and the update, delete and merge hooks under `features/tags/`, and `useTags` gains a `withUsage` option (included in the query key). `TagPicker`, `ArtistCharacterFilter` and `SaveFilterModal` consume `useTags` without `withUsage` and are unaffected, but are call sites to re-check.
- **Tests**: unit tests for the tag usecase and handler; integration tests for the tag repository (usage count including zero and user scoping, count with `sort=recent`, merge with overlapping artpieces, cross-user rejection); frontend tests colocated with each new hook and component.
- **Dependencies**: none new.
