// ratingHandler.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RatingSystemHandler ...
type RatingHandler struct {
	db *gorm.DB
}

// NewRatingSystemHandler ...
func NewRatingHandler(db *gorm.DB) *RatingHandler {
	return &RatingHandler{db: db}
}

// GetAllRatingsHTML handler function for GET method
func (rsh *RatingHandler) GetAllRatingsHTML(w http.ResponseWriter, r *http.Request) {
	ratingSystems := []model.RatingSystem{}
	err := rsh.db.Preload("Country").Preload("RatingValues").Find(&ratingSystems).Error
	if err != nil {
		rsh.handleError(w, err)
		return
	}

	var responses []model.RatingSystem
	for _, rs := range ratingSystems {
		var ratingSystemResponse model.RatingSystem
		ratingSystemResponse.RatingSystemID = rs.RatingSystemID
		ratingSystemResponse.CountryID = rs.CountryID
		ratingSystemResponse.Description = rs.Description
		ratingSystemResponse.Country.CountryCode = rs.Country.CountryCode
		ratingSystemResponse.Country.CountryName = rs.Country.CountryName

		var ratingValues []model.RatingValue
		for _, rv := range rs.RatingValues {
			var ratingValueResponse model.RatingValue
			ratingValueResponse.RatingValueID = rv.RatingValueID
			ratingValueResponse.RatingSystemID = rv.RatingSystemID
			ratingValueResponse.Value = rv.Value
			ratingValueResponse.MinAge = rv.MinAge
			ratingValueResponse.Description = rv.Description

			ratingValues = append(ratingValues, ratingValueResponse)
		}

		ratingSystemResponse.RatingValues = ratingValues

		responses = append(responses, ratingSystemResponse)
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:   "Ratings",
		Heading: "Ratings List",
		Ratings: responses,
	}

	RenderTemplate(w, "static/html/rating.html", data)
}

// GetAllRatingSystems handler function for GET method
func (rsh *RatingHandler) GetAllRatingSystems(w http.ResponseWriter, r *http.Request) {
	ratingSystems := []model.RatingSystem{}
	err := rsh.db.Preload("Country").Preload("RatingValues").Find(&ratingSystems).Error
	if err != nil {
		rsh.handleError(w, err)
		return
	}
	rsh.encodeJSONResponse(w, ratingSystems)
}

// GetRatingSystemById handler function for GET method
func (rsh *RatingHandler) GetRatingSystemById(w http.ResponseWriter, r *http.Request) {
	ratingSystemId, err := strconv.ParseInt(mux.Vars(r)["ratingSystemId"], 10, 64)
	if err != nil {
		rsh.handleError(w, err)
		return
	}
	ratingSystem := &model.RatingSystem{}
	err = rsh.db.Preload("Country").Preload("RatingValues").Where("rating_system_id = ?", ratingSystemId).First(ratingSystem).Error
	if err != nil {
		rsh.handleError(w, err)
		return
	}
	rsh.encodeJSONResponse(w, ratingSystem)
}

// handleError ...
func (rsh *RatingHandler) handleError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{"status": false, "message": err.Error()}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// encodeJSONResponse ...
func (rsh *RatingHandler) encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
