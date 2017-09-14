package graphite

import (
	"log"
	"sync"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/fgrosse/graphigo"
)

const (
	MeasurementNode            = "node"     // Measurement for per-node statistics
	MeasurementGlobal          = "global"   // Measurement for summarized global statistics
	CounterMeasurementFirmware = "firmware" // Measurement for firmware statistics
	CounterMeasurementModel    = "model"    // Measurement for model statistics
)

type Connection struct {
	database.Connection
	client graphigo.Client
	points chan []graphigo.Metric
	wg     sync.WaitGroup
}

type Config map[string]interface{}

func (c Config) Address() string {
	return c["address"].(string)
}

func (c Config) Prefix() string {
	return c["prefix"].(string)
}

func (c Config) Enable() bool {
	return c["enable"].(bool)
}

func Connect(configuration interface{}) (database.Connection, error) {
	var config Config

	config = configuration.(map[string]interface{})

	if !config.Enable() {
		return nil, nil
	}

	con := &Connection{
		client: graphigo.Client{
			Address: config.Address(),
			Prefix:  config.Prefix(),
		},
		points: make(chan []graphigo.Metric, 1000),
	}

	if err := con.client.Connect(); err != nil {
		return nil, err
	}

	con.wg.Add(1)
	go con.addWorker()

	return con, nil
}

func (c *Connection) Close() {
	close(c.points)
	if c.client.Connection != nil {
		c.client.Close()
	}
}

func (c *Connection) addWorker() {
	defer c.wg.Done()
	defer c.Close()
	for point := range c.points {
		err := c.client.SendAll(point)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func (c *Connection) addPoint(point []graphigo.Metric) {
	c.points <- point
}

func init() {
	database.RegisterAdapter("graphite", Connect)
}
