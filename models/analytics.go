package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type UserPhishingScore struct {
	Id              int64      `json:"id" gorm:"primary_key"`
	UserId          int64      `json:"user_id"`
	Email           string     `json:"email" gorm:"index"`
	Department      string     `json:"department"`
	Domain          string     `json:"domain" gorm:"index"`
	Score           int        `json:"score" gorm:"default:0"`
	TimesClicked    int        `json:"times_clicked" gorm:"default:0"`
	TimesOpened     int        `json:"times_opened" gorm:"default:0"`
	TimesReported   int        `json:"times_reported" gorm:"default:0"`
	LastCampaignId  int64      `json:"last_campaign_id"`
	LastCampaign    string     `json:"last_campaign"`
	LastActivity    time.Time  `json:"last_activity"`
	RiskLevel       RiskLevel  `json:"risk_level" gorm:"default:'low'"`
	FirstSeen       time.Time  `json:"first_seen"`
	LastUpdated     time.Time  `json:"last_updated"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (UserPhishingScore) TableName() string {
	return "user_phishing_scores"
}

type DepartmentMetrics struct {
	Id                  int64     `json:"id" gorm:"primary_key"`
	Domain              string    `json:"domain" gorm:"index"`
	Department          string    `json:"department"`
	TotalUsers          int       `json:"total_users" gorm:"default:0"`
	TotalEmails         int       `json:"total_emails" gorm:"default:0"`
	ClickRate           float64   `json:"click_rate" gorm:"default:0.0"`
	OpenRate            float64   `json:"open_rate" gorm:"default:0.0"`
	ReportRate          float64   `json:"report_rate" gorm:"default:0.0"`
	AvgTimeToClick      float64   `json:"avg_time_to_click" gorm:"default:0.0"`
	AvgTimeToOpen       float64   `json:"avg_time_to_open" gorm:"default:0.0"`
	AvgTimeToReport     float64   `json:"avg_time_to_report" gorm:"default:0.0"`
	LowRiskUsers        int       `json:"low_risk_users" gorm:"default:0"`
	MediumRiskUsers     int       `json:"medium_risk_users" gorm:"default:0"`
	HighRiskUsers       int       `json:"high_risk_users" gorm:"default:0"`
	CriticalRiskUsers   int       `json:"critical_risk_users" gorm:"default:0"`
	TotalCampaigns      int       `json:"total_campaigns" gorm:"default:0"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (DepartmentMetrics) TableName() string {
	return "department_metrics"
}

type CampaignAnalytics struct {
	Id                  int64     `json:"id" gorm:"primary_key"`
	CampaignId          int64     `json:"campaign_id" gorm:"unique_index"`
	CampaignName        string    `json:"campaign_name"`
	TotalEmails         int       `json:"total_emails" gorm:"default:0"`
	EmailsSent          int       `json:"emails_sent" gorm:"default:0"`
	EmailsOpened        int       `json:"emails_opened" gorm:"default:0"`
	LinksClicked        int       `json:"links_clicked" gorm:"default:0"`
	FormsSubmitted      int       `json:"forms_submitted" gorm:"default:0"`
	EmailsReported      int       `json:"emails_reported" gorm:"default:0"`
	Errors              int       `json:"errors" gorm:"default:0"`
	ClickRate           float64   `json:"click_rate" gorm:"default:0.0"`
	OpenRate            float64   `json:"open_rate" gorm:"default:0.0"`
	ReportRate          float64   `json:"report_rate" gorm:"default:0.0"`
	ConversionRate      float64   `json:"conversion_rate" gorm:"default:0.0"`
	MedianTimeToClick   float64   `json:"median_time_to_click" gorm:"default:0.0"`
	MedianTimeToOpen    float64   `json:"median_time_to_open" gorm:"default:0.0"`
	LowRiskCount        int       `json:"low_risk_count" gorm:"default:0"`
	MediumRiskCount     int       `json:"medium_risk_count" gorm:"default:0"`
	HighRiskCount       int       `json:"high_risk_count" gorm:"default:0"`
	CriticalRiskCount   int       `json:"critical_risk_count" gorm:"default:0"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (CampaignAnalytics) TableName() string {
	return "campaign_analytics"
}

type TimeAnalytics struct {
	Id                  int64     `json:"id" gorm:"primary_key"`
	CampaignId          int64     `json:"campaign_id" gorm:"index"`
	AvgTimeToClick      float64   `json:"avg_time_to_click" gorm:"default:0.0"`
	AvgTimeToOpen       float64   `json:"avg_time_to_open" gorm:"default:0.0"`
	AvgTimeToReport     float64   `json:"avg_time_to_report" gorm:"default:0.0"`
	MinTimeToClick      float64   `json:"min_time_to_click" gorm:"default:0.0"`
	MaxTimeToClick      float64   `json:"max_time_to_click" gorm:"default:0.0"`
	ClicksWithin1Min    int       `json:"clicks_within_1min" gorm:"default:0"`
	Clicks1To5Min       int       `json:"clicks_1_to_5min" gorm:"default:0"`
	Clicks5To30Min      int       `json:"clicks_5_to_30min" gorm:"default:0"`
	Clicks30To60Min     int       `json:"clicks_30_to_60min" gorm:"default:0"`
	ClicksAfter1Hour    int       `json:"clicks_after_1hour" gorm:"default:0"`
	CreatedAt           time.Time `json:"created_at"`
}

func (TimeAnalytics) TableName() string {
	return "time_analytics"
}

type DashboardSummary struct {
	TotalCampaigns      int64              `json:"total_campaigns"`
	TotalUsers          int64              `json:"total_users"`
	TotalEmailsSent     int64              `json:"total_emails_sent"`
	OverallClickRate    float64            `json:"overall_click_rate"`
	OverallOpenRate     float64            `json:"overall_open_rate"`
	OverallReportRate   float64            `json:"overall_report_rate"`
	RiskDistribution    RiskDistribution   `json:"risk_distribution"`
	TopDepartments      []DepartmentMetrics `json:"top_departments"`
	RecentCampaigns    []CampaignSummary   `json:"recent_campaigns"`
	TimeSeriesData     []TimeSeriesPoint  `json:"time_series_data"`
}

type RiskDistribution struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type TimeSeriesPoint struct {
	Date        string  `json:"date"`
	Clicked     int     `json:"clicked"`
	Opened      int     `json:"opened"`
	Reported    int     `json:"reported"`
}

func GetUserPhishingScores(userId int64) ([]UserPhishingScore, error) {
	var scores []UserPhishingScore
	err := db.Where("user_id = ?", userId).Order("score DESC").Find(&scores).Error
	return scores, err
}

func GetUserPhishingScoreByEmail(email string) (*UserPhishingScore, error) {
	var score UserPhishingScore
	err := db.Where("email = ?", email).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func UpsertUserPhishingScore(score *UserPhishingScore) error {
	var existing UserPhishingScore
	err := db.Where("email = ?", score.Email).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		score.FirstSeen = time.Now().UTC()
		score.LastUpdated = time.Now().UTC()
		return db.Create(score).Error
	}
	existing.Score = score.Score
	existing.TimesClicked = score.TimesClicked
	existing.TimesOpened = score.TimesOpened
	existing.TimesReported = score.TimesReported
	existing.RiskLevel = score.RiskLevel
	existing.LastUpdated = time.Now().UTC()
	return db.Save(&existing).Error
}

func GetDepartmentMetrics(userId int64) ([]DepartmentMetrics, error) {
	var metrics []DepartmentMetrics
	err := db.Order("total_users DESC").Find(&metrics).Error
	return metrics, err
}

func GetCampaignAnalytics(campaignId int64) (*CampaignAnalytics, error) {
	var analytics CampaignAnalytics
	err := db.Where("campaign_id = ?", campaignId).First(&analytics).Error
	if err != nil {
		return nil, err
	}
	return &analytics, nil
}

func GetAllCampaignAnalytics() ([]CampaignAnalytics, error) {
	var analytics []CampaignAnalytics
	err := db.Order("campaign_id DESC").Find(&analytics).Error
	return analytics, err
}

func UpsertCampaignAnalytics(a *CampaignAnalytics) error {
	var existing CampaignAnalytics
	err := db.Where("campaign_id = ?", a.CampaignId).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		a.CreatedAt = time.Now().UTC()
		a.UpdatedAt = time.Now().UTC()
		return db.Create(a).Error
	}
	existing.EmailsSent = a.EmailsSent
	existing.EmailsOpened = a.EmailsOpened
	existing.LinksClicked = a.LinksClicked
	existing.EmailsReported = a.EmailsReported
	existing.ClickRate = a.ClickRate
	existing.OpenRate = a.OpenRate
	existing.ReportRate = a.ReportRate
	existing.UpdatedAt = time.Now().UTC()
	return db.Save(&existing).Error
}

func GetTimeAnalytics(campaignId int64) (*TimeAnalytics, error) {
	var analytics TimeAnalytics
	err := db.Where("campaign_id = ?", campaignId).First(&analytics).Error
	if err != nil {
		return nil, err
	}
	return &analytics, nil
}

func GetDashboardSummary(userId int64) (*DashboardSummary, error) {
	var analytics []CampaignAnalytics
	err := db.Order("campaign_id DESC").Find(&analytics).Error
	if err != nil {
		return nil, err
	}

	var totalSent, totalOpened, totalClicked, totalReported int
	for _, a := range analytics {
		totalSent += a.EmailsSent
		totalOpened += a.EmailsOpened
		totalClicked += a.LinksClicked
		totalReported += a.EmailsReported
	}

	var totalUsers int64
	db.Model(&UserPhishingScore{}).Count(&totalUsers)

	summary := &DashboardSummary{
		TotalCampaigns:     int64(len(analytics)),
		TotalUsers:         totalUsers,
		TotalEmailsSent:    int64(totalSent),
		OverallClickRate:   0,
		OverallOpenRate:    0,
		OverallReportRate:  0,
		RiskDistribution:   RiskDistribution{},
		TopDepartments:     []DepartmentMetrics{},
		RecentCampaigns:    []CampaignSummary{},
		TimeSeriesData:     []TimeSeriesPoint{},
	}

	if totalSent > 0 {
		summary.OverallClickRate = float64(totalClicked) / float64(totalSent) * 100
		summary.OverallOpenRate = float64(totalOpened) / float64(totalSent) * 100
		summary.OverallReportRate = float64(totalReported) / float64(totalSent) * 100
	}

	var deptMetrics []DepartmentMetrics
	db.Order("total_users DESC").Limit(5).Find(&deptMetrics)
	summary.TopDepartments = deptMetrics

	var low, medium, high, critical int64
	db.Model(&UserPhishingScore{}).Where("risk_level = ?", "low").Count(&low)
	db.Model(&UserPhishingScore{}).Where("risk_level = ?", "medium").Count(&medium)
	db.Model(&UserPhishingScore{}).Where("risk_level = ?", "high").Count(&high)
	db.Model(&UserPhishingScore{}).Where("risk_level = ?", "critical").Count(&critical)
	summary.RiskDistribution = RiskDistribution{
		Low:      int(low),
		Medium:   int(medium),
		High:     int(high),
		Critical: int(critical),
	}

	return summary, nil
}

func CalculateCampaignStats(c Campaign) CampaignStats {
	var stats CampaignStats
	stats.Total = int64(len(c.Results))
	for _, r := range c.Results {
		switch r.Status {
		case "Email Sent":
			stats.EmailsSent++
		case "Email Opened":
			stats.OpenedEmail++
		case "Clicked Link":
			stats.ClickedLink++
		case "Submitted Data":
			stats.SubmittedData++
		case "Email Reported":
			stats.EmailReported++
		}
		if r.Status == "Error Sending Email" || r.Status == "Error" {
			stats.Error++
		}
	}
	return stats
}

func CalculateRiskLevel(score int) RiskLevel {
	if score >= 75 {
		return RiskCritical
	} else if score >= 50 {
		return RiskHigh
	} else if score >= 25 {
		return RiskMedium
	}
	return RiskLow
}

func CalculatePhishingScore(opens, clicks, reports int) int {
	score := 0
	score += opens * 10
	score += clicks * 25
	score -= reports * 20
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}
