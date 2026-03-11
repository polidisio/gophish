-- +goose Up
ALTER TABLE results ADD COLUMN status VARCHAR(255) DEFAULT 'sent';
ALTER TABLE results ALTER COLUMN status DROP NOT NULL;

-- +goose Down
ALTER TABLE results ALTER COLUMN status SET NOT NULL;
