package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
		RootCmd.AddCommand(versionCMD)
	}
}
