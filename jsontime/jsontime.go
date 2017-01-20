package jsontime

import (
	"errors"
	"time"
)

// TimeFormat of JSONTime
const TimeFormat = "2006-01-02T15:04:05-0700"

//Time struct of JSONTime
type Time struct {
	time time.Time
}

// Now current Time
func Now() Time {
	return Time{time.Now()}
}

//MarshalJSON to bytearray
func (t Time) MarshalJSON() ([]byte, error) {
	stamp := `"` + t.time.Format(TimeFormat) + `"`
	return []byte(stamp), nil
}

// UnmarshalJSON from bytearray
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid jsontime")
	}
	if nativeTime, err := time.Parse(TimeFormat, string(data[1:len(data)-1])); err == nil {
		t.time = nativeTime
	}
	return
}

// GetTime normal
func (t Time) GetTime() time.Time {
	return t.time
}

// Unix of this time
func (t Time) Unix() int64 {
	return t.time.Unix()
}

// IsZero is time zero?
func (t Time) IsZero() bool {
	return t.time.IsZero()
}

// Add given Duration to this time
func (t Time) Add(d time.Duration) Time {
	return Time{time: t.time.Add(d)}
}

// After is this time after the given?
func (t Time) After(u Time) bool {
	return t.time.After(u.GetTime())
}

// Before is this time before the given?
func (t Time) Before(u Time) bool {
	return t.time.Before(u.GetTime())
}
