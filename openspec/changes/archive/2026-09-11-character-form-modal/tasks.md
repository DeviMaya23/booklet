## 1. API Hooks

- [x] 1.1 Create `useCreateCharacter` mutation hook (`POST /characters`)
- [x] 1.2 Create `useUpdateCharacter` mutation hook (`PUT /characters/:id`)
- [x] 1.3 Create `useInitAvatarUpload` mutation hook (`POST /characters/:id/avatar/init`)
- [x] 1.4 Create `useCompleteAvatarUpload` mutation hook (`POST /characters/:id/avatar/:uploadID/complete`)

## 2. ResourceCard — onClick support

- [x] 2.1 Add optional `onClick` prop to `ResourceCard`; invoke it on card body click, not on dropdown trigger interaction

## 3. CharacterFormModal component

- [x] 3.1 Create `CharacterFormModal` component scaffold (Dialog, shared open/edit-mode props, `Character` type import)
- [x] 3.2 Implement avatar picker area: placeholder icon when no avatar, image preview when avatar exists or file picked, click-to-open file input restricted to `image/jpeg,image/png`
- [x] 3.3 Implement avatar local preview via `URL.createObjectURL`; revoke object URL on modal close/reset
- [x] 3.4 Implement "remove avatar" button in edit mode (sets local `avatarCleared` state; no network call until Save)
- [x] 3.5 Implement name input (required) and notes textarea
- [x] 3.6 Implement footer: Cancel + Save for create mode; Delete (destructive, left) + Cancel + Save for edit mode
- [x] 3.7 Implement create submit handler: `POST /characters` → if avatar picked: init → R2 PUT → complete → on avatar failure: `DELETE /characters/:id` + error toast + keep modal open
- [x] 3.8 Implement edit submit handler: avatar changed → init → R2 PUT → complete first; avatar cleared → `DELETE /characters/:id/avatar` first; then `PUT /characters/:id`; invalidate query on success
- [x] 3.9 Implement delete flow: Delete button opens `AlertDialog`; on confirm calls `DELETE /characters/:id`, closes modal, refetches list

## 4. CharactersGrid — wire modal

- [x] 4.1 Add modal open state and selected character state to `CharactersGrid`
- [x] 4.2 Pass `onClick` to `ResourceCard` to open modal in edit mode with the clicked character
- [x] 4.3 Wire `+ New` button (in the characters page) to open modal in create mode
- [x] 4.4 Render `CharacterFormModal` in `CharactersGrid` (or characters page) with open/onOpenChange/character props

## 5. Build & lint

- [x] 5.1 Run `npm run build` and fix any type errors
- [x] 5.2 Run `npm run lint` and fix any lint findings
