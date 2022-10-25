-- Deploy microbin:paste to pg
-- requires: paste_type

BEGIN;

SET
client_min_messages = 'warning';

CREATE TABLE microbin.paste
(
    id         SERIAL PRIMARY KEY,
    name       text                NOT NULL,
    type       microbin.paste_type NOT NULL,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    content    BYTEA               NOT NULL,
    expiration INTERVAL            NOT NULL,
    metadata   JSON NULL
);

COMMIT;
