package respondd

/**
 * This database type is for injecting into another yanic instance.
 */
import (
	"bufio"
	"compress/flate"
	"encoding/json"
	"net"
	"time"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	config Config
	conn   net.Conn
}

type Config map[string]interface{}

func (c Config) Type() string {
	return c["type"].(string)
}

func (c Config) Address() string {
	return c["address"].(string)
}

func init() {
	database.RegisterAdapter("respondd", Connect)
}

func Connect(configuration map[string]interface{}) (database.Connection, error) {
	var config Config
	config = configuration

	conn, err := net.Dial(config.Type(), config.Address())
	if err != nil {
		return nil, err
	}

	return &Connection{conn: conn, config: config}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	res := &data.ResponseData{
		NodeInfo:   node.Nodeinfo,
		Statistics: node.Statistics,
		Neighbours: node.Neighbours,
	}

	writer := bufio.NewWriterSize(conn.conn, 8192)

	flater, err := flate.NewWriter(writer, flate.BestCompression)
	if err != nil {
		log.Errorf("[database-yanic] could not create flater: %s", err)
		return
	}
	defer flater.Close()
	err = json.NewEncoder(flater).Encode(res)
	if err != nil {
		nodeid := "unknown"
		if node.Nodeinfo != nil && node.Nodeinfo.NodeID != "" {
			nodeid = node.Nodeinfo.NodeID
		}
		log.WithField("node_id", nodeid).Errorf("[database-yanic] could not encode node: %s", err)
		return
	}
	err = flater.Flush()
	if err != nil {
		log.Errorf("[database-yanic] could not compress: %s", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Errorf("[database-yanic] could not send: %s", err)
	}
}

func (conn *Connection) InsertLink(link *runtime.Link, time time.Time) {
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string, domain string) {
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
}

func (conn *Connection) Close() {
	conn.conn.Close()
}
