## Purpose

Defines the `AvatarCropDialog` component and `cropCanvas` utility used to crop a source image to a square JPEG before it is used as a character avatar.

## Requirements

### Requirement: AvatarCropDialog — square crop with zoom and pan
The system SHALL provide an `AvatarCropDialog` component that renders a square crop editor over a source image. The editor SHALL use `react-easy-crop` with `aspect={1}` and `cropShape="rect"`. A zoom slider below the crop view SHALL allow the user to zoom between 0.5× and 3×. The dialog footer SHALL contain Cancel and Crop buttons.

On Crop confirmation, the component SHALL:
1. Compute the `PixelCrop` from `react-easy-crop`'s `onCropComplete` callback
2. Export the cropped region to a square JPEG via canvas, bounded to 1024 × 1024 px at quality 0.9
3. Invoke an `onCrop(file: File)` callback with the resulting `File`
4. Clean up the source object URL

On Cancel, the component SHALL invoke an `onCancel()` callback without producing a file.

#### Scenario: Crop confirmed produces square JPEG File
- **WHEN** the user positions the crop area and clicks Crop
- **THEN** a `File` of type `image/jpeg` is produced from the cropped region (square, ≤ 1024 px per side) and passed to `onCrop`

#### Scenario: Crop confirmed cleans up source object URL
- **WHEN** the user clicks Crop
- **THEN** the source object URL is revoked after the canvas export completes

#### Scenario: Cancel aborts crop without producing a file
- **WHEN** the user clicks Cancel in the crop dialog
- **THEN** `onCancel` is called, no File is produced, and the source object URL is revoked

#### Scenario: Zoom slider adjusts crop magnification
- **WHEN** the user moves the zoom slider
- **THEN** the crop view zooms between 0.5× and 3×, updating the preview in real time

### Requirement: cropCanvas utility
A `getCroppedFile(sourceUrl: string, pixelCrop: PixelCrop, mimeType: 'image/jpeg'): Promise<File>` utility function SHALL exist at `frontend/src/features/characters/lib/cropCanvas.ts`. It SHALL:
1. Load `sourceUrl` into an `HTMLImageElement`
2. Draw the `pixelCrop` region onto an offscreen canvas sized to `min(pixelCrop.width, 1024) × min(pixelCrop.height, 1024)`
3. Return a `File` produced by `canvas.toBlob` at JPEG quality 0.9

#### Scenario: Returns a File for a valid crop region
- **WHEN** called with a valid source URL and a non-zero pixel crop
- **THEN** the function resolves to a `File` of type `image/jpeg`

#### Scenario: Output is bounded to 1024 × 1024
- **WHEN** the pixel crop region exceeds 1024 px on either axis
- **THEN** the canvas output is scaled down to fit within 1024 × 1024 while preserving the square aspect ratio
