## MODIFIED Requirements

### Requirement: New Commission button
The page SHALL display a "+ New Commission" button in the header. Clicking the button opens the `CommissionFormModal` in create mode.

#### Scenario: Button present
- **WHEN** a user views the commissions page
- **THEN** a "+ New Commission" button is visible in the page header

#### Scenario: Button opens create modal
- **WHEN** a user clicks the "+ New Commission" button
- **THEN** the `CommissionFormModal` opens in create mode with all fields empty

---

## ADDED Requirements

### Requirement: Row action controls
Each commission row SHALL display a brush icon and a trash icon at the end of the row (after the Last Contact / Time Taken column). The brush icon opens the `CommissionFormModal` in edit mode for that commission. The trash icon opens the `DeleteCommissionDialog` for that commission.

#### Scenario: Brush icon opens edit modal
- **WHEN** a user clicks the brush icon on a commission row
- **THEN** the `CommissionFormModal` opens in edit mode pre-filled with that commission's data

#### Scenario: Trash icon opens delete dialog
- **WHEN** a user clicks the trash icon on a commission row
- **THEN** the `DeleteCommissionDialog` opens for that commission
