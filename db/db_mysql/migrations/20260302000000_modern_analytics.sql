-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- Migration for modern dashboard analytics (MySQL)

-- User Phishing Score table
CREATE TABLE IF NOT EXISTS user_phishing_scores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    score INT DEFAULT 0,
    times_clicked INT DEFAULT 0,
    times_opened INT DEFAULT 0,
    times_reported INT DEFAULT 0,
    last_campaign DATETIME,
    risk_level VARCHAR(20) DEFAULT 'low',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Department Metrics table
CREATE TABLE IF NOT EXISTS department_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    domain VARCHAR(255) NOT NULL,
    department VARCHAR(100) NOT NULL,
    total_users INT DEFAULT 0,
    total_emails INT DEFAULT 0,
    click_rate REAL DEFAULT 0.0,
    open_rate REAL DEFAULT 0.0,
    report_rate REAL DEFAULT 0.0,
    avg_time_to_click REAL DEFAULT 0.0,
    avg_time_to_open REAL DEFAULT 0.0,
    avg_time_to_report REAL DEFAULT 0.0,
    low_risk_users INT DEFAULT 0,
    medium_risk_users INT DEFAULT 0,
    high_risk_users INT DEFAULT 0,
    critical_risk_users INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_domain (domain),
    INDEX idx_department (department)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Time Analytics table
CREATE TABLE IF NOT EXISTS time_analytics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL,
    avg_time_to_click REAL DEFAULT 0.0,
    avg_time_to_open REAL DEFAULT 0.0,
    avg_time_to_report REAL DEFAULT 0.0,
    min_time_to_click REAL DEFAULT 0.0,
    max_time_to_click REAL DEFAULT 0.0,
    clicks_within_1min INT DEFAULT 0,
    clicks_1_to_5min INT DEFAULT 0,
    clicks_5_to_30min INT DEFAULT 0,
    clicks_30_to_60min INT DEFAULT 0,
    clicks_after_1hour INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Campaign Analytics table
CREATE TABLE IF NOT EXISTS campaign_analytics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL UNIQUE,
    total_emails BIGINT DEFAULT 0,
    emails_sent BIGINT DEFAULT 0,
    emails_opened BIGINT DEFAULT 0,
    links_clicked BIGINT DEFAULT 0,
    forms_submitted BIGINT DEFAULT 0,
    emails_reported BIGINT DEFAULT 0,
    errors BIGINT DEFAULT 0,
    click_rate REAL DEFAULT 0.0,
    open_rate REAL DEFAULT 0.0,
    report_rate REAL DEFAULT 0.0,
    conversion_rate REAL DEFAULT 0.0,
    median_time_to_click REAL DEFAULT 0.0,
    median_time_to_open REAL DEFAULT 0.0,
    low_risk_count INT DEFAULT 0,
    medium_risk_count INT DEFAULT 0,
    high_risk_count INT DEFAULT 0,
    critical_risk_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS campaign_analytics;
DROP TABLE IF EXISTS time_analytics;
DROP TABLE IF EXISTS department_metrics;
DROP TABLE IF EXISTS user_phishing_scores;
