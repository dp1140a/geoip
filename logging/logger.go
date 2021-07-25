package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func InitLogger() {
	config := InitConfig()
	var logLevel = log.InfoLevel
	if config.TraceLogging == true {
		logLevel = log.TraceLevel
	}

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat:   time.RFC3339,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			//r, _ := regexp.Compile("[^.]+$")
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: false,
	})
	log.SetLevel(logLevel)
	log.SetReportCaller(true)
	log.Info("Logging started")
	dir, _ := filepath.Split(config.LogFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("Log Error: ", err, ": ", config.LogFile)
		}
	}
	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Log Error: ", err, " ", config.LogFile)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}
