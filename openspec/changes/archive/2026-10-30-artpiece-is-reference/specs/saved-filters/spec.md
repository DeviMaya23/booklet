# Spec Delta

## ADDED Requirements

### Requirement: filter_payload includes is_reference
The `filter_payload` object stored in a saved filter SHALL support an `is_reference` boolean field. When the field is absent on read it SHALL be treated as false. The field is accepted on create and update without additional validation beyond type.

#### Scenario: Save filter with is_reference true
- **WHEN** a user creates a saved filter with filter_payload containing is_reference: true
- **THEN** the system stores the payload and returns the record with filter_payload.is_reference: true

#### Scenario: Save filter with is_reference false
- **WHEN** a user creates a saved filter with filter_payload containing is_reference: false
- **THEN** the system stores the payload and returns the record with filter_payload.is_reference: false

#### Scenario: Saved filter without is_reference — treated as false
- **WHEN** a saved filter's filter_payload does not contain an is_reference field
- **THEN** the system treats is_reference as false when reading that filter
