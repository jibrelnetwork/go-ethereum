CREATE TABLE headers (
    block_number bigint UNIQUE,
    block_hash varchar(70),
    fields jsonb
);

CREATE TABLE bodies (
    block_number bigint UNIQUE,
    block_hash varchar(70),
    fields jsonb
);

CREATE TABLE pending_transactions (
    tx_hash varchar(70) UNIQUE,
    status varchar,
    fields jsonb
);

CREATE TABLE receipts (
    block_number bigint UNIQUE,
    block_hash varchar(70),
    fields jsonb
);

CREATE TABLE accounts (
    block_number bigint,
    block_hash varchar(70),
    address varchar(45),
    fields jsonb,
    UNIQUE(block_number, address)
);

CREATE TABLE rewards (
    block_number bigint,
    block_hash varchar(70),
    address varchar(45),
    fields jsonb,
    UNIQUE(block_number, address)
);

CREATE TABLE internal_transactions (
    block_number bigint,
    type varchar(20),
    timestamp bigint,
    fields jsonb
);

CREATE INDEX ON internal_transactions(block_number);

CREATE INDEX ON pending_transactions(status);

CREATE INDEX ON internal_transactions(type);

CREATE INDEX ON internal_transactions(block_number);
