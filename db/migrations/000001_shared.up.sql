CREATE TABLE IF NOT EXISTS countries (
    id      UUID PRIMARY KEY,
    name    TEXT NOT NULL,
    continent TEXT NOT NULL
);