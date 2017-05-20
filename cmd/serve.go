package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	allDatabase "github.com/FreifunkBremen/yanic/database/all"
	"github.com/FreifunkBremen/yanic/output"
	allOutput "github.com/FreifunkBremen/yanic/output/all"
	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/FreifunkBremen/yanic/webserver"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Runs the yanic server",
	Example: "yanic serve --config /etc/yanic.toml",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()

		connections, err := allDatabase.Connect(config.Database.Connection)
		if err != nil {
			panic(err)
		}
		database.Start(connections, config)
		defer database.Close(connections)

		nodes = runtime.NewNodes(config)
		nodes.Start()

		outputs, err := allOutput.Register(config.Nodes.Output)
		if err != nil {
			panic(err)
		}
		output.Start(outputs, nodes, config)
		defer output.Close()

		if config.Webserver.Enable {
			log.Println("starting webserver on", config.Webserver.Bind)
			srv := webserver.New(config.Webserver.Bind, config.Webserver.Webroot)
			go webserver.Start(srv)
			defer srv.Close()
		}

		if config.Respondd.Enable {
			// Delaying startup to start at a multiple of `duration` since the zero time.
			if duration := config.Respondd.Synchronize.Duration; duration > 0 {
				now := time.Now()
				delay := duration - now.Sub(now.Truncate(duration))
				log.Printf("delaying %0.1f seconds", delay.Seconds())
				time.Sleep(delay)
			}

			collector = respond.NewCollector(connections, nodes, config.Respondd.Interfaces, config.Respondd.Port)
			collector.Start(config.Respondd.CollectInterval.Duration)
			defer collector.Close()
		}

		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Println("received", sig)

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "config.toml", "Path to configuration file")
}
