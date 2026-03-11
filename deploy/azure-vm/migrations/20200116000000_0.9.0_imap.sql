-- +goose Up
CREATE TABLE IF NOT EXISTS "imap" (
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
    host VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    inbox_folder VARCHAR(255) DEFAULT 'INBOX',
    last_run_mailbox VARCHAR(255),
    last_run TIMESTAMP,
    enabled BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS "imap";
