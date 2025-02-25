
CREATE TABLE user_auth (
    id SERIAL PRIMARY KEY,
    user_code uuid NOT NULL,
    user_name VARCHAR NOT NULL UNIQUE,
    token_hash TEXT NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS user_code_index
    on user_auth (user_code);

