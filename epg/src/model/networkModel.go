// Network Model
package model

import (
	"time"

	config "epg/src/config"

	"gorm.io/gorm"
)

type Network struct {
	NetworkID       uint      `gorm:"column:network_id;primaryKey;autoIncrement" json:"networkID"`
	CountryID       uint      `gorm:"column:country_id;not null" json:"countryID"`
	TimezoneID      uint      `gorm:"column:timezone_id;not null" json:"timezoneID"`
	ServiceID       uint      `gorm:"column:service_id;not null" json:"serviceID"`
	Description     string    `gorm:"column:description;not null;type:text" json:"description"`
	StartTime       time.Time `json:"startTime"`
	FinishTime      time.Time `json:"finishTime"`
	CridDescription string    `gorm:"column:crid_description;type:text" json:"cridDescription"`
	Country         Country   `gorm:"foreignKey:CountryID;references:CountryID" json:"-"`
	Timezone        Timezone  `gorm:"foreignKey:TimezoneID;references:TimeZoneID"`
	CountryCode     string    `gorm:"-" json:"countryCode"`
	CountryName     string    `gorm:"-" json:"countryName"`
	TimezoneName    string    `gorm:"-" json:"TimezoneName"`
	StandardOffset  int       `gorm:"-" json:"StandardOffset"`
	DSTOffset       int       `gorm:"-" json:"DSTOffset"`
}

// PopulateInitialNetworkValues populates the database with initial network values called from init()
func PopulateInitialNetworkValues(db *gorm.DB) error {
	// Loop through each network
	for _, networkConfig := range config.Config.Network {
		// Create a new country instance
		country := &Country{
			CountryCode: networkConfig.CountryCode,
		}

		// Check if the country already exists.
		err := db.First(country, "country_code = ?", networkConfig.CountryCode).Error
		if err != nil {
			return err
		}

		// Get the country ID
		countryID := country.CountryID

		// Lookup the TimezoneID based on the TimezoneName and CountryCode
		var timezone Timezone
		err = db.Where("country_code = ? AND timezone_name = ?", networkConfig.CountryCode, networkConfig.TimezoneName).First(&timezone).Error
		if err != nil {
			return err
		}

		// Create a new network instance
		network := &Network{
			CountryID:       countryID,
			TimezoneID:      timezone.TimeZoneID,
			ServiceID:       networkConfig.ServiceID,
			Description:     networkConfig.Description,
			StartTime:       networkConfig.StartTime,
			FinishTime:      networkConfig.FinishTime,
			CridDescription: networkConfig.CridDescription,
			CountryCode:     country.CountryCode,
			CountryName:     country.CountryName,
			TimezoneName:    timezone.TimezoneName,
		}

		// Insert the network instance into the database
		err = db.Create(network).Error
		if err != nil {
			return err
		}
	}

	return nil
}
