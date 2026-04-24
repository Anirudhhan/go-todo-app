BEGIN;

CREATE TYPE user_role AS ENUM ('user', 'admin');

ALTER TABLE users
    ADD COLUMN role user_role NOT NULL DEFAULT 'user';

ALTER TABLE users
    ADD COLUMN suspended_at TIMESTAMPTZ NULL;

COMMIT;