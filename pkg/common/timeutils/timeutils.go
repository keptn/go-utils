package timeutils

import (
	"fmt"
	"time"
)

const keptnTimeFormatISO8601 = "2006-01-02T15:04:05.000Z"

const defaultEvaluationTimeframe = "5m"

// GetKeptnTimeStamp formats a given timestamp into the format used by
// Keptn which is following the ISO 8601 standard
func GetKeptnTimeStamp(timestamp time.Time) string {
	return timestamp.Format(keptnTimeFormatISO8601)
}

// GetStartEndTime parses the provided start date, end date and/or timeframe
func GetStartEndTime(startDatePoint, endDatePoint, timeframe, timeFormat string) (*time.Time, *time.Time, error) {
	if timeFormat == "" {
		timeFormat = keptnTimeFormatISO8601
	}
	var err error
	// input validation
	if startDatePoint != "" && endDatePoint == "" {
		// if a start date is set, but no end date is set, we require the timeframe to be set
		if timeframe == "" {
			errMsg := "no timeframe or end date provided"

			return nil, nil, fmt.Errorf(errMsg)
		}
	}
	if endDatePoint != "" && timeframe != "" {
		// can not use end date and timeframe at the same time
		errMsg := "You can not use 'end' together with 'timeframe'"

		return nil, nil, fmt.Errorf(errMsg)
	}
	if endDatePoint != "" && startDatePoint == "" {
		errMsg := "start date is required when using an end date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	// parse timeframe
	var timeframeDuration time.Duration
	if timeframe != "" {
		timeframeDuration, err = time.ParseDuration(timeframe)
	} else {
		timeframeDuration, err = time.ParseDuration(defaultEvaluationTimeframe)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse provided timeframe: %s", err.Error())
	}

	// initialize default values for end and start time
	end := time.Now().UTC()
	start := time.Now().UTC().Add(-timeframeDuration)

	// Parse start date
	if startDatePoint != "" {
		start, err = time.Parse(timeFormat, startDatePoint)

		if err != nil {
			return nil, nil, err
		}
	}

	// Parse end date
	if endDatePoint != "" {
		end, err = time.Parse(timeFormat, endDatePoint)

		if err != nil {
			return nil, nil, err
		}
	}

	// last but not least: if a start date and a timeframe is provided, we set the end date to start date + timeframe
	if startDatePoint != "" && endDatePoint == "" && timeframe != "" {
		end = start.Add(timeframeDuration)
	}

	// ensure end date is greater than start date
	diff := end.Sub(start).Minutes()

	if diff < 1 {
		errMsg := "end date must be at least 1 minute after start date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	return &start, &end, nil
}
