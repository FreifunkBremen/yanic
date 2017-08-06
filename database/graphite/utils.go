package graphite

import (
	"regexp"
)

var reInvalidChars = regexp.MustCompile("(?i)[^a-z0-9\\-]")

func replaceInvalidChars(name string) string {
	return reInvalidChars.ReplaceAllString(name, "_")
}
