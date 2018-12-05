CREATE TABLE headers (
    block_number bigint,
    block_hash varchar(70) UNIQUE,
    fields jsonb
);

CREATE TABLE bodies (
    block_number bigint,
    block_hash varchar(70) UNIQUE,
    fields jsonb
);

CREATE TABLE pending_transactions (
    tx_hash varchar(70) UNIQUE,
    status varchar,
    fields jsonb
);

CREATE TABLE receipts (
    block_number bigint,
    block_hash varchar(70) UNIQUE,
    fields jsonb
);

CREATE TABLE accounts (
    id bigserial primary key,
    block_number bigint,
    block_hash varchar(70),
    address varchar(45),
    fields jsonb,
    UNIQUE(block_hash, address)
);

CREATE TABLE rewards (
    id bigserial primary key,
    block_number bigint,
    block_hash varchar(70) UNIQUE,
    address varchar(45),
    fields jsonb
);

CREATE TABLE internal_transactions (
    id bigserial primary key,
    block_number bigint,
    block_hash varchar(70),
    parent_tx_hash varchar(70),
    index bigint,
    type varchar(20),
    timestamp bigint,
    fields jsonb
);

CREATE TABLE chain_splits (
    id bigserial primary key,
    common_block_number bigint,
    common_block_hash varchar(70),
    drop_length bigint,
    drop_block_hash varchar(70),
    add_length bigint,
    add_block_hash varchar(70),
    node_id varchar(70)
);

CREATE TABLE reorgs (
    id bigserial primary key,
    block_number bigint,
    block_hash varchar(70),
    header jsonb,
    reinserted boolean,
    node_id varchar(70)
);


CREATE INDEX ON internal_transactions(block_number);

CREATE UNIQUE INDEX internal_transactions_uniq ON internal_transactions(block_hash, parent_tx_hash, index, type);

CREATE INDEX ON pending_transactions(status);

CREATE INDEX ON internal_transactions(type);

CREATE INDEX ON internal_transactions(block_hash);

CREATE INDEX ON rewards(block_number);

CREATE INDEX ON rewards(block_hash);

CREATE INDEX ON accounts(block_number);

CREATE INDEX ON accounts(block_hash);

CREATE INDEX ON receipts(block_number);

CREATE INDEX ON receipts(block_hash);

CREATE INDEX ON bodies(block_number);

CREATE INDEX ON bodies(block_hash);

CREATE INDEX ON headers(block_number);

CREATE INDEX ON headers(block_hash);

CREATE INDEX ON chain_splits(common_block_number);

CREATE INDEX ON reorgs(block_number);

CREATE INDEX ON reorgs(block_hash);

CREATE INDEX ON reorgs(reinserted);
