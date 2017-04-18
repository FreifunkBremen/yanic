package influxdb

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/influxdata/influxdb/models"
)

// InsertGlobals implementation of database
func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {
	conn.addPoint(MeasurementGlobal, nil, GlobalStatsFields(stats), time)
	conn.addCounterMap(CounterMeasurementModel, stats.Models, time)
	conn.addCounterMap(CounterMeasurementFirmware, stats.Firmwares, time)
}

// GlobalStatsFields returns fields for InfluxDB
func GlobalStatsFields(stats *runtime.GlobalStats) map[string]interface{} {
	return map[string]interface{}{
		"nodes":          stats.Nodes,
		"gateways":       stats.Gateways,
		"clients.total":  stats.Clients,
		"clients.wifi":   stats.ClientsWifi,
		"clients.wifi24": stats.ClientsWifi24,
		"clients.wifi5":  stats.ClientsWifi5,
	}
}

// Saves the values of a CounterMap in the database.
// The key are used as 'value' tag.
// The value is used as 'counter' field.
func (conn *Connection) addCounterMap(name string, m runtime.CounterMap, t time.Time) {
	for key, count := range m {
		conn.addPoint(
			name,
			models.Tags{
				models.Tag{Key: []byte("value"), Value: []byte(key)},
			},
			models.Fields{"count": count},
			t,
		)
	}
}
