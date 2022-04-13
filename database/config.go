package database

import "yanic/lib/duration"

type Config struct {
	DeleteInterval duration.Duration `toml:"delete_interval"` // Delete stats of nodes every n minutes
	DeleteAfter    duration.Duration `toml:"delete_after"`    // Delete stats of nodes till now-deletetill n minutes
	Connection     map[string]interface{}
}
