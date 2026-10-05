## 1. Migrations

- [x] 1.1 Write up migration: `CREATE TABLE commissions` (id, user_id, artist_id, status, price, paid, paid_date, finish_date, notes, created_at, updated_at) and `CREATE TABLE commission_characters` (commission_id, character_id, PK); include down migration
- [x] 1.2 Write up migration: `ALTER TABLE artpieces ADD COLUMN commission_id UUID NULL REFERENCES commissions(id) ON DELETE SET NULL`; include down migration

## 2. Domain

- [x] 2.1 Add `domain.Commission` model with GORM tags: all scalar columns, `Artist *Artist`, `Characters []Character` (many2many commission_characters), `Artpieces []Artpiece` (foreignKey:CommissionID)
- [x] 2.2 Add `CommissionID *uuid.UUID` column field to `domain.Artpiece`

## 3. Repository interfaces & implementations

- [x] 3.1 Create `usecase/commission_repository.go`: define `CommissionRepository` interface (Create, GetByID, List, Update, Delete), `CommissionArtistRepository`, `CommissionCharacterRepository`, `CommissionArtpieceRepository` interfaces, and params/errors types (`CreateCommissionParams`, `UpdateCommissionParams`, `ErrArtpieceAlreadyAttached`, valid status consts)
- [x] 3.2 Create `repository/commission_repository.go`: implement `CommissionRepository` — `GetByID` preloads Artist, Characters, and Artpieces with their CoverFile; `List` ordered by created_at DESC; `Update` replaces character set via GORM Association Replace
- [x] 3.3 Add `GetByIDsAndUserID`, `BulkUpdateCommissionID`, and `GetArtpiecesForCommission` to `artpiece_repository.go` and expose them through the `ArtpieceRepository` interface in `usecase/artpiece_repository.go`

## 4. Commission usecase

- [x] 4.1 Create `usecase/commission_usecase.go` with `CommissionUsecase` struct wiring `CommissionRepository`, `CommissionArtistRepository`, `CommissionCharacterRepository`, `CommissionArtpieceRepository`, and `Transactor`
- [x] 4.2 Implement `Create`: validate artist ownership, character ownership, artpiece ownership + conflict check (`ErrArtpieceAlreadyAttached`); create commission; attach artpieces in transaction
- [x] 4.3 Implement `GetByID`, `List`, `Delete`
- [x] 4.4 Implement `Update`: validate artist/character ownership, replace character set via commission repo
- [x] 4.5 Implement `AttachArtpieces` (POST bulk): validate ownership + conflict for each artpiece ID, bulk update commission_id in transaction; no-op for artpieces already on this commission
- [x] 4.6 Implement `DetachArtpieces` (DELETE bulk): bulk set commission_id = NULL for artpiece IDs that belong to this commission; no-op for others
- [x] 4.7 Implement `ReplaceArtpieces` (PUT): validate ownership + conflict on incoming set, reconcile attach/detach in a single transaction (same pattern as `artpiece_usecase.ReplaceFiles`)

## 5. Commission handler

- [x] 5.1 Create `handler/commission_handler.go`: define `CommissionUsecase` interface, `CommissionHandler` struct, request/response types (`createCommissionRequest`, `updateCommissionRequest`, `artpieceSummary`, `commissionResponse`); `commissionResponse` includes artpieces as `[]artpieceSummary{id, thumbnail_url}`
- [x] 5.2 Implement `CreateCommission`, `GetCommissionByID`, `ListCommissions`, `UpdateCommission`, `DeleteCommission` — presign thumbnail URLs for artpiece summaries in Get/Create response
- [x] 5.3 Implement `AttachArtpieces`, `DetachArtpieces`, `ReplaceArtpieces` handlers
- [x] 5.4 Wire `CommissionHandler` in `main.go`: instantiate repository, usecase, handler; register routes under the `protected` group

## 6. Artpiece response update

- [x] 6.1 Add `CommissionID *string` to `artpieceResponse` and update `toArtpieceResponse` to map `a.CommissionID`

## 7. Bruno files

- [x] 7.1 Create Bruno request files for commission CRUD: Create Commission, Get Commission, List Commissions, Update Commission, Delete Commission
- [x] 7.2 Create Bruno request files for attachment endpoints: Attach Artpieces, Detach Artpieces, Replace Artpieces

## 8. Unit tests

- [x] 8.1 Create `usecase/commission_usecase_test.go` with fakes for all commission repo interfaces; write Create tests: success with artpieces, artist not owned, character not owned, artpiece not owned, artpiece already attached to another commission
- [x] 8.2 Write Update tests: success, artist not owned, character not owned, commission not found
- [x] 8.3 Write AttachArtpieces tests: success, artpiece not owned, artpiece already attached to another commission, artpiece already on this commission (no-op)
- [x] 8.4 Write DetachArtpieces tests: success, artpiece not attached to this commission (no-op)
- [x] 8.5 Write ReplaceArtpieces tests: success full replace, empty set clears all, artpiece not owned rejects without changes, artpiece belonging to another commission rejects without changes

## 9. Lint

- [x] 9.1 Run `golangci-lint run ./...` from the backend directory and fix all reported issues
