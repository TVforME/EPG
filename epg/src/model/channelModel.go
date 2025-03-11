// channel model
package model

import (
	"time"

	config "epg/src/config"

	"gorm.io/gorm"
)

// Channel represents a channel with endpoint responses.
type Channel struct {
	ChannelID           uint      `gorm:"primaryKey;autoIncrement" json:"channelID"`
	NetworkID           uint      `gorm:"not null" json:"networkID"`
	Description         string    `gorm:"type:text;not null" json:"description"`
	BroadcastStartTime  time.Time `gorm:"not null" json:"broadcastStartTime"`
	BroadcastFinishTime time.Time `gorm:"not null" json:"broadcastFinishTime"`
	ServiceID           uint      `gorm:"not null" json:"serviceID"`
	ServiceVPid         uint      `gorm:"not null" json:"serviceVPid"`
	ServiceAPid         uint      `gorm:"not null" json:"serviceAPid"`
	AuthorityMeta       *string   `gorm:"type:text" json:"authorityMeta"`
	LogoName            *string   `gorm:"type:text" json:"logoName"`
	Network             Network   `gorm:"foreignKey:NetworkID;constraint:OnDelete:CASCADE" json:"-"`
	Events              []Event   `gorm:"foreignKey:ChannelID" json:"events"`
	NetworkName         string    `gorm:"-" json:"networkName"`
}

// PopulateInitialChannelValues populates the database with initial channel values called from init()
func PopulateInitialChannelValues(db *gorm.DB) error {
	// Loop through each channel
	for _, channelConfig := range config.Config.Channels {
		// Create a new channel instance
		channel := &Channel{
			NetworkID:           channelConfig.NetworkID,
			Description:         channelConfig.Description,
			BroadcastStartTime:  channelConfig.BroadcastStartTime,
			BroadcastFinishTime: channelConfig.BroadcastFinishTime,
			ServiceID:           channelConfig.ServiceID,
			ServiceVPid:         channelConfig.ServiceVPid,
			ServiceAPid:         channelConfig.ServiceAPid,
			AuthorityMeta:       &channelConfig.AuthorityMeta,
			LogoName:            &channelConfig.LogoName,
		}

		// Insert the channel instance into the database
		err := db.Create(channel).Error
		if err != nil {
			return err
		}
	}

	return nil
}
