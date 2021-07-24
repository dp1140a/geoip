package httplogging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LOGGING_HTTP = "logging.http"
	LOG_DIR      = "/var/log/geoip"
	HTTP_LOG     = "geoip_http.log"
)

type HttpLoggingConfig struct {
	Enabled bool
	StdOut  bool
	FileOut bool
	LogFile string
}

func InitConfig() (config *HttpLoggingConfig) {
	config = &HttpLoggingConfig{
		Enabled: true,
		StdOut:  true,
		FileOut: false,
		LogFile: fmt.Sprintf("%v/%v.log", LOG_DIR, HTTP_LOG),
	}

	h := viper.Sub(LOGGING_HTTP)
	if h != nil {
		err := h.Unmarshal(config)
		if err != nil {
			log.Error("HTTP Config Error: ", err.Error())
			return config
		}
	}

	return config
}

func String() string {
	return `    [logging.http]
        enabled=true
        stdout=true
        fileout=true
        logfile="/var/log/geoip/geoip_http.log"
`
}
