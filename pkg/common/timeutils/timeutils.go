package timeutils

import "time"

const keptnTimeFormat = "2006-01-02T15:04:05.000Z"

func GetKeptnTimeStamp(timestamp time.Time) string {
	return timestamp.Format(keptnTimeFormat)
}
