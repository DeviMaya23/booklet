# Spec Delta

## REMOVED Requirements

### Requirement: File response includes thumbnail_gen_state
**Reason**: `thumbnail_gen_state` is removed from the database and domain model as part of moving thumbnail generation client-side. The field no longer exists to expose.
**Migration**: Clients must not rely on `thumbnail_gen_state` in file API responses. Use `thumbnail_url` (null when no thumbnail is available) to determine thumbnail presence.
