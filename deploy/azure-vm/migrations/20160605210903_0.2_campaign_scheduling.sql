-- +goose Up
ALTER TABLE campaigns ADD COLUMN send_by_date TIMESTAMP;
ALTER TABLE campaigns ADD COLUMN scheduled_count INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE campaigns DROP COLUMN IF EXISTS scheduled_count;
ALTER TABLE campaigns DROP COLUMN IF EXISTS send_by_date;
