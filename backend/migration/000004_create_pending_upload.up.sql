CREATE TABLE pending_uploads (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    r2_key TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    title TEXT,
    artist_name TEXT,
    artist_link TEXT,
    notes TEXT,
    character_ids JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_pending_uploads_user_id ON pending_uploads (user_id);
