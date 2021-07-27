package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func InitLogger() {
	config := InitLoggingConfig()
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
			repopath := fmt.Sprintf("%s/src/github.com/bob", os.Getenv("GOPATH"))
			filename := strings.Replace(f.File, repopath, "", -1)
			r, _ := regexp.Compile(`[^\/]+\/[^\/]+$`)
			daFunc := strings.Split(f.Function, ".")
			return "", fmt.Sprintf("%s:%d[%s()]", r.Find([]byte(filename)), f.Line, daFunc[len(daFunc)-1])
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
