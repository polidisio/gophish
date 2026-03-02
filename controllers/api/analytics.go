package api

import (
	"net/http"
	"strconv"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

// DashboardSummary returns the overall dashboard summary
// GET /api/analytics/dashboard
func (as *Server) DashboardSummary(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get user's campaigns
		campaigns, err := models.GetCampaigns(userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Get all results for these campaigns
		var allResults []models.Result
		for _, c := range campaigns {
			cr, err := models.GetCampaignResults(c.Id, userID)
			if err != nil {
				continue
			}
			allResults = append(allResults, cr.Results...)
		}
		
		// Get all events
		var allEvents []models.Event
		for _, c := range campaigns {
			events, err := models.GetEvents(c.Id)
			if err != nil {
				continue
			}
			allEvents = append(allEvents, events...)
		}
		
		// Calculate dashboard summary
		analyticsService := models.NewAnalyticsService()
		summary := analyticsService.CalculateDashboardSummary(campaigns, allResults, allEvents)
		
		JSONResponse(w, summary, http.StatusOK)
	}
}

// CampaignAnalytics returns analytics for a specific campaign
// GET /api/analytics/campaign/{id}
func (as *Server) CampaignAnalytics(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		vars := mux.Vars(r)
		campaignID, err := strconv.ParseInt(vars["id"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid campaign ID"}, http.StatusBadRequest)
			return
		}
		
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get campaign
		campaign, err := models.GetCampaign(campaignID, userID)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
			return
		}
		
		// Get results and events
		cr, err := models.GetCampaignResults(campaignID, userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		events, err := models.GetEvents(campaignID)
		if err != nil {
			log.Error(err)
			// Don't fail, just continue without events
			events = []models.Event{}
		}
		
		// Calculate analytics
		analyticsService := models.NewAnalyticsService()
		analytics := analyticsService.CalculateCampaignAnalytics(campaign, cr.Results, events)
		
		JSONResponse(w, analytics, http.StatusOK)
	}
}

// UserScores returns phishing susceptibility scores for all users
// GET /api/analytics/users/scores
func (as *Server) UserScores(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get user's campaigns
		campaigns, err := models.GetCampaigns(userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Get all results
		var allResults []models.Result
		for _, c := range campaigns {
			cr, err := models.GetCampaignResults(c.Id, userID)
			if err != nil {
				continue
			}
			allResults = append(allResults, cr.Results...)
		}
		
		// Calculate scores
		analyticsService := models.NewAnalyticsService()
		scores := analyticsService.CalculateUserScores(allResults)
		
		JSONResponse(w, scores, http.StatusOK)
	}
}

// UserScore returns phishing score for a specific user
// GET /api/analytics/users/{email}/score
func (as *Server) UserScore(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		vars := mux.Vars(r)
		email := vars["email"]
		
		if email == "" {
			JSONResponse(w, models.Response{Success: false, Message: "Email required"}, http.StatusBadRequest)
			return
		}
		
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get user's campaigns
		campaigns, err := models.GetCampaigns(userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Get results for this email
		var userResults []models.Result
		for _, c := range campaigns {
			cr, err := models.GetCampaignResults(c.Id, userID)
			if err != nil {
				continue
			}
			for _, r := range cr.Results {
				if r.Email == email {
					userResults = append(userResults, r)
				}
			}
		}
		
		// Calculate single user score
		score := models.UserPhishingScore{
			Email: email,
		}
		for _, r := range userResults {
			switch r.Status {
			case models.EventClicked:
				score.TimesClicked++
			case models.EventOpened:
				score.TimesOpened++
			}
			if r.Reported {
				score.TimesReported++
			}
		}
		score.CalculateScore()
		
		JSONResponse(w, score, http.StatusOK)
	}
}

// DepartmentMetrics returns metrics grouped by department
// GET /api/analytics/departments
func (as *Server) DepartmentMetrics(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get user's campaigns
		campaigns, err := models.GetCampaigns(userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Get all results
		var allResults []models.Result
		for _, c := range campaigns {
			cr, err := models.GetCampaignResults(c.Id, userID)
			if err != nil {
				continue
			}
			allResults = append(allResults, cr.Results...)
		}
		
		// Calculate department metrics
		analyticsService := models.NewAnalyticsService()
		metrics := analyticsService.CalculateDepartmentMetrics(allResults)
		
		JSONResponse(w, metrics, http.StatusOK)
	}
}

// TimeAnalytics returns time-based analytics for a campaign
// GET /api/analytics/campaign/{id}/time
func (as *Server) TimeAnalytics(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		vars := mux.Vars(r)
		campaignID, err := strconv.ParseInt(vars["id"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid campaign ID"}, http.StatusBadRequest)
			return
		}
		
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get results and events
		cr, err := models.GetCampaignResults(campaignID, userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		events, err := models.GetEvents(campaignID)
		if err != nil {
			events = []models.Event{}
		}
		
		// Calculate time analytics
		analyticsService := models.NewAnalyticsService()
		timeAnalytics := analyticsService.CalculateTimeAnalytics(events, cr.Results)
		
		JSONResponse(w, timeAnalytics, http.StatusOK)
	}
}

// RiskDistribution returns the risk distribution across all users
// GET /api/analytics/risk/distribution
func (as *Server) RiskDistribution(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get user's campaigns
		campaigns, err := models.GetCampaigns(userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Get all results
		var allResults []models.Result
		for _, c := range campaigns {
			cr, err := models.GetCampaignResults(c.Id, userID)
			if err != nil {
				continue
			}
			allResults = append(allResults, cr.Results...)
		}
		
		// Calculate scores
		analyticsService := models.NewAnalyticsService()
		scores := analyticsService.CalculateUserScores(allResults)
		
		// Count by risk level
		distribution := map[string]int{
			"low":      0,
			"medium":   0,
			"high":     0,
			"critical": 0,
		}
		
		for _, s := range scores {
			distribution[s.RiskLevel]++
		}
		
		JSONResponse(w, distribution, http.StatusOK)
	}
}

// TopAtRiskUsers returns the users with highest phishing susceptibility
// GET /api/analytics/users/at-risk
func (as *Server) TopAtRiskUsers(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		userID := ctx.Get(r, "user_id").(int64)
		
		// Get limit from query param (default 10)
		limit := 10
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}
		
		// Get user's campaigns
		campaigns, err := models.GetCampaigns(userID)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Get all results
		var allResults []models.Result
		for _, c := range campaigns {
			cr, err := models.GetCampaignResults(c.Id, userID)
			if err != nil {
				continue
			}
			allResults = append(allResults, cr.Results...)
		}
		
		// Calculate scores
		analyticsService := models.NewAnalyticsService()
		scores := analyticsService.CalculateUserScores(allResults)
		
		// Return top at-risk (highest scores)
		if len(scores) > limit {
			scores = scores[:limit]
		}
		
		JSONResponse(w, scores, http.StatusOK)
	}
}
