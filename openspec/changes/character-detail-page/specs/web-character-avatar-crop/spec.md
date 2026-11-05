# Spec Delta

## MODIFIED Requirements

### Requirement: Clipboard paste as avatar source
When the character form modal is open, or the character detail page is in edit mode, the user SHALL be able to paste an image from the clipboard (Ctrl+V / Cmd+V) to initiate the avatar crop flow. The pasted image SHALL be routed through the same crop dialog as file-picker and drag-and-drop uploads.

The paste SHALL be ignored when the event originates from a focused `<input>` or `<textarea>` element, so that normal text-paste behavior in the name and notes fields is unaffected.

#### Scenario: Paste image when no text field is focused
- **WHEN** the character form modal is open or the detail page is in edit mode, no `<input>` or `<textarea>` is focused, and the user presses Ctrl+V (or Cmd+V) with an `image/*` item on the clipboard
- **THEN** the crop dialog opens with the pasted image as its source

#### Scenario: Paste image to replace an existing avatar
- **WHEN** the detail page is in edit mode with an existing avatar displayed and the user pastes an image from the clipboard
- **THEN** the crop dialog opens, and on confirm the new cropped file replaces the previous avatar preview

#### Scenario: Paste while a text field is focused
- **WHEN** the name input or notes textarea has focus and the user presses Ctrl+V
- **THEN** the paste event is not intercepted; the text field handles it normally and the crop dialog does not open

#### Scenario: Paste non-image content
- **WHEN** the user presses Ctrl+V and the clipboard contains no `image/*` item
- **THEN** the crop dialog does not open and no avatar change occurs

#### Scenario: Paste in view mode is ignored
- **WHEN** the detail page is in view mode and the user pastes an image
- **THEN** the crop dialog does not open
