# Design

## Context

See `proposal.md — Why` for motivation. Relevant current state:

- Both `POST /files` (init file upload) and `POST /characters/:id/avatar/init` call `storage.GeneratePresignedPutURL(ctx, key, contentType, ttl)`. Neither passes a size.
- `GeneratePresignedPutURL` in `r2.go` calls `presign.PresignPutObject` with `PutObjectInput{Bucket, Key, ContentType}`. No `ContentLength` is set, so R2 does not enforce a size cap on the resulting URL.
- The FE `uploadFiles()` in `ArtpieceFilesInput` and `AddFilesModal` are near-identical and share the core logic via `uploadFile()` in `lib/upload.ts`. Avatar upload goes through a separate path in `CharacterFormModal.stageFileForCrop()`.
- `fileCategory.ts` already maps every extension to a `FileCategory`. This is the natural anchor for the FE size constants.

## Goals / Non-Goals

**Goals:**
- FE rejects oversized files immediately (before any network call) with a clear toast.
- BE validates `file_size` on init and returns 422 if the limit is exceeded.
- The presigned PUT URL is signed with `ContentLength` so R2 enforces the exact size even if FE validation is bypassed.
- Size limits are consistent between FE and BE via a cross-reference comment (both sides are the source of truth for their own constant; no runtime sync).

**Non-Goals:**
- Switching from presigned PUT to presigned POST — not needed since the exact file size is available at init time.
- A shared config file or build-time codegen to eliminate the constant duplication — justified but out of scope.
- BE enforcement on the thumbnail upload URL — thumbnails are generated client-side at bounded dimensions (1024×1024 JPEG) and will never be large.

## Decisions

### D1 — Size limit constants live in `fileCategory.ts` (FE) and a new `pkg/filesize` package (BE)

**FE:** Add a `FILE_SIZE_LIMITS` map (`Record<FileCategory, number>`) alongside the existing `CATEGORY_META` in `fileCategory.ts`. Category is already derived from extension there, so the limit lookup is a natural one-liner: `FILE_SIZE_LIMITS[getFileCategory(extension).category]`.

**BE:** New `pkg/filesize` package with a single exported function `MaxBytesForExtension(ext string) int64` backed by a flat `extension → int64` map. The category abstraction is not needed on the BE; only the limit value matters. Unknown extensions fall back to a conservative default (100 MB).

A comment in both files points to the other: _"Keep in sync with backend/pkg/filesize — limits must agree or BE will reject files that passed FE validation."_ / _"Keep in sync with frontend/src/features/files/lib/fileCategory.ts FILE_SIZE_LIMITS — same values."_

**Alternative considered:** A shared JSON file consumed by both build pipelines. Cleaner in theory, but adds a build dependency and is overkill for ~8 constants that change rarely.

### D2 — FE size check is a utility function in `upload.ts`, called by both `ArtpieceFilesInput` and `AddFilesModal`

Both components already call `uploadFile()` from `upload.ts`. A `filterBySize(files: File[]): { accepted: File[]; rejected: File[] }` helper (also in `upload.ts` or a new `lib/filesize.ts`) keeps the logic in one place. Each component handles the rejection toast itself (they already own their own toast calls).

Avatar check stays in `CharacterFormModal.stageFileForCrop()` — it's a single file, separate from the batch upload path, so no shared utility is needed.

### D3 — `GeneratePresignedPutURL` gains a `contentLength int64` parameter

The storage interface and `r2.go` implementation both gain the parameter. It is passed as `ContentLength` on `PutObjectInput`. All call sites (file upload usecase, character usecase) must provide it.

**Alternative considered:** A new `GeneratePresignedPutURLWithSize` function to avoid touching existing call sites. Rejected — there's only two call sites and keeping one signature is simpler.

### D4 — `file_size` is a required field on both init endpoints

`initiateFileUploadRequest` in the file upload handler and `initAvatarUploadRequest` in the character handler both gain `FileSizeBytes int64 json:"file_size" validate:"required,min=1"`. The handler passes it to the usecase. The usecase calls `filesize.MaxBytesForExtension(ext)` (or the equivalent for avatar, which is always an image) and returns a typed error if exceeded. The handler maps that error to HTTP 422.

**FE side:** `useInitFileUpload` mutation params gain `fileSizeBytes`. `uploadFile()` in `upload.ts` passes `file.size`. Avatar upload in `useInitAvatarUpload` gains the same param; `CharacterFormModal` passes the source file's size before cropping.

Note on avatar: the BE receives the source file size (before crop), not the cropped output size. The crop output is always a small JPEG (1024×1024 at 0.9 quality), so the `ContentLength` signed into the presigned URL is for the source — but the actual PUT body is the cropped file. This means ContentLength enforcement will misfire for avatar uploads (source ≠ crop output size). To avoid this, the FE should pass `croppedFile.size` (available in `onCrop` before the upload) to the avatar init call, not the source size. The BE signs ContentLength for the cropped file; the FE size check on the source guards against loading an oversized image into the canvas.

This splits the avatar into two checks:
- FE source check (in `stageFileForCrop`): guards canvas memory. Uses the image category limit on the source file.
- BE init + ContentLength: uses `croppedFile.size`, which is always small. The BE validates against the image limit, which the crop output will trivially pass.

## Risks / Trade-offs

- **Duplicated constants can drift** → Mitigated by cross-reference comment; the consequence is a user seeing a confusing 422 from BE after passing FE validation, which is noticeable and will surface quickly.
- **ContentLength mismatch on avatars if crop output size is not sent** → Mitigated by the design decision above (send crop size, not source size, to avatar init).
- **`aws-sdk-go-v2` `ContentLength` on presigned PUT enforced by R2** → R2 is S3-compatible and enforces signed `Content-Length` headers. If a future R2 change breaks this, FE validation remains as the user-facing guard.
