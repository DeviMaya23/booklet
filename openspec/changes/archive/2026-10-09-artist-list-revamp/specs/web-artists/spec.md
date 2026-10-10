## MODIFIED Requirements

### Requirement: Artists list display
The artists page at `/app/artists` SHALL fetch the authenticated user's artists from `GET /artists` and display them in a table with columns: Name, Links, Notes, and an edit action. Each row SHALL display the artist's `name` in bold, a links chip group, the artist's `notes` truncated with ellipsis (or "—" if absent), and an edit icon button.

#### Scenario: Artists load and display
- **WHEN** an authenticated user navigates to `/app/artists`
- **THEN** the page SHALL fetch `GET /artists` and render one table row per artist

#### Scenario: Empty list
- **WHEN** the user has no artists
- **THEN** the page SHALL render an empty table body with no rows

---

## ADDED Requirements

### Requirement: Artist link chips
Each artist row's Links cell SHALL display the artist's links as chips, showing the hostname of each URL (e.g. `bsky.app` for `https://bsky.app/@artist`). Chips SHALL be ordered with the primary link first. A maximum of 3 chips SHALL be visible; if the artist has more than 3 links, a non-interactive `+N` chip SHALL be shown after the third, where N is the count of hidden links. If the artist has no links, the cell SHALL display the text "No links" in a muted style. Each chip SHALL be an anchor that opens the full URL in a new tab.

The primary link chip SHALL be visually distinguished with a star icon and an outlined border. All other chips SHALL render in a filled/muted style without a star.

If a URL cannot be parsed by `new URL()`, the chip SHALL display the raw URL string as its label.

#### Scenario: Chips show hostname only
- **WHEN** an artist has links
- **THEN** each visible chip SHALL display only the hostname portion of the URL, not the full URL

#### Scenario: Primary link chip is first and distinguished
- **WHEN** an artist has a primary link alongside other links
- **THEN** the primary link chip SHALL be rendered first and SHALL show a star icon with an outlined border style

#### Scenario: Maximum 3 chips shown with overflow badge
- **WHEN** an artist has more than 3 links
- **THEN** exactly 3 chips SHALL be visible and a non-interactive `+N` chip SHALL indicate the number of hidden links

#### Scenario: No links placeholder
- **WHEN** an artist has no links
- **THEN** the Links cell SHALL render "No links" in muted text instead of chips

#### Scenario: Chip opens URL in new tab
- **WHEN** a user clicks a link chip
- **THEN** the full URL SHALL open in a new browser tab

---

### Requirement: Artist link chip tooltip
Hovering a link chip SHALL show a tooltip containing the full URL on the first line and a descriptor on the second line. For the primary link the descriptor SHALL read "Main link · opens in a new tab". For all other links the descriptor SHALL read "Opens in a new tab".

#### Scenario: Tooltip on primary link chip
- **WHEN** a user hovers the primary link chip
- **THEN** a tooltip SHALL appear showing the full URL and "Main link · opens in a new tab"
- **Verified by:** test — hover interaction asserting tooltip content

#### Scenario: Tooltip on non-primary link chip
- **WHEN** a user hovers a non-primary link chip
- **THEN** a tooltip SHALL appear showing the full URL and "Opens in a new tab"
- **Verified by:** test — hover interaction asserting tooltip content

---

### Requirement: Artist notes truncation
The Notes column SHALL display the artist's `notes` value. When the text is too long to fit the available column width it SHALL be clipped with a trailing ellipsis. When the artist has no notes, the cell SHALL display "—".

#### Scenario: Notes truncated when too long
- **WHEN** an artist's notes text exceeds the available column width
- **THEN** the text SHALL be clipped with a trailing "…"
- **Verified by:** manual smoke — CSS overflow not assertable in jsdom

#### Scenario: No notes placeholder
- **WHEN** an artist has no notes
- **THEN** the Notes cell SHALL display "—"

---

## REMOVED Requirements

### Requirement: Copy artist link
**Reason**: Replaced by link chips that open directly in a new tab. The copy-link affordance is no longer needed as the primary link is now a clickable anchor.
**Migration**: Users who previously copied the primary link by clicking the icon can now click the primary chip directly to open it, or copy from the browser address bar.
