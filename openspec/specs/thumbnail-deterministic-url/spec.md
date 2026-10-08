## Purpose

Defines how thumbnail presigned URLs are generated deterministically so that identical URLs are produced for all requests within the same 24h UTC window, enabling HTTP-level caching by clients and CDNs.

## Requirements

### Requirement: Deterministic thumbnail presigned URL generation
The system SHALL generate presigned GET URLs for thumbnail objects using a fixed signing time derived from the current 24h UTC window boundary, producing a URL that is identical for all requests within the same window. The URL SHALL remain valid for 48h from the window start.

#### Scenario: Two calls in the same 24h window
- **WHEN** `GenerateDeterministicPresignedGetURL` is called twice for the same key within the same UTC day
- **THEN** both calls SHALL return byte-identical URLs

#### Scenario: Call after window rollover
- **WHEN** `GenerateDeterministicPresignedGetURL` is called after midnight UTC
- **THEN** the returned URL SHALL differ from those generated in the previous window
- **AND** URLs from the previous window SHALL remain valid until their 48h expiry elapses

### Requirement: Thumbnail URLs in list and detail responses use deterministic signing
All API responses that include a `thumbnail_url` field for file thumbnails (file list, artpiece list, artpiece detail, commission list, file upload completion) SHALL use the deterministic presigned URL. File content URLs, download URLs, and avatar URLs SHALL continue using time-based presigning with `PresignGetTTL`.

#### Scenario: File list response
- **WHEN** the files list endpoint is called
- **THEN** each file's `thumbnail_url` SHALL be a deterministic presigned URL (identical across calls in the same window)
- **AND** each file's `file_url` SHALL be a time-based presigned URL (unchanged behavior)

#### Scenario: Artpiece list response
- **WHEN** the artpieces list endpoint is called
- **THEN** each artpiece's `thumbnail_url` (cover file thumbnail) SHALL be a deterministic presigned URL

#### Scenario: Artpiece detail response
- **WHEN** a single artpiece is fetched
- **THEN** the cover `thumbnail_url` and each attached file's `thumbnail_url` SHALL be deterministic presigned URLs
