package api

import (
	"fmt"
	"github.com/dp1140a/geoip/api/httplogging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_HTTP   = "http"
	SSL_DIR = "/etc/ssl"
)

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
		TLSCert:              fmt.Sprintf("%v/geoip.crt", SSL_DIR),
		TLSKey:               fmt.Sprintf("%v/geoip.key", SSL_DIR),
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
    tlsCert = "/etc/ssl/geoip.crt"
    tlsKey = "/etc/ssl/geoip.key"
    enableCORS = true
    jwtSecret="There is a mouse in my house"
`
}
