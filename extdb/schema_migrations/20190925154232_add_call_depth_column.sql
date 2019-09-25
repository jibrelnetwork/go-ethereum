-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE internal_transactions ADD COLUMN IF NOT EXISTS call_depth INTEGER;
ALTER TABLE internal_transactions ALTER COLUMN call_depth SET DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
