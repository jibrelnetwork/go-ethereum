-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS headers (
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    fields jsonb
);

CREATE TABLE IF NOT EXISTS bodies (
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    fields jsonb
);

CREATE TABLE IF NOT EXISTS pending_transactions (
    tx_hash varchar(70) UNIQUE NOT NULL,
    status varchar,
    fields jsonb
);

CREATE TABLE IF NOT EXISTS receipts (
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    fields jsonb
);

CREATE TABLE IF NOT EXISTS accounts (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    address varchar(45) NOT NULL,
    fields jsonb,
    UNIQUE(block_hash, address)
);

CREATE TABLE IF NOT EXISTS rewards (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    address varchar(45) NOT NULL,
    fields jsonb
);

CREATE TABLE IF NOT EXISTS internal_transactions (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    parent_tx_hash varchar(70) NOT NULL,
    index bigint NOT NULL,
    type varchar(20) NOT NULL,
    timestamp bigint NOT NULL,
    fields jsonb
);

CREATE TABLE IF NOT EXISTS chain_splits (
    id bigserial primary key,
    common_block_number bigint NOT NULL,
    common_block_hash varchar(70) NOT NULL,
    drop_length bigint NOT NULL,
    drop_block_hash varchar(70) NOT NULL,
    add_length bigint NOT NULL,
    add_block_hash varchar(70) NOT NULL,
    node_id varchar(70)
);

CREATE TABLE IF NOT EXISTS reorgs (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    header jsonb,
    reinserted boolean NOT NULL,
    node_id varchar(70)
);


CREATE INDEX IF NOT EXISTS internal_transactions_block_number_idx ON internal_transactions(block_number);

CREATE UNIQUE INDEX IF NOT EXISTS internal_transactions_uniq ON internal_transactions(block_hash, parent_tx_hash, index, type);

CREATE INDEX IF NOT EXISTS internal_transactions_type_idx ON internal_transactions(type);

CREATE INDEX IF NOT EXISTS internal_transactions_block_hash_idx ON internal_transactions(block_hash);

CREATE INDEX IF NOT EXISTS pending_transactions_status_idx ON pending_transactions(status);

CREATE INDEX IF NOT EXISTS rewards_block_number_idx ON rewards(block_number);

CREATE INDEX IF NOT EXISTS rewards_block_hash_idx ON rewards(block_hash);

CREATE INDEX IF NOT EXISTS accounts_block_number_idx ON accounts(block_number);

CREATE INDEX IF NOT EXISTS accounts_block_hash_idx ON accounts(block_hash);

CREATE INDEX IF NOT EXISTS receipts_block_number_idx ON receipts(block_number);

CREATE INDEX IF NOT EXISTS receipts_block_hash_idx ON receipts(block_hash);

CREATE INDEX IF NOT EXISTS bodies_block_number_idx ON bodies(block_number);

CREATE INDEX IF NOT EXISTS bodies_block_hash_idx ON bodies(block_hash);

CREATE INDEX IF NOT EXISTS headers_block_number_idx ON headers(block_number);

CREATE INDEX IF NOT EXISTS headers_block_hash_idx ON headers(block_hash);

CREATE INDEX IF NOT EXISTS chain_splits_common_block_number_idx ON chain_splits(common_block_number);

CREATE INDEX IF NOT EXISTS reorgs_block_number_idx ON reorgs(block_number);

CREATE INDEX IF NOT EXISTS reorgs_block_hash_idx ON reorgs(block_hash);

CREATE INDEX IF NOT EXISTS reorgs_reinserted_idx ON reorgs(reinserted);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
