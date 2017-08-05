package graphite

import (
	"regexp"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/fgrosse/graphigo"
)

func (c *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {
	c.addPoint(GlobalStatsFields(stats))
	c.addCounterMap(CounterMeasurementModel, stats.Models, time)
	c.addCounterMap(CounterMeasurementFirmware, stats.Firmwares, time)
}

func GlobalStatsFields(stats *runtime.GlobalStats) []graphigo.Metric {
	return []graphigo.Metric{
		{Name: MeasurementGlobal + ".nodes", Value: stats.Nodes},
		{Name: MeasurementGlobal + ".gateways", Value: stats.Gateways},
		{Name: MeasurementGlobal + ".clients.total", Value: stats.Clients},
		{Name: MeasurementGlobal + ".clients.wifi", Value: stats.ClientsWifi},
		{Name: MeasurementGlobal + ".clients.wifi24", Value: stats.ClientsWifi24},
		{Name: MeasurementGlobal + ".clients.wifi5", Value: stats.ClientsWifi5},
	}
}

func (c *Connection) addCounterMap(name string, m runtime.CounterMap, t time.Time) {
	var fields []graphigo.Metric
	re := regexp.MustCompile("(?i)[^a-z0-9\\-]")
	for key, count := range m {
		fields = append(fields, graphigo.Metric{Name: name + `.` + re.ReplaceAllString(key, "_") + `.count`, Value: count, Timestamp: t})
	}
	c.addPoint(fields)
}
