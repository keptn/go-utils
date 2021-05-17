package timeutils

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetStartEndTime(t *testing.T) {
	const userFriendlyTimeFormat = "2006-01-02T15:04:05"
	tests := []struct {
		name      string
		args      GetStartEndTimeParams
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{
			name: "start and end date provided - return those",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Round(time.Minute).UTC().Format(keptnTimeFormatISO8601),
				EndDate:   time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format(keptnTimeFormatISO8601),
				Timeframe: "",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "start and end date provided (different time format) - return those",
			args: GetStartEndTimeParams{
				StartDate:  time.Now().Round(time.Minute).UTC().Format(userFriendlyTimeFormat),
				EndDate:    time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format(userFriendlyTimeFormat),
				Timeframe:  "",
				TimeFormat: userFriendlyTimeFormat,
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "start and timeframe - return startdate and startdate + timeframe",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Round(time.Minute).UTC().Format(keptnTimeFormatISO8601),
				EndDate:   "",
				Timeframe: "10m",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(10 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "only timeframe provided - return time.Now - timeframe and time.Now",
			args: GetStartEndTimeParams{
				StartDate: "",
				EndDate:   "",
				Timeframe: "10m",
			},
			wantStart: time.Now().UTC().Add(-10 * time.Minute).Round(time.Minute).UTC(),
			wantEnd:   time.Now().UTC().Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "startDate > endDate provided - return error",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				EndDate:   time.Now().UTC().Format(keptnTimeFormatISO8601),
				Timeframe: "",
			},
			wantErr: true,
		},
		{
			name: "startDate, endDate and timeframe provided - return error",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				EndDate:   time.Now().UTC().Format(keptnTimeFormatISO8601),
				Timeframe: "5m",
			},
			wantErr: true,
		},
		{
			name: "startDate provided, but no endDate or timeframe - return error",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				EndDate:   "",
				Timeframe: "",
			},
			wantErr: true,
		},
		{
			name: "endDate provided, but no startDate or timeframe - return error",
			args: GetStartEndTimeParams{
				StartDate: "",
				EndDate:   time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				Timeframe: "",
			},
			wantErr: true,
		},
		{
			name: "invalid timeframe string - return error",
			args: GetStartEndTimeParams{
				StartDate: "",
				EndDate:   "",
				Timeframe: "xyz",
			},
			wantErr: true,
		},
		{
			name: "invalid timeframe string - return error",
			args: GetStartEndTimeParams{
				StartDate: "",
				EndDate:   "",
				Timeframe: "xym",
			},
			wantErr: true,
		},
		{
			name: "invalid startDate string - return error",
			args: GetStartEndTimeParams{
				StartDate: "abc",
				EndDate:   "",
				Timeframe: "5m",
			},
			wantErr: true,
		},
		{
			name: "invalid endDate string - return error",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				EndDate:   "abc",
				Timeframe: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := GetStartEndTime(tt.args)
			if tt.wantErr {
				require.NotNil(t, err)
				require.Nil(t, gotStart)
				require.Nil(t, gotEnd)
			} else {
				require.Nil(t, err)
				require.WithinDuration(t, *gotStart, tt.wantStart, time.Minute)
				require.WithinDuration(t, *gotEnd, tt.wantEnd, time.Minute)
			}
		})
	}
}
