package cmd

import (
	"log"
	"net"
	"time"

	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/spf13/cobra"
)

var wait int

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:     "query <interface> <destination>",
	Short:   "Sends a query on the interface to the destination and waits for a response",
	Example: `yanic query wlan0 "fe80::eade:27ff:dead:beef"`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		iface := args[0]
		dstAddress := net.ParseIP(args[1])

		log.Printf("Sending request address=%s iface=%s", dstAddress, iface)

		nodes := runtime.NewNodes(&runtime.NodesConfig{})

		sitesDomains := make(map[string][]string)

		collector := respond.NewCollector(nil, nodes, sitesDomains, []string{iface}, 0)
		defer collector.Close()
		collector.SendPacket(dstAddress)

		time.Sleep(time.Second * time.Duration(wait))

		for id, data := range nodes.List {
			log.Printf("%s: %+v", id, data)
		}
	},
}

func init() {
	RootCmd.AddCommand(queryCmd)
	queryCmd.Flags().IntVar(&wait, "wait", 1, "Seconds to wait for a response")
}
