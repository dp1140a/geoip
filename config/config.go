package config

import (
	"encoding/json"
	"fmt"
	"github.com/dp1140a/geoip/api"
	"github.com/dp1140a/geoip/api/geo"
	"github.com/dp1140a/geoip/api/httplogging"
	"github.com/dp1140a/geoip/logging"
	"github.com/dp1140a/geoip/version"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	_PRINT_CONFIG = "printConfig"
)

/**
Initialize Config via Viper.
Read in Config File
Set Defaults
*/
func InitConfig(configFile string) (err error) {
	//Get Current Working Directory so we know where we are
	dir, _ := os.Getwd()
	log.Info("Current Working Directory is: ", dir)

	//Set the config File
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")                          // Local dir
		viper.AddConfigPath("/etc/" + version.APP_NAME)   // Looks in etc/appname
		viper.AddConfigPath("$HOME/." + version.APP_NAME) // Looks in HOME
	}

	//Set ENV Prefix
	viper.SetEnvPrefix(version.APP_NAME)
	viper.AutomaticEnv() //Reads in Environ Vars that match
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; Use defaults
			log.Info("No Config File Found.  Continuing with defaults.")
			err = nil
		} else if e, ok := err.(viper.ConfigParseError); ok {
			//Config Parsing Error
			log.Errorf("error parsing config file: %v", e)
			return err
		} else if e, ok := err.(viper.UnsupportedConfigError); ok {
			//Config Parsing Error
			log.Errorf("upsupported config type: %v", e)
			return err
		} else {
			// Other Error
			log.Errorf("fatal config error: %v", err)
			return err
		}
	} else {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	// Set Defaults
	setDefaults()
	if viper.GetBool(_PRINT_CONFIG) {
		PrintConfig()
	}

	//Watch config file for changes and apply
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", e.Name)
	})

	return err
}

func setDefaults() {
	viper.SetDefault(_PRINT_CONFIG, false)
}

func PrintConfig() {
	// Marshal the map into a JSON string.
	config := viper.AllSettings()
	cJSON, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(cJSON))
}

func String() string {
	return `printConfig = true
`
}

func Generate() {
	builder := strings.Builder{}
	builder.WriteString(String())
	builder.WriteString(logging.String())
	builder.WriteString(httplogging.String())
	builder.WriteString(api.String())
	builder.WriteString(geo.String())

	// Convert Builder to String and print it.
	result := builder.String()
	fmt.Println(result)
}
