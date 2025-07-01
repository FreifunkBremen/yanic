package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bdlm/log"
	"github.com/spf13/cobra"

	allDatabase "github.com/FreifunkBremen/yanic/database/all"
	allOutput "github.com/FreifunkBremen/yanic/output/all"
	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/FreifunkBremen/yanic/webserver"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Runs the yanic server",
	Example: "yanic serve --config /etc/yanic.toml",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()

		err := allDatabase.Start(config.Database)
		if err != nil {
			log.WithError(err).Panic("could not connect to database")
		}
		defer allDatabase.Close()

		nodes = runtime.NewNodes(&config.Nodes)
		nodes.Start()

		err = allOutput.Start(nodes, config.Nodes)
		if err != nil {
			log.WithError(err).Panicf("error on init outputs")
		}
		defer allOutput.Close()

		if config.Webserver.Enable {
			log.WithField("address", config.Webserver.Bind).Info("starting webserver")
			srv := webserver.New(config.Webserver, nodes)
			go webserver.Start(srv)
			defer func() {
				if err := srv.Close(); err != nil {
					log.WithError(err).Panic("could not stop webserver")
				}
			}()
		} else if prom := config.Webserver.Prometheus; prom != nil && prom.Enable {
			log.Error("to enable prometheus exporter, please enable webserver ")
		}

		if config.Respondd.Enable {
			// Delaying startup to start at a multiple of `duration` since the zero time.
			if duration := config.Respondd.Synchronize.Duration; duration > 0 {
				now := time.Now()
				delay := duration - now.Sub(now.Truncate(duration))
				log.Infof("delaying %0.1f seconds", delay.Seconds())
				time.Sleep(delay)
			}

			collector = respond.NewCollector(allDatabase.Conn, nodes, &config.Respondd)
			collector.Start(config.Respondd.CollectInterval.Duration)
			defer collector.Close()
		}

		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Infof("received %s", sig)

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "config.toml", "Path to configuration file")
}
