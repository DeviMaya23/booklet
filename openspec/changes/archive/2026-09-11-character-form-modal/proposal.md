## Why

The characters page can list and delete characters but has no way to create or edit them. The `+ New` button is currently a no-op, and cards have no edit entry point — making the feature incomplete for day-to-day use.

## What Changes

- The `+ New` button on the characters page opens a `CharacterFormModal` in create mode
- Clicking a character card opens the modal in edit mode (pre-populated with existing data)
- `CharacterFormModal` supports: avatar upload (jpg/png, presigned R2 flow), name (required), notes (blurb), and delete/cancel/save actions
- On create, if avatar upload fails after the character is created, the character is deleted and an error is surfaced — the whole sequence is atomic from the user's perspective
- `ResourceCard` gains an `onClick` prop for card-body click; the existing `(...)` dropdown delete remains for quick delete without opening the modal
- Moodboard (folder linking) is out of scope for this proposal

## Capabilities

### New Capabilities

_(none — all backend endpoints are already implemented and specced)_

### Modified Capabilities

- `web-characters`: Requirements change for the `+ New` button (was no-op, now opens modal), card click behavior (new), and the full create/edit/delete flow via modal

## Impact

- **Frontend only** — no backend, database, or API changes
- New hooks: `useCreateCharacter`, `useUpdateCharacter`, `useInitAvatarUpload`, `useCompleteAvatarUpload`
- New component: `CharacterFormModal`
- Modified component: `CharactersGrid` (wires modal open state), `ResourceCard` (adds `onClick` prop)
- Consumes existing endpoints: `POST /characters`, `PUT /characters/:id`, `DELETE /characters/:id`, `POST /characters/:id/avatar/init`, `POST /characters/:id/avatar/:uploadID/complete`
