CREATE TABLE IF NOT EXISTS files(
    id          UUID PRIMARY KEY,
    bucket      TEXT NOT NULL,
    object_key  TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    type        TEXT NOT NULL
);