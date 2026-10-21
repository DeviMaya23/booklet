# Design

## Context

`CharacterFormModal` currently supports two paths into `stageFileForCrop`: a file-picker click and a drag-and-drop onto the avatar button. Both produce a `File` and call `stageFileForCrop(file)`, which creates an object URL and opens `AvatarCropDialog`. The crop dialog, canvas utility, and upload flow are unchanged by this feature.

The modal has two text inputs — a name `<Input>` and a notes `<Textarea>` — that must continue to receive normal Ctrl+V text paste without triggering the avatar crop flow.

## Goals / Non-Goals

**Goals:**
- Intercept clipboard paste of an image inside the open modal and route it into the existing crop flow.
- Leave text paste in `<input>` and `<textarea>` elements unaffected.

**Non-Goals:**
- Handling clipboard paste outside the character form modal.
- Accepting formats other than the browser's native `image/*` clipboard items.
- Adding input-size validation (the existing file-picker and drag-drop paths have none; parity is correct here).

## Decisions

### Where to attach the paste listener

**Decision**: Attach an `onPaste` handler to the `<DialogContent>` element (the modal's root DOM node), rather than to `document` or `window`.

**Rationale**: Scoping to the dialog container means the listener is active only while the modal is mounted and naturally handles the stacking case where `AvatarCropDialog` is open on top (events from the crop dialog do not bubble into the character form's handler in a meaningful way).

**Alternative considered**: `useEffect` + `document.addEventListener('paste', ...)` with a guard on `open`. This works but requires explicit cleanup and makes the listener active even if the dialog is visually covered by other dialogs.

### Suppressing paste in text fields

**Decision**: In the `onPaste` handler, check `event.target` — if `instanceof HTMLInputElement` or `instanceof HTMLTextAreaElement`, return early without extracting image data.

**Rationale**: This is the minimal check. The modal has exactly one `<Input>` (name) and one `<Textarea>` (notes); both are standard HTML elements. Checking by instance type is reliable and avoids coupling to class names or refs.

**Alternative considered**: Checking the first clipboard item's `kind === 'file'` and type before the target check, to skip the guard when only text is in the clipboard. Unnecessary — text pastes into an input have no image item anyway, so the guard would never fire in that case. The simpler `instanceof` check covers it.

### Extracting the image from ClipboardData

```
const item = Array.from(event.clipboardData?.items ?? [])
               .find(i => i.type.startsWith('image/'))
const file = item?.getAsFile()
if (file) stageFileForCrop(file)
```

`getAsFile()` returns a `File` with a synthetic name like `image.png`; the output is always renamed to `avatar.jpg` by `getCroppedFile`, so the source name is irrelevant.

## Risks / Trade-offs

- **Large clipboard images**: A clipboard item can in theory be a very large image (e.g., a multi-MP screenshot from a Retina display). The crop dialog loads the full source into a browser `Image` element, which may use significant memory. This risk is pre-existing and identical to the drag-drop path; no new guard is introduced here.
- **`clipboardData` unavailability**: In some embedded or restricted browser contexts `clipboardData` may be null. The optional-chaining guard (`clipboardData?.items`) makes this a silent no-op rather than a crash.
