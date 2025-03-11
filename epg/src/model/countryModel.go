// countries model
package model

import (
	"strings"

	"gorm.io/gorm"
)

// Country represents a country
type Country struct {
	CountryID   uint       `gorm:"primaryKey;autoIncrement"`
	CountryCode string     `gorm:"not null;type:char(2);unique"`
	CountryName string     `gorm:"not null;type:varchar(50)"`
	Region      string     `gorm:"not null;type:varchar(50)"`
	Timezones   []Timezone `gorm:"foreignKey:CountryCode;references:CountryCode"`
}

// LoadFromCSV loads countries from a CSV file
func (c *Country) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new country instance
		country := &Country{
			CountryCode: strings.TrimSpace(record[0]),
			CountryName: strings.TrimSpace(record[1]),
			Region:      strings.TrimSpace(record[2]),
		}

		// Insert the country instance into the database
		err = db.Create(country).Error
		if err != nil {
			return err
		}
	}

	return nil
}
