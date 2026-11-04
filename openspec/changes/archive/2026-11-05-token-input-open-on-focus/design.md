# Design

## Context

`TokenInput` filters suggestions at render time. With empty text it returns the full unselected list when `openOnFocus` is set, otherwise nothing (`TokenInput.tsx`, the `filtered` computation). `TagPicker` sets `openOnFocus`; every Characters field uses `TokenInput` directly without it. Both lists come from `sort=recent` queries that return the whole list, so the cap must be applied client-side. See proposal.md for motivation.

## Goals / Non-Goals

**Goals:**
- One cap implementation in `TokenInput`, shared by Tags and Characters.
- Characters on focus in the two artpiece forms.

**Non-Goals:**
- Server-side limits or pagination.
- Changing the filter bar or folder picker.
- Capping typed-text matches.

## Decisions

- **Cap in `TokenInput`, not per wrapper.** The list is built there, so one change covers every `openOnFocus` user. Alternative: slice in `TagPicker` and a new character wrapper; rejected as duplicated logic and a new wrapper layer.
- **Slice, not scroll.** Cap of 10 as a module constant in `TokenInput` (no new prop). Recency order makes the first rows the useful ones, and a slice avoids rendering hundreds of rows. Alternative: max-height scroll; rejected by the user.
- **Cap applies only to empty text.** Once the user types, results are as filtered. Typed queries are narrow by nature and capping would hide valid matches.
- **Hint row is presentational.** Not counted in the keyboard-navigation row count and has no mouse handler, so arrows, Enter and clicks ignore it. It shows only when the empty-text list was truncated.
- **Selected items are excluded before slicing**, so selections do not consume the cap.
- **Opt in per call site** for the two artpiece Characters fields (`openOnFocus`), leaving the filter bar and folder picker unchanged.

## Risks / Trade-offs

- [Hint row could look like an option] → style it muted, non-hover, no pointer cursor.
- [User unaware more characters exist] → hint row says to type to search.
