package runtime

import "github.com/FreifunkBremen/yanic/lib/duration"

type NodesConfig struct {
	StatePath    string            `toml:"state_path"`
	SaveInterval duration.Duration `toml:"save_interval"` // Save nodes periodically
	OfflineAfter duration.Duration `toml:"offline_after"` // Set node to offline if not seen within this period
	PruneAfter   duration.Duration `toml:"prune_after"`   // Remove nodes after n days of inactivity
	PingCount    int               `toml:"ping_count"`    // send x pings to verify if node is offline (for disable count < 1)
	PingTimeout  duration.Duration `toml:"ping_timeout"`  // timeout of sending ping to a node
	Output       map[string]interface{}
}
