package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_LOGGING     = "logging"
	LOG_FILE     = "geoip.log"
	LOG_DIR      = "/var/log/geoip"
	TRACELOGGING = true
)

type LoggingConfig struct {
	TraceLogging bool
	LogFile      string
}

func InitLoggingConfig() (config *LoggingConfig) {
	config = &LoggingConfig{
		TraceLogging: TRACELOGGING,
		LogFile:      fmt.Sprintf("%v/%v.log", LOG_DIR, LOG_FILE),
	}

	l := viper.Sub(_LOGGING)
	if l != nil {
		err := l.Unmarshal(config)
		if err != nil {
			log.Error("Logging Config Error: ", err.Error())
			return config
		}
	}

	return config

}

func String() string {
	return `
[logging]
    tracelogging = ` + fmt.Sprintf("%t", TRACELOGGING) + `
    logfile = "` + fmt.Sprintf("%s/%s", LOG_DIR, LOG_FILE) + `"
`
}
