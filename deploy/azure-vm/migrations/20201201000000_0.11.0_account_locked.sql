-- +goose Up
ALTER TABLE users ADD COLUMN account_locked BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS account_locked;
