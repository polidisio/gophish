-- +migrate Up
-- Migration for modern dashboard analytics

-- User Phishing Score table
CREATE TABLE IF NOT EXISTS user_phishing_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,
    score INTEGER DEFAULT 0,
    times_clicked INTEGER DEFAULT 0,
    times_opened INTEGER DEFAULT 0,
    times_reported INTEGER DEFAULT 0,
    last_campaign DATETIME,
    risk_level VARCHAR(20) DEFAULT 'low',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Department Metrics table
CREATE TABLE IF NOT EXISTS department_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain VARCHAR(255) NOT NULL,
    department VARCHAR(100) NOT NULL,
    total_users INTEGER DEFAULT 0,
    total_emails INTEGER DEFAULT 0,
    click_rate REAL DEFAULT 0.0,
    open_rate REAL DEFAULT 0.0,
    report_rate REAL DEFAULT 0.0,
    avg_time_to_click REAL DEFAULT 0.0,
    avg_time_to_open REAL DEFAULT 0.0,
    avg_time_to_report REAL DEFAULT 0.0,
    low_risk_users INTEGER DEFAULT 0,
    medium_risk_users INTEGER DEFAULT 0,
    high_risk_users INTEGER DEFAULT 0,
    critical_risk_users INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Time Analytics table
CREATE TABLE IF NOT EXISTS time_analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    campaign_id INTEGER NOT NULL,
    avg_time_to_click REAL DEFAULT 0.0,
    avg_time_to_open REAL DEFAULT 0.0,
    avg_time_to_report REAL DEFAULT 0.0,
    min_time_to_click REAL DEFAULT 0.0,
    max_time_to_click REAL DEFAULT 0.0,
    clicks_within_1min INTEGER DEFAULT 0,
    clicks_1_to_5min INTEGER DEFAULT 0,
    clicks_5_to_30min INTEGER DEFAULT 0,
    clicks_30_to_60min INTEGER DEFAULT 0,
    clicks_after_1hour INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Campaign Analytics table
CREATE TABLE IF NOT EXISTS campaign_analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    campaign_id INTEGER NOT NULL UNIQUE,
    total_emails INTEGER DEFAULT 0,
    emails_sent INTEGER DEFAULT 0,
    emails_opened INTEGER DEFAULT 0,
    links_clicked INTEGER DEFAULT 0,
    forms_submitted INTEGER DEFAULT 0,
    emails_reported INTEGER DEFAULT 0,
    errors INTEGER DEFAULT 0,
    click_rate REAL DEFAULT 0.0,
    open_rate REAL DEFAULT 0.0,
    report_rate REAL DEFAULT 0.0,
    conversion_rate REAL DEFAULT 0.0,
    median_time_to_click REAL DEFAULT 0.0,
    median_time_to_open REAL DEFAULT 0.0,
    low_risk_count INTEGER DEFAULT 0,
    medium_risk_count INTEGER DEFAULT 0,
    high_risk_count INTEGER DEFAULT 0,
    critical_risk_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_phishing_scores;
DROP TABLE IF EXISTS department_metrics;
DROP TABLE IF EXISTS time_analytics;
DROP TABLE IF EXISTS campaign_analytics;
