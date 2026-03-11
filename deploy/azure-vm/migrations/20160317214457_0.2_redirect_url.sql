-- +goose Up
ALTER TABLE pages ADD COLUMN redirect_url VARCHAR(255);

-- +goose Down
ALTER TABLE pages DROP COLUMN IF EXISTS redirect_url;
