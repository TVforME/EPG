// eventHandler.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// EventHandler ...
type EventHandler struct {
	db *gorm.DB
}

// NewEventHandler ...
func NewEventHandler(db *gorm.DB) *EventHandler {
	return &EventHandler{db: db}
}

// GetAllEvents handler function for GET method
func (eh *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events := []model.Event{}
	err := eh.db.Preload("Channel").Preload("Category").Preload("Genre").Preload("EventRatings").Find(&events).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}
	eh.encodeJSONResponse(w, events)
}

// GetEventById handler function for GET method
func (eh *EventHandler) GetEventById(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["eventId"], 10, 64)
	if err != nil {
		eh.handleError(w, err)
		return
	}
	event := &model.Event{}
	err = eh.db.Preload("Channel").Preload("Category").Preload("Genre").Preload("EventRatings").Where("event_id = ?", eventId).First(event).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}
	eh.encodeJSONResponse(w, event)
}

// http://localhost:8080/eventbytime/2025-02-09%2000:00:00
// GetEventByTime handler function for GET method
func (eh *EventHandler) GetEventByTime(w http.ResponseWriter, r *http.Request) {
	timeStr := mux.Vars(r)["time"]
	timeLayout := "2006-01-02 15:04:05"
	timeCreated, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		eh.handleError(w, err)
		return
	}

	events := []model.Event{}
	err = eh.db.Preload("Channel").
		Preload("Category").
		Preload("Genre").
		Preload("Genre.GenreColor").
		Preload("EventRatings").
		Preload("EventRatings.RatingValue").
		Preload("EventRatings.RatingValue.RatingSystem").
		Preload("EventRatings.RatingValue.RatingSystem.Country").
		Where("created_at >= ?", timeCreated).
		Find(&events).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}

	eh.encodeJSONResponse(w, events)
}

// CreateEvent handler function for POST method
func (eh *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	event := &model.Event{}
	err := json.NewDecoder(r.Body).Decode(event)
	if err != nil {
		eh.handleError(w, err)
		return
	}
	err = eh.db.Create(event).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}
	eh.encodeJSONResponse(w, event)
}

// UpdateEvent handler function for PUT method
func (eh *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	event := &model.Event{}
	err := json.NewDecoder(r.Body).Decode(event)
	if err != nil {
		eh.handleError(w, err)
		return
	}
	eventId, err := strconv.ParseInt(mux.Vars(r)["eventId"], 10, 64)
	if err != nil {
		eh.handleError(w, err)
		return
	}
	event.EventID = uint(eventId)
	err = eh.db.Save(event).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}
	eh.encodeJSONResponse(w, event)
}

// DeleteEvent handler function for DELETE method
func (eh *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["eventId"], 10, 64)
	if err != nil {
		eh.handleError(w, err)
		return
	}
	event := &model.Event{}
	err = eh.db.Where("event_id = ?", eventId).First(event).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}
	err = eh.db.Delete(event).Error
	if err != nil {
		eh.handleError(w, err)
		return
	}
	eh.encodeJSONResponse(w, map[string]interface{}{"message": "Event deleted successfully"})
}

// handleError ...
func (eh *EventHandler) handleError(w http.ResponseWriter, err error) {
	msg := map[string]interface{}{"status": false, "message": err.Error()}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// encodeJSONResponse ...
func (eh *EventHandler) encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
