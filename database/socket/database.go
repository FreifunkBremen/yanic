package socket

/**
 * This database type is just for,
 * - debugging without a influxconn
 * - example for other developers for new databases
 */
import (
	"log"
	"net"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	listener net.Listener
	clients  map[net.Addr]net.Conn
}

func init() {
	database.RegisterAdapter("socket", Connect)
}

func Connect(configuration interface{}) (database.Connection, error) {
	config := configuration.(map[string]interface{})

	if !config["enable"].(bool) {
		return nil, nil
	}

	ln, err := net.Listen(config["type"].(string), config["address"].(string))
	if err != nil {
		return nil, err
	}
	conn := &Connection{listener: ln, clients: make(map[net.Addr]net.Conn)}
	go conn.handleSocketConnection(ln)

	log.Println("[socket-database] listen on: ", ln.Addr())

	return conn, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	conn.sendJSON(EventMessage{Event: "insert_node", Body: node})
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {
	conn.sendJSON(EventMessage{Event: "insert_globals", Body: stats})
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	conn.sendJSON(EventMessage{Event: "prune_nodes"})
}

func (conn *Connection) Close() {
	for _, c := range conn.clients {
		c.Close()
	}
	conn.listener.Close()
}
