package jsontime

import (
	"time"
)

const TimeFormat = "2006-01-02T15:04:05"

type Time time.Time

func Now() Time {
	return Time(time.Now())
}

func (t Time) MarshalJSON() ([]byte, error) {
	stamp := `"` + time.Time(t).Format(TimeFormat) + `"`
	return []byte(stamp), nil
}

func (t Time) UnmarshalJSON(data []byte) (err error) {
	if nativeTime, err := time.Parse(TimeFormat, string(data)); err == nil {
		t = Time(nativeTime)
	}
	return
}
