-- +goose Up
-- Entra ID SSO configuration table

CREATE TABLE IF NOT EXISTS entra_id_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    enabled INTEGER DEFAULT 0,
    client_id VARCHAR(255),
    client_secret VARCHAR(512),
    tenant_id VARCHAR(255),
    redirect_uri VARCHAR(512),
    scopes VARCHAR(512) DEFAULT 'openid profile email User.Read',
    admin_only INTEGER DEFAULT 0,
    auto_create_users INTEGER DEFAULT 0,
    sync_departments INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Add Entra ID and Department to users table
ALTER TABLE users ADD COLUMN entra_id VARCHAR(255);
ALTER TABLE users ADD COLUMN department VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_users_entra_id ON users(entra_id);
CREATE INDEX IF NOT EXISTS idx_users_department ON users(department);

-- +goose Down
DROP INDEX IF EXISTS idx_users_entra_id;
DROP INDEX IF EXISTS idx_users_department;
ALTER TABLE users DROP COLUMN department;
ALTER TABLE users DROP COLUMN entra_id;
DROP TABLE IF EXISTS entra_id_settings;
