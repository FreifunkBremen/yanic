package influxdb

import (
	"sync"
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

const TEST_SITE = "ffxx"

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes(), []string{TEST_SITE})

	assert := assert.New(t)

	// check SITE_GLOBAL fields
	fields := GlobalStatsFields(stats[runtime.GLOBAL_SITE])
	assert.EqualValues(3, fields["nodes"])

	fields = GlobalStatsFields(stats[TEST_SITE])
	assert.EqualValues(1, fields["nodes"])

	conn := &Connection{
		points: make(chan *client.Point),
	}

	global := 0
	globalSite := 0
	model := 0
	modelSite := 0
	firmware := 0
	firmwareSite := 0
	autoupdater := 0
	autoupdaterSite := 0
	wg := sync.WaitGroup{}
	wg.Add(9)
	go func() {
		for p := range conn.points {
			switch p.Name() {
			case MeasurementGlobal:
				global++
				break
			case "global_site":
				globalSite++
				break
			case CounterMeasurementModel:
				model++
				break
			case "model_site":
				modelSite++
				break
			case CounterMeasurementFirmware:
				firmware++
				break
			case "firmware_site":
				firmwareSite++
				break
			case CounterMeasurementAutoupdater:
				autoupdater++
				break
			case "autoupdater_site":
				autoupdaterSite++
				break
			default:
				assert.Equal("invalid p.Name found", p.Name())
			}
			wg.Done()
		}
	}()
	for site, stat := range stats {
		conn.InsertGlobals(stat, time.Now(), site)
	}
	wg.Wait()
	assert.Equal(1, global)
	assert.Equal(1, globalSite)
	assert.Equal(2, model)
	assert.Equal(1, modelSite)
	assert.Equal(1, firmware)
	assert.Equal(0, firmwareSite)
	assert.Equal(2, autoupdater)
	assert.Equal(1, autoupdaterSite)
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.Config{})

	nodeData := &runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.NodeInfo{
			NodeID: "abcdef012345",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
			System: data.System{},
		},
	}
	nodeData.Nodeinfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodeData.Nodeinfo.Software.Autoupdater.Enabled = true
	nodeData.Nodeinfo.Software.Autoupdater.Branch = "stable"
	nodes.AddNode(nodeData)

	nodes.AddNode(&runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 2,
			},
		},
		Nodeinfo: &data.NodeInfo{
			NodeID: "112233445566",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.NodeInfo{
			NodeID: "0xdeadbeef0x",
			VPN:    true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
			System: data.System{
				SiteCode: TEST_SITE,
			},
		},
	})

	return nodes
}
