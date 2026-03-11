-- +goose Up
ALTER TABLE events ADD COLUMN details TEXT;

-- +goose Down
ALTER TABLE events DROP COLUMN IF EXISTS details;
