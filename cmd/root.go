package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	timestamps bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "yanic",
	Short: "Yet another node info collector",
	Long:  `A respondd client that fetches, stores and publishes information about a Freifunk network.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().BoolVar(&timestamps, "timestamps", false, "Enables timestamps for log output")
}

func initConfig() {
	if timestamps {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}
