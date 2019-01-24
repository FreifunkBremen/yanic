package influxdb

import (
	"sync"
	"testing"
	"time"

	"github.com/influxdata/influxdb1-client/v2"
	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

const (
	TEST_SITE   = "ffhb"
	TEST_DOMAIN = "city"
)

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes(), map[string][]string{TEST_SITE: {TEST_DOMAIN}})

	assert := assert.New(t)

	// check SITE_GLOBAL fields
	fields := GlobalStatsFields(stats[runtime.GLOBAL_SITE][runtime.GLOBAL_DOMAIN])
	assert.EqualValues(3, fields["nodes"])

	fields = GlobalStatsFields(stats[TEST_SITE][runtime.GLOBAL_DOMAIN])
	assert.EqualValues(2, fields["nodes"])
	fields = GlobalStatsFields(stats[TEST_SITE][TEST_DOMAIN])

	assert.EqualValues(1, fields["nodes"])

	conn := &Connection{
		points: make(chan *client.Point),
	}

	global := 0
	globalSite := 0
	globalDomain := 0

	model := 0
	modelSite := 0
	modelDomain := 0

	firmware := 0
	firmwareSite := 0
	firmwareDomain := 0

	autoupdater := 0
	autoupdaterSite := 0
	autoupdaterDomain := 0

	wg := sync.WaitGroup{}
	wg.Add(15)
	go func() {
		for p := range conn.points {
			switch p.Name() {
			case MeasurementGlobal:
				global++
			case "global_site":
				globalSite++
			case "global_site_domain":
				globalDomain++

			case CounterMeasurementModel:
				model++
			case "model_site":
				modelSite++
			case "model_site_domain":
				modelDomain++

			case CounterMeasurementFirmware:
				firmware++
			case "firmware_site":
				firmwareSite++
			case "firmware_site_domain":
				firmwareDomain++

			case CounterMeasurementAutoupdater:
				autoupdater++
			case "autoupdater_site":
				autoupdaterSite++
			case "autoupdater_site_domain":
				autoupdaterDomain++

			default:
				assert.Equal("invalid p.Name found", p.Name())
			}
			wg.Done()
		}
	}()
	for site, domains := range stats {
		for domain, stat := range domains {
			conn.InsertGlobals(stat, time.Now(), site, domain)
		}
	}
	wg.Wait()
	assert.Equal(1, global)
	assert.Equal(1, globalSite)
	assert.Equal(1, globalDomain)

	assert.Equal(2, model)
	assert.Equal(2, modelSite)
	assert.Equal(1, modelDomain)

	assert.Equal(1, firmware)
	assert.Equal(1, firmwareSite)
	assert.Equal(0, firmwareDomain)

	assert.Equal(2, autoupdater)
	assert.Equal(2, autoupdaterSite)
	assert.Equal(1, autoupdaterDomain)
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.NodesConfig{})

	nodeData := &runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.Nodeinfo{
			NodeID: "abcdef012345",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
			System: data.System{
				SiteCode: TEST_SITE,
			},
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
		Nodeinfo: &data.Nodeinfo{
			NodeID: "112233445566",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "0xdeadbeef0x",
			VPN:    true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
			System: data.System{
				SiteCode:   TEST_SITE,
				DomainCode: TEST_DOMAIN,
			},
		},
	})

	return nodes
}
