package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

func (as *Server) registerAnalyticsRoutes(router *mux.Router) {
	router.HandleFunc("/analytics/summary", as.AnalyticsSummary)
	router.HandleFunc("/analytics/user-scores", as.UserPhishingScores)
	router.HandleFunc("/analytics/user-scores/{email}", as.UserPhishingScoreByEmail)
	router.HandleFunc("/analytics/departments", as.DepartmentMetrics)
	router.HandleFunc("/analytics/campaigns", as.CampaignAnalyticsList)
	router.HandleFunc("/analytics/campaigns/{id:[0-9]+}", as.CampaignAnalyticsDetail)
	router.HandleFunc("/analytics/time/{campaign_id:[0-9]+}", as.TimeAnalytics)
}

func (as *Server) AnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	userId := ctx.Get(r, "user_id").(int64)
	summary, err := models.GetDashboardSummary(userId)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, summary, http.StatusOK)
}

func (as *Server) UserPhishingScores(w http.ResponseWriter, r *http.Request) {
	userId := ctx.Get(r, "user_id").(int64)
	scores, err := models.GetUserPhishingScores(userId)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, scores, http.StatusOK)
}

func (as *Server) UserPhishingScoreByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	score, err := models.GetUserPhishingScoreByEmail(email)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusNotFound)
		return
	}
	JSONResponse(w, score, http.StatusOK)
}

func (as *Server) DepartmentMetrics(w http.ResponseWriter, r *http.Request) {
	userId := ctx.Get(r, "user_id").(int64)
	metrics, err := models.GetDepartmentMetrics(userId)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, metrics, http.StatusOK)
}

func (as *Server) CampaignAnalyticsList(w http.ResponseWriter, r *http.Request) {
	analytics, err := models.GetAllCampaignAnalytics()
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, analytics, http.StatusOK)
}

func (as *Server) CampaignAnalyticsDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 64)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	analytics, err := models.GetCampaignAnalytics(id)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusNotFound)
		return
	}
	JSONResponse(w, analytics, http.StatusOK)
}

func (as *Server) TimeAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId, err := strconv.ParseInt(vars["campaign_id"], 0, 64)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	analytics, err := models.GetTimeAnalytics(campaignId)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusNotFound)
		return
	}
	JSONResponse(w, analytics, http.StatusOK)
}

func (as *Server) UpdateUserScore(w http.ResponseWriter, r *http.Request) {
	var score models.UserPhishingScore
	err := json.NewDecoder(r.Body).Decode(&score)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	score.Score = models.CalculatePhishingScore(score.TimesOpened, score.TimesClicked, score.TimesReported)
	score.RiskLevel = models.CalculateRiskLevel(score.Score)
	err = models.UpsertUserPhishingScore(&score)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, models.Response{Success: true, Message: "Score updated"}, http.StatusOK)
}
