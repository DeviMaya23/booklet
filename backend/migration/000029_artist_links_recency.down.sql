DROP TABLE IF EXISTS artist_links;

ALTER TABLE artists
    DROP COLUMN IF EXISTS last_used_at,
    ADD COLUMN artist_link TEXT NULL;
