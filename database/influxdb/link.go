package influxdb

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	models "github.com/influxdata/influxdb/models"
)

// InsertLink adds a link data point
func (conn *Connection) InsertLink(link *runtime.Link, t time.Time) {
	tags := models.Tags{}
	tags.SetString("source", link.SourceID)
	tags.SetString("target", link.TargetID)

	conn.addPoint(MeasurementLink, tags, models.Fields{"tq": float32(link.TQ) / 2.55}, t)
}
