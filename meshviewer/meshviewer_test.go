package meshviewer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
)

func TestNewMeshviewer(t *testing.T) {
	assert := assert.New(t)

	stats := NewStatistics(&data.Statistics{
		Clients: &data.Clients{Total: 32},
		Memory: &data.Memory{
			Total: 2,
			Free:  1,
		},
	})
	assert.Equal(0.5, stats.MemoryUsage, "Wrong calculated memory")
	assert.Equal(uint32(32), stats.Clients, "Wrong client count with given total")

	stats = NewStatistics(&data.Statistics{
		Clients: &data.Clients{
			Wifi24: 3,
			Wifi5:  4,
		},
		Memory: &data.Memory{
			Total: 0,
			Free:  1,
		},
	})
	assert.Equal(1.0, stats.MemoryUsage, "Wrong calculated memory during divide by zero")
	assert.Equal(uint32(7), stats.Clients, "Wrong client count without total and wifi from batman")
}
