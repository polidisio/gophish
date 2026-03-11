-- +goose Up
-- Analytics tables for modern dashboard (PostgreSQL version)

CREATE TABLE IF NOT EXISTS user_phishing_scores (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    department VARCHAR(100),
    domain VARCHAR(255),
    score INTEGER DEFAULT 0,
    times_clicked INTEGER DEFAULT 0,
    times_opened INTEGER DEFAULT 0,
    times_reported INTEGER DEFAULT 0,
    last_campaign_id BIGINT,
    last_campaign VARCHAR(255),
    last_activity TIMESTAMP,
    risk_level VARCHAR(20) DEFAULT 'low',
    first_seen TIMESTAMP,
    last_updated TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_scores_email ON user_phishing_scores(email);
CREATE INDEX IF NOT EXISTS idx_user_scores_domain ON user_phishing_scores(domain);
CREATE INDEX IF NOT EXISTS idx_user_scores_risk ON user_phishing_scores(risk_level);

CREATE TABLE IF NOT EXISTS department_metrics (
    id SERIAL PRIMARY KEY,
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
    total_campaigns INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dept_domain ON department_metrics(domain);
CREATE INDEX IF NOT EXISTS idx_dept_name ON department_metrics(department);

CREATE TABLE IF NOT EXISTS campaign_analytics (
    id SERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL UNIQUE,
    campaign_name VARCHAR(255),
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_campaign_analytics_id ON campaign_analytics(campaign_id);

CREATE TABLE IF NOT EXISTS time_analytics (
    id SERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL,
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_time_analytics_campaign ON time_analytics(campaign_id);

-- +goose Down
DROP TABLE IF EXISTS time_analytics;
DROP TABLE IF EXISTS campaign_analytics;
DROP TABLE IF EXISTS department_metrics;
DROP TABLE IF EXISTS user_phishing_scores;
