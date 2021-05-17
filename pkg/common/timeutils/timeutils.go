package timeutils

import (
	"errors"
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
	var timeFormat string
	if params.TimeFormat == "" {
		timeFormat = keptnTimeFormatISO8601
	} else {
		timeFormat = params.TimeFormat
	}
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

	// initialize default values for end and start time
	end := time.Now().UTC()
	start := time.Now().UTC().Add(-timeframeDuration)

	// Parse start date
	if params.StartDate != "" {
		start, err = time.Parse(timeFormat, params.StartDate)

		if err != nil {
			return nil, nil, err
		}
	}

	// Parse end date
	if params.EndDate != "" {
		end, err = time.Parse(timeFormat, params.EndDate)

		if err != nil {
			return nil, nil, err
		}
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
