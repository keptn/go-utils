package timeutils

import (
	"errors"
	"fmt"
	"time"
)

const KeptnTimeFormatISO8601 = "2006-01-02T15:04:05.000Z"
const fallbackTimeformat = "2006-01-02T15:04:05"
const defaultEvaluationTimeframe = "5m"

// GetKeptnTimeStamp formats a given timestamp into the format used by
// Keptn which is following the ISO 8601 standard
func GetKeptnTimeStamp(timestamp time.Time) string {
	return timestamp.Format(KeptnTimeFormatISO8601)
}

// ParseTimestamp tries to parse the given timestamp using the ISO8601 format (e.g. '2006-01-02T15:04:05.000Z')
// if this is not possible, the fallback format '2006-01-02T15:04:05.000000000Z' will be used, then RFC3339
//If this fails as well, an error is returned

func ParseTimestamp(timestamp string) (*time.Time, error) {
	parsedTime, err := time.Parse(KeptnTimeFormatISO8601, timestamp)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			parsedTime, err = time.Parse(fallbackTimeformat, timestamp)
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

	/*if params.TimeFormat == "" {
		timeFormat = KeptnTimeFormatISO8601

	} else {
		timeFormat = params.TimeFormat
	}*/

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

	// initialize default values for end and start time
	end := time.Now().UTC()
	start := time.Now().UTC().Add(-timeframeDuration)
	var temp *time.Time
	// Parse start date
	if params.StartDate != "" {
		temp, err = ParseTimestamp(params.StartDate)
		if err != nil {
			return nil, nil, err
		}
		start = *temp
	}

	// Parse end date
	if params.EndDate != "" {
		temp, err = ParseTimestamp(params.EndDate)
		if err != nil {
			return nil, nil, err
		}
		end = *temp
	}

	// last but not least: if a start date and a timeframe is provided, we set the end date to start date + timeframe
	if params.StartDate != "" && params.EndDate == "" && params.Timeframe != "" {
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
