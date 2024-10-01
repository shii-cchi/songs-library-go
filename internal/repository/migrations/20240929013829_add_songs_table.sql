-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link TEXT,
    CONSTRAINT unique_group_song UNIQUE ("group", song)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songs;
-- +goose StatementEnd
