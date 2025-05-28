-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE gender_enum AS ENUM ('male', 'female');
CREATE TABLE IF NOT EXISTS users(
                                    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                    name TEXT NOT NULL,
                                    surname TEXT NOT NULL,
                                    patronymic TEXT NOT NULL DEFAULT '',
                                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    age INT NOT NULL,
                                    gender gender_enum NOT NULL,
                                    nationality TEXT NOT NULL ,
                                    version BIGINT NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
