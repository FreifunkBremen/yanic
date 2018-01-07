package runtime

import "github.com/FreifunkBremen/yanic/lib/duration"

type NodesConfig struct {
	StatePath    string            `toml:"state_path"`
	SaveInterval duration.Duration `toml:"save_interval"` // Save nodes periodically
	OfflineAfter duration.Duration `toml:"offline_after"` // Set node to offline if not seen within this period
	PruneAfter   duration.Duration `toml:"prune_after"`   // Remove nodes after n days of inactivity
	Output       map[string]interface{}
}
