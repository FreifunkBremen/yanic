package cmd

import (
	"github.com/bdlm/log"
	"github.com/spf13/cobra"

	allDatabase "github.com/FreifunkBremen/yanic/database/all"
	"github.com/FreifunkBremen/yanic/rrd"
	"github.com/FreifunkBremen/yanic/runtime"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import <file.rrd> <site> <domain>",
	Short:   "Imports global statistics from the given RRD files, requires InfluxDB",
	Example: "yanic import --config /etc/yanic.toml olddata.rrd global global",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		site := args[1]
		domain := args[2]
		config := loadConfig()

		err := allDatabase.Start(config.Database)
		if err != nil {
			log.Panicf("could not connect to database: %s", err)
		}
		defer allDatabase.Close()

		log.Infof("importing RRD from %s", path)

		for ds := range rrd.Read(path) {
			allDatabase.Conn.InsertGlobals(
				&runtime.GlobalStats{
					Nodes:   uint32(ds.Nodes),
					Clients: uint32(ds.Clients),
				},
				ds.Time,
				site,
				domain,
			)
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&configPath, "config", "c", "config.toml", "Path to configuration file")
}
