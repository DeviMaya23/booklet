## Context

Three domains (artist, character, image) have update endpoints implemented as `PATCH`. The PATCH approach required pointer-based optional fields (`*string`, `**uuid.UUID`, `Patch[T]`) at every layer — handler request structs, usecase params, and repository logic — to distinguish "field not sent" from "field set to null." This complexity was introduced unnecessarily: the app is single-user and every edit form loads the full entity before submission, so the client always holds complete state.

The image domain partially worked around this by introducing a `Patch[T]` generic type for `artist_id`, but this pattern would need to spread to every nullable field across all three domains if PATCH were kept.

Current handler tests for all three update handlers are written against PATCH. The FE has a working update flow only for artists (`useUpdateArtist` + `ArtistFormModal`); character and image update UIs don't exist yet.

## Goals / Non-Goals

**Goals:**
- Replace `PATCH /artists/:id`, `PATCH /characters/:id`, and `PATCH /images/:id` with `PUT` equivalents
- Simplify all update params structs to plain nullable fields (no pointer-of-pointer, no `Patch[T]`)
- Simplify repository update logic to write all received fields unconditionally
- Update `useUpdateArtist` hook and `ArtistFormModal` to reflect PUT semantics
- Update handler tests to reflect new method and full-body semantics

**Non-Goals:**
- Implementing character or image update UIs (no update hooks or forms exist for those yet — they'll be authored under PUT semantics when built)
- Removing the `Patch[T]` type itself (it can be deleted as a cleanup, but this change doesn't depend on anything that reuses it)
- Changing create endpoints

## Decisions

### 1. Full replace on every PUT call — no field-level diffing

The handler receives the full body, passes all fields to the usecase, the usecase passes them to the repo, and the repo writes them. No `if param != nil` guards for nullable fields — just write what you receive.

The only guard that remains is for `name`: it is required and non-nullable, so the validate tag (`required`) handles that at the handler boundary.

**Alternative considered**: keep `*T` params and treat `nil` as "use existing value." Rejected — this is the same PATCH semantics under a different name, and it doesn't solve the clearing problem.

### 2. Empty string → `null` conversion happens at the FE, not the BE

The BE accepts `null` for nullable fields; it does not treat empty string as null. The FE converts an empty input (`""`) to `null` before serializing to JSON. This keeps the API contract clean and consistent with the domain model (nullable columns, not empty-string columns).

### 3. `folder_ids` and `character_ids` retain replace-all semantics — send `[]` to clear

These are collection fields managed via association replacement. An empty array `[]` means "remove all associations." `null` (field absent) is not valid under PUT — the body must always include these fields. The validation tag on the handler will require them.

**Rationale**: replace-all was already working correctly under PATCH because the "absent = skip" ambiguity doesn't apply to collections the same way — an empty array is unambiguously "clear." Keeping this behavior is consistent and requires no change to repository logic for these fields.

### 4. `thumbnail_r2_path` stays out of the image PUT body

This field is written exclusively by the thumbnail worker. Exposing it in the user-facing PUT body would be a security and correctness concern. It remains on its own internal path.

### 5. Avatar endpoints for character are unchanged

Avatar is a multi-step async upload flow with its own endpoints (`POST /avatar/uploads`, `POST /avatar/uploads/:id/complete`, `DELETE /avatar`). It is not a form field and cannot be included in a PUT body. These endpoints are untouched.

### 6. `Patch[T]` removed from image handler

`image_handler.go` currently uses `Patch[string]` for `artist_id` and has accompanying logic to convert it to `**uuid.UUID`. With PUT, `artist_id` in the request body is a plain `*string` (nullable). The handler converts it to `*uuid.UUID` and passes it directly. The `**uuid.UUID` double-pointer in `UpdateImageParams` becomes `*uuid.UUID`.

## Risks / Trade-offs

- **Breaking API change** → The method change from PATCH to PUT is a breaking contract change. Since the only consumer is the first-party frontend (no public API), this is low-risk and coordinated in the same change. Bruno collection files must be updated alongside.

- **Handler tests rewritten for PUT** → Existing handler tests send PATCH requests and test partial-field behavior. They will be rewritten to send full bodies. Test coverage should remain equivalent; no new scenarios introduced.

## Migration Plan

No data migration required. The change is purely to the API method and request/response handling. BE and FE ship together in one branch; there is no interim compatibility period needed.

Rollback: revert the branch. No database state is affected.
