package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bdlm/log"
	"github.com/spf13/cobra"

	"github.com/FreifunkBremen/yanic/respond/daemon"
)

// serveCmd represents the serve command
var responddCMD = &cobra.Command{
	Use:     "respondd",
	Short:   "Runs a respond daemon",
	Example: "yanic respondd --config /etc/respondd.toml",
	Run: func(cmd *cobra.Command, args []string) {
		daemon := &respondd.Daemon{}
		if err := ReadConfigFile(configPath, daemon); err != nil {
			log.Panicf("unable to load config file: %s", err)
		}

		go daemon.Start()

		log.Info("respondd daemon started")
		// Wait for INT/TERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Infof("received %s", sig)

	},
}

func init() {
	RootCmd.AddCommand(responddCMD)
	responddCMD.Flags().StringVarP(&configPath, "config", "c", "config-respondd.toml", "Path to configuration file")
}
