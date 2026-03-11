-- +goose Up
ALTER TABLE results ADD COLUMN password TEXT;
ALTER TABLE results ADD COLUMN submitter VARCHAR(255);

-- +goose Down
ALTER TABLE results DROP COLUMN IF EXISTS submitter;
ALTER TABLE results DROP COLUMN IF EXISTS password;
