package cmd

import (
	"context"
	"github.com/dp1140a/geoip/api"
	"github.com/dp1140a/geoip/config"
	"github.com/dp1140a/geoip/logging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
)

var ConfigLocation string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//cobra.OnInitialize(config.InitConfig, logging.InitLogger)
		if ConfigLocation != "" {
			log.Infof("Config flag set to %v", ConfigLocation)
		}
		err := config.InitConfig(ConfigLocation)
		if err != nil {
			log.Fatal("unable to configure. shutting down.")
		}
		logging.InitLogger()
		ctx, cancel := context.WithCancel(context.Background())
		s, err := api.NewServer(ctx)
		if err != nil {
			log.Error("Error starting server. ", err)
			cancel()
		}

		// Setup clean shutdown on ^C
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			sig := <-ch
			log.Warn("Signal caught. Shutting down... Reason: ", sig)
			cancel()
		}()

		var wg sync.WaitGroup

		//Go Run It
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			s.ServeAPI()
		}()

		wg.Wait()
		log.Println("Server gracefully stopped")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&ConfigLocation, "config", "c", "", "location of the config file to use")
}
