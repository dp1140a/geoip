package influx

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_INFLUX = "influxdb"
)

type Config struct {
	Url    string
	Token  string
	Org    string
	Bucket string
}

func InitConfig() (config *Config, err error) {
	config = &Config{
		Url:    "http://localhost:8086",
		Token:  "myToken",
		Org:    "myOrg",
		Bucket: "myBucket",
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
	url="http://localhost:8086"
	token="myToken"
	org="myOrg"
	bucket="myBucket"
`
}
