package update

import (
	"fmt"
	"github.com/maxmind/geoipupdate/v4/pkg/geoipupdate"
	"github.com/maxmind/geoipupdate/v4/pkg/geoipupdate/database"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

const (
	URL      = "https://updates.maxmind.com"
	LOCKFILE = ".geoipupdate.lock"
)

type GeoIpUpdater struct {
	client *http.Client
	config *geoipupdate.Config
}

func NewGeoIPUpdater(config *geoipupdate.Config) (gu *GeoIpUpdater) {
	config.LockFile = filepath.Join(config.DatabaseDirectory, LOCKFILE)
	config.URL = URL
	config.Verbose = true
	gu = &GeoIpUpdater{
		config: config,
	}
	gu.client = geoipupdate.NewClient(gu.config)

	return gu
}

/**
https://github.com/maxmind/geoipupdate/blob/main/cmd/geoipupdate/main.go
*/
func (gu GeoIpUpdater) Update() error {
	dbReader := database.NewHTTPDatabaseReader(gu.client, gu.config)
	log.Info("Starting GeoIp Database Update")
	fmt.Println(gu.config)
	for _, editionID := range gu.config.EditionIDs {
		log.Info("Updating ", editionID)
		filename, err := geoipupdate.GetFilename(gu.config, editionID, gu.client)
		if err != nil {
			return errors.Wrapf(err, "error retrieving filename for %s", editionID)
		}
		filePath := filepath.Join(gu.config.DatabaseDirectory, filename)
		dbWriter, err := database.NewLocalFileDatabaseWriter(filePath, gu.config.LockFile, gu.config.Verbose)
		if err != nil {
			return errors.Wrapf(err, "error creating database writer for %s", editionID)
		}
		if err := dbReader.Get(dbWriter, editionID); err != nil {
			return errors.WithMessagef(err, "error while getting database for %s", editionID)
		}
	}
	log.Info("Completed GeoIp DataBase Update")
	return nil
}
