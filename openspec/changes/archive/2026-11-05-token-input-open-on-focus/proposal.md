## Why

In the artpiece forms the Tags field lists suggestions as soon as it is focused, but the Characters field shows nothing until the user types. The two fields share the same chip-input component; only the Tags picker opts in to listing on focus. Tags' on-focus list is also uncapped, so a user with hundreds of tags gets a dropdown that renders every tag and runs off the modal.

## What Changes

- Cap the on-focus (empty text) suggestion list in the shared chip input at the 10 first suggestions in the order given. When more exist, show a non-selectable hint row telling the user to type to search. This applies to every field that lists on focus, so Tags gets the cap too.
- Turn on list-on-focus for the Characters field in the artpiece create modal (reached from the Inbox and from the Artpieces page) and in the artpiece detail view's edit mode.
- Typed-text filtering is unchanged and remains uncapped. The artist/character/tag filter bar and the folder picker are unchanged (still list only after typing).

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `web-tag-picker`: the on-focus tag list is capped, with a hint row when truncated.
- `web-artpiece-create-modal`: the Characters field lists characters on focus, capped the same way.
- `web-artpiece-detail`: the Characters field in edit mode lists characters on focus, capped the same way.

## Impact

- `frontend/src/components/TokenInput.tsx`: cap and hint row for the on-focus list.
- `frontend/src/features/files/components/ArtpieceFormModal.tsx`: Characters field opts in (covers Inbox and Artpieces page entry points, which both render this modal).
- `frontend/src/features/artpieces/components/ArtpieceDetailView.tsx`: Characters field opts in.
- `frontend/src/components/TokenInput.test.tsx` and the form/detail tests that cover the Characters field.
- No backend, API, or Bookleaf changes. Call sites checked: `TokenInput` is also used by `ArtistCharacterFilter` and `FolderPicker`; neither opts in, so neither changes.
