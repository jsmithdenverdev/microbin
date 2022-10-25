-- Deploy microbin:paste_type to pg
-- requires: appschema

BEGIN;

CREATE TYPE paste_type AS ENUM ('file', 'text', 'url');

COMMIT;
