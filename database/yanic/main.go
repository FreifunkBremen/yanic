package yanic

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	config Config
	conn   *net.UDPConn
}

type Config map[string]interface{}

func (c Config) Enable() bool {
	return c["enable"].(bool)
}
func (c Config) Address() string {
	return c["address"].(string)
}

func init() {
	database.RegisterAdapter("yanic", Connect)
}

func Connect(configuration interface{}) (database.Connection, error) {
	var config Config
	config = configuration.(map[string]interface{})
	if !config.Enable() {
		return nil, nil
	}
	udpAddr, err := net.ResolveUDPAddr("udp", config.Address())
	if err != nil {
		log.Panicf("Invalid yanic address: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Panicf("Unable to dial yanic: %s", err)
	}
	return &Connection{config: config, conn: conn}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	buf := bytes.Buffer{}
	writer := bufio.NewWriter(&buf)
	deflater, err := flate.NewWriter(writer, flate.DefaultCompression)

	err = json.NewEncoder(deflater).Encode(&data.ResponseData{
		Statistics: node.Statistics,
		NodeInfo:   node.Nodeinfo,
	})
	if err != nil {
		panic(err)
	}
	deflater.Close()
	writer.Flush()

	conn.conn.Write(buf.Bytes())
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {

}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
}

func (conn *Connection) Close() {
	conn.conn.Close()
}
