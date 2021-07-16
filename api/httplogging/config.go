package httplogging

import (
	"github.com/dp1140a/geoip/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LOGGING_HTTP = "logging.http"
	_LOG_FILE    = "logfile"
)

/**
  		enabled=true
        stdout=true
        file=true
        logfile="log/d5_http.log"
*/
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
		LogFile: "./log/" + version.APP_NAME + ".log",
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
        logfile="log/geoip_http.log"
`
}
