// countryHandler.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CountryHandler ...
type CountryHandler struct {
	db *gorm.DB
}

// NewCountryHandler ...
func NewCountryHandler(db *gorm.DB) *CountryHandler {
	return &CountryHandler{db: db}
}

// GetAllCountries handler function for GET method
func (ch *CountryHandler) GetAllCountries(w http.ResponseWriter, r *http.Request) {
	countries := []model.Country{}
	err := ch.db.Preload("Timezones").Find(&countries).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	ch.encodeJSONResponse(w, countries)
}

// GetCountryById handler function for GET method
func (ch *CountryHandler) GetCountryById(w http.ResponseWriter, r *http.Request) {
	countryCode := mux.Vars(r)["countryCode"]
	country := &model.Country{}
	err := ch.db.Preload("Timezones").Where("country_code = ?", countryCode).First(country).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	ch.encodeJSONResponse(w, country)
}

// GetCountryByCode handler function for GET method
func (ch *CountryHandler) GetCountryByCode(w http.ResponseWriter, r *http.Request) {
	countryCode := mux.Vars(r)["countryCode"]
	countries := []model.Country{}
	err := ch.db.Preload("Timezones").Where("country_code = ?", countryCode).Find(&countries).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	ch.encodeJSONResponse(w, countries)
}

// CreateCountry handler function for POST method
func (ch *CountryHandler) CreateCountry(w http.ResponseWriter, r *http.Request) {
	country := &model.Country{}
	err := json.NewDecoder(r.Body).Decode(country)
	if err != nil {
		ch.handleError(w, err)
		return
	}
	err = ch.db.Create(country).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	ch.encodeJSONResponse(w, country)
}

// UpdateCountry handler function for PUT method
func (ch *CountryHandler) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	country := &model.Country{}
	err := json.NewDecoder(r.Body).Decode(country)
	if err != nil {
		ch.handleError(w, err)
		return
	}
	countryCode := mux.Vars(r)["countryCode"]
	country.CountryCode = countryCode
	err = ch.db.Save(country).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	ch.encodeJSONResponse(w, country)
}

// DeleteCountry handler function for DELETE method
func (ch *CountryHandler) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	countryCode := mux.Vars(r)["countryCode"]
	country := &model.Country{}
	err := ch.db.Where("country_code = ?", countryCode).First(country).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	err = ch.db.Delete(country).Error
	if err != nil {
		ch.handleError(w, err)
		return
	}
	ch.encodeJSONResponse(w, map[string]interface{}{"message": "Country deleted successfully"})
}

// handleError ...
func (ch *CountryHandler) handleError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{"status": false, "message": err.Error()}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// encodeJSONResponse ...
func (ch *CountryHandler) encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
