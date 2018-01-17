package respond

import "github.com/FreifunkBremen/yanic/lib/duration"

type Config struct {
	Enable          bool                  `toml:"enable"`
	Synchronize     duration.Duration     `toml:"synchronize"`
	Interfaces      []string              `toml:"interfaces"`
	Sites           map[string]SiteConfig `toml:"sites"`
	Port            int                   `toml:"port"`
	CollectInterval duration.Duration     `toml:"collect_interval"`
}

func (c *Config) SitesDomains() (result map[string][]string) {
	result = make(map[string][]string)
	for site, siteConfig := range c.Sites {
		result[site] = siteConfig.Domains
	}
	return
}

type SiteConfig struct {
	Domains []string `toml:"domains"`
}
