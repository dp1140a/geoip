package geo

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	_GEO           = "geoip"
	DATA_DIR       = "data/geoip"
	DB_NAME        = "GeoLite2-City.mmdb"
	EDITION        = "GeoLite2-City"
	DBREFRESH      = "24h"
	VERBOSE_UPDATE = false
	LICENSE_KEY    = "testing"
	ACCOUNT_ID     = 123
)

type GeoIpConfig struct {
	AccountId         int
	DatabaseDirectory string
	DatabaseName      string
	EditionIDs        []string
	LicenseKey        string
	RefreshDuration   time.Duration
	VerboseUpdate     bool
}

func InitGeoIpConfig() (config *GeoIpConfig, err error) {
	dbrt, _ := time.ParseDuration(DBREFRESH)

	config = &GeoIpConfig{
		AccountId:         ACCOUNT_ID,
		DatabaseDirectory: DATA_DIR,
		DatabaseName:      DB_NAME,
		EditionIDs:        []string{EDITION},
		LicenseKey:        LICENSE_KEY,
		RefreshDuration:   dbrt,
		VerboseUpdate:     VERBOSE_UPDATE,
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
    RefreshDuration = "` + DBREFRESH + `"
    accountId = ` + fmt.Sprintf("%d", ACCOUNT_ID) + `
    DatabaseDirectory = "` + DATA_DIR + `"
    DatabaseName = "` + DB_NAME + `"
    LicenseKey = "` + LICENSE_KEY + `"
    Verbose = ` + fmt.Sprintf("%t", VERBOSE_UPDATE) + `
`
}
