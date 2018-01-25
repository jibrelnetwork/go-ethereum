CREATE TABLE block_headers (
    block_number integer,
    block_hash varchar,
    fields jsonb
);

CREATE TABLE block_bodies (
    block_number integer,
    block_hash varchar,
    fields jsonb
);