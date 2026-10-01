CREATE TABLE files (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    artpiece_id UUID REFERENCES artpieces(id) ON DELETE SET NULL,
    file_r2_path TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    thumbnail_r2_path TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_files_user_id ON files (user_id);
