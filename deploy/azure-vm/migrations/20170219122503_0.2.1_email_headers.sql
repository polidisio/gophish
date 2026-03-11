-- +goose Up
ALTER TABLE results ADD COLUMN headers TEXT;

-- +goose Down
ALTER TABLE results DROP COLUMN IF EXISTS headers;
