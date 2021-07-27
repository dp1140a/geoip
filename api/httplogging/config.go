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
	ENABLED      = true
	STDOUT       = true
	FILEOUT      = false
)

type HttpLoggingConfig struct {
	Enabled bool
	StdOut  bool
	FileOut bool
	LogFile string
}

func InitHttpLoggingConfig() (config *HttpLoggingConfig) {
	config = &HttpLoggingConfig{
		Enabled: ENABLED,
		StdOut:  STDOUT,
		FileOut: FILEOUT,
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
        enabled = ` + fmt.Sprintf("%t", ENABLED) + `
        stdout = ` + fmt.Sprintf("%t", STDOUT) + `
        fileout = ` + fmt.Sprintf("%t", FILEOUT) + `
        logfile = "` + fmt.Sprintf("%s/%s", LOG_DIR, HTTP_LOG) + `"
`
}
