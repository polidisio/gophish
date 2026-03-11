-- +goose Up
ALTER TABLE campaigns ADD COLUMN send_by_date TIMESTAMP;

-- +goose Down
ALTER TABLE campaigns DROP COLUMN IF EXISTS send_by_date;
