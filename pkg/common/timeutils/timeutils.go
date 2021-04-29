package timeutils

import "time"

const keptnTimeFormatISO8601 = "2006-01-02T15:04:05.000Z"

// GetKeptnTimeStamp formats a given timestamp into the format used by
// Keptn which is following the ISO 8601 standard
func GetKeptnTimeStamp(timestamp time.Time) string {
	return timestamp.Format(keptnTimeFormatISO8601)
}
