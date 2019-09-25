-- +goose Up
-- SQL in this section is executed when the migration is applied.
DROP INDEX IF EXISTS internal_transactions_uniq;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
