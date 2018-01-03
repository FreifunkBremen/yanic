package all

import (
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

var conn database.Connection
var quit chan struct{}

func Start(config runtime.DatabaseConfig) (err error) {
	conn, err = Connect(config.Connection)
	if err != nil {
		return
	}
	quit = make(chan struct{})
	go deleteWorker(config.DeleteInterval.Duration, config.DeleteAfter.Duration)
	return
}

func Close() {
	close(quit)
	conn.Close()
	quit = nil
}

// prunes node-specific data periodically
func deleteWorker(deleteInterval time.Duration, deleteAfter time.Duration) {
	ticker := time.NewTicker(deleteInterval)
	for {
		select {
		case <-ticker.C:
			conn.PruneNodes(deleteAfter)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
