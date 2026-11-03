ALTER TABLE characters
    ADD COLUMN last_used_at TIMESTAMPTZ NULL;

UPDATE characters SET last_used_at = created_at;
