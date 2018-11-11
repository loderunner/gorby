package main

import (
	"time"
)

// MaxUnixTime is the maximum time that can be represented in Go
// as UNIX time. See https://stackoverflow.com/a/32620397/769262
// for details.
var MaxUnixTime = time.Unix(1<<63-62135596801, 999999999)

func copyMap(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for k, v := range src {
		dst[k] = make([]string, len(v))
		copy(dst[k], v)
	}
	return dst
}
