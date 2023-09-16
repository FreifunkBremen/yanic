package influxdb

import (
	"context"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	influxdbAPI "github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/bdlm/log"
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
)

type Connection struct {
	database.Connection
	config   Config
	client   influxdb.Client
	writeAPI map[string]influxdbAPI.WriteAPI
}

type Config map[string]interface{}

func (c Config) Address() string {
	return c["address"].(string)
}
func (c Config) Token() string {
	if d, ok := c["token"]; ok {
		return d.(string)
	}
	log.Panic("influxdb2 - no token given")
	return ""
}
func (c Config) Organization() string {
	if d, ok := c["organization_id"]; ok {
		return d.(string)
	}
	return ""
}
func (c Config) Bucket(measurement string) string {
	logger := log.WithFields(map[string]interface{}{
		"organization_id": c.Organization(),
		"address":         c.Address(),
		"measurement":     measurement,
	})
	if d, ok := c["buckets"]; ok {
		dMap := d.(map[string]interface{})
		if d, ok := dMap[measurement]; ok {
			bucket := d.(string)
			logger.WithField("bucket", bucket).Info("get bucket for writeapi")
			return bucket
		}
		if d, ok := c["bucket_default"]; ok {
			bucket := d.(string)
			logger.WithField("bucket", bucket).Info("get bucket for writeapi")
			return bucket
		}
	}
	if d, ok := c["bucket_default"]; ok {
		bucket := d.(string)
		logger.WithField("bucket", bucket).Info("get bucket for writeapi")
		return bucket
	}
	logger.Panic("no bucket found for measurement")
	return ""
}
func (c Config) Tags() map[string]string {
	if c["tags"] != nil {
		tags := make(map[string]string)
		for k, v := range c["tags"].(map[string]interface{}) {
			tags[k] = v.(string)
		}
		return tags
	}
	return nil
}

func init() {
	database.RegisterAdapter("influxdb2", Connect)
}
func Connect(configuration map[string]interface{}) (database.Connection, error) {
	config := Config(configuration)

	// Make client
	client := influxdb.NewClientWithOptions(config.Address(), config.Token(), influxdb.DefaultOptions().SetBatchSize(batchMaxSize))

	ok, err := client.Ping(context.Background())
	if !ok || err != nil {
		return nil, err
	}

	writeAPI := map[string]influxdbAPI.WriteAPI{
		MeasurementLink:               client.WriteAPI(config.Organization(), config.Bucket(MeasurementLink)),
		MeasurementNode:               client.WriteAPI(config.Organization(), config.Bucket(MeasurementNode)),
		MeasurementDHCP:               client.WriteAPI(config.Organization(), config.Bucket(MeasurementDHCP)),
		MeasurementGlobal:             client.WriteAPI(config.Organization(), config.Bucket(MeasurementGlobal)),
		CounterMeasurementFirmware:    client.WriteAPI(config.Organization(), config.Bucket(CounterMeasurementFirmware)),
		CounterMeasurementModel:       client.WriteAPI(config.Organization(), config.Bucket(CounterMeasurementModel)),
		CounterMeasurementAutoupdater: client.WriteAPI(config.Organization(), config.Bucket(CounterMeasurementAutoupdater)),
	}

	db := &Connection{
		config:   config,
		client:   client,
		writeAPI: writeAPI,
	}

	return db, nil
}

// Close all connection and clean up
func (conn *Connection) Close() {
	for _, api := range conn.writeAPI {
		api.Flush()
	}
	conn.client.Close()
}
