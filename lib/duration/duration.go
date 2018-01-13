package duration

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Duration is a TOML datatype
// A duration string is a possibly signed sequence of
// decimal numbers  and a unit suffix,
// such as "300s", "1.5h" or "5d".
// Valid time units are "s", "m", "h", "d", "w".
type Duration struct {
	time.Duration
}

// UnmarshalTOML parses a duration string.
func (d *Duration) UnmarshalText(data []byte) error {

	// " + int + unit + "
	if len(data) < 2 {
		return fmt.Errorf("invalid duration: \"%s\"", data)
	}

	unit := data[len(data)-1]
	value, err := strconv.Atoi(string(data[:len(data)-1]))
	if err != nil {
		return errors.Wrapf(err, "unable to parse duration \"%s\"", data)
	}

	switch unit {
	case 's':
		d.Duration = time.Duration(value) * time.Second
	case 'm':
		d.Duration = time.Duration(value) * time.Minute
	case 'h':
		d.Duration = time.Duration(value) * time.Hour
	case 'd':
		d.Duration = time.Duration(value) * time.Hour * 24
	case 'w':
		d.Duration = time.Duration(value) * time.Hour * 24 * 7
	case 'y':
		d.Duration = time.Duration(value) * time.Hour * 24 * 365
	default:
		return fmt.Errorf("invalid duration unit \"%s\"", string(unit))
	}

	return nil
}
