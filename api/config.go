package api

import (
	"github.com/dp1140a/geoip/api/httplogging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_HTTP = "http"
)

/**
[http]
    port = "8081"
    host = "localhost"
    useHttps = false
    tlsMinVersion = "1.2"
    httpTLSStrictCiphers = false
    tlsCert = "config/d5-test.crt"
    tlsKey = "config/d5-test.key"
    enableCORS = true
    jwtSecret="There is a mouse in my house"
*/
type Config struct {
	Hostname             string
	Port                 string
	UseHttps             bool
	TLSMinVersion        string
	HttpTLSStrictCiphers bool
	TLSCert              string
	TLSKey               string
	EnableCORS           bool
	JWTSecret            string
	LoggingConfig        *httplogging.HttpLoggingConfig //api/logging/config
}

func InitConfig() (config *Config, err error) {
	config = &Config{
		Hostname:             "0.0.0.0",
		Port:                 "8081",
		UseHttps:             false,
		TLSMinVersion:        "1.2",
		HttpTLSStrictCiphers: false,
		TLSCert:              "config/d5-test.crt",
		TLSKey:               "config/d5-test.key",
		EnableCORS:           false,
		JWTSecret:            "There is a mouse in my house",
		LoggingConfig:        httplogging.InitConfig(),
	}
	h := viper.Sub(_HTTP)
	if h != nil {
		err := h.Unmarshal(config)
		if err != nil {
			log.Error("HTTP Config Error: ", err.Error())
			return nil, err
		}
	}
	return config, nil
}

func String() string {
	return `
[http]
    port = "8081"
    host = "localhost"
    useHttps = false
    tlsMinVersion = "1.2"
    httpTLSStrictCiphers = false
    tlsCert = "config/d5-test.crt"
    tlsKey = "config/d5-test.key"
    enableCORS = true
    jwtSecret="There is a mouse in my house"
`
}
