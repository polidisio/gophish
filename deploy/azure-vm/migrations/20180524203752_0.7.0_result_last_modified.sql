-- +goose Up
ALTER TABLE results ALTER COLUMN modified_date TYPE TIMESTAMP;
ALTER TABLE results ADD COLUMN modified_date TIMESTAMP;

-- +goose Down
-- No rollback needed
