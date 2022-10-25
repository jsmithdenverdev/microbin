-- Deploy microbin:paste_content_base64 to pg
-- requires: paste

BEGIN;

SET client_min_messages = 'warning';

create or replace function microbin.paste_content_base64(paste_row microbin.paste) returns text
    stable
    language sql
as
$$
select encode("content", 'base64')
from microbin.paste p
$$;

alter function microbin.paste_content_base64(microbin.paste) owner to postgres;

COMMIT;
