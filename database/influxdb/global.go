package influxdb

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/influxdata/influxdb1-client/models"
)

// InsertGlobals implementation of database
func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string, domain string) {
	tags := models.Tags{}

	measurementGlobal := MeasurementGlobal
	counterMeasurementModel := CounterMeasurementModel
	counterMeasurementFirmware := CounterMeasurementFirmware
	counterMeasurementAutoupdater := CounterMeasurementAutoupdater

	if site != runtime.GLOBAL_SITE {
		tags.Set([]byte("site"), []byte(site))

		measurementGlobal += "_site"
		counterMeasurementModel += "_site"
		counterMeasurementFirmware += "_site"
		counterMeasurementAutoupdater += "_site"
	}
	if domain != runtime.GLOBAL_DOMAIN {
		tags.Set([]byte("domain"), []byte(domain))

		measurementGlobal += "_domain"
		counterMeasurementModel += "_domain"
		counterMeasurementFirmware += "_domain"
		counterMeasurementAutoupdater += "_domain"
	}

	conn.addPoint(measurementGlobal, tags, GlobalStatsFields(stats), time)
	conn.addCounterMap(counterMeasurementModel, stats.Models, time, site, domain)
	conn.addCounterMap(counterMeasurementFirmware, stats.Firmwares, time, site, domain)
	conn.addCounterMap(counterMeasurementAutoupdater, stats.Autoupdater, time, site, domain)
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
func (conn *Connection) addCounterMap(name string, m runtime.CounterMap, t time.Time, site string, domain string) {
	for key, count := range m {
		conn.addPoint(
			name,
			models.Tags{
				models.Tag{Key: []byte("value"), Value: []byte(key)},
				models.Tag{Key: []byte("site"), Value: []byte(site)},
				models.Tag{Key: []byte("domain"), Value: []byte(domain)},
			},
			models.Fields{"count": count},
			t,
		)
	}
}
