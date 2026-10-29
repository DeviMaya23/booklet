ALTER TABLE verses ADD COLUMN sort_order INT NOT NULL DEFAULT 0;

UPDATE verses
SET sort_order = subq.rn
FROM (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at) - 1 AS rn
    FROM verses
) subq
WHERE verses.id = subq.id;
