package models

import (
	"math"
	"sort"
	"time"
)

// AnalyticsService provides methods to calculate analytics
type AnalyticsService struct{}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// CalculateUserScores calculates phishing scores for all users
func (a *AnalyticsService) CalculateUserScores(results []Result) []UserPhishingScore {
	// Group results by email
	emailStats := make(map[string]*UserPhishingScore)
	
	for _, r := range results {
		if _, exists := emailStats[r.Email]; !exists {
			emailStats[r.Email] = &UserPhishingScore{
				Email:         r.Email,
				TimesClicked:  0,
				TimesOpened:   0,
				TimesReported: 0,
			}
		}
		
		switch r.Status {
		case EventClicked:
			emailStats[r.Email].TimesClicked++
		case EventOpened:
			emailStats[r.Email].TimesOpened++
		}
		
		if r.Reported {
			emailStats[r.Email].TimesReported++
		}
	}
	
	// Calculate scores
	var scores []UserPhishingScore
	for _, s := range emailStats {
		s.CalculateScore()
		scores = append(scores, *s)
	}
	
	// Sort by score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	
	return scores
}

// CalculateTimeAnalytics calculates time-based analytics from events
func (a *AnalyticsService) CalculateTimeAnalytics(events []Event, results []Result) TimeAnalytics {
	ta := TimeAnalytics{}
	
	// Create lookup for send dates
	sendDates := make(map[string]time.Time)
	for _, r := range results {
		sendDates[r.Email] = r.SendDate
	}
	
	var clickTimes, openTimes []float64
	
	for _, e := range events {
		sendTime, ok := sendDates[e.Email]
		if !ok {
			continue
		}
		
		minutesToAction := e.Time.Sub(sendTime).Minutes()
		
		switch e.Message {
		case EventClicked:
			clickTimes = append(clickTimes, minutesToAction)
			// Bucket the click time
			switch {
			case minutesToAction <= 1:
				ta.ClicksWithin1Min++
			case minutesToAction <= 5:
				ta.Clicks1To5Min++
			case minutesToAction <= 30:
				ta.Clicks5To30Min++
			case minutesToAction <= 60:
				ta.Clicks30To60Min++
			default:
				ta.ClicksAfter1Hour++
			}
		case EventOpened:
			openTimes = append(openTimes, minutesToAction)
		}
	}
	
	// Calculate averages
	if len(clickTimes) > 0 {
		ta.AvgTimeToClick = calculateMean(clickTimes)
		ta.MinTimeToClick = calculateMin(clickTimes)
		ta.MaxTimeToClick = calculateMax(clickTimes)
	}
	
	if len(openTimes) > 0 {
		ta.AvgTimeToOpen = calculateMean(openTimes)
	}
	
	return ta
}

// CalculateDepartmentMetrics calculates metrics by department
func (a *AnalyticsService) CalculateDepartmentMetrics(results []Result) []DepartmentMetrics {
	deptStats := make(map[string]*DepartmentMetrics)
	
	for _, r := range results {
		department := GetDepartmentFromEmail(r.Email)
		domain := extractDomain(r.Email)
		
		key := department
		if _, exists := deptStats[key]; !exists {
			deptStats[key] = &DepartmentMetrics{
				Domain:     domain,
				Department: department,
			}
		}
		
		dm := deptStats[key]
		dm.TotalUsers++
		dm.TotalEmails++
		
		switch r.Status {
		case EventClicked:
			dm.ClickRate++
		case EventOpened:
			dm.OpenRate++
		}
		
		if r.Reported {
			dm.ReportRate++
		}
	}
	
	// Calculate percentages
	var metrics []DepartmentMetrics
	for _, dm := range deptStats {
		if dm.TotalEmails > 0 {
			dm.ClickRate = (dm.ClickRate / float64(dm.TotalEmails)) * 100
			dm.OpenRate = (dm.OpenRate / float64(dm.TotalEmails)) * 100
			dm.ReportRate = (dm.ReportRate / float64(dm.TotalEmails)) * 100
		}
		metrics = append(metrics, *dm)
	}
	
	// Sort by click rate descending
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ClickRate > metrics[j].ClickRate
	})
	
	return metrics
}

