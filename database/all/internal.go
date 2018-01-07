package all

import (
	"sync"
	"time"

	"github.com/FreifunkBremen/yanic/database"
)

var conn database.Connection
var wg = sync.WaitGroup{}
var quit chan struct{}

func Start(config database.Config) (err error) {
	conn, err = Connect(config.Connection)
	if err != nil {
		return
	}
	quit = make(chan struct{})
	wg.Add(1)
	go deleteWorker(config.DeleteInterval.Duration, config.DeleteAfter.Duration)
	return
}

func Close() {
	close(quit)
	wg.Wait()
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
			wg.Done()
			return
		}
	}
}
