// events model
package model

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	EventID             uint           `gorm:"primaryKey;autoIncrement"`
	ChannelID           uint           `gorm:"not null;index:idx_events_channel_start"`
	StartTime           time.Time      `gorm:"not null;index:idx_events_channel_start"`
	EndTime             time.Time      `gorm:"not null"`
	Title               string         `gorm:"type:text;not null;index:idx_events_title"`
	ShortDescription    *string        `gorm:"column:short_description;type:text"`
	ExtendedDescription *string        `gorm:"column:extended_description;type:text"`
	GenreID             uint           `gorm:"not null"`
	CategoryID          uint           `gorm:"not null"`
	CreatedAt           time.Time      `gorm:"type:datetime;default:current_timestamp;not null"`
	UpdatedAt           time.Time      `gorm:"type:datetime;default:current_timestamp;not null"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	Channel             Channel        `gorm:"foreignKey:ChannelID;constraint:OnDelete:CASCADE"`
	Category            Category       `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Genre               Genre          `gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE"`
	EventRatings        []EventRating  `gorm:"foreignKey:EventID"`
}

type EventRating struct {
	EventID       uint        `gorm:"primaryKey"`
	RatingValueID uint        `gorm:"primaryKey"`
	Event         Event       `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	RatingValue   RatingValue `gorm:"foreignKey:RatingValueID;constraint:OnDelete:CASCADE"`
}

// EventTemplate represents an event template
type EventTemplate struct {
	Title         string
	StartMinute   int
	CategoryID    uint
	GenreID       uint
	RatingValueID uint
}

// PopulateInitialEvents populates the database with initial events from CSV template
func (e *Event) PopulateInitialEvents(db *gorm.DB, channelIDs []uint, startTime time.Time) error {
	// Read events template from CSV file
	eventsTemplate, err := readEventsTemplate("./csv/events_template.csv")
	if err != nil {
		return err
	}

	// Loop through each channel
	for _, channelID := range channelIDs {
		// Loop through each event in the template
		for _, eventTemplate := range eventsTemplate {
			// Calculate the start and finish times for the event
			eventStartTime := startTime.Add(time.Duration(eventTemplate.StartMinute) * time.Minute)
			eventFinishTime := eventStartTime.Add(15 * time.Minute)

			// Create a new event instance
			event := &Event{
				ChannelID:  channelID,
				StartTime:  eventStartTime,
				EndTime:    eventFinishTime,
				Title:      eventTemplate.Title,
				CategoryID: eventTemplate.CategoryID,
				GenreID:    eventTemplate.GenreID,
			}

			// Check if the event already exists
			var existingEvent Event
			err = db.Where("start_time = ? AND channel_id = ?", eventStartTime, channelID).First(&existingEvent).Error
			if err == nil {
				// Event already exists, skip it
				continue
			}

			// Insert the event instance into the database
			err = db.Create(event).Error
			if err != nil {
				return err
			}

			// Create a new event rating instance
			eventRating := &EventRating{
				EventID:       event.EventID,
				RatingValueID: eventTemplate.RatingValueID,
			}

			// Insert the event rating instance into the database
			err = db.Create(eventRating).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// readEventsTemplate reads the events template from a CSV file
func readEventsTemplate(filename string) ([]EventTemplate, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row and discard it
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Create event templates from the records
	var eventTemplates []EventTemplate
	for _, record := range records {
		eventTemplate := EventTemplate{
			Title:         record[0],
			StartMinute:   atoi(record[1]),
			CategoryID:    uint(atoi(record[2])),
			GenreID:       uint(atoi(record[3])),
			RatingValueID: uint(atoi(record[4])),
		}
		eventTemplates = append(eventTemplates, eventTemplate)
	}

	return eventTemplates, nil
}

// atoi converts a string to an integer, if negative then return 0
func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	if i < 0 {
		return 0
	}

	return i
}
