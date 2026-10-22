# Frontend Constraints

Critical constraints that contradict common defaults. Full conventions in `/CONVENTIONS.md`.

## UI Primitives

- UI components use **Base UI** (`@base-ui/react`) as the primitive layer, not Radix.
- **Never** import from `@radix-ui/*` packages.
- When in doubt about a component API, check Base UI docs — not Radix docs.
- Exception for `ContextMenuItem`: use `onClick`, not `onSelect`.

## Linting

- Never suppress ESLint findings with `eslint-disable` comments — fix the underlying issue.
- Exception allowed: `react-hooks/exhaustive-deps` (add a comment explaining why).
