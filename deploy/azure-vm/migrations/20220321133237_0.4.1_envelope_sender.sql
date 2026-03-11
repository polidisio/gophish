-- +goose Up
ALTER TABLE smtp ADD COLUMN from_address VARCHAR(255) DEFAULT '';

-- +goose Down
ALTER TABLE smtp DROP COLUMN IF EXISTS from_address;
