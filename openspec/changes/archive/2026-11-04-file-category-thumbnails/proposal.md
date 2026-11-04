# Proposal

## Why

Non-image files (video, audio, layered art project files, 3D models, archives) render as a generic paper icon on a flat muted tile in every file grid. At a glance, a `.clip` file and a `.zip` file look identical, which makes the inbox and artpiece file grids hard to scan — especially since these are exactly the file types artists upload most often alongside reference images.

## What Changes

- **Backend — `extension` field**: A nullable `extension` column is added to `files` and `pending_file_uploads`. The FE sends the file's extension at upload-initiation time (derived from the original filename before name-stripping); the value is stored through `pending_file_uploads` and copied to `files` on completion. All file API responses include `extension`.
- **Backend — migration**: One migration adds the `extension` column to both tables.
- **Frontend — `fileCategory.ts` utility**: A new pure function `getFileCategory(mimeType, extension)` maps a file to one of eight categories (Video, Audio, Layered, 3D, Archive, Document, Image, Fallback). Extension takes precedence for proprietary formats; mime-type prefix is used for standard types.
- **Frontend — `FileThumb` component**: A new component renders the three-branch fallback chain: (1) generated image thumbnail, (2) `Skeleton` while `thumbnail_gen_state === 'pending'`, (3) category tile with icon, optional color tint, and uppercase extension label pill. Custom user thumbnails (future, not built) are reserved as branch 0.
- **Frontend — CSS tint tokens**: Four category hues (Video, Audio, Layered, 3D) are defined as CSS custom properties per theme block in `index.css`. Archive, Document, Image, and Fallback use `bg-muted` unchanged.
- **Frontend — existing tile sites updated**: `FileTile`, `ArtpieceDetailFileGrid`, `ArtpieceFilesInput`, and `FileTileEditOverlay` replace their inline fallback with `<FileThumb>`.
- **Frontend — upload paths send extension**: Every `initUpload` call is updated to also send the `extension` field derived from the original filename.

## Capabilities

### New Capabilities

- `web-file-category-thumb`: FileThumb component — category-aware fallback tile with icon, tint, and extension label shown when no thumbnail is available.

### Modified Capabilities

- `file-upload`: Initiate-upload request now accepts and stores an optional `extension` field; complete-upload carries it from the pending record to the file row.
- `file-crud`: All file responses now include the `extension` field.
- `artpiece-files`: Artpiece file ref response now includes the `extension` field.
- `web-file-inbox`: Inbox drag-and-drop upload sends `extension`; tile rendering uses `FileThumb`.
- `web-file-inbox-upload-modal`: Modal upload sends `extension`.
- `web-artpiece-detail`: Edit-mode file upload sends `extension`.

## Impact

- **BE — migration**: `000030_add_extension_to_files` — ALTER TABLE on `files` and `pending_file_uploads`.
- **BE — domain**: `domain.File` and `domain.PendingFileUpload` gain `Extension *string`.
- **BE — usecase**: `InitiateFileUploadParams` gains `Extension *string`; `CompleteUpload` copies it from pending to file.
- **BE — handler**: `initiateFileUploadRequest`, `fileResponse` (file upload handler), and `fileRef` (artpiece handler) gain `extension`.
- **FE — types**: `File` (inbox API) and `ArtpieceFile` (artpiece API) gain `extension: string | null`.
- **FE — new files**: `features/files/lib/fileCategory.ts`, `features/files/components/FileThumb.tsx`.
- **FE — updated files**: `index.css` (CSS vars), `FileTile`, `ArtpieceDetailFileGrid`, `ArtpieceFilesInput`, `FileTileEditOverlay`, `FileInboxGrid`, `AddFilesModal`/`ArtpieceFilesInput` upload call sites.
- **No breaking changes**: `extension` is additive to all API responses; existing consumers that ignore the field are unaffected.
