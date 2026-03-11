-- +goose Up
ALTER TABLE events ALTER COLUMN time TYPE TIMESTAMP;
ALTER TABLE templates ALTER COLUMN modified_date TYPE TIMESTAMP;
ALTER TABLE pages ALTER COLUMN modified_date TYPE TIMESTAMP;
ALTER TABLE groups ALTER COLUMN modified_date TYPE TIMESTAMP;
ALTER TABLE campaigns ALTER COLUMN created_date TYPE TIMESTAMP;
ALTER TABLE campaigns ALTER COLUMN completed_date TYPE TIMESTAMP;

-- +goose Down
-- No rollback needed for timezone changes
