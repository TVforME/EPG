package epg

import (
	"encoding/json"
	"os"
	"time"
)

type NetworkConfig struct {
	ServiceID       uint      `json:"service_id"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"start_time"`
	FinishTime      time.Time `json:"finish_time"`
	CountryCode     string    `json:"country_code"`
	TimezoneName    string    `json:"timezone_name"`
	CridDescription string    `json:"crid_description"`
}

type ChannelConfig struct {
	Description         string    `json:"description"`
	BroadcastStartTime  time.Time `json:"broadcast_start_time"`
	BroadcastFinishTime time.Time `json:"broadcast_finish_time"`
	ServiceID           uint      `json:"service_id"`
	ServiceVPid         uint      `json:"service_vpid"`
	ServiceAPid         uint      `json:"service_apid"`
	AuthorityMeta       string    `json:"authority_meta"`
	LogoName            string    `json:"logo_name"`
	NetworkID           uint      `json:"network_id"`
}

var Config struct {
	DbType      string          `json:"dbtype"`
	Dbname      string          `json:"dbname"`
	BindPort    int             `json:"bindport"`
	LoadBalance bool            `json:"load_balance"`
	Network     []NetworkConfig `json:"network"`
	Channels    []ChannelConfig `json:"channels"`
}

func init() {
	cfgJson, err := os.ReadFile("src/config/config.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(cfgJson, &Config)
}
