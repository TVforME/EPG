// genre model join to genre colors
package model

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// GenreColor represents a genre color
type GenreColor struct {
	ColorID      uint   `gorm:"column:color_id;primaryKey;autoIncrement"`
	NibbleLevel1 uint8  `gorm:"column:nibble_level_1;not null;unique"`
	ColorHex     string `gorm:"column:color_hex;not null;type:text"`
}

// Genre represents event genre
type Genre struct {
	GenreID      uint        `gorm:"column:genre_id;primaryKey;autoIncrement" json:"genreID"`
	NibbleLevel1 uint8       `gorm:"column:nibble_level_1;not null;index:idx_genres_nibble_level_1" json:"nibbleLevel1"`
	NibbleLevel2 uint8       `gorm:"column:nibble_level_2;not null" json:"nibbleLevel2"`
	Description  string      `gorm:"column:description;not null;type:text;index:idx_genres_description" json:"description"`
	GenreColor   *GenreColor `gorm:"foreignKey:NibbleLevel1;references:NibbleLevel1" json:"-"`
	ColorHex     string      `gorm:"-" json:"colorHex"`
}

// LoadFromCSV loads genre colors from a CSV file
func (gc *GenreColor) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new genre color instance
		nibbleLevel1, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return err
		}
		genreColor := &GenreColor{
			NibbleLevel1: uint8(nibbleLevel1),
			ColorHex:     record[1],
		}

		// Insert the genre color instance into the database
		err = db.Create(genreColor).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadFromCSV loads genres from a CSV file
func (g *Genre) LoadFromCSV(db *gorm.DB, filename string) error {
	// Load the CSV records
	records, err := loadCSVRecords(filename)
	if err != nil {
		return err
	}

	// Insert the records into the database
	for _, record := range records {
		// Create a new genre instance
		nibbleLevel1, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return err
		}
		nibbleLevel2, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			return err
		}
		genre := &Genre{
			NibbleLevel1: uint8(nibbleLevel1),
			NibbleLevel2: uint8(nibbleLevel2),
			Description:  strings.TrimSpace(record[2]),
		}

		// Insert the genre instance into the database
		err = db.Create(genre).Error
		if err != nil {
			return err
		}
	}

	return nil
}
