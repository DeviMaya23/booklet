# Proposal

## Why

Users currently have to leave booklet and navigate to Bookleaf to see what images are inside a folder they've linked to a character. This change adds the ability to fetch a folder's image list directly from booklet, so folder contents can be surfaced in-app (initially to support the character create/edit modal revamp).

## What Changes

- New Bookleaf internal endpoint (implemented by Bookleaf): `GET /internal/users/:userID/folders/:folderID/images`
- New method on the Bookleaf HTTP client: `GetFolderImages(ctx, userID, folderID)` returning a list of images with `image_id` and `thumbnail_url`
- New booklet endpoint: `GET /folders/:folderID/images` — authenticates the caller, proxies to Bookleaf, and returns the image list
- New frontend query hook `useGetFolderImages(folderID)` that calls the new endpoint

## Capabilities

### New Capabilities

- `folder-images`: Defines the rules for fetching and exposing the image list of a Bookleaf public folder, including the Bookleaf client method, the booklet proxy endpoint, and the frontend query hook.

### Modified Capabilities

_(none — existing folder listing and character folder behaviour is unchanged)_

## Impact

- `backend/internal/bookleaf/client.go` — new `GetFolderImages` method and `FolderImage` / `FolderImageList` types
- `backend/internal/usecase/folder_usecase.go` — new `GetFolderImages` usecase method; `BookleafClient` interface extended
- `backend/internal/handler/folder_handler.go` — new `GetFolderImages` handler; `FolderUsecase` interface extended
- `backend/cmd/server/main.go` — route registration for `GET /folders/:folderID/images`
- `frontend/src/features/characters/api/useGetFolderImages.ts` — new query hook
- `collection/folders/` — new Bruno file for the endpoint
