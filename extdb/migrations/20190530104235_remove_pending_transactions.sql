-- +goose Up
-- SQL in this section is executed when the migration is applied.
DELETE FROM pending_transactions WHERE id>0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
