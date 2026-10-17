# Design

## Context

See `proposal.md — Why` for motivation.

The current fallback for all file tiles is a single `<FileIcon>` on `bg-muted`, rendered inline in four separate places: `FileTile` (inbox), `ArtpieceDetailFileGrid` (artpiece detail), `ArtpieceFilesInput` (form uploader), and `FileTileEditOverlay` (edit panel). There is no shared component.

The `File` API type has `mime_type` but no `extension`. For common web types (`image/*`, `video/*`, `audio/*`, `application/pdf`), mime is sufficient to determine category. For proprietary art formats (`.clip`, `.kra`, `.sai2`, `.vrm`, `.procreate`, etc.) the browser reports `application/octet-stream`, making mime alone useless for category matching. An `extension` field sent at upload time resolves this.

The project has four themes: `warm` (default), `lumen`, `sunless` (dark), `root` (alias of warm). Dark uses the `[data-theme="sunless"]` selector.

## Goals / Non-Goals

**Goals:**
- Add `extension` to the upload pipeline (FE → pending_file_uploads → files → API responses)
- Category-map files to one of eight categories using extension-first, mime-fallback logic
- Render a category tile (icon, tint for 4 categories, label) when no thumbnail is available
- Consolidate the four inline fallback rendering sites into a single `FileThumb` component
- Theme-aware tint tokens that require no component changes when new themes are added

**Non-Goals:**
- Custom user-uploaded thumbnails (reserved as branch 0 but not built)
- Server-side mime sniffing or extension validation
- Retroactively populating `extension` on existing file rows (nulls are acceptable; pre-existing files show no label)
- Changing the label shown on the existing `ArtpieceDetailFileGrid` corner chip (the `bg-black/60` mime label is a separate, unrelated display element — not replaced by this change)

## Decisions

### Extension is client-sent, not server-derived

**Decision**: The FE derives the extension from `File.name` before stripping it and sends it in the initiate-upload request body. The server stores it verbatim with no validation or normalization beyond what the existing `name` field undergoes.

**Rationale**: The server has no reliable way to sniff extensions for proprietary binary formats. The client has the original filename. Normalizing to lowercase happens on the FE before sending so the stored value is consistent.

**Alternative rejected**: Server-side sniffing using magic bytes. Too complex, requires a new dependency, and still wouldn't solve all proprietary format cases.

### CSS custom properties for tint tokens, not hardcoded Tailwind classes

**Decision**: Define `--file-cat-{category}-bg` and `--file-cat-{category}-fg` CSS custom properties inside each theme block in `index.css`. The component references them via inline `style`. Only the four colored categories (Video, Audio, Layered, 3D) need vars; the other four use `bg-muted` and `text-muted-foreground`.

**Rationale**: Tailwind's `dark:` variant only targets `sunless`. If tints were expressed as `bg-slate-100 dark:bg-slate-800`, adding a third theme (`lumen`) would require changing component code. CSS vars in theme blocks are additive — a new theme just adds new values.

**Alternative rejected**: Hardcoded `dark:` Tailwind classes per category. Doesn't scale to more than two themes.

### Percentage-based icon sizing, no container queries

**Decision**: The icon container uses `w-[28%] aspect-square` (relative to the parent `w-full h-full` tile). The tile itself is already responsive via CSS grid column counts. The label uses a fixed `text-[10px]` floor.

**Rationale**: Container queries add complexity (`@container` class, `cqw` units) for a problem already solved by the grid's responsive column sizing. Since every tile is `aspect-square` and fills its grid cell, `w-[28%]` scales identically at every tile size in the grid.

**Alternative rejected**: `@container` + `cqw`. Works but is unnecessary given the grid's existing responsiveness.

### FileThumb lives in `features/files/components/`

**Decision**: `FileThumb.tsx` is a feature-local component under `features/files/`. All four consuming sites already live in this feature. It is not promoted to `components/ui/` or top-level `components/`.

**Rationale**: The conventions promote to shared layers only when a second feature actually imports the component. All current consumers are in the `files` feature; promotion is premature.

### `fileCategory.ts` placement and interface

**Decision**: `getFileCategory(extension: string | null)` lives in `features/files/lib/fileCategory.ts`. Mime type is not a parameter. It returns a plain object:
```ts
{
  icon: LucideIcon,
  label: string | null,  // extension.toUpperCase(), or null when extension is null
  tinted: boolean,        // true for Video, Audio, Layered, 3D
  cssVar: string | null   // '--file-cat-video' etc., null for neutral categories
}
```
Extension is the only signal. Null extension maps to Fallback with no label.

**The existing `mimeTypeLabel.ts` files are not changed** — they serve `ArtpieceDetailFileGrid`'s corner chip and `FileViewer`, which are separate display elements with different rendering.

### Extension → category mapping

| Category | Icon | Extensions |
|---|---|---|
| Video | Film | mp4, mov, webm, mkv, m4v, avi, wmv, flv, ogv, ts |
| Audio | Music | mp3, wav, flac, m4a, ogg, aac, wma, opus, aiff, alac |
| Layered / project | Layers | psd, clip, kra, sai, sai2, xcf, procreate, ai, afphoto, psb, csp, afdesign, pdn, ase, aseprite |
| 3D / rigged model | Box | blend, fbx, obj, glb, gltf, vrm, moc3, cmo3, vox, dae, stl |
| Archive | FileArchive | zip, rar, 7z, tar, gz, tgz, bz2, xz, zst |
| Document | FileText | pdf, txt, md, doc, docx, rtf, odt, pages |
| Image | FileImage | jpg, jpeg, png, gif, svg, webp, tiff, tif, bmp, heic, heif, avif, ico |
| Fallback | File | everything else |

### Label source

**Decision**: The label pill shows `extension.toUpperCase()` when extension is non-null. Files uploaded before this change have null extension and show no label pill (fallback icon, no label). Mime type is not used as a label source.

**Rationale**: Mime type for `application/` types is unreliable and often verbose (e.g. `application/vnd.openxmlformats-officedocument.wordprocessingml.document`). The extension derived from the original filename is always user-readable. Files with null extension already have a working generic icon — the same as before this change — so there is no regression.

## Risks / Trade-offs

- **Pre-existing files have null extension** → They will show no label pill. This is acceptable; the icon still renders correctly for types identifiable via mime (images, video, audio). Only proprietary formats uploaded before this change show a plain fallback icon — same as today, no regression.
- **Browser-reported extension may be wrong** → A user could rename `photo.jpg` to `photo.mp4`. The tile would show a Video icon for a JPEG. Since this is a single-user personal library, this is an acceptable edge case. The server does not validate.
- **Long extensions truncate in the pill** → Extensions like `procreate` (9 chars) are long. The pill should truncate with `max-w` and `truncate` to avoid layout issues at small tile sizes.

## Migration Plan

Migration `000030_add_extension_to_files`:
```sql
ALTER TABLE files ADD COLUMN extension text;
ALTER TABLE pending_file_uploads ADD COLUMN extension text;
```

Both columns are nullable. No backfill is needed — existing rows keep `extension = NULL` and render without a label. The migration is safe to run with zero downtime (additive column, no default, no constraints).

Rollback: drop the columns. No FE change is needed for rollback since `extension: null` already renders correctly (no pill shown).
