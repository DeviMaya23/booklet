CREATE TABLE image_metadata (
    file_id UUID PRIMARY KEY REFERENCES files(id) ON DELETE CASCADE,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL
);
