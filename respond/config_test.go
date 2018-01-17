package respond

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSitesDomainsConfigTransform(t *testing.T) {
	assert := assert.New(t)
	c := Config{
		Sites: map[string]SiteConfig{
			"ffhb": {Domains: []string{"city"}},
		},
	}
	result := c.SitesDomains()
	assert.Len(result, 1)
	assert.Contains(result, "ffhb")

	domains := result["ffhb"]

	assert.Len(domains, 1)
	assert.Equal("city", domains[0])
}
