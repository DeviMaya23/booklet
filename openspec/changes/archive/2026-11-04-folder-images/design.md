# Design

## Context

The bookleaf HTTP client (`internal/bookleaf/client.go`) already handles public folder listing via one method. The folder handler and usecase follow a thin proxy pattern: authenticate, delegate to the bookleaf client, return the result without transformation. This change adds a second method and endpoint on that same path.

The Bookleaf internal endpoint (`GET /internal/users/:userID/folders/:folderID/images`) is new — implemented by the Bookleaf agent against the agreed API contract. The response includes `image_id` and `thumbnail_url` per image.

## Goals / Non-Goals

**Goals:**
- Add `GetFolderImages` to the bookleaf client
- Expose `GET /folders/:folderID/images` in booklet, following the existing folder proxy pattern
- Provide a `useGetFolderImages(folderID)` hook for the frontend

**Non-Goals:**
- UI changes to the character modal (separate change)
- Caching or storing presigned URLs server-side
- Fetching full-resolution images (thumbnails only)

## Decisions

### Proxy without transformation
Booklet returns `{ "images": [{ "image_id", "thumbnail_url" }] }` exactly as Bookleaf provides it, consistent with how `GET /folders` passes through `folder_list` unchanged.

**Alternative considered**: map to a local response struct.
**Rejected**: no transformation needed; adds churn for no benefit.

### Path param for folderID, user ID from JWT
The endpoint is `GET /folders/:folderID/images`. The user ID is derived from the authenticated JWT (same as `GetPublicFolders`), not accepted from the caller.

**Why**: consistent with how all other user-scoped proxy endpoints work; prevents callers from impersonating other users.

### No server-side caching of presigned URLs
The thumbnail URLs are presigned and have a TTL set by Bookleaf. Booklet does not cache them — each call to `GET /folders/:folderID/images` results in a live Bookleaf call that returns fresh URLs.

**Why**: the primary consumer is a modal (opened on demand). The browser loads images into memory on render; the URL TTL is irrelevant as long as the page stays mounted. Caching would add complexity (invalidation, storage) for negligible benefit in this access pattern.

### folderID as `string` in the handler (not `uuid.UUID`)
The path param is validated via `uuid.Parse` in the handler before use, per project conventions for path params where `c.Bind` does not apply. The bookleaf client accepts it as a string since no downstream logic requires the typed UUID form.

## Risks / Trade-offs

- **Bookleaf availability**: the endpoint adds a new synchronous dependency on the Bookleaf API. A slow or unavailable Bookleaf will surface as 500 to the caller. Mitigation: consistent with existing `GET /folders` behaviour; no special handling needed.
- **URL expiry on long-lived pages**: if a page displaying folder images is kept open beyond the presigned URL TTL and then re-renders (e.g., React re-mounts the component), images will fail to load. Mitigation: out of scope for this change — handled when the UI is designed.
