package api

import (
	"fmt"
	"github.com/dp1140a/geoip/api/httplogging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_HTTP         = "http"
	SSL_DIR       = "/etc/ssl"
	JWT_SECRET    = "There is a mouse in my house"
	HOST          = "localhost"
	PORT          = "8081"
	USEHTTPS      = false
	TLSMINVERSION = "1.2"
	STRICTCIPHERS = false
	ENABLECORS    = false
)

type ServerConfig struct {
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

func InitServerConfig() (config *ServerConfig, err error) {
	config = &ServerConfig{
		Hostname:             HOST,
		Port:                 PORT,
		UseHttps:             false,
		TLSMinVersion:        TLSMINVERSION,
		HttpTLSStrictCiphers: STRICTCIPHERS,
		TLSCert:              fmt.Sprintf("%v/geoip.crt", SSL_DIR),
		TLSKey:               fmt.Sprintf("%v/geoip.key", SSL_DIR),
		EnableCORS:           ENABLECORS,
		JWTSecret:            JWT_SECRET,
		LoggingConfig:        httplogging.InitHttpLoggingConfig(),
	}
	h := viper.Sub(_HTTP)
	if h != nil {
		err := h.Unmarshal(config)
		if err != nil {
			log.Panic("Exiting! Server Config Error: ", err.Error())
			return nil, err
		}
	}
	return config, nil
}

func String() string {
	return `
[http]
    host = "` + HOST + `"
    port = "` + PORT + `"
    useHttps = ` + fmt.Sprintf("%t", USEHTTPS) + `
    tlsMinVersion = "` + TLSMINVERSION + `"
    httpTLSStrictCiphers = ` + fmt.Sprintf("%t", STRICTCIPHERS) + `
    tlsCert = "` + fmt.Sprintf("%v/geoip.crt", SSL_DIR) + `"
    tlsKey = "` + fmt.Sprintf("%v/geoip.key", SSL_DIR) + `"
    enableCORS = ` + fmt.Sprintf("%t", ENABLECORS) + `
    jwtSecret = "` + JWT_SECRET + `"
`
}
