## Context

The images module was the original upload/display system. It has been superseded by artpieces + files, which handle the same use case with a better model. The app is pre-launch with no live user data, making this a straight deletion with no migration concerns beyond dropping the DB tables.

The images module is not fully isolated — it is woven into the character layer (`GetCharacterImages`, `imageRepo` in `characterUsecase`) and the user account-deletion flow (`DeleteAllUserData` in `userRepository`). These call sites must be cleaned up as part of this change.

## Goals / Non-Goals

**Goals:**
- Remove all code, tests, migrations, routes, and frontend pages belonging to the images module
- Remove `GetCharacterImages` endpoint and its character-layer plumbing
- Drop `images`, `image_characters`, and `pending_uploads` tables
- Leave the codebase compiling cleanly with all remaining tests passing

**Non-Goals:**
- Replacing `GetCharacterImages` (future change against artpieces/files)
- Migrating or preserving existing image data in R2 (app is pre-launch)
- Changing how artpieces, files, or characters work

## Decisions

**Delete `GetCharacterImages` entirely, no stub.**
The frontend never calls `GET /characters/:id/images`. No clients will break. A 410 stub would be dead code immediately. A future implementation will be a clean addition against the artpieces/files model.

**Single drop migration for all three tables.**
`images`, `image_characters`, and `pending_uploads` are all owned by the images module. They have no dependencies from the surviving modules. One migration is cleaner than three.

**Keep `fakeImageRepository` removal scoped to `fakes_test.go`.**
`character_usecase_test.go` uses `&fakeImageRepository{}` as a constructor argument. Once `imageRepo` is removed from `NewCharacterUsecase`, all those call sites must be updated too — they pass it as a positional arg, so the compiler will catch every instance.

**Remove image/pending-upload cleanup from `DeleteAllUserData` entirely.**
Once the tables are gone, this code would panic. No partial data concern: the drop migration removes everything at once.

## Risks / Trade-offs

- [Worker registration] `GenerateThumbnailWorker` and `PurgeExpiredUploadsWorker` are registered in `main.go`. Removing them while any queued jobs reference their `Kind()` strings (`generate_thumbnail`, `purge_expired_uploads`) would cause River to log errors. Non-issue pre-launch; noting for awareness.
- [api.test.ts] Tests use `/images` as a fixture URL for `apiFetch` behaviour tests (auth headers, maintenance mode). These are unrelated to the images feature — the URL just needs to be any valid string. Swapping to `/artpieces` is safe.

## Migration Plan

1. Apply code changes (BE + FE)
2. Run `golangci-lint` and `npm run build` / `npm run lint` to confirm clean compile
3. Write and apply `000023_drop_images_module` migration:
   ```sql
   DROP TABLE IF EXISTS image_characters;
   DROP TABLE IF EXISTS pending_uploads;
   DROP TABLE IF EXISTS images;
   ```
   Down migration recreates the three tables from their original create migrations.
4. No rollback concern pre-launch.

## Open Questions

_None._
