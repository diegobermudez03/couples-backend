CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    gender          TEXT NOT NULL,
    birth_date      DATE NOT NULL,
    created_at      TIMESTAMP  NOT NULL,
    language        TEXT NOT NULL,
    active          BOOLEAN NOT NULL,
    country_id      UUID REFERENCES countries(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS couples(
    id          UUID PRIMARY KEY, 
    he_id       UUID REFERENCES users(id) NOT NULL,
    she_id      UUID REFERENCES users(id) NOT NULL,
    relation_start  DATE NOT NULL,
    end_date        DATE
);

CREATE TABLE IF NOT EXISTS temp_couples(
    user_id     UUID REFERENCES users(id),
    code        INTEGER NOT NULL UNIQUE,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    PRIMARY KEY(user_id)
);

CREATE TABLE IF NOT EXISTS points(
    id              UUID PRIMARY KEY,
    points          INTEGER NOT NULL,
    starting_day    DATE NOT NULL,
    ending_day      DATE NOT NULL,
    user_id         UUID REFERENCES users(id),
    couple_id       UUID REFERENCES couples(id)
);

CREATE TABLE IF NOT EXISTS couple_levels(
    id          UUID PRIMARY KEY,
    level_name  TEXT NOT NULL,
    level_description TEXT NOT NULL,
    min_month_points  INTEGER NOT NULL,
    image_url         TEXT NOT NULL
);