package logging

/**
 * This database type is just for,
 * - debugging without a influxconn
 * - example for other developers for new databases
 */
import (
	"fmt"
	"log"
	"os"
	"time"

	logger "github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	config Config
	file   *os.File
}

type Config map[string]interface{}

func (c Config) Path() string {
	return c["path"].(string)
}

func init() {
	database.RegisterAdapter("logging", Connect)
}

func Connect(configuration map[string]interface{}) (database.Connection, error) {
	config := Config(configuration)

	file, err := os.OpenFile(config.Path(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return &Connection{config: config, file: file}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	conn.log("InsertNode: [", node.Statistics.NodeID, "] clients: ", node.Statistics.Clients.Total)
}

func (conn *Connection) InsertLink(link *runtime.Link, time time.Time) {
	conn.log("InsertLink: ", link)
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string, domain string) {
	conn.log("InsertGlobals: [", time.String(), "] site: ", site, " domain: ", domain, ", nodes: ", stats.Nodes, ", clients: ", stats.Clients, " models: ", len(stats.Models))
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	conn.log("PruneNodes")
}

func (conn *Connection) Close() {
	conn.log("Close")
	if err := conn.file.Close(); err != nil {
		logger.WithError(err).Error("unable to close connection")
	}
}

func (conn *Connection) log(v ...interface{}) {
	_, err := fmt.Fprintf(conn.file, "[%s] %v", time.Now().String(), v)
	if err != nil {
		log.Println(err)
	}
}
