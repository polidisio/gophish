-- +goose Up
ALTER TABLE users ADD COLUMN password_changed TIMESTAMP;
ALTER TABLE users ADD COLUMN password_change_required BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS password_change_required;
ALTER TABLE users DROP COLUMN IF EXISTS password_changed;
