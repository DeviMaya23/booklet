## Context

Characters can be tagged with Bookleaf folder IDs via the `character_folders` join table, but today there is no FE UI for this and no BE validation — folder IDs are stored blindly. The `character_folders` table has `(character_id, folder_id)` but no `folder_name`. Bookleaf's `GetPublicFolders` API already exists and is wired through the bookleaf client and folder usecase; `GET /folders` is already an endpoint. The character update handler accepts `folder_ids` in the request body but passes them straight to the repository.

The change needs to:
1. Add `folder_name` to `character_folders` so existing links can render without a live Bookleaf call.
2. Thread Bookleaf validation into character create/update — validate, enrich with names, and persist atomically.
3. Change the character API response shape to expose `folders [{id, name}]` instead of `folder_ids []`.
4. Build the folder picker UI in the character modal.

## Goals / Non-Goals

**Goals:**
- Users can add/remove Bookleaf folder links on characters via a picker in the modal
- Existing folder links render their names without requiring Bookleaf to be reachable at modal open time
- Bookleaf validates every folder ID on character save; absent IDs are silently dropped, Bookleaf unreachability fails the save
- `folder_name` stays fresh: every successful save writes the latest name from Bookleaf

**Non-Goals:**
- Surfacing dropped folder IDs to the user (logged internally only, deferred)
- Exposing folder contents (images) in Booklet
- Any change to folder visibility rules on Bookleaf

## Decisions

### D1: Bookleaf validation in the usecase, not the handler

The validation logic (call Bookleaf → build id→name map → filter → enrich) lives in `characterUsecase`, not the handler. The handler extracts `idpSubject` from the JWT context and passes it through `CreateCharacterParams` / `UpdateCharacterParams`. The usecase then calls `BookleafClient.GetPublicFolders(ctx, idpSubject)`.

**Why:** Keeps the handler thin (bind/validate/delegate). The usecase already owns the character lifecycle; Bookleaf validation is a business rule, not a transport concern. Mirrors how `folderUsecase` already wraps the client.

**Alternative considered:** Call `GetPublicFolders` in the handler and pass the validated set to the usecase. Rejected — it puts business logic in the transport layer and makes the handler harder to unit-test.

### D2: Usecase enriches folders; repository accepts `[]domain.CharacterFolder`

`characterUsecase` builds `[]domain.CharacterFolder{FolderID, FolderName}` after Bookleaf validation. `updateWithFolders` in the repository is updated to accept this pre-enriched slice instead of `[]uuid.UUID`.

**Why:** The repository's job is persistence, not enrichment. Keeping name resolution in the usecase avoids leaking external API knowledge into the repo layer. The repo just inserts whatever it receives.

### D3: Fail closed on Bookleaf unavailability

If `GetPublicFolders` returns any error, the entire character create/update is aborted and a 500 is returned. No partial success.

**Why:** Silent partial saves (e.g. character fields saved but folders untouched) would leave the data in an inconsistent state relative to user intent. Failing closed and asking the user to retry is safer than a split outcome. FE can handle this with a generic error toast.

**Alternative considered:** Fall back to the current persisted folder set on Bookleaf failure. Rejected — this hides the failure and lets stale links accumulate without re-validation.

### D4: `folder_name` always refreshed on successful save (write-through)

Every successful save writes the latest `folder_name` from Bookleaf's response, even if the name hasn't changed. No dirty-checking.

**Why:** Names can change on Bookleaf at any time. Write-through is simple and guarantees the persisted name is at most one save stale. The cost (a redundant string write per row per save) is negligible.

### D5: API response changes `folder_ids: []string` to `folders: [{id, name}]`

The character response type changes from a flat `folder_ids []string` to `folders []FolderResponse` where each entry carries both `id` and `name`.

**Why:** The FE needs names to render existing folder chips without calling Bookleaf. Adding a parallel `folder_names` field would be redundant and confusing. The existing `folder_ids` field is only consumed by the FE `Character` type, which can be updated in lockstep.

### D6: FE tracks folder changes as an explicit diff

The modal holds `addedIDs` and `removedIDs` sets against the original folder list from the character's GET response. The final submitted `folder_ids` is computed as: original IDs + addedIDs − removedIDs.

**Why:** If `GetPublicFolders` fails at modal-open time, the picker is disabled but the user's existing folders still render (from persisted names). The diff model means the user can still remove an existing folder or save unrelated field changes — the final intent is always computable from the original list, regardless of whether the public-folders fetch succeeded.

### D7: `FolderPicker` is a standalone component, not a generic multi-select

A dedicated `FolderPicker` component wraps the chip display + search dropdown + tooltip. It takes `selected: {id, name}[]`, `available: {id, name}[]`, `onAdd`, `onRemove`, and `disabled`.

**Why:** The folder picker has domain-specific behaviour (no add option, Bookleaf sourcing tooltip). Reusing a generic multi-select would require too many props or context leakage. A focused component is easier to test and extend.

## Risks / Trade-offs

- **Bookleaf is a hard dependency on every character save that includes folders** → Because character update uses PUT semantics (full field replacement), a non-empty `folder_ids` list always triggers a Bookleaf call — even if the user made no folder changes this session. This is accepted: the call serves double duty as a consistency sweep, silently dropping any folders that have since gone private or been deleted on Bookleaf. The trade-off is that Bookleaf unavailability blocks those saves; fail-closed behaviour makes the failure predictable.
- **Migration adds `folder_name NOT NULL`** → Existing rows have no name. Migration must provide a default (empty string `''`) for existing rows; the usecase will refresh names on the next save.
- **BREAKING response shape change** → Any caller relying on `folder_ids` will break. Currently only the FE `Character` type consumes it; FE is updated in this same change. No external API consumers are known.
- **`idpSubject` not previously used in character handler** → If the JWT context is somehow missing it (should be impossible given middleware always sets it), the handler returns 401. This is the correct fail-safe.

## Migration Plan

1. Deploy migration adding `folder_name text not null default ''` to `character_folders`
2. Deploy BE changes (validation logic, response shape) — backward-incompatible response change, so FE must deploy together or after
3. Deploy FE changes (updated `Character` type, folder picker, modal wiring)
4. Existing `character_folders` rows will have `folder_name = ''` until the user next saves the character; names are refreshed on every successful save from that point on
