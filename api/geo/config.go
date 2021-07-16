package geo

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	_GEO = "geoip"
)

type Config struct {
	AccountId         int
	DatabaseDirectory string
	DatabaseName      string
	EditionIDs        []string
	LicenseKey        string
	RefreshDuration   time.Duration
	Verbose           bool
}

func InitConfig() (config *Config, err error) {
	dbrt, _ := time.ParseDuration("1h")

	config = &Config{
		AccountId:         123,
		DatabaseDirectory: "/usr/share/GeoIP",
		DatabaseName:      "GeoLite2-City.mmdb",
		EditionIDs:        []string{"GeoLite2-City"},
		LicenseKey:        "testing",
		RefreshDuration:   dbrt,
		Verbose:           false,
	}
	h := viper.Sub(_GEO)
	if h != nil {
		err := h.Unmarshal(config)
		if err != nil {
			log.Error("Geo Config Error: ", err.Error())
			return nil, err
		}
	}
	return config, nil
}

func String() string {
	return `
[geoip]
    geoDB="GeoIP2-city.mmdb"
    dbRefreshTime="1h"
    accountId="YOURACCTNUMBER"
    DatabaseDirectory="data"
    DatabaseName="GeoLite2-City.mmdb"
    LicenseKey="YOURLICENSEKEY"
`
}
