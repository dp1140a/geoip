package influx

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_INFLUX = "influxdb"
	DB_URL  = "http://localhost:8086"
	TOKEN   = "mytoken"
	ORG     = "myOrg"
	BUCKET  = "myBucket"
)

type InfluxDBConfig struct {
	Url    string
	Token  string
	Org    string
	Bucket string
}

func InitInfluxDBConfig() (config *InfluxDBConfig, err error) {
	config = &InfluxDBConfig{
		Url:    DB_URL,
		Token:  TOKEN,
		Org:    ORG,
		Bucket: BUCKET,
	}

	h := viper.Sub(_INFLUX)
	if h != nil {
		err := h.Unmarshal(config)
		if err != nil {
			log.Error("InfluxDB Config Error: ", err.Error())
			return nil, err
		}
	}
	return config, nil
}

func String() string {
	return `
[influxdb]
    url = "` + DB_URL + `"
    token = "` + TOKEN + `"
    org = "` + ORG + `"
    bucket = "` + BUCKET + `"
`
}
