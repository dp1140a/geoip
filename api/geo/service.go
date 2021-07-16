package geo

import (
	"context"
	"fmt"
	"github.com/dp1140a/geoip/api/geo/update"
	influxwriter "github.com/dp1140a/geoip/api/influx"
	"github.com/dp1140a/geoip/models"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	influx "github.com/influxdata/influxdb/models"
	"github.com/maxmind/geoipupdate/v4/pkg/geoipupdate"
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	SERVICE_NAME      = "GeoService"
	BLOCK_MEASUREMENT = "ip_block"
	TIME_LAYOUT       = "Jul 13 15:20:54"
)

type GeoIPService struct {
	name           string
	Config         *Config
	database       *geoip2.Reader
	Updater        *update.GeoIpUpdater
	InfluxDbWriter *influxwriter.InfluxDbWriter
}

func NewGeoIPService(ctx context.Context) (geoService models.Service) {
	log.Info("Initializing Service ", SERVICE_NAME)
	config, err := InitConfig()
	if err != nil {
		log.Error("NewGeoIPService Config error: ", err)
	}
	updater := update.NewGeoIPUpdater(mapConfig(config))

	db, err := openDB(filepath.Join(config.DatabaseDirectory, config.DatabaseName))
	if err != nil {
		log.Warn("No Geo DB found. Fetching")
		err = updater.Update()
		if err != nil {
			log.Error("Cannot retrieve Geo DB.  Shutting down")
			os.Exit(1)
		}
		db, _ = openDB(filepath.Join(config.DatabaseDirectory, config.DatabaseName))
	}

	writer, err := influxwriter.NewInfluxDBWriter()
	if err != nil {
		log.Error("Cannot get db connection ", err)
		return nil
	}

	gs := GeoIPService{
		Config:         config,
		database:       db,
		name:           SERVICE_NAME,
		Updater:        updater,
		InfluxDbWriter: writer,
	}

	ticker := time.NewTicker(config.RefreshDuration)

	userInput := make(chan string)

	go readInput(userInput)

	go func(updater *update.GeoIpUpdater, input chan<- string) {
		for {
			select {
			case <-ticker.C:
				log.Info("Updating GeoIp Database")
				err := updater.Update()
				if err != nil {
					log.Error(err)
				}
			case userAnswer := <-userInput:
				log.Info(userAnswer)
				err := updater.Update()
				if err != nil {
					log.Error(err)
				}

			case <-ctx.Done():
				log.Warn("Stopping updater.")
				gs.shutdown()
				ticker.Stop()
				return
			}
		}
	}(updater, userInput)

	return gs
}

func (gs GeoIPService) GetName() (name string) {
	return gs.name
}

func openDB(dbFile string) (*geoip2.Reader, error) {
	//If no DB go get one
	db, err := geoip2.Open(dbFile)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	log.Info("DB found at ", dbFile)
	//defer db.Close()
	return db, nil
}

func (gs GeoIPService) shutdown() {
	gs.database.Close()
}

func mapConfig(config *Config) (gconfig *geoipupdate.Config) {
	return &geoipupdate.Config{
		AccountID:         config.AccountId,
		DatabaseDirectory: config.DatabaseDirectory,
		LicenseKey:        config.LicenseKey,
		EditionIDs:        config.EditionIDs,
		Verbose:           config.Verbose,
	}
}

func (gs GeoIPService) locate(ipaddress string) (record *geoip2.City, err error) {
	log.Debug("Locating ", ipaddress)
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipaddress)
	record, err = gs.database.City(ip)
	if err != nil {
		log.Error("Error locating ip. ", err)
		return &geoip2.City{}, err
	}

	return record, nil
}

/**
{
  "action": "block",
  "data_length": "0",
  "dest_ip": "173.160.205.9",
  "dest_port": "22",
  "direction": "in",
  "flags": "none",
  "id": "54321",
  "iface": "em0",
  "ip_ver": "4",
  "length": "44",
  "offset": "0",
  "pid": "45720",
  "program": "filterlog",
  "reason": "match",
  "rule": "102",
  "src_ip": "198.98.56.149",
  "src_port": "44397",
  "timestamp": "Jul  8 09:58:38",
  "tos": "0x20",
  "tracker": "1521923716",
  "ttl": "231"
}

*/
func (gs GeoIPService) mapPointsWithLookup(body []byte) (err error) {
	var newPoints []*write.Point
	points, err := influx.ParsePointsWithPrecision(body, time.Now().UTC(), "n")
	if err != nil {
		log.Error("Error parsing line protocol. ", err)
	}

	for _, point := range points {
		fields, _ := point.Fields()
		location, err := gs.locate(fmt.Sprintf("%v", fields["src_ip"]))
		if err != nil {

		}

		city := ""
		if location.City.Names == nil {
			city = "unknown"
		} else {
			city = location.City.Names["en"]
		}

		newPoint := influxdb2.NewPointWithMeasurement(BLOCK_MEASUREMENT).
			AddTag("iface", fmt.Sprintf("%v", fields["iface"])).
			AddField("dest_ip", fields["dest_ip"]).
			AddField("dest_port", fields["dest_port"]).
			AddField("rule", fields["rule"]).
			AddField("src_ip", fields["src_ip"]).
			AddField("src_port", fields["src_port"]).
			AddField("country", location.Country.IsoCode).
			AddField("city", city).
			AddField("latitude", location.Location.Latitude).
			AddField("longitude", location.Location.Longitude).
			SetTime(point.Time())

		newPoints = append(newPoints, newPoint)
		//fmt.Println(point2Str(newPoint))
	}
	log.Infof("Appending %v new Points", len(newPoints))

	err = gs.InfluxDbWriter.WritePoints(newPoints)
	if err != nil {
		return err
	}

	return nil
}

func readInput(ch chan<- string) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var b = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		ch <- string(b)
	}
}

func point2Str(point *write.Point) string {
	builder := strings.Builder{}
	builder.WriteString(point.Name())
	builder.WriteString(",")
	var tags []string
	for _, tag := range point.TagList() {
		tags = append(tags, fmt.Sprintf("%s=%s", tag.Key, tag.Value))
	}
	builder.WriteString(strings.Join(tags, ","))
	builder.WriteString(" ")
	var fields []string
	for _, field := range point.FieldList() {
		fields = append(fields, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}
	builder.WriteString(strings.Join(fields, ","))
	builder.WriteString(" ")
	builder.WriteString(fmt.Sprintf("%v", point.Time().UnixNano()))

	// Convert Builder to String and print it.
	return builder.String()
}
