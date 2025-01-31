CREATE TABLE IF NOT EXISTS users_auth(
    id                  UUID PRIMARY KEY,
    email               TEXT,
    hash                TEXT,
    oauth_provider      TEXT,
    oauth_id            TEXT,
    created_at          TIMESTAMP NOT NULL,
    user_id             UUID REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS sessions(
    id                  UUID PRIMARY KEY,
    token               TEXT NOT NULL,
    device              TEXT,
    os                  TEXT,
    expires_at          TIMESTAMP NOT NULL,
    created_at          TIMESTAMP NOT NULL,
    last_used           TIMESTAMP NOT NULL,
    user_auth_id        UUID REFERENCES users_auth(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS admin_sessions(
    id          UUID PRIMARY KEY,
    token       TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL
);