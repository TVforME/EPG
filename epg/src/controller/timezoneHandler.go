package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// TimezoneHandler ...
type TimezoneHandler struct {
	db *gorm.DB
}

// NewTimezoneHandler ...
func NewTimezoneHandler(db *gorm.DB) *TimezoneHandler {
	return &TimezoneHandler{db: db}
}

// GetAllTimezones handler function for GET method
func (th *TimezoneHandler) GetAllTimezones(w http.ResponseWriter, r *http.Request) {
	timezones := []model.Timezone{}
	err := th.db.Preload("Country").Find(&timezones).Error
	if err != nil {
		th.handleError(w, err)
		return
	}
	th.encodeJSONResponse(w, timezones)
}

// GetTimezoneById handler function for GET method
func (th *TimezoneHandler) GetTimezoneById(w http.ResponseWriter, r *http.Request) {
	countryCode := mux.Vars(r)["countryCode"]
	timezoneName := mux.Vars(r)["timezoneName"]
	timezone := &model.Timezone{}
	err := th.db.Preload("Country").Where("country_code = ? AND timezone_name = ?", countryCode, timezoneName).First(timezone).Error
	if err != nil {
		th.handleError(w, err)
		return
	}
	th.encodeJSONResponse(w, timezone)
}

// handleError ...
func (th *TimezoneHandler) handleError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{"status": false, "message": err.Error()}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// encodeJSONResponse ...
func (th *TimezoneHandler) encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
