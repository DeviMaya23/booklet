# Spec Delta

## ADDED Requirements

### Requirement: Source image size validation before crop
The system SHALL validate the source image's size against the image category size limit before opening the crop dialog. If the source image exceeds the limit, the crop dialog SHALL NOT open and the system SHALL display a toast identifying the limit. This applies to all entry points: file picker, drag-and-drop, and clipboard paste.

#### Scenario: Oversized source image rejected
- **WHEN** the user selects, drops, or pastes a source image that exceeds the image category size limit
- **THEN** the crop dialog does not open and a toast is shown naming the size limit

#### Scenario: Source image within limit
- **WHEN** the user selects, drops, or pastes a source image within the image category size limit
- **THEN** the crop dialog opens normally
