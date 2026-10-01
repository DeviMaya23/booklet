CREATE TABLE pending_file_uploads (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    r2_key TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    artpiece_id UUID REFERENCES artpieces(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
