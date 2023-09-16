package influxdb

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"

	"github.com/bdlm/log"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

// InsertGlobals implementation of database
func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string, domain string) {
	measurementGlobal := MeasurementGlobal
	counterMeasurementModel := CounterMeasurementModel
	counterMeasurementFirmware := CounterMeasurementFirmware
	counterMeasurementAutoupdater := CounterMeasurementAutoupdater

	if site != runtime.GLOBAL_SITE {
		measurementGlobal += "_site"
		counterMeasurementModel += "_site"
		counterMeasurementFirmware += "_site"
		counterMeasurementAutoupdater += "_site"
	}
	if domain != runtime.GLOBAL_DOMAIN {
		measurementGlobal += "_domain"
		counterMeasurementModel += "_domain"
		counterMeasurementFirmware += "_domain"
		counterMeasurementAutoupdater += "_domain"
	}
	p := influxdb.NewPoint(measurementGlobal,
		conn.config.Tags(),
		map[string]interface{}{
			"nodes":          stats.Nodes,
			"gateways":       stats.Gateways,
			"clients.total":  stats.Clients,
			"clients.wifi":   stats.ClientsWifi,
			"clients.wifi24": stats.ClientsWifi24,
			"clients.wifi5":  stats.ClientsWifi5,
			"clients.owe":    stats.ClientsOWE,
			"clients.owe24":  stats.ClientsOWE24,
			"clients.owe5":   stats.ClientsOWE5,
		},
		time)

	if site != runtime.GLOBAL_SITE {
		p = p.AddTag("site", site)
	}
	if domain != runtime.GLOBAL_DOMAIN {
		p = p.AddTag("domain", domain)
	}
	conn.writeAPI[MeasurementGlobal].WritePoint(p)

	conn.addCounterMap(CounterMeasurementModel, counterMeasurementModel, stats.Models, time, site, domain)
	conn.addCounterMap(CounterMeasurementFirmware, counterMeasurementFirmware, stats.Firmwares, time, site, domain)
	conn.addCounterMap(CounterMeasurementAutoupdater, counterMeasurementAutoupdater, stats.Autoupdater, time, site, domain)
}

// Saves the values of a CounterMap in the database.
// The key are used as 'value' tag.
// The value is used as 'counter' field.
func (conn *Connection) addCounterMap(writeName, name string, m runtime.CounterMap, t time.Time, site string, domain string) {
	writeAPI, ok := conn.writeAPI[writeName]
	if !ok {
		log.WithField("writeName", writeName).Panic("no writeAPI found for countermap")
	}
	for key, count := range m {
		p := influxdb.NewPoint("stat",
			conn.config.Tags(),
			map[string]interface{}{
				"count": count,
			},
			t).
			AddTag("site", site).
			AddTag("domain", domain).
			AddTag("value", key)
		writeAPI.WritePoint(p)
	}
}
