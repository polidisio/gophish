-- +goose Up
ALTER TABLE smtp ADD COLUMN ignore_cert BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE smtp DROP COLUMN IF EXISTS ignore_cert;
