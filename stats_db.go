package main

import (
	"log"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/FreifunkBremen/RespondCollector/data"
)

const (
	batchWaiting = time.Second * 5
)

type StatsDb struct {
	points chan *client.Point
	wg     sync.WaitGroup
	client client.Client
}

func NewStatsDb() *StatsDb {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: config.Influxdb.Username,
		Password: config.Influxdb.Password,
	})

	if err != nil {
		panic(err)
	}

	db := &StatsDb{
		client: c,
		points: make(chan *client.Point, 500),
	}

	// start worker
	db.wg.Add(1)
	go db.worker()

	return db
}

func (c *StatsDb) Add(stats *data.Statistics) {
	tags := map[string]string{
		"nodeid": stats.NodeId,
	}
	fields := map[string]interface{}{
		"load":              stats.LoadAverage,
		"idletime":          uint64(stats.Idletime),
		"uptime":            uint64(stats.Uptime),
		"processes.running": stats.Processes.Running,
		"clients.wifi":      stats.Clients.Wifi,
		"clients.wifi24":    stats.Clients.Wifi24,
		"clients.wifi5":     stats.Clients.Wifi5,
		"clients.total":     stats.Clients.Total,
		"memory.buffers":    stats.Memory.Buffers,
		"memory.cached":     stats.Memory.Cached,
		"memory.free":       stats.Memory.Free,
		"memory.total":      stats.Memory.Total,
	}

	if t := stats.Traffic.Rx; t != nil {
		fields["traffic.rx.bytes"] = uint64(t.Bytes)
		fields["traffic.rx.packets"] = t.Packets
	}
	if t := stats.Traffic.Tx; t != nil {
		fields["traffic.tx.bytes"] = uint64(t.Bytes)
		fields["traffic.tx.packets"] = t.Packets
		fields["traffic.tx.dropped"] = t.Dropped
	}
	if t := stats.Traffic.Forward; t != nil {
		fields["traffic.forward.bytes"] = uint64(t.Bytes)
		fields["traffic.forward.packets"] = t.Packets
	}
	if t := stats.Traffic.MgmtRx; t != nil {
		fields["traffic.mgmt_rx.bytes"] = uint64(t.Bytes)
		fields["traffic.mgmt_rx.packets"] = t.Packets
	}
	if t := stats.Traffic.MgmtTx; t != nil {
		fields["traffic.mgmt_tx.bytes"] = uint64(t.Bytes)
		fields["traffic.mgmt_tx.packets"] = t.Packets
	}

	point, err := client.NewPoint("node", tags, fields, time.Now())
	if err != nil {
		panic(err)
	}
	c.points <- point
}

func (c *StatsDb) Close() {
	close(c.points)
	c.wg.Wait()
	c.client.Close()
}

// stores data points in batches into the influxdb
func (c *StatsDb) worker() {
	bpConfig := client.BatchPointsConfig{
		Database:  config.Influxdb.Database,
		Precision: "m",
	}

	var bp client.BatchPoints
	var err error
	var dirty, closed bool
	var batchStarted time.Time

	for !closed {
		// create new batch points?
		if bp == nil {
			if bp, err = client.NewBatchPoints(bpConfig); err != nil {
				panic(err)
			}
		}

		// wait for new points
		select {
		case point, ok := <-c.points:
			if ok {
				bp.AddPoint(point)
				if !dirty {
					batchStarted = time.Now()
					dirty = true
				}
			} else {
				closed = true
			}
		case <-time.After(time.Second):
			// nothing
		}

		// write batch now?
		if (dirty && batchStarted.Add(batchWaiting).Before(time.Now())) || closed {
			log.Println("saving", len(bp.Points()), "points")

			if err = c.client.Write(bp); err != nil {
				panic(err)
			}
			dirty = false
			bp = nil
		}
	}

	c.wg.Done()
}
