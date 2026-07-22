CREATE TABLE character_folders (
    character_id UUID NOT NULL REFERENCES characters(id),
    folder_id    UUID NOT NULL,
    PRIMARY KEY (character_id, folder_id)
);
