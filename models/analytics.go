package models

import (
	"time"
)

// UserPhishingScore represents the phishing susceptibility score for a user
type UserPhishingScore struct {
	ID          int64   `json:"id" gorm:"primary_key"`
	UserID      int64   `json:"user_id"`
	Email       string  `json:"email"`
	Score       int     `json:"score"` // 0-100
	TimesClicked int    `json:"times_clicked"`
	TimesOpened  int    `json:"times_opened"`
	TimesReported int   `json:"times_reported"`
	LastCampaign  *time.Time `json:"last_campaign"`
	RiskLevel    string    `json:"risk_level"` // low, medium, high, critical
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CalculateRiskLevel determines the risk level based on score
func (u *UserPhishingScore) CalculateRiskLevel() string {
	switch {
	case u.Score >= 80:
		return "critical"
	case u.Score >= 60:
		return "high"
	case u.Score >= 40:
		return "medium"
	default:
		return "low"
	}
}

// CalculateScore computes the susceptibility score based on behavior
func (u *UserPhishingScore) CalculateScore() int {
	// Base score starts at 0 (safe)
	score := 0

	// Clicking adds 25 points per incident
	score += u.TimesClicked * 25

	// Opening without clicking adds 10 points per incident
	score += u.TimesOpened * 10

	// Reporting is positive - reduces score by 15 points each time
	score -= u.TimesReported * 15

	// Clamp score between 0 and 100
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	u.Score = score
	u.RiskLevel = u.CalculateRiskLevel()

	return score
}

// DepartmentMetrics represents aggregated metrics for a department/domain
type DepartmentMetrics struct {
	ID          int64   `json:"id" gorm:"primary_key"`
	Domain      string  `json:"domain"`
	Department  string  `json:"department"`
	TotalUsers  int     `json:"total_users"`
	TotalEmails int     `json:"total_emails"`
	
	// Rates (percentages)
	ClickRate   float64 `json:"click_rate"`
	OpenRate    float64 `json:"open_rate"`
	ReportRate  float64 `json:"report_rate"`
	
	// Average times (in minutes)
	AvgTimeToClick  float64 `json:"avg_time_to_click"`
	AvgTimeToOpen  float64 `json:"avg_time_to_open"`
	AvgTimeToReport float64 `json:"avg_time_to_report"`
	
	// Risk distribution
	LowRiskUsers    int `json:"low_risk_users"`
	MediumRiskUsers int `json:"medium_risk_users"`
	HighRiskUsers   int `json:"high_risk_users"`
	CriticalRiskUsers int `json:"critical_risk_users"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TimeAnalytics contains time-based metrics
type TimeAnalytics struct {
	ID                int64     `json:"id" gorm:"primary_key"`
	CampaignID        int64     `json:"campaign_id"`
	
	// Average time metrics (in minutes)
	AvgTimeToClick   float64   `json:"avg_time_to_click"`
	AvgTimeToOpen    float64   `json:"avg_time_to_open"`
	AvgTimeToReport  float64   `json:"avg_time_to_report"`
	MinTimeToClick   float64   `json:"min_time_to_click"`
	MaxTimeToClick   float64   `json:"max_time_to_click"`
	
	// Distribution buckets (number of clicks in each time range)
	ClicksWithin1Min   int `json:"clicks_within_1min"`
	Clicks1To5Min      int `json:"clicks_1_to_5min"`
	Clicks5To30Min     int `json:"clicks_5_to_30min"`
	Clicks30To60Min    int `json:"clicks_30_to_60min"`
	ClicksAfter1Hour   int `json:"clicks_after_1hour"`
	
	CreatedAt time.Time `json:"created_at"`
}

// CampaignAnalytics contains comprehensive analytics for a campaign
type CampaignAnalytics struct {
	ID                    int64   `json:"id" gorm:"primary_key"`
	CampaignID            int64   `json:"campaign_id"`
	
	// Basic stats
	TotalEmails          int64   `json:"total_emails"`
	EmailsSent           int64   `json:"emails_sent"`
	EmailsOpened         int64   `json:"emails_opened"`
	LinksClicked         int64   `json:"links_clicked"`
	FormsSubmitted       int64   `json:"forms_submitted"`
	EmailsReported       int64   `json:"emails_reported"`
	Errors               int64   `json:"errors"`
	
	// Advanced metrics
	ClickRate            float64 `json:"click_rate"`      // clicks / sent
	OpenRate             float64 `json:"open_rate"`       // opened / sent
	ReportRate           float64 `json:"report_rate"`     // reported / sent
	ConversionRate       float64 `json:"conversion_rate"`  // submitted / clicked
	
	// Time metrics
	MedianTimeToClick    float64 `json:"median_time_to_click"`   // minutes
	MedianTimeToOpen     float64 `json:"median_time_to_open"`    // minutes
	
	// User risk distribution at campaign end
	LowRiskCount         int     `json:"low_risk_count"`
	MediumRiskCount      int     `json:"medium_risk_count"`
	HighRiskCount        int     `json:"high_risk_count"`
	CriticalRiskCount    int     `json:"critical_risk_count"`
	
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CalculateRates computes the percentage rates for the campaign
func (c *CampaignAnalytics) CalculateRates() {
	if c.EmailsSent > 0 {
		c.ClickRate = float64(c.LinksClicked) / float64(c.EmailsSent) * 100
		c.OpenRate = float64(c.EmailsOpened) / float64(c.EmailsSent) * 100
		c.ReportRate = float64(c.EmailsReported) / float64(c.EmailsSent) * 100
	}
	
	if c.LinksClicked > 0 {
		c.ConversionRate = float64(c.FormsSubmitted) / float64(c.LinksClicked) * 100
	}
}

// TrendData represents a single point in a time series
type TrendData struct {
	Date          time.Time `json:"date"`
	Value         float64   `json:"value"`
	Label         string    `json:"label"`
}

// DashboardSummary is the main dashboard data structure
type DashboardSummary struct {
	// Overview stats
	TotalCampaigns     int     `json:"total_campaigns"`
	ActiveCampaigns    int     `json:"active_campaigns"`
	CompletedCampaigns int     `json:"completed_campaigns"`
	
	// Aggregate metrics
	TotalEmailsSent    int64   `json:"total_emails_sent"`
	TotalClicks        int64   `json:"total_clicks"`
	TotalOpens         int64   `json:"total_opens"`
	TotalReports       int64   `json:"total_reports"`
	
	// Overall rates
	GlobalClickRate    float64 `json:"global_click_rate"`
	GlobalOpenRate     float64 `json:"global_open_rate"`
	GlobalReportRate   float64 `json:"global_report_rate"`
	
	// Risk distribution
	LowRiskUsers       int     `json:"low_risk_users"`
	MediumRiskUsers    int     `json:"medium_risk_users"`
	HighRiskUsers      int     `json:"high_risk_users"`
	CriticalRiskUsers  int     `json:"critical_risk_users"`
	
	// Top at-risk users
	TopAtRiskUsers     []UserPhishingScore `json:"top_at_risk_users"`
	
	// Department summary
	DepartmentStats     []DepartmentMetrics `json:"department_stats"`
	
	// Trend data
	ClickTrend         []TrendData `json:"click_trend"`
	OpenTrend          []TrendData `json:"open_trend"`
	ReportTrend        []TrendData `json:"report_trend"`
	
	// Time period
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}

// CalculateDashboard computes all dashboard metrics
func (d *DashboardSummary) CalculateDashboard() {
	totalSent := d.TotalEmailsSent
	if totalSent > 0 {
		d.GlobalClickRate = float64(d.TotalClicks) / float64(totalSent) * 100
		d.GlobalOpenRate = float64(d.TotalOpens) / float64(totalSent) * 100
		d.GlobalReportRate = float64(d.TotalReports) / float64(totalSent) * 100
	}
}

// GetDepartmentFromEmail extracts department from email domain
func GetDepartmentFromEmail(email string) string {
	// Simple logic - could be extended to use company directory
	domain := ""
	if i := find(email, "@"); i != -1 {
		domain = email[i+1:]
	}
	
	// Common department patterns
	switch {
	case contains(domain, "engineering") || contains(domain, "dev") || contains(domain, "tech"):
		return "Engineering"
	case contains(domain, "sales"):
		return "Sales"
	case contains(domain, "marketing"):
		return "Marketing"
	case contains(domain, "hr") || contains(domain, "human resources"):
		return "HR"
	case contains(domain, "finance") || contains(domain, "accounting"):
		return "Finance"
	case contains(domain, "it") || contains(domain, "support"):
		return "IT"
	case contains(domain, "management") || contains(domain, "exec") || contains(domain, "ceo"):
		return "Executive"
	default:
		return "General"
	}
}

func find(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func contains(s, substr string) bool {
	return find(s, substr) != -1
}
