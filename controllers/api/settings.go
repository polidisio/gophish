package api

import (
	"encoding/json"
	"net/http"

	"github.com/gophish/gophish/models"
	log "github.com/gophish/gophish/logger"
	"github.com/gorilla/mux"
)

func (as *Server) registerSettingsRoutes(router *mux.Router) {
	router.HandleFunc("/settings/entra", as.GetEntraIDSettings).Methods("GET")
	router.HandleFunc("/settings/entra", as.SaveEntraIDSettings).Methods("POST")
}

func (as *Server) GetEntraIDSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := models.GetEntraIDSettings()
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	settings.ClientSecret = ""
	JSONResponse(w, settings, http.StatusOK)
}

func (as *Server) SaveEntraIDSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusBadRequest)
		return
	}

	var settings models.EntraIDSettings
	err := json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Error decoding JSON Request"}, http.StatusBadRequest)
		return
	}

	err = models.SaveEntraIDSettings(&settings)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}

	settings.ClientSecret = ""
	JSONResponse(w, settings, http.StatusOK)
}
