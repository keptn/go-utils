package timeutils

import (
	"errors"
	"fmt"
	"time"
)

const KeptnTimeFormatISO8601 = "2006-01-02T15:04:05.000Z"
const fallbackTimeFormat = "2006-01-02T15:04:05"
const defaultEvaluationTimeframe = "5m"

// GetKeptnTimeStamp formats a given timestamp into the format used by
// Keptn which is following the ISO 8601 standard
func GetKeptnTimeStamp(timestamp time.Time) string {
	return timestamp.Format(KeptnTimeFormatISO8601)
}

// ParseTimestamp tries to parse the given timestamp using the ISO8601 format (e.g. '2006-01-02T15:04:05.000Z').
// If this is not possible, the fallback format RFC3339  will be used, then '2006-01-02T15:04:05'.
// If this fails as well, an error is returned

func ParseTimestamp(timestamp string) (*time.Time, error) {
	parsedTime, err := time.Parse(KeptnTimeFormatISO8601, timestamp)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			parsedTime, err = time.Parse(fallbackTimeFormat, timestamp)
			if err != nil {
				return nil, err
			}

		}
	}
	return &parsedTime, nil
}

type GetStartEndTimeParams struct {
	StartDate  string
	EndDate    string
	Timeframe  string
	TimeFormat string
}

func (params *GetStartEndTimeParams) Validate() error {
	if params.StartDate != "" && params.EndDate == "" {
		// if a start date is set, but no end date is set, we require the timeframe to be set
		if params.Timeframe == "" {
			return fmt.Errorf("no timeframe or end date provided")
		}
	}
	if params.EndDate != "" && params.Timeframe != "" {
		return fmt.Errorf("'end' and 'timeframe' are mutually exclusive")
	}
	if params.EndDate != "" && params.StartDate == "" {
		return errors.New("start date is required when using an end date")
	}
	return nil
}

// GetStartEndTime parses the provided start date, end date and/or timeframe
func GetStartEndTime(params GetStartEndTimeParams) (*time.Time, *time.Time, error) {
	//var timeFormat string
	var err error

	// input validation
	if err := params.Validate(); err != nil {
		return nil, nil, err
	}

	// parse timeframe
	var timeframeDuration time.Duration
	if params.Timeframe != "" {
		timeframeDuration, err = time.ParseDuration(params.Timeframe)
	} else {
		timeframeDuration, err = time.ParseDuration(defaultEvaluationTimeframe)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse provided timeframe: %s", err.Error())
	}

	// calculate possible start and end time
	now := time.Now().UTC()
	calcStart := now.Add(-timeframeDuration)

	// initialize default values for end and start time pointer
	var end, start *time.Time
	end = &now
	start = &calcStart

	// Parse start date
	if params.StartDate != "" {
		start, err = ParseTimestamp(params.StartDate)
		if err != nil {
			return nil, nil, err
		}
	}

	// Parse end date
	if params.EndDate != "" {
		end, err = ParseTimestamp(params.EndDate)
		if err != nil {
			return nil, nil, err
		}
	}

	// last but not least: if a start date and a timeframe is provided, we set the end date to start date + timeframe
	if params.StartDate != "" && params.EndDate == "" && params.Timeframe != "" {
		*end = start.Add(timeframeDuration)
	}

	// ensure end date is greater than start date
	diff := end.Sub(*start).Minutes()

	if diff < 1 {
		errMsg := "end date must be at least 1 minute after start date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	return start, end, nil
}
