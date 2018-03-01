CREATE TABLE headers (
    block_number bigint UNIQUE,
    block_hash varchar UNIQUE,
    fields jsonb
);

CREATE TABLE bodies (
    block_number bigint UNIQUE,
    block_hash varchar UNIQUE,
    fields jsonb
);

CREATE TABLE pending_transactions (
    tx_hash varchar UNIQUE,
    status varchar,
    fields jsonb
);

CREATE TABLE receipts (
    block_number bigint UNIQUE,
    block_hash varchar UNIQUE,
    fields jsonb
);

CREATE TABLE accounts (
    block_number bigint,
    block_hash varchar,
    address varchar,
    fields jsonb,
    UNIQUE(block_number, address),
    UNIQUE(block_hash, address)
);

CREATE TABLE rewards (
    block_number bigint,
    block_hash varchar,
    address varchar,
    fields jsonb,
    UNIQUE(block_number, address),
    UNIQUE(block_hash, address)
);

CREATE TABLE internal_transactions (
    block_number bigint,
    type varchar,
    timestamp bigint,
    fields jsonb,
    UNIQUE(block_number, timestamp)
);

CREATE INDEX ON accounts(address);

CREATE INDEX ON rewards(address);

CREATE INDEX ON pending_transactions(status);

CREATE INDEX ON internal_transactions(type);

CREATE INDEX ON internal_transactions(block_number);
