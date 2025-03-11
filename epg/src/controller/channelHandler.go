// channelHandler.go
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// ChannelHandler ...
type ChannelHandler struct {
	db *gorm.DB
}

// NewChannelHandler ...
func NewChannelHandler(db *gorm.DB) *ChannelHandler {
	return &ChannelHandler{db: db}
}

// GetAllChannelsHTML handler function for GET method
func (ch *ChannelHandler) GetAllChannelsHTML(w http.ResponseWriter, r *http.Request) {
	channels := []model.Channel{}
	err := ch.db.Preload("Network").Preload("Network.Country").Preload("Events").Find(&channels).Error
	if err != nil {
		HandleHtmlError(w, err)
		return
	}

	var responses []model.Channel
	for _, c := range channels {
		// Check if the Network association is loaded
		networkName := ""
		if c.NetworkID != 0 {
			var network model.Network
			err := ch.db.First(&network, c.NetworkID).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					HandleHtmlError(w, errors.New("network not found"))
				} else {
					HandleHtmlError(w, err)
				}
				return
			}
			networkName = network.Description
		}

		// Create a new Channel struct with the NetworkName
		newChannel := model.Channel{
			ChannelID:           c.ChannelID,
			Description:         c.Description,
			BroadcastStartTime:  c.BroadcastStartTime,
			BroadcastFinishTime: c.BroadcastFinishTime,
			ServiceID:           c.ServiceID,
			ServiceVPid:         c.ServiceVPid,
			ServiceAPid:         c.ServiceAPid,
			AuthorityMeta:       c.AuthorityMeta,
			LogoName:            c.LogoName,
			NetworkID:           c.NetworkID,
			Network:             c.Network,
			NetworkName:         networkName,
			Events:              c.Events,
		}

		responses = append(responses, newChannel)
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:    "Channel",
		Heading:  "Channel List",
		Channels: responses,
	}

	RenderTemplate(w, "static/html/channel.html", data)
}

// GetChannelByIdHTML handles GET requests for retrieving a channel by its ID.
func (ch *ChannelHandler) GetChannelByIdHTML(w http.ResponseWriter, r *http.Request) {
	channelId, err := strconv.Atoi(mux.Vars(r)["channelId"])
	if err != nil {
		HandleHtmlError(w, err)
		return
	}

	var channel model.Channel
	// Retrieve channel with preloading Network and Events.
	if err = ch.db.Preload("Network").Preload("Network.Country").Preload("Events").Where("channel_id = ?", channelId).First(&channel).Error; err != nil {
		HandleHtmlError(w, err)
		return
	}

	if channel.NetworkID != 0 {
		var network model.Network
		err := ch.db.First(&network, channel.NetworkID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				HandleHtmlError(w, errors.New("channel not found"))
			} else {
				HandleHtmlError(w, err)
			}
			return
		}

		channel.NetworkName = network.Description
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:    "Channel",
		Heading:  "Channel ByID",
		Channels: []model.Channel{channel},
	}

	RenderTemplate(w, "static/html/channel.html", data)
}
