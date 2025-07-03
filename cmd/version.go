package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/FreifunkBremen/yanic/webserver"
)

// VERSION is set at build time
var VERSION string

// versionCMD to print version
var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "print version of yanic",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("yanic version: %s\n", VERSION)
	},
}

func init() {
	if VERSION != "" {
		webserver.VERSION = VERSION
		RootCmd.AddCommand(versionCMD)
	}
}
