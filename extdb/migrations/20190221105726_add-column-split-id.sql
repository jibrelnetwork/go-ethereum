-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE reorgs ADD COLUMN IF NOT EXISTS split_id INTEGER;

CREATE INDEX IF NOT EXISTS reorgs_split_id_idx ON reorgs(split_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
