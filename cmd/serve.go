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

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//cobra.OnInitialize(config.InitConfig, logging.InitLogger)
		err := config.InitConfig()
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

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			sig := <-ch
			log.Warn("Signal caught. Shutting down... Reason: ", sig)
			cancel()

		}()

		var wg sync.WaitGroup

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
}
