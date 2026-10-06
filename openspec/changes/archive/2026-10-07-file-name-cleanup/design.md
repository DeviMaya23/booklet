## Context

File uploads in the app send `file.name` (the browser `File` object's name, e.g. `pic.jpg`) as the `name` field to the backend. The backend stores it verbatim. When the download-all feature is wired to use `file.Name` for zip entry filenames, it appends an extension derived from MIME type — producing `pic.jpg.jpg`. A secondary issue is that no validation prevents names containing characters that are illegal on Windows/macOS (`/ \ : * ? " < > |`).

Four upload entry points all have the same pattern:
- `AddFilesModal.tsx`
- `ArtpieceFilesInput.tsx` (used in the Create Artpiece modal)
- `FileInboxGrid.tsx` (inbox DnD)
- `ArtpieceDetailView.tsx` (edit-mode DnD)

## Goals / Non-Goals

**Goals:**
- Strip the file extension from `file.name` before it is used as the `name` in any upload call or pre-filled name input
- Validate file name inputs against a forbidden-character rule and block saves when invalid
- Share both utilities across all four entry points with no duplication

**Non-Goals:**
- Backend changes — the `name` field is already optional and stored as-is; no schema or handler changes needed
- Retroactively fixing names on existing records (app is not live)
- Sanitizing names silently (show an error instead of auto-stripping forbidden chars)
- Validating the notes field

## Decisions

### Strip extension at upload call sites, not in `useInitFileUpload`

The strip could live inside the `useInitFileUpload` mutation hook so it applies automatically everywhere. However:
- The name input is pre-filled from `file.name` before the mutation is called (it appears in the uploading placeholder row and then the uploaded row). If stripping only happens inside the hook, the UI would show the full extension until the upload completes.
- Stripping at the call site means both the `name` passed to the hook and the value shown in the input are consistent from the start.

**Decision**: export a `stripFileExtension(name: string): string` utility from `frontend/src/lib/files.ts` and apply it at every call site when building the name to send and to pre-fill.

### Forbidden character validation on blur, not on keystroke

Keystroke-level validation blocks users mid-type when they paste a value with a forbidden char. On blur, the user has finished editing and the error is actionable. Save is blocked until the error clears.

**Forbidden set**: `/ \ : * ? " < > |` — the union of Windows and macOS illegal filename characters. This matches the existing `SanitizeFilename` regex on the backend (`[/\\:*?"<>|]+`).

Additional rules:
- Blank name (empty or whitespace-only after trim) → error: "Name is required"
- Name longer than 255 characters → error: "Name must be 255 characters or fewer"

**Decision**: export `validateFileName(name: string): string | null` from `frontend/src/lib/files.ts`, returning an error string or null. Call it on blur across all name inputs that persist to the backend.

### Where to apply validation

Name validation applies wherever a name input calls `PUT /files/:id` on blur:
- `AddFilesModal.tsx` — already has inline error display; add the rule alongside the existing API error handling
- `ArtpieceFilesInput.tsx` — add inline error display (same pattern as `AddFilesModal`)
- `FileInboxGrid.tsx` — the inbox DnD flow does not have an editable name input at upload time; the name can be edited later via the tile edit overlay (`FileTileEditOverlay.tsx`)

The inbox grid's tile edit overlay is out of scope for this change — it's a separate surface. This change adds validation only to the upload-time name inputs in `AddFilesModal` and `ArtpieceFilesInput`.

`ArtpieceDetailView.tsx`'s upload path has no inline name input (names are edited in the tile overlay), so validation does not apply there either.

### `stripFileExtension` rules

```
"pic.jpg"         → "pic"
"my.art.png"      → "my.art"    (last dot only)
".gitignore"      → ".gitignore" (leading dot, leave alone: idx <= 0)
"no-extension"    → "no-extension"
"trailing."       → "trailing"  (strip trailing dot)
```

The function trims a trailing dot after stripping. An empty result after stripping (e.g. `"."`) falls back to the original.

## Risks / Trade-offs

- **Existing records**: names already stored with extensions will not be fixed. Acceptable — the app is not live and old records can be cleaned manually if needed.
- **User renames a file to include extension**: after initial strip, the name input is editable and the user can type `pic.jpg` themselves. The validation rules do not forbid `.`, so they can re-add the extension. This is intentional — the strip is a convenience default, not a hard constraint.
- **Shared utility in `frontend/src/lib/files.ts`**: this file may not exist yet. Creating it is a one-line utility file — low risk.
