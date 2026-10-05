## ADDED Requirements

### Requirement: Attach artpieces to a commission
The system SHALL allow an authenticated user to attach one or more artpieces to a commission in a single request. Every artpiece in the request must be owned by the same user as the commission. An artpiece that already belongs to a different commission SHALL be rejected.

#### Scenario: Successful bulk attach
- **WHEN** a user posts a list of artpiece_ids that all belong to them and have no commission_id set
- **THEN** the system sets commission_id on each artpiece and returns the updated commission with HTTP 200

#### Scenario: Artpiece not owned by user
- **WHEN** a user attaches an artpiece_id that does not belong to them
- **THEN** the system returns HTTP 422 with an ownership error and attaches nothing

#### Scenario: Artpiece already attached to a different commission
- **WHEN** a user attaches an artpiece_id whose commission_id points to a different commission
- **THEN** the system returns HTTP 422 with a conflict error and attaches nothing

#### Scenario: Artpiece already attached to the same commission
- **WHEN** a user attaches an artpiece_id that is already attached to this commission
- **THEN** the system treats it as a no-op for that artpiece and returns HTTP 200

---

### Requirement: Detach artpieces from a commission
The system SHALL allow an authenticated user to detach one or more artpieces from a commission in a single request. Detaching sets commission_id to NULL on the artpieces; it does not delete them.

#### Scenario: Successful bulk detach
- **WHEN** a user deletes a list of artpiece_ids that are attached to the commission
- **THEN** the system sets commission_id to NULL on each artpiece and returns HTTP 200

#### Scenario: Artpiece not attached to this commission
- **WHEN** a user detaches an artpiece_id that is not attached to this commission
- **THEN** the system treats it as a no-op for that artpiece and returns HTTP 200

#### Scenario: Commission not found or belongs to another user
- **WHEN** a user detaches from a commission that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: Replace artpiece set on a commission
The system SHALL allow an authenticated user to replace the full set of artpieces on a commission in a single transaction. The system reconciles server-side: artpieces in the new set but not the current set are attached; artpieces in the current set but not the new set are detached; artpieces in both are left unchanged.

#### Scenario: Successful full replace
- **WHEN** a user puts a complete list of artpiece_ids, all owned by the user and none attached to a different commission
- **THEN** the system reconciles the artpiece set in one transaction and returns the updated commission with HTTP 200

#### Scenario: Empty replace clears all artpieces
- **WHEN** a user puts an empty artpiece_ids array
- **THEN** the system detaches all current artpieces and returns the updated commission with HTTP 200

#### Scenario: Artpiece not owned by user
- **WHEN** the replacement set includes an artpiece_id not owned by the user
- **THEN** the system returns HTTP 422 and makes no changes

#### Scenario: Artpiece already attached to a different commission
- **WHEN** the replacement set includes an artpiece_id whose commission_id points to a different commission
- **THEN** the system returns HTTP 422 and makes no changes

---

### Requirement: Commission GET returns artpiece summaries
The system SHALL include in the commission GET response a list of attached artpieces, each containing the artpiece id and a presigned cover thumbnail URL (null if the artpiece has no cover or the cover has no thumbnail).

#### Scenario: Commission with attached artpieces
- **WHEN** a user fetches a commission that has attached artpieces
- **THEN** each artpiece in the response has an id and a thumbnail_url (null if no thumbnail)

#### Scenario: Commission with no artpieces
- **WHEN** a user fetches a commission with no attached artpieces
- **THEN** the artpieces array in the response is empty
