## Context

The characters page already has list and delete working. All backend endpoints for create, update, and avatar upload are implemented. This change is frontend-only: wire the existing endpoints to a new modal, following the pattern already established by `ArtistFormModal`.

The avatar upload flow is a 2-step presigned URL sequence: the client calls `init` to get a presigned PUT URL, uploads the file directly to R2, then calls `complete` to commit. This flow runs after character creation, so create + avatar is inherently a multi-step sequence.

## Goals / Non-Goals

**Goals:**
- Ship a `CharacterFormModal` for create and edit with avatar upload, name, and notes
- Keep the create+avatar sequence atomic from the user's perspective (failure = rollback + error)
- Stay consistent with the `ArtistFormModal` pattern for dialog structure and footer layout

**Non-Goals:**
- Moodboard (folder linking) — deferred to a follow-up proposal
- `is_public` toggle — not exposed in the UI for now
- Any backend changes

## Decisions

### Avatar upload is local-preview only until Save

When the user picks an image file, the modal displays it immediately using `URL.createObjectURL(file)`. The file is not uploaded until the user presses Save. This keeps the flow simple and avoids initiating an upload the user might cancel.

**Alternative considered**: upload eagerly on file pick (fire init+upload immediately, complete on Save). Rejected because it generates orphaned pending upload records on Cancel and complicates cleanup.

### Create sequence: character first, avatar second — fail whole on avatar error

On Save with a new avatar:
1. `POST /characters` → get `characterID`
2. `POST /characters/:id/avatar/init` → get `{ id, upload_url }`
3. `PUT <upload_url>` (direct to R2, no auth header, `Content-Type` set to file mime type)
4. `POST /characters/:id/avatar/:uploadID/complete`
5. If steps 2–4 fail: `DELETE /characters/:id` to clean up, then surface error toast

The character must exist before avatar init because the endpoint takes `:id`. This means a small window where the character exists without an avatar if the avatar steps fail. The cleanup call handles this — the user sees a single "Failed to create character" toast and the modal stays open for retry.

On Save without a new avatar (create with no image): just `POST /characters`, done.

### Edit sequence: avatar upload is independent of metadata save

On Save in edit mode:
- If the user changed the avatar: run the init → upload → complete sequence first, then `PUT /characters/:id` for metadata. Both steps show a single pending state.
- If the user removed the avatar: call `DELETE /characters/:id/avatar`, then `PUT /characters/:id`.
- If only metadata changed: just `PUT /characters/:id`.

Avatar and metadata are separate concerns on the backend. Running avatar first means the PUT always sees the latest `avatar_r2_path` state, even though the response from PUT carries the presigned URL we'll display.

**Alternative considered**: run PUT first, then avatar. Rejected because on failure between steps the displayed data would be inconsistent (metadata saved, avatar not updated).

### `ResourceCard` gets an `onClick` prop — card body triggers edit

The `(...)` dropdown keeps its Delete option (quick delete without opening modal). The card body gets an `onClick` that opens the modal in edit mode. The `onEditClick` and `onClick` are the same handler at the `CharactersGrid` level.

### No `useDeleteAvatar` mutation hook

`DELETE /characters/:id/avatar` is only called from inside the modal's save handler (edit mode, avatar cleared). It doesn't need its own React Query mutation — a plain `apiFetch` call inside the submit handler is sufficient and avoids coupling the cache invalidation logic to an action that is always followed by a `PUT`.

### Avatar "remove" is a local state change only

In edit mode, the user can remove the current avatar by clicking an X on the preview. This sets a local `avatarCleared` boolean. The actual `DELETE /characters/:id/avatar` call only happens on Save. This is consistent with how new avatar picks work — nothing hits the network until Save.

## Risks / Trade-offs

**Create cleanup can fail**: If `POST /characters` succeeds but avatar upload fails, and the subsequent `DELETE /characters/:id` also fails (network drop), the user has an orphaned characterwith no avatar. They can delete it manually from the list. Acceptable risk given the rarity of back-to-back failures.

**Presigned URL expiry**: The presigned PUT URL from `init` expires in 15 minutes. If the user sits on the modal for 15+ minutes after picking an image, the upload will fail at step 3. This surfaces as an error toast; the user can Save again to restart. Not worth adding expiry tracking in the UI.

**Object URL leak**: `URL.createObjectURL` results should be revoked when the modal closes to avoid memory leaks. Revoke in the modal's close/reset handler.
