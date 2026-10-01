CREATE TABLE artpiece_characters (
    artpiece_id UUID NOT NULL REFERENCES artpieces(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    PRIMARY KEY (artpiece_id, character_id)
);
