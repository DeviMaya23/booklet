# Tasks

## 1. Bookleaf client

- [x] 1.1 Add `FolderImage` and `FolderImageList` types to `internal/bookleaf/client.go` and implement `GetFolderImages(ctx, userID, folderID string)` calling `GET <host>/internal/users/:userID/folders/:folderID/images` with the internal secret header; verify it compiles
- [x] 1.2 Add unit tests in `internal/bookleaf/` for `GetFolderImages`: successful decode returns the image list, non-200 response returns an error containing the status code

## 2. Backend usecase

- [x] 2.1 Extend the `BookleafClient` interface in `internal/usecase/folder_usecase.go` with `GetFolderImages(ctx, userID, folderID string)`, add `GetFolderImages(ctx, idpSubject, folderID string)` method to `folderUsecase`, and write unit tests covering: successful result forwarded from client, client error returned as error to caller

## 3. Backend handler

- [x] 3.1 Extend `FolderUsecase` interface in `internal/handler/folder_handler.go` with `GetFolderImages(ctx, idpSubject, folderID string)`, implement `GetFolderImages` handler (parse and validate `:folderID` path param with `uuid.Parse`, call usecase, return 200 with image list or 404/500 on error), and write handler unit tests covering: success returns 200 with image list, Bookleaf 404 returns 404, upstream error returns 500
- [x] 3.2 Register `GET /folders/:folderID/images` route in `backend/cmd/server/main.go` and verify the server builds

## 4. Bruno collection

- [x] 4.1 Create `collection/folders/get_folder_images.bru` for `GET /folders/:folderID/images` with a path param variable, consistent with the existing folder collection files

## 5. Frontend hook

- [x] 5.1 Create `frontend/src/features/characters/api/useGetFolderImages.ts` exporting a `useGetFolderImages(folderID: string | undefined)` hook; it calls `GET /folders/:folderID/images`, is enabled only when `folderID` is a non-empty string, sets `retry: false`, and returns the typed image list; verify it type-checks

## 6. Quality checks

- [x] 6.1 Run `golangci-lint run ./...` from the `backend/` directory and fix any issues
- [x] 6.2 Run `npm run build` and `npm run lint` from the `frontend/` directory and fix any issues

## Workflow follow-up

- Archive the change once review is complete.
