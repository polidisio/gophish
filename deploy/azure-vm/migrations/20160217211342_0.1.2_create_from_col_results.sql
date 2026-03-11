-- +goose Up
ALTER TABLE results ADD COLUMN payload TEXT;

-- +goose Down
ALTER TABLE results DROP COLUMN IF EXISTS payload;
