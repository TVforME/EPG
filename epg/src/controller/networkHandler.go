// / networkHandler.go
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"epg/src/model"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

// NetworkHandler ...
type NetworkHandler struct {
	db *gorm.DB
}

// NewNetworkHandler ...
func NewNetworkHandler(db *gorm.DB) *NetworkHandler {
	return &NetworkHandler{db: db}
}

// GetAllNetworksHTML handler function for GET method
func (nh *NetworkHandler) GetAllNetworksHTML(w http.ResponseWriter, r *http.Request) {
	networks := []model.Network{}
	err := nh.db.Preload("Country").Preload("Country.Timezones").Find(&networks).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleHtmlError(w, errors.New("country not found"))
		} else {
			HandleHtmlError(w, err)
		}
		return
	}

	var responses []model.Network
	for _, n := range networks {
		// Find the selected timezone for the network
		var selectedTimezoneName *model.Timezone
		for _, tz := range n.Country.Timezones {
			if tz.TimeZoneID == n.TimezoneID {
				selectedTimezoneName = &tz
				break
			}
		}

		var timezoneName string
		if selectedTimezoneName != nil {
			timezoneName = selectedTimezoneName.TimezoneName
		} else {
			timezoneName = "Unknown"
		}

		responses = append(responses, model.Network{
			NetworkID:       n.NetworkID,
			CountryID:       n.CountryID,
			ServiceID:       n.ServiceID,
			Description:     n.Description,
			StartTime:       n.StartTime,
			FinishTime:      n.FinishTime,
			CridDescription: n.CridDescription,
			Country:         n.Country,
			CountryCode:     n.Country.CountryCode,
			CountryName:     n.Country.CountryName,
			TimezoneName:    timezoneName,
			StandardOffset:  selectedTimezoneName.StandardOffset,
			DSTOffset:       selectedTimezoneName.DSTOffset,
		})
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:   "Network",
		Heading: "Network List",
		Network: responses,
	}

	RenderTemplate(w, "static/html/network.html", data)
}

// GetNetworkById handler function for GET method requests for retrieving a channel by its ID.
func (nh *NetworkHandler) GetNetworkByIdHTML(w http.ResponseWriter, r *http.Request) {
	networkId, err := strconv.ParseInt(mux.Vars(r)["networkId"], 10, 64)
	if err != nil {
		HandleHtmlError(w, err)
		return
	}
	network := &model.Network{}
	err = nh.db.Preload("Country").Preload("Country.Timezones").First(network, "network_id = ?", networkId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleHtmlError(w, errors.New("network not found"))
		} else {
			HandleHtmlError(w, err)
		}
		return
	}

	// Find the selected timezone for the network
	var selectedTimezoneName *model.Timezone
	for _, tz := range network.Country.Timezones {
		if tz.TimeZoneID == network.TimezoneID {
			selectedTimezoneName = &tz
			break
		}
	}

	var timezoneName string
	if selectedTimezoneName != nil {
		timezoneName = selectedTimezoneName.TimezoneName
	} else {
		timezoneName = "Unknown"
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:   "Network",
		Heading: "Network ById",
		Network: []model.Network{
			{
				NetworkID:       network.NetworkID,
				CountryID:       network.CountryID,
				ServiceID:       network.ServiceID,
				Description:     network.Description,
				StartTime:       network.StartTime,
				FinishTime:      network.FinishTime,
				CridDescription: network.CridDescription,
				Country:         network.Country,
				CountryCode:     network.Country.CountryCode,
				CountryName:     network.Country.CountryName,
				TimezoneName:    timezoneName,
				StandardOffset:  selectedTimezoneName.StandardOffset,
				DSTOffset:       selectedTimezoneName.DSTOffset,
			},
		},
	}
	RenderTemplate(w, "static/html/network.html", data)
}

/*
// GetAllNetworksHTML handler function for GET method
func (nh *NetworkHandler) GetAllNetworksHTML(w http.ResponseWriter, r *http.Request) {
    networks := []model.Network{}
    err := nh.loadNetworks(nh.db, &networks, nil)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            HandleHtmlError(w, errors.New("country not found"))
        } else {
            HandleHtmlError(w, err)
        }
        return
    }

    data := nh.preparePageData("Network", "Network List", nh.prepareNetworkData(networks...))
    RenderTemplate(w, "static/html/network.html", data)
}

// GetNetworkById handler function for GET method requests for retrieving a channel by its ID.
func (nh *NetworkHandler) GetNetworkByIdHTML(w http.ResponseWriter, r *http.Request) {
    networkId, err := strconv.ParseInt(mux.Vars(r)["networkId"], 10, 64)
    if err != nil {
        HandleHtmlError(w, err)
        return
    }
    network := &model.Network{}
    err = nh.loadNetworks(nh.db, []*model.Network{network}, "network_id = ?", networkId).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            HandleHtmlError(w, errors.New("network not found"))
        } else {
            HandleHtmlError(w, err)
        }
        return
    }

    data := nh.preparePageData("Network", "Network ById", nh.prepareNetworkData(*network))
    RenderTemplate(w, "static/html/network.html", data)
}

// loadNetworks loads networks from the database with preloaded Country and Timezones
func (nh *NetworkHandler) loadNetworks(db *gorm.DB, networks *[]model.Network, query interface{}) error {
    return db.Preload("Country").Preload("Country.Timezones").Find(networks, query).Error
}

// extractTimezoneInfo extracts timezone information from a network
func (nh *NetworkHandler) extractTimezoneInfo(network model.Network) (string, *model.Timezone) {
    for _, tz := range network.Country.Timezones {
        if tz.TimeZoneID == network.TimezoneID {
            return tz.TimezoneName, &tz
        }
    }
    return "Unknown", nil
}

// prepareNetworkData prepares network data for template rendering
func (nh *NetworkHandler) prepareNetworkData(networks ...model.Network) []model.Network {
    var responses []model.Network
    for _, n := range networks {
        timezoneName, selectedTimezoneName := nh.extractTimezoneInfo(n)
        responses = append(responses, model.Network{
            NetworkID:       n.NetworkID,
            CountryID:       n.CountryID,
            ServiceID:       n.ServiceID,
            Description:     n.Description,
            StartTime:       n.StartTime,
            FinishTime:      n.FinishTime,
            CridDescription: n.CridDescription,
            Country:         n.Country,
            CountryCode:     n.Country.CountryCode,
            CountryName:     n.Country.CountryName,
            TimezoneName:    timezoneName,
            StandardOffset:  selectedTimezoneName.StandardOffset,
            DSTOffset:       selectedTimezoneName.DSTOffset,
        })
    }
    return responses
}

// preparePageData prepares page data for template rendering
func (nh *NetworkHandler) preparePageData(title, heading string, networks []model.Network) PageData {
    return PageData{
        Title:   title,
        Heading: heading,
        Network: networks,
    }
}

*/
