-- +goose Up
CREATE TABLE IF NOT EXISTS "webhooks" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    url VARCHAR(1000),
    secret VARCHAR(255),
    is_active BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS "webhooks";
