package influxdb

import (
	"log"
	"sync"
	"time"

	"github.com/influxdata/influxdb1-client/models"
	"github.com/influxdata/influxdb1-client/v2"

	"github.com/FreifunkBremen/yanic/database"
)

const (
	MeasurementLink               = "link"        // Measurement for per-link statistics
	MeasurementNode               = "node"        // Measurement for per-node statistics
	MeasurementDHCP               = "dhcp"        // Measurement for DHCP server statistics
	MeasurementGlobal             = "global"      // Measurement for summarized global statistics
	CounterMeasurementFirmware    = "firmware"    // Measurement for firmware statistics
	CounterMeasurementModel       = "model"       // Measurement for model statistics
	CounterMeasurementAutoupdater = "autoupdater" // Measurement for autoupdater
	batchMaxSize                  = 1000
	batchTimeout                  = 5 * time.Second
)

type Connection struct {
	database.Connection
	config Config
	client client.Client
	points chan *client.Point
	wg     sync.WaitGroup
}

type Config map[string]interface{}

func (c Config) Address() string {
	return c["address"].(string)
}
func (c Config) Database() string {
	return c["database"].(string)
}
func (c Config) Username() string {
	return c["username"].(string)
}
func (c Config) Password() string {
	return c["password"].(string)
}
func (c Config) InsecureSkipVerify() bool {
	if d, ok := c["insecure_skip_verify"]; ok {
		return d.(bool)
	}
	return false
}
func (c Config) Tags() map[string]interface{} {
	if c["tags"] != nil {
		return c["tags"].(map[string]interface{})
	}
	return nil
}

func init() {
	database.RegisterAdapter("influxdb", Connect)
}
func Connect(configuration map[string]interface{}) (database.Connection, error) {
	var config Config
	config = configuration

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               config.Address(),
		Username:           config.Username(),
		Password:           config.Password(),
		InsecureSkipVerify: config.InsecureSkipVerify(),
	})

	if err != nil {
		return nil, err
	}

	_, _, err = c.Ping(time.Millisecond * 50)
	if err != nil {
		return nil, err
	}

	db := &Connection{
		config: config,
		client: c,
		points: make(chan *client.Point, batchMaxSize),
	}

	db.wg.Add(1)
	go db.addWorker()

	return db, nil
}

func (conn *Connection) addPoint(name string, tags models.Tags, fields models.Fields, t ...time.Time) {
	if configTags := conn.config.Tags(); configTags != nil {
		for tag, valueInterface := range configTags {
			if value, ok := valueInterface.(string); ok && tags.Get([]byte(tag)) == nil {
				tags.SetString(tag, value)
			} else {
				log.Println(name, "could not saved configured value of tag", tag)
			}
		}
	}
	point, err := client.NewPoint(name, tags.Map(), fields, t...)
	if err != nil {
		panic(err)
	}
	conn.points <- point
}

// Close all connection and clean up
func (conn *Connection) Close() {
	close(conn.points)
	conn.wg.Wait()
	conn.client.Close()
}

// stores data points in batches into the influxdb
func (conn *Connection) addWorker() {
	bpConfig := client.BatchPointsConfig{
		Database:  conn.config.Database(),
		Precision: "m",
	}

	var bp client.BatchPoints
	var err error
	var writeNow, closed bool
	timer := time.NewTimer(batchTimeout)

	for !closed {
		// wait for new points
		select {
		case point, ok := <-conn.points:
			if ok {
				if bp == nil {
					// create new batch
					timer.Reset(batchTimeout)
					if bp, err = client.NewBatchPoints(bpConfig); err != nil {
						log.Fatal(err)
					}
				}
				bp.AddPoint(point)
			} else {
				closed = true
			}
		case <-timer.C:
			if bp == nil {
				timer.Reset(batchTimeout)
			} else {
				writeNow = true
			}
		}

		// write batch now?
		if bp != nil && (writeNow || closed || len(bp.Points()) >= batchMaxSize) {
			log.Println("saving", len(bp.Points()), "points")

			if err = conn.client.Write(bp); err != nil {
				log.Print(err)
			}
			writeNow = false
			bp = nil
		}
	}
	timer.Stop()
	conn.wg.Done()
}
