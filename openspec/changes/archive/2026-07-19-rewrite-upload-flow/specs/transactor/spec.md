## ADDED Requirements

### Requirement: Transactor
The system SHALL provide a `Transactor` that allows the usecase layer to coordinate writes across multiple repositories in a single atomic database transaction, without exposing a database handle to the usecase.

The `Transactor` interface SHALL be defined in the `usecase` package:
```
InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
```

The concrete implementation SHALL store the active transaction in the context under a package-private key. Repository methods SHALL read the transaction from context if present, falling back to the live database connection if not. This means the same repository methods work correctly both inside and outside a transaction.

#### Scenario: Operations inside a transaction are atomic
- **WHEN** two repository writes are performed inside `InTransaction` and one fails
- **THEN** both writes are rolled back; no partial state is committed to the database

#### Scenario: Repository methods work outside a transaction
- **WHEN** a repository method is called with a context that contains no active transaction
- **THEN** the method executes against the live database connection with no change in behavior
