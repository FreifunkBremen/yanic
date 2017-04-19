package logging

/**
 * This database type is just for,
 * - debugging without a influxconn
 * - example for other developers for new databases
 */
import (

	"log"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
	"gopkg.in/olivere/elastic.v5"


	"context"

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
func (c Config) Path() string {
	return c["path"].(string)
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
	     elastic.SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"),
	     elastic.SetBasicAuth("user", "secret"))



	if err != nil {
		// Handle error
		panic(err)
	}


	return &Connection{config: config, client: client}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	log("InsertNode: [", node.Statistics.NodeID, "] clients: ", node.Statistics.Clients.Total)

	_, err = conn.client.Index()
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

	_, err = conn.client.Index()
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
