package logging

import (
	"fmt"
	"github.com/dp1140a/geoip/version"
	"github.com/spf13/viper"
)

const (
	_LOGGING       = "logging"
	_TRACE_LOGGING = "tracelogging"
	_LOG_FILE      = "logfile"
)

/**
  tracelogging = true
  logfile="log/d5.log"
*/
type Config struct {
	TraceLogging bool
	LogFile      string
}

func InitConfig() (config *Config, err error) {
	l := viper.Sub(_LOGGING)
	fmt.Println("App_name: " + version.APP_NAME)
	if l == nil {
		config = &Config{
			TraceLogging: false,
			LogFile:      "./log/" + version.APP_NAME + ".log",
		}
		return config, err
	} else {
		setViperDefaults()
		config = &Config{}
		config.TraceLogging = l.GetBool(_TRACE_LOGGING)
		config.LogFile = l.GetString(_LOG_FILE)
		return config, err
	}
}

func setViperDefaults() {
	lv := viper.Sub(_LOGGING)
	lv.SetDefault(_TRACE_LOGGING, false)
	lv.SetDefault(_LOG_FILE, "./log/"+version.APP_NAME+".log")
}

func String() string {
	return `
[logging]
    tracelogging = true
    logfile="log/geoip.log"
`
}
