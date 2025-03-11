// ratingValueHandler.go
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

// RatingValueHandler ...
type RatingValueHandler struct {
	db *gorm.DB
}

// NewRatingValueHandler ...
func NewRatingValueHandler(db *gorm.DB) *RatingValueHandler {
	return &RatingValueHandler{db: db}
}

// GetAllRatingValues handler function for GET method
func (rvh *RatingValueHandler) GetAllRatingValues(w http.ResponseWriter, r *http.Request) {
	ratingValues := []model.RatingValue{}
	err := rvh.db.Preload("RatingSystem").Find(&ratingValues).Error
	if err != nil {
		rvh.handleError(w, err)
		return
	}
	rvh.encodeJSONResponse(w, ratingValues)
}

// GetRatingValueById handler function for GET method
func (rvh *RatingValueHandler) GetRatingValueById(w http.ResponseWriter, r *http.Request) {
	ratingValueId, err := strconv.ParseInt(mux.Vars(r)["ratingValueId"], 10, 64)
	if err != nil {
		rvh.handleError(w, err)
		return
	}
	ratingValue := &model.RatingValue{}
	err = rvh.db.Preload("RatingSystem.Country").Preload("RatingSystem").Where("rating_value_id = ?", ratingValueId).First(ratingValue).Error
	if err != nil {
		rvh.handleError(w, err)
		return
	}
	rvh.encodeJSONResponse(w, ratingValue)
}

// handleError ...
func (rvh *RatingValueHandler) handleError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{"status": false, "message": err.Error()}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// encodeJSONResponse ...
func (rvh *RatingValueHandler) encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
