package all

import (
	_ "github.com/FreifunkBremen/yanic/database/graphite"
	_ "github.com/FreifunkBremen/yanic/database/influxdb"
	_ "github.com/FreifunkBremen/yanic/database/influxdb2"
	_ "github.com/FreifunkBremen/yanic/database/logging"
	_ "github.com/FreifunkBremen/yanic/database/respondd"
)
