-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS token_holders (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    token_address varchar(45) NOT NULL,
    holder_address varchar(45) NOT NULL,
    balance numeric NOT NULL,
    UNIQUE(block_hash, token_address, holder_address)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
