ALTER TABLE pending_uploads
    DROP COLUMN artist_id,
    ADD COLUMN artist_name TEXT,
    ADD COLUMN artist_link TEXT;

ALTER TABLE images
    DROP COLUMN artist_id,
    ADD COLUMN artist_name TEXT,
    ADD COLUMN artist_link TEXT;

DROP INDEX IF EXISTS idx_artists_user_id;
DROP TABLE IF EXISTS artists;
