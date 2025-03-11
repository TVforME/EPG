// timezone model
package model

import (
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Timezone represents a timezone
type Timezone struct {
	TimeZoneID     uint      `gorm:"primaryKey;autoIncrement" json:"timezone_id"`
	CountryCode    string    `gorm:"primaryKey;type:char(2)" json:"country_code"`
	TimezoneName   string    `gorm:"primaryKey;type:varchar(50)" json:"timezone_name"`
	StandardOffset int       `gorm:"not null"`
	DSTOffset      int       `gorm:"not null"`
	DSTStartMonth  int       `gorm:"not null"`
	DSTStartDay    int       `gorm:"not null"`
	DSTStartTime   time.Time `gorm:"not null"`
	DSTEndMonth    int       `gorm:"not null"`
	DSTEndDay      int       `gorm:"not null"`
	DSTEndTime     time.Time `gorm:"not null"`
	IsDefault      bool      `gorm:"default:false"`
	Country        Country   `gorm:"foreignKey:CountryCode;references:CountryCode"`
}

// LoadFromCSV loads timezones from a CSV file
func (t *Timezone) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new timezone instance
		var dstStartMonth int
		var dstStartDay int
		var dstStartTime time.Time
		var dstEndMonth int
		var dstEndDay int
		var dstEndTime time.Time
		if record[4] != "NULL" {
			parts := strings.Split(record[4], " ")
			dateParts := strings.Split(parts[0], "/")
			dstStartMonth, _ = strconv.Atoi(dateParts[1])
			dstStartDay, _ = strconv.Atoi(dateParts[0])
			dstStartTime, _ = time.Parse("15:04", parts[1])
		}
		if record[5] != "NULL" {
			parts := strings.Split(record[5], " ")
			dateParts := strings.Split(parts[0], "/")
			dstEndMonth, _ = strconv.Atoi(dateParts[1])
			dstEndDay, _ = strconv.Atoi(dateParts[0])
			dstEndTime, _ = time.Parse("15:04", parts[1])
		}
		timezone := &Timezone{
			CountryCode:    strings.TrimSpace(record[0]),
			TimezoneName:   strings.TrimSpace(record[1]),
			StandardOffset: parseInt(strings.TrimSpace(record[2])),
			DSTOffset:      parseInt(strings.TrimSpace(record[3])),
			DSTStartMonth:  dstStartMonth,
			DSTStartDay:    dstStartDay,
			DSTStartTime:   dstStartTime,
			DSTEndMonth:    dstEndMonth,
			DSTEndDay:      dstEndDay,
			DSTEndTime:     dstEndTime,
			IsDefault:      parseBool(strings.TrimSpace(record[6])),
		}

		// Insert the timezone instance into the database
		err = db.Create(timezone).Error
		if err != nil {
			return err
		}
	}

	return nil
}
