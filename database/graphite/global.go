package graphite

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/fgrosse/graphigo"
)

func (c *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string, domain string) {
	measurementGlobal := MeasurementGlobal
	counterMeasurementModel := CounterMeasurementModel
	counterMeasurementFirmware := CounterMeasurementFirmware
	counterMeasurementAutoupdater := CounterMeasurementAutoupdater

	if site != runtime.GLOBAL_SITE {
		measurementGlobal += "_" + site
		counterMeasurementModel += "_" + site
		counterMeasurementFirmware += "_" + site
		counterMeasurementAutoupdater += "_" + site
	}

	if domain != runtime.GLOBAL_DOMAIN {
		measurementGlobal += "_" + domain
		counterMeasurementModel += "_" + domain
		counterMeasurementFirmware += "_" + domain
		counterMeasurementAutoupdater += "_" + domain
	}

	c.addPoint(GlobalStatsFields(measurementGlobal, stats))
	c.addCounterMap(counterMeasurementModel, stats.Models, time)
	c.addCounterMap(counterMeasurementFirmware, stats.Firmwares, time)
	c.addCounterMap(counterMeasurementAutoupdater, stats.Autoupdater, time)
}

func GlobalStatsFields(name string, stats *runtime.GlobalStats) []graphigo.Metric {
	return []graphigo.Metric{
		{Name: name + ".nodes", Value: stats.Nodes},
		{Name: name + ".gateways", Value: stats.Gateways},
		{Name: name + ".clients.total", Value: stats.Clients},
		{Name: name + ".clients.wifi", Value: stats.ClientsWifi},
		{Name: name + ".clients.wifi24", Value: stats.ClientsWifi24},
		{Name: name + ".clients.wifi5", Value: stats.ClientsWifi5},
	}
}

func (c *Connection) addCounterMap(name string, m runtime.CounterMap, t time.Time) {
	var fields []graphigo.Metric
	for key, count := range m {
		fields = append(fields, graphigo.Metric{Name: name + `.` + replaceInvalidChars(key) + `.count`, Value: count, Timestamp: t})
	}
	c.addPoint(fields)
}
