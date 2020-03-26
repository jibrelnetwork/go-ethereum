-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE pending_transactions ADD COLUMN IF NOT EXISTS replaced_by_hash varchar(70);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
