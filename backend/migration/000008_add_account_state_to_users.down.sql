ALTER TABLE users
    ADD COLUMN is_pending_deletion BOOLEAN NOT NULL DEFAULT false;

UPDATE users SET is_pending_deletion = true WHERE account_state = 'pending_deletion';

ALTER TABLE users
    DROP COLUMN account_state,
    DROP COLUMN purged_at;
