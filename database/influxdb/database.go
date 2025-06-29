package influxdb

import (
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/bdlm/log"
	"github.com/influxdata/influxdb1-client/models"
	client "github.com/influxdata/influxdb1-client/v2"

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
	config := Config(configuration)

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               config.Address(),
		Username:           config.Username(),
		Password:           config.Password(),
		InsecureSkipVerify: config.InsecureSkipVerify(),
		Timeout:            batchTimeout,
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

func sanitizeValues(tags models.Tags) models.Tags {
	// https://docs.influxdata.com/influxdb/v2/reference/syntax/line-protocol/
	// Line protocol does not support the newline character \n in tag or field values.
	// To be safe, remove all non-printable characters and spaces except ASCII space and U+0020.
	for _, tag := range tags {
		cleaned_value := strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return ' '
		}, string(tag.Value))

		if cleaned_value != string(tag.Value) {
			tags.SetString(string(tag.Key), cleaned_value)
		}
	}
	return tags
}

func (conn *Connection) addPoint(name string, tags models.Tags, fields models.Fields, t ...time.Time) {
	if configTags := conn.config.Tags(); configTags != nil {
		for tag, valueInterface := range configTags {
			if value, ok := valueInterface.(string); ok && tags.Get([]byte(tag)) == nil {
				tags.SetString(tag, value)
			} else {
				log.WithFields(map[string]interface{}{
					"name": name,
					"tag":  tag,
				}).Warnf("could not save tag configuration on point")
			}
		}
	}

	tags = sanitizeValues(tags)

	point, err := client.NewPoint(name, tags.Map(), fields, t...)
	if err != nil {
		log.Panicf("could not save points: %s", err)
	}
	conn.points <- point
}

// Close all connection and clean up
func (conn *Connection) Close() {
	close(conn.points)
	conn.wg.Wait()
	if err := conn.client.Close(); err != nil {
		log.WithError(err).Error("during close connection")
	}
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
						log.WithError(err).Fatal("not able to create new batch for points")
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
			log.WithField("count", len(bp.Points())).Info("saving points")

			if err = conn.client.Write(bp); err != nil {
				log.WithError(err).Error("not able to write batch of points")
			}
			writeNow = false
			bp = nil
		}
	}
	timer.Stop()
	conn.wg.Done()
}
