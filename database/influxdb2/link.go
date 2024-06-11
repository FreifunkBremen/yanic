package influxdb

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

// InsertLink adds a link data point
func (conn *Connection) InsertLink(link *runtime.Link, t time.Time) {
	p := influxdb.NewPoint(MeasurementLink,
		conn.config.Tags(),
		map[string]interface{}{
			"tq": link.TQ * 100,
		},
		t).
		AddTag("source.id", link.SourceID).
		AddTag("source.addr", link.SourceAddress).
		AddTag("target.id", link.TargetID).
		AddTag("target.addr", link.TargetAddress)
	if link.SourceHostname != "" {
		p.AddTag("source.hostname", link.SourceHostname)
	}
	if link.TargetHostname != "" {
		p.AddTag("target.hostname", link.TargetHostname)
	}
	conn.writeAPI[MeasurementLink].WritePoint(p)
}
