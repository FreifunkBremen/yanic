package all

import (
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	list []database.Connection
}

func Connect(configuration interface{}) (database.Connection, error) {
	var list []database.Connection
	allConnection := configuration.(map[string][]interface{})
	for dbType, conn := range database.Adapters {
		dbConfigs := allConnection[dbType]
		for _, config := range dbConfigs {
			connected, err := conn(config)
			if err != nil {
				return nil, err
			}
			if connected == nil {
				continue
			}
			list = append(list, connected)
		}
	}
	return &Connection{list: list}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	for _, item := range conn.list {
		item.InsertNode(node)
	}
}

func (conn *Connection) InsertLink(link *runtime.Link, time time.Time) {
	for _, item := range conn.list {
		item.InsertLink(link, time)
	}
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string) {
	for _, item := range conn.list {
		item.InsertGlobals(stats, time, site)
	}
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	for _, item := range conn.list {
		item.PruneNodes(deleteAfter)
	}
}

func (conn *Connection) Close() {
	for _, item := range conn.list {
		item.Close()
	}
}
