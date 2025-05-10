BEGIN;

DROP TABLE IF EXISTS users, pastes, pastes_passwords;
DROP TABLE IF EXISTS token;

DROP TYPE IF EXISTS access_type;
DROP EXTENSION IF EXISTS pgcrypto;

DROP INDEX IF EXISTS idx_pastes_passwords;

COMMIT;