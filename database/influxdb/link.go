package influxdb

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	models "github.com/influxdata/influxdb/models"
)

// InsertLink adds a link data point
func (conn *Connection) InsertLink(link *runtime.Link, t time.Time) {
	tags := models.Tags{}
	tags.SetString("source.id", link.SourceID)
	tags.SetString("source.mac", link.SourceMAC)
	tags.SetString("target.id", link.TargetID)
	tags.SetString("target.mac", link.TargetMAC)

	conn.addPoint(MeasurementLink, tags, models.Fields{"tq": float32(link.TQ) / 2.55}, t)
}
