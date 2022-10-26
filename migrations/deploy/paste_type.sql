-- Deploy microbin:paste_type to pg

BEGIN;

CREATE TYPE paste_type AS ENUM ('file', 'text', 'url');

COMMIT;
