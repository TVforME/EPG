// eventRatingHandler.go
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

// EventRatingHandler ...
type EventRatingHandler struct {
	db *gorm.DB
}

// NewEventRatingHandler ...
func NewEventRatingHandler(db *gorm.DB) *EventRatingHandler {
	return &EventRatingHandler{db: db}
}

// GetAllEventRatings handler function for GET method
func (erh *EventRatingHandler) GetAllEventRatings(w http.ResponseWriter, r *http.Request) {
	eventRatings := []model.EventRating{}
	err := erh.db.Preload("Event").Preload("RatingValue").Find(&eventRatings).Error
	if err != nil {
		erh.handleError(w, err)
		return
	}
	erh.encodeJSONResponse(w, eventRatings)
}

// GetEventRatingById handler function for GET method
func (erh *EventRatingHandler) GetEventRatingById(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["eventId"], 10, 64)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	ratingValueId, err := strconv.ParseInt(mux.Vars(r)["ratingValueId"], 10, 64)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	eventRating := &model.EventRating{}
	err = erh.db.Preload("Event").Preload("RatingValue").Where("event_id = ? AND rating_value_id = ?", eventId, ratingValueId).First(eventRating).Error
	if err != nil {
		erh.handleError(w, err)
		return
	}
	erh.encodeJSONResponse(w, eventRating)
}

// CreateEventRating handler function for POST method
func (erh *EventRatingHandler) CreateEventRating(w http.ResponseWriter, r *http.Request) {
	eventRating := &model.EventRating{}
	err := json.NewDecoder(r.Body).Decode(eventRating)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	err = erh.db.Create(eventRating).Error
	if err != nil {
		erh.handleError(w, err)
		return
	}
	erh.encodeJSONResponse(w, eventRating)
}

// UpdateEventRating handler function for PUT method
func (erh *EventRatingHandler) UpdateEventRating(w http.ResponseWriter, r *http.Request) {
	eventRating := &model.EventRating{}
	err := json.NewDecoder(r.Body).Decode(eventRating)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	eventId, err := strconv.ParseInt(mux.Vars(r)["eventId"], 10, 64)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	ratingValueId, err := strconv.ParseInt(mux.Vars(r)["ratingValueId"], 10, 64)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	eventRating.EventID = uint(eventId)
	eventRating.RatingValueID = uint(ratingValueId)
	err = erh.db.Save(eventRating).Error
	if err != nil {
		erh.handleError(w, err)
		return
	}
	erh.encodeJSONResponse(w, eventRating)
}

// DeleteEventRating handler function for DELETE method
func (erh *EventRatingHandler) DeleteEventRating(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["eventId"], 10, 64)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	ratingValueId, err := strconv.ParseInt(mux.Vars(r)["ratingValueId"], 10, 64)
	if err != nil {
		erh.handleError(w, err)
		return
	}
	eventRating := &model.EventRating{}
	err = erh.db.Where("event_id = ? AND rating_value_id = ?", eventId, ratingValueId).First(eventRating).Error
	if err != nil {
		erh.handleError(w, err)
		return
	}
	err = erh.db.Delete(eventRating).Error
	if err != nil {
		erh.handleError(w, err)
		return
	}
	erh.encodeJSONResponse(w, map[string]interface{}{"message": "Event rating deleted successfully"})
}

// handleError ...
func (erh *EventRatingHandler) handleError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{"status": false, "message": err.Error()}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// encodeJSONResponse ...
func (erh *EventRatingHandler) encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
