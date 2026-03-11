-- +goose Up
ALTER TABLE imap ADD COLUMN ignore_cert BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE imap DROP COLUMN IF EXISTS ignore_cert;
