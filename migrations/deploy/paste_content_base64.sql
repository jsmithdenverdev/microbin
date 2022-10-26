-- Deploy microbin:paste_content_base64 to pg
-- requires: paste

BEGIN;

SET
client_min_messages = 'warning';

CREATE
OR REPLACE FUNCTION paste_content_base64(paste_row paste) RETURNS TEXT
    STABLE
    LANGUAGE SQL
AS
$$
SELECT encode("content", 'base64')
FROM paste P
    $$;

ALTER FUNCTION paste_content_base64(paste) OWNER TO postgres;

COMMIT;
