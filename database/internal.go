package database

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

var quit chan struct{}

// Start workers of database
// WARNING: Do not override this function
//  you should use New()
func Start(conn Connection, config *runtime.Config) {
	quit = make(chan struct{})
	go deleteWorker(conn, config.Database.DeleteInterval.Duration, config.Database.DeleteAfter.Duration)
}

func Close(conn Connection) {
	if quit != nil {
		close(quit)
	}
	if conn != nil {
		conn.Close()
	}
}

// prunes node-specific data periodically
func deleteWorker(conn Connection, deleteInterval time.Duration, deleteAfter time.Duration) {
	ticker := time.NewTicker(deleteInterval)
	for {
		select {
		case <-ticker.C:
			conn.DeleteNode(deleteAfter)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
