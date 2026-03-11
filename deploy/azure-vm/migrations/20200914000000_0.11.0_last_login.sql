-- +goose Up
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS last_login;
