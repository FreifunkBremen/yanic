package cmd

import (
	"log"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/database/all"
	"github.com/FreifunkBremen/yanic/rrd"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import <file.rrd> <site>",
	Short:   "Imports global statistics from the given RRD files, requires InfluxDB",
	Example: "yanic import --config /etc/yanic.toml olddata.rrd global",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		site := args[1]
		config := loadConfig()

		connections, err := all.Connect(config.Database.Connection)
		if err != nil {
			panic(err)
		}
		database.Start(connections, config)
		defer database.Close(connections)

		log.Println("importing RRD from", path)

		for ds := range rrd.Read(path) {
			connections.InsertGlobals(
				&runtime.GlobalStats{
					Nodes:   uint32(ds.Nodes),
					Clients: uint32(ds.Clients),
				},
				ds.Time,
				site,
			)
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&configPath, "config", "c", "config.toml", "Path to configuration file")
}
