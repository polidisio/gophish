-- +migrate Up
-- OAuth Providers table
CREATE TABLE IF NOT EXISTS oauth_providers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(100) NOT NULL,
    redirect_url VARCHAR(500) NOT NULL,
    scopes VARCHAR(500) NOT NULL,
    enabled TINYINT(1) DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Azure AD Users table
CREATE TABLE IF NOT EXISTS azure_ad_users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    given_name VARCHAR(100),
    surname VARCHAR(100),
    job_title VARCHAR(100),
    department VARCHAR(100),
    office_location VARCHAR(255),
    photo LONGBLOB,
    manager_id VARCHAR(255),
    last_synced_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_email (email),
    INDEX idx_department (department)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Sync Status table
CREATE TABLE IF NOT EXISTS sync_status (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    provider_id BIGINT NOT NULL,
    started_at DATETIME NOT NULL,
    completed_at DATETIME,
    total_users INT DEFAULT 0,
    new_users INT DEFAULT 0,
    updated_users INT DEFAULT 0,
    deleted_users INT DEFAULT 0,
    failed_users INT DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (provider_id) REFERENCES oauth_providers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS sync_status;
DROP TABLE IF EXISTS azure_ad_users;
DROP TABLE IF EXISTS oauth_providers;
