package output

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

var quit chan struct{}

// Start workers of database
// WARNING: Do not override this function
//  you should use New()
func Start(output Output, config *runtime.Config) {
	quit = make(chan struct{})
	go saveWorker(output, config.Nodes.SaveInterval.Duration)
}

func Close() {
	if quit != nil {
		close(quit)
	}
}

// save periodically to output
func saveWorker(output Output, saveInterval time.Duration) {
	ticker := time.NewTicker(saveInterval)
	for {
		select {
		case <-ticker.C:
			output.Save()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
