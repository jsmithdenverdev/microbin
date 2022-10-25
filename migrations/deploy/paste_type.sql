-- Deploy microbin:paste_type to pg
-- requires: appschema

BEGIN;

CREATE TYPE microbin.paste_type AS ENUM ('file', 'text', 'url');

COMMIT;
