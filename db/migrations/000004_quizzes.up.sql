CREATE TABLE IF NOT EXISTS quiz_categories(
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    image_id    UUID REFERENCES FILES(id) NOT NULL,
    active      BOOLEAN NOT NULL
);