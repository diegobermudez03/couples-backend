CREATE TABLE IF NOT EXISTS strategic_type_answers(
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS quiz_categories(
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    image_id    UUID REFERENCES FILES(id) NOT NULL,
    active      BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS quizzes(
    id              UUID PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL,
    language_code   TEXT NOT NULL,
    published       BOOLEAN NOT NULL,
    active          BOOLEAN NOT NULL,
    image_id        UUID REFERENCES files(id),
    created_at      TIMESTAMP NOT NULL,
    creator_id      UUID REFERENCES users(id),
    category_id     UUID REFERENCES quiz_categories(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS quiz_questions(
    id                  UUID PRIMARY KEY,
    ordering               INTEGER NOT NULL,
    question            TEXT NOT NULL,
    question_type       TEXT NOT NULL,
    options_json        JSONB NOT NULL,
    quiz_id             UUID REFERENCES quizzes(id) NOT NULL,
    strategic_answer_id    UUID REFERENCES strategic_type_answers(id)
);


CREATE TABLE IF NOT EXISTS user_answers(
    id              UUID PRIMARY KEY,
    user_id         UUID REFERENCES users(id) NOT NULL,
    question_id     UUID REFERENCES quiz_questions(id) NOT NULL,
    own_answer      TEXT NOT NULL,
    partner_answer  TEXT NOT NULL,
    answered_at     TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS quizzes_played(
    id              UUID PRIMARY KEY,
    quiz_id         UUID REFERENCES quizzes(id) NOT NULL,
    user_id         UUID REFERENCES userS(id) NOT NULL,
    shared          BOOLEAN NOT NULL,
    score           INTEGER,
    completed_at    TIMESTAMP,
    started_at      TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS challenges(
    id              UUID PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL,
    published       BOOLEAN NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    creator_id      UUID REFERENCES users(id) NOT NULL
);


CREATE TABLE IF NOT EXISTS challenges_played(
    id              UUID PRIMARY KEY,
    started_at      TIMESTAMP NOT NULL,
    completed_at    TIMESTAMP,
    score           INTEGER,
    shared          BOOLEAN NOT NULL,
    challenge_id    UUID REFERENCES challenges(id) NOT NULL,
    user_id         UUID REFERENCES users(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS chall_questions(
    id              UUID PRIMARY KEY,
    ordering           INTEGER NOT NULL,
    question        TEXT NOT NULL,
    question_type   TEXT NOT NULL,
    options_json    JSONB NOT NULL,
    challenge_id    UUID REFERENCES challenges(id) NOT NULL
);