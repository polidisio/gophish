-- +goose Up
ALTER TABLE results ADD COLUMN sent_at TIMESTAMP;
ALTER TABLE results ADD COLUMN opened_at TIMESTAMP;
ALTER TABLE results ADD COLUMN clicked_at TIMESTAMP;
ALTER TABLE results ADD COLUMN submitted_at TIMESTAMP;

-- +goose Down
ALTER TABLE results DROP COLUMN IF EXISTS submitted_at;
ALTER TABLE results DROP COLUMN IF EXISTS clicked_at;
ALTER TABLE results DROP COLUMN IF EXISTS opened_at;
ALTER TABLE results DROP COLUMN IF EXISTS sent_at;
