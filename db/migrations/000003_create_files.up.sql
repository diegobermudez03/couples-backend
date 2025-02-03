CREATE TABLE IF NOT EXISTS files(
    id          UUID PRIMARY KEY,
    bucket      TEXT NOT NULL,
    grouping       TEXT NOT NULL,
    object_key  TEXT NOT NULL,
    url         TEXT,
    public      BOOLEAN NOT NULL,
    type        TEXT NOT NULL
);