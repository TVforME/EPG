// event category model
package model

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Category represents a category
type Category struct {
	CategoryID  uint    `gorm:"primaryKey;autoIncrement" json:"categoryID"`
	Description string  `gorm:"not null;type:text;unique" json:"description"`
	Events      []Event `gorm:"foreignKey:CategoryID" json:"-"`
}

// LoadFromCSV loads categories from a CSV file
func (c *Category) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new category instance
		categoryID, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return err
		}
		category := &Category{
			CategoryID:  uint(categoryID),
			Description: strings.TrimSpace(record[1]),
		}

		// Insert the category instance into the database
		err = db.Create(category).Error
		if err != nil {
			return err
		}
	}

	return nil
}
