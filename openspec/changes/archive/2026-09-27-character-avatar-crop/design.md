## Context

`CharacterFormModal` currently accepts any picked image file and sets it directly as `localFile`, which is then uploaded as-is via presigned R2 URL. Avatars are always rendered as squares in the UI, so unframed source images produce poor results. The image modal already has DnD and helper text; the character modal has neither.

No backend changes are required — the presigned upload pipeline already accepts JPEG/PNG and is agnostic to image dimensions.

## Goals / Non-Goals

**Goals:**
- Intercept every file pick (click or drop) and route it through a square crop dialog before `localFile` is set
- Produce a square JPEG output, bounded to 1024 × 1024 px, from the crop confirmation
- Add DnD and helper text to the avatar area, matching the image modal
- Accept any source file size — the crop editor is the tool for handling oversized or oddly-dimensioned inputs

**Non-Goals:**
- Backend changes of any kind
- Rotation in the crop editor
- Crop shapes other than square
- Changing when or how the avatar upload fires (timing follows existing flow)

## Decisions

### Library: `react-easy-crop` over `react-image-crop`

`react-easy-crop` provides the zoom + pan gesture model the design calls for (zoom slider, drag-to-reposition), whereas `react-image-crop` uses a resize-handle model that does not support zoom. `react-easy-crop` is MIT, ~50 KB, and has no transitive dependencies. It returns a `PixelCrop` (`{ x, y, width, height }` relative to the source image at natural resolution), which feeds directly into a `<canvas>` export.

### Canvas export via `cropCanvas.ts` utility

A small `getCroppedFile(source: string, pixelCrop: PixelCrop, mimeType: 'image/jpeg'): Promise<File>` utility handles the canvas work:
1. Load source URL into an `HTMLImageElement`
2. Draw the cropped region onto an offscreen canvas sized to `min(pixelCrop.width, 1024) × min(pixelCrop.height, 1024)`
3. Call `canvas.toBlob(cb, 'image/jpeg', 0.9)` to produce the output
4. Wrap the Blob in a `File` with the mime type

This produces a bounded JPEG regardless of source dimensions or format. The original source object URL is revoked after export.

### Crop dialog as a separate `Dialog` component (`AvatarCropDialog`)

The crop state (zoom level, crop area, source object URL) is entirely transient and exists only during the cropping step. Isolating it in `AvatarCropDialog` keeps `CharacterFormModal` clean and avoids resetting crop state on unrelated re-renders. The dialog mounts when a source file is staged and unmounts (cleaning up the source object URL) when the user confirms or cancels.

The dialog renders:
- `react-easy-crop` view with `aspect={1}` and `cropShape="rect"`
- A zoom range slider (0.5× – 3×) below the crop canvas
- Cancel and Crop buttons in the footer

### Output MIME type: always JPEG

The output is always `image/jpeg` regardless of the source format. PNG sources go through canvas and come out as JPEG. This keeps the upload type consistent and avoids large PNG outputs for what are essentially square avatar thumbnails. The `mimeType` passed to `useInitAvatarUpload` is hardcoded to `image/jpeg` after crop.

### No raw file size gate

Any picked or dropped file is accepted and passed to the crop dialog. The crop editor is the mechanism for handling any source image — the user shouldn't need to pre-process in another app. Output dimensions are bounded at 1024 × 1024 by the canvas export regardless of input size.

## Risks / Trade-offs

- **Canvas output size is not explicitly checked** — bounded to 1024 × 1024 @ JPEG 0.9, typical output is 80–200 KB. A pathologically high-detail 1024 px² image could be larger, but still well under the R2 upload limit.
- **JPEG conversion loses transparency** — PNG sources with alpha become JPEG with a white background in canvas. Acceptable for avatar use.
- **`react-easy-crop` is a new dependency** — it is small and stable, but adds to the bundle. No lighter alternative provides zoom+pan with the required output model.
- **Crop dialog blocks the main modal** — there is no way to "skip" crop and use the raw file. This is intentional per the requirements.
