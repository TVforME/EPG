// rating system model
package model

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// RatingValue represents a rating value
type RatingValue struct {
	RatingValueID  uint          `gorm:"primaryKey;autoIncrement" json:"RatingValueID"`
	RatingSystemID uint          `gorm:"not null" json:"RatingSystemID"`
	Value          string        `gorm:"not null;type:text" json:"Value"`
	MinAge         uint          `gorm:"not null" json:"MinAge"`
	Description    string        `gorm:"type:text" json:"Description"`
	RatingSystem   RatingSystem  `gorm:"foreignKey:RatingSystemID"`
	EventRatings   []EventRating `gorm:"foreignKey:RatingValueID"`
}

// RatingSystem represents a rating system
type RatingSystem struct {
	RatingSystemID uint          `gorm:"primaryKey;autoIncrement" json:"RatingSystemID"`
	CountryID      uint          `gorm:"not null" json:"CountryID"`
	Description    string        `gorm:"not null;type:varchar(100)" json:"Description"`
	Country        Country       `gorm:"foreignKey:CountryID;references:CountryID"`
	RatingValues   []RatingValue `gorm:"foreignKey:RatingSystemID" json:"ratingValues"`
}

// LoadFromCSV loads rating systems from a CSV file
func (rs *RatingSystem) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new rating system instance
		ratingSystemID, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return err
		}
		country := &Country{}
		err = db.Where("country_code = ?", strings.TrimSpace(record[1])).First(country).Error
		if err != nil {
			return err
		}
		ratingSystem := &RatingSystem{
			RatingSystemID: uint(ratingSystemID),
			CountryID:      country.CountryID,
			Description:    strings.TrimSpace(record[2]),
		}

		// Insert the rating system instance into the database
		err = db.Create(ratingSystem).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadFromCSV loads rating values from a CSV file
func (rv *RatingValue) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new rating value instance
		ratingValueID, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		ratingSystemID, err := strconv.Atoi(record[1])
		if err != nil {
			return err
		}
		minAge, err := strconv.Atoi(record[3])
		if err != nil {
			return err
		}
		ratingValue := &RatingValue{
			RatingValueID:  uint(ratingValueID),
			RatingSystemID: uint(ratingSystemID),
			Value:          record[2],
			MinAge:         uint(minAge),
			Description:    record[4],
		}

		// Insert the rating value instance into the database
		err = db.Create(ratingValue).Error
		if err != nil {
			return err
		}
	}

	return nil
}
