package main

import (
	"log"
	"sync"
	"time"

	"github.com/ffdo/node-informant/gluon-collector/data"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	saveInterval = time.Second * 5
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

func (c *StatsDb) Add(stats *data.StatisticsStruct) {
	tags := map[string]string{
		"nodeid": stats.NodeId,
	}
	fields := map[string]interface{}{
		"load":                   stats.LoadAverage,
		"processes.running":      stats.Processes.Running,
		"clients.wifi":           stats.Clients.Wifi,
		"clients.total":          stats.Clients.Total,
		"traffic.forward":        stats.Traffic.Forward,
		"traffic.rx":             stats.Traffic.Rx,
		"traffic.tx":             stats.Traffic.Tx,
		"traffic.mgmt.rx":        stats.Traffic.MgmtRx,
		"traffic.mgmt.tx":        stats.Traffic.MgmtTx,
		"traffic.memory.buffers": stats.Memory.Buffers,
		"traffic.memory.cached":  stats.Memory.Cached,
		"traffic.memory.free":    stats.Memory.Free,
		"traffic.memory.total":   stats.Memory.Total,
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

func (c *StatsDb) worker() {
	lastSent := time.Now()
	bpConfig := client.BatchPointsConfig{
		Database:  config.Influxdb.Database,
		Precision: "m",
	}

	var bp client.BatchPoints
	var err error
	var abort bool
	var dirty bool

	for {
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
				dirty = true
			} else {
				abort = true
			}
		case <-time.After(time.Second):
			// nothing
		}

		// write now?
		if dirty && (abort || lastSent.Add(saveInterval).Before(time.Now())) {
			log.Println("saving", len(bp.Points()), "points")

			if err := c.client.Write(bp); err != nil {
				panic(err)
			}
			lastSent = time.Now()
			dirty = false
			bp = nil
		}

		if abort {
			break
		}
	}

	c.wg.Done()
}
