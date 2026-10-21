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

### Requirement: Clipboard paste as avatar source
When the character form modal is open, the user SHALL be able to paste an image from the clipboard (Ctrl+V / Cmd+V) to initiate the avatar crop flow. The pasted image SHALL be routed through the same crop dialog as file-picker and drag-and-drop uploads.

The paste SHALL be ignored when the event originates from a focused `<input>` or `<textarea>` element, so that normal text-paste behavior in the name and notes fields is unaffected.

#### Scenario: Paste image when no text field is focused
- **WHEN** the character form modal is open, no `<input>` or `<textarea>` is focused, and the user presses Ctrl+V (or Cmd+V) with an `image/*` item on the clipboard
- **THEN** the crop dialog opens with the pasted image as its source

#### Scenario: Paste image to replace an existing avatar
- **WHEN** the character form modal is open with an existing avatar displayed and the user pastes an image from the clipboard
- **THEN** the crop dialog opens, and on confirm the new cropped file replaces the previous avatar preview

#### Scenario: Paste while a text field is focused
- **WHEN** the name input or notes textarea has focus and the user presses Ctrl+V
- **THEN** the paste event is not intercepted; the text field handles it normally and the crop dialog does not open

#### Scenario: Paste non-image content
- **WHEN** the user presses Ctrl+V and the clipboard contains no `image/*` item
- **THEN** the crop dialog does not open and no avatar change occurs
