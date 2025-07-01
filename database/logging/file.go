package logging

/**
 * This database type is just for,
 * - debugging without a influxconn
 * - example for other developers for new databases
 */
import (
	"log/slog"
	"os"
	"time"

	logger "github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	config Config
	file   *os.File
	logger *slog.Logger
}

type Config map[string]interface{}

func (c Config) Path() string {
	return c["path"].(string)
}
func (c Config) Type() string {
	if t, ok := c["type"].(string); ok {
		return t
	}
	return "text"
}

func init() {
	database.RegisterAdapter("logging", Connect)
}

func Connect(configuration map[string]interface{}) (database.Connection, error) {
	config := Config(configuration)

	file, err := os.OpenFile(config.Path(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	var handler slog.Handler
	switch config.Type() {
	case "text":
		handler = slog.NewTextHandler(file, nil)
	case "json":
		handler = slog.NewJSONHandler(file, nil)
	}

	return &Connection{
		config: config,
		logger: slog.New(handler),
	}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	conn.logger.Info("InsertNode",
		slog.Group("node",
			slog.Group("statistics",
				slog.String("node_id", node.Statistics.NodeID),
				slog.Any("clients", node.Statistics.Clients.Total),
			),
		),
	)
}

func (conn *Connection) InsertLink(link *runtime.Link, insertTimestamp time.Time) {
	conn.logger.Info("InsertLink",
		slog.Any("link", link),
		slog.Time("insert_timestamp", insertTimestamp),
	)
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, insertTimestamp time.Time, site string, domain string) {
	conn.logger.Info("InsertGlobals",
		slog.Group("stats",
			slog.Any("nodes", stats.Nodes),
			slog.Any("clients", stats.Clients),
			slog.Int("models", len(stats.Models)),
		),
		slog.Time("insert_timestamp", insertTimestamp),
		slog.String("site", site),
		slog.String("domain", domain),
	)
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	conn.logger.Info("PruneNodes",
		slog.Duration("delete_after", deleteAfter),
	)
}

func (conn *Connection) Close() {
	conn.logger.Info("Close")
	if err := conn.file.Close(); err != nil {
		logger.WithError(err).Error("unable to close connection")
	}
}
