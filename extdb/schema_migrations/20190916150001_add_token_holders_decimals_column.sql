-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE token_holders ADD COLUMN IF NOT EXISTS decimals smallint;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
