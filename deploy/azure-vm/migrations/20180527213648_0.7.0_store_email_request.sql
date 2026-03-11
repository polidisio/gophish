-- +goose Up
ALTER TABLE results ADD COLUMN request TEXT;

-- +goose Down
ALTER TABLE results DROP COLUMN IF EXISTS request;
