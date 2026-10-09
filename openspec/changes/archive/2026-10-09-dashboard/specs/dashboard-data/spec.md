## ADDED Requirements

### Requirement: Dashboard endpoint
The system SHALL expose `GET /dashboard` as a protected endpoint returning all data needed to render the dashboard page. The response SHALL be user-scoped (only data belonging to the authenticated user). The endpoint SHALL run all queries sequentially and return a single JSON object.

#### Scenario: Authenticated user fetches dashboard data
- **WHEN** an authenticated user sends `GET /dashboard`
- **THEN** the system SHALL respond with HTTP 200 and a JSON body containing `recent_artpieces`, `in_progress`, and `housekeeping` fields

#### Scenario: Unauthenticated request is rejected
- **WHEN** an unauthenticated request is sent to `GET /dashboard`
- **THEN** the system SHALL respond with HTTP 401

### Requirement: Recent artpieces
The `recent_artpieces` field SHALL contain the five most recently created artpieces for the user, ordered by `created_at` descending. Each entry SHALL include: `id`, `title`, `artist_name` (nullable), and `thumbnail_url` (nullable, deterministic presigned URL). If the user has fewer than five artpieces, all are returned.

#### Scenario: Returns up to five artpieces ordered by created_at
- **WHEN** the user has more than five artpieces
- **THEN** `recent_artpieces` SHALL contain exactly five entries, ordered newest first by `created_at`

#### Scenario: Returns all artpieces when fewer than five exist
- **WHEN** the user has fewer than five artpieces
- **THEN** `recent_artpieces` SHALL contain all of them

#### Scenario: Thumbnail URL is null when no cover file is set
- **WHEN** an artpiece has no cover file
- **THEN** its `thumbnail_url` SHALL be null

### Requirement: In-progress commissions
The `in_progress` field SHALL contain all commissions with `status` of `waitlist` or `wip`, ordered by `created_at` ascending (oldest first). Each entry SHALL include: `id`, `title`, `artist_name` (nullable), `artist_link` (nullable), `status`, `paid`, `last_contacted_at` (nullable ISO timestamp), and `created_at`.

#### Scenario: Returns only waitlist and wip commissions
- **WHEN** the user has commissions in all three statuses
- **THEN** `in_progress` SHALL contain only commissions with `status` of `waitlist` or `wip`

#### Scenario: Returns empty array when no in-progress commissions exist
- **WHEN** the user has no commissions with status waitlist or wip
- **THEN** `in_progress` SHALL be an empty array

### Requirement: Housekeeping data
The `housekeeping` field SHALL be an object with four arrays. Items in each array SHALL contain only `id` and `title`. All four queries are always run regardless of whether any items exist.

`commissions_no_artist`: all commissions where `artist_id` is NULL, across all statuses.
`artpieces_no_artist`: all artpieces where `artist_id` is NULL.
`done_no_artpieces`: all commissions where `status = 'done'` and no artpieces are attached to that commission.
`artpieces_no_files`: all artpieces where no files exist in the `files` table with a matching `artpiece_id` (NOT EXISTS subquery).

#### Scenario: Commission with no artist appears in commissions_no_artist
- **WHEN** a commission has a null artist_id regardless of its status
- **THEN** it SHALL appear in `housekeeping.commissions_no_artist`

#### Scenario: Done commission with artpieces does not appear in done_no_artpieces
- **WHEN** a done commission has at least one artpiece attached
- **THEN** it SHALL NOT appear in `housekeeping.done_no_artpieces`

#### Scenario: Artpiece with files does not appear in artpieces_no_files
- **WHEN** an artpiece has at least one file in the files table
- **THEN** it SHALL NOT appear in `housekeeping.artpieces_no_files`

#### Scenario: All housekeeping arrays are empty when nothing needs fixing
- **WHEN** there are no housekeeping issues
- **THEN** all four arrays SHALL be empty
