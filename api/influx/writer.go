package influx

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	SERVICE_NAME = "InfluxDBWriter"
)

type InfluxDbWriter struct {
	Client api.WriteAPI
}

func NewInfluxDBWriter() (influxWriter *InfluxDbWriter, err error) {
	log.Info("Initializing Service ", SERVICE_NAME)
	config, err := InitInfluxDBConfig()
	if err != nil {
		log.Error("NewGeoIPService Config error: ", err)
	}
	client := influxdb2.NewClient(config.Url, config.Token)
	ok, err := client.Ready(context.TODO())
	if err != nil && !ok {
		return nil, errors.New("Cannot connect to InfluxDB.")
	}
	writer := client.WriteAPI(config.Org, config.Bucket)
	return &InfluxDbWriter{Client: writer}, nil

}

func (w *InfluxDbWriter) WritePoints(points []*write.Point) (err error) {

	writeErrors := 0
	// Get errors channel
	errorsCh := w.Client.Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			log.Errorf("InfluxDB write error: %s\n", err.Error())
			writeErrors++
		}
	}()
	// write some points
	for _, p := range points {
		// write asynchronously
		w.Client.WritePoint(p)
	}
	// Force all unwritten data to be sent
	w.Client.Flush()

	if writeErrors > 0 {
		return fmt.Errorf("Error writing %d lines of %d total points.  Check log for more information", writeErrors,
			len(points))
	} else {
		return nil
	}
}
