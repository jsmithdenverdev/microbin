-- Verify microbin:paste on pg

BEGIN;

SELECT id
     , name
     , type
     , created_at
     , updated_at
     , content
     , expiration
FROM microbin.paste
WHERE FALSE;
ROLLBACK;
