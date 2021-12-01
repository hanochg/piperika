package utils

import (
	"time"
)

func PipelinesTimeParser(timeStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	//e.g - "2021-11-09T06:09:08.239Z"
	return time.Parse(layout, timeStr)
}
