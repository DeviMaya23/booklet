# Spec Delta

## ADDED Requirements

### Requirement: Edit mode — is_reference checkbox
The detail view in edit mode SHALL render a "This is a reference" checkbox. The checkbox reflects the artpiece's current `is_reference` value when edit mode is entered, and its state is included when the user saves.

#### Scenario: Checkbox initialized from artpiece state
- **WHEN** the user enters edit mode on an artpiece where is_reference is true
- **THEN** the "This is a reference" checkbox is checked

#### Scenario: Checkbox initialized unchecked for normal artpiece
- **WHEN** the user enters edit mode on an artpiece where is_reference is false
- **THEN** the "This is a reference" checkbox is unchecked

#### Scenario: Save with checkbox checked
- **WHEN** the user checks "This is a reference" and saves
- **THEN** PUT /artpieces/:id is called with is_reference: true

#### Scenario: Save with checkbox unchecked
- **WHEN** the user unchecks "This is a reference" and saves
- **THEN** PUT /artpieces/:id is called with is_reference: false

#### Scenario: Unsaved changes — checkbox contributes to dirty check
- **WHEN** the user changes only the is_reference checkbox and the route blocker evaluates unsaved changes
- **THEN** the blocker treats the form as having unsaved changes
