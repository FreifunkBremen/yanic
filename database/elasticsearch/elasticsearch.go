package logging

/**
 * This database type is just for,
 * - debugging without a influxconn
 * - example for other developers for new databases
 */
import (
	"context"
	"log"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	config Config
	client *elastic.Client
}

type Config map[string]interface{}

func (c Config) Enable() bool {
	return c["enable"].(bool)
}
func (c Config) Host() string {
	return c["host"].(string)
}
func (c Config) Username() string {
	return c["username"].(string)
}
func (c Config) Password() string {
	return c["password"].(string)
}
func (c Config) IndexPrefix() string {
	return c["index_prefix"].(string)
}
func (c Config) UpdateTemplates() bool {
	return c["update_templates"].(bool)
}

func init() {
	database.RegisterAdapter("elasticsearch", Connect)
}

func Connect(configuration interface{}) (database.Connection, error) {
	var config Config
	config = configuration.(map[string]interface{})
	if !config.Enable() {
		return nil, nil
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(config.Host()),
		elastic.SetBasicAuth(config.Username(), config.Password()))

	if err != nil {
		// Handle error
		panic(err)
	}

	return &Connection{config: config, client: client}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	log.Print("InsertNode: [", node.Statistics.NodeID, "] clients: ", node.Statistics.Clients.Total)

	_, err := conn.client.Index().
		Index("ffhb").
		Type("node").
		BodyJson(node).
		Refresh("true").
		Do(context.Background())

	if err != nil {
		// Handle error
		panic(err)
	}

}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {
	log.Print("InsertGlobals: [", time.String(), "] nodes: ", stats.Nodes, ", clients: ", stats.Clients, " models: ", len(stats.Models))

	_, err := conn.client.Index().
		Index("ffhb").
		Type("globals").
		BodyJson(stats).
		Refresh("true").
		Do(context.Background())

	if err != nil {
		// Handle error
		panic(err)
	}

}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	log.Print("PruneNodes")

	// TODO
}

func (conn *Connection) Close() {

	log.Print("Closing connection, stop client")
	conn.client.Stop()
}