// CalculateCampaignAnalytics calculates comprehensive analytics for a campaign
func (a *AnalyticsService) CalculateCampaignAnalytics(campaign Campaign, results []Result, events []Event) CampaignAnalytics {
	ca := CampaignAnalytics{
		CampaignID:   campaign.Id,
		TotalEmails:  int64(len(results)),
		EmailsSent:   0,
		EmailsOpened: 0,
		LinksClicked: 0,
	}
	
	sendDates := make(map[string]time.Time)
	var clickTimes, openTimes []float64
	
	for _, r := range results {
		// Skip if email wasn't sent
		if r.Status == "" || r.Status == Error {
			continue
		}
		
		ca.EmailsSent++
		sendDates[r.Email] = r.SendDate
		
		switch r.Status {
		case EventOpened:
			ca.EmailsOpened++
		case EventClicked:
			ca.LinksClicked++
		case EventDataSubmit:
			ca.LinksClicked++
			ca.FormsSubmitted++
		case EventReported:
			ca.EmailsReported++
		default:
			if r.Status == Error || r.Status == StatusRetry {
				ca.Errors++
			}
		}
	}
	
	// Calculate time metrics
	for _, e := range events {
		if e.CampaignId != campaign.Id {
			continue
		}
		
		sendTime, ok := sendDates[e.Email]
		if !ok {
			continue
		}
		
		minutesToAction := e.Time.Sub(sendTime).Minutes()
		
		switch e.Message {
		case EventClicked:
			clickTimes = append(clickTimes, minutesToAction)
		case EventOpened:
			openTimes = append(openTimes, minutesToAction)
		}
	}
	
	if len(clickTimes) > 0 {
		ca.AvgTimeToClick = calculateMean(clickTimes)
		ca.MedianTimeToClick = calculateMedian(clickTimes)
	}
	
	if len(openTimes) > 0 {
		ca.AvgTimeToOpen = calculateMean(openTimes)
		ca.MedianTimeToOpen = calculateMedian(openTimes)
	}
	
	// Calculate rates
	ca.CalculateRates()
	
	return ca
}

// CalculateDashboardSummary calculates the overall dashboard summary
func (a *AnalyticsService) CalculateDashboardSummary(campaigns []Campaign, results []Result, events []Event) DashboardSummary {
	ds := DashboardSummary{
		TotalCampaigns:     len(campaigns),
		ActiveCampaigns:   0,
		CompletedCampaigns: 0,
	}
	
	for _, c := range campaigns {
		switch c.Status {
		case "In Progress", "Queued":
			ds.ActiveCampaigns++
		case "Completed":
			ds.CompletedCampaigns++
		}
	}
	
	// Aggregate all results
	for _, r := range results {
		ds.TotalEmailsSent++
		
		switch r.Status {
		case EventClicked:
			ds.TotalClicks++
		case EventOpened:
			ds.TotalOpens++
		}
		
		if r.Reported {
			ds.TotalReports++
		}
	}
	
	// Calculate user risk distribution
	userScores := a.CalculateUserScores(results)
	for _, s := range userScores {
		switch s.RiskLevel {
		case "low":
			ds.LowRiskUsers++
		case "medium":
			ds.MediumRiskUsers++
		case "high":
			ds.HighRiskUsers++
		case "critical":
			ds.CriticalRiskUsers++
		}
	}
	
	// Get top 10 at-risk users
	if len(userScores) > 10 {
		ds.TopAtRiskUsers = userScores[:10]
	} else {
		ds.TopAtRiskUsers = userScores
	}
	
	// Calculate department stats
	ds.DepartmentStats = a.CalculateDepartmentMetrics(results)
	
	// Calculate dashboard rates
	ds.CalculateDashboard()
	
	return ds
}

// Helper functions
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func calculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func extractDomain(email string) string {
	atIndex := find(email, "@")
	if atIndex == -1 {
		return ""
	}
	return email[atIndex+1:]
}

// Standard deviation calculation
func calculateStdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	mean := calculateMean(values)
	variance := 0.0
	
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	
	return math.Sqrt(variance / float64(len(values)))
}
