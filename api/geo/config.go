package geo

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	_GEO     = "geoip"
	DATA_DIR = "data/geoip"
	DB_NAME  = "GeoLite2-City.mmdb"
	EDITION  = "GeoLite2-City"
	DBRT     = "1d"
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
		DatabaseDirectory: DATA_DIR,
		DatabaseName:      DB_NAME,
		EditionIDs:        []string{EDITION},
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
    RefreshDuration="1d"
    accountId=1234
    DatabaseDirectory="data/geoip"
    DatabaseName="GeoLite2-City.mmdb"
    LicenseKey="YOURLICENSEKEY"
`
}
