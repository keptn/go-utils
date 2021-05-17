package timeutils

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetStartEndTime(t *testing.T) {
	type args struct {
		startDatePoint string
		endDatePoint   string
		timeframe      string
		timeFormat     string
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{
			name: "start and end date provided - return those",
			args: args{
				startDatePoint: time.Now().Round(time.Minute).UTC().Format(keptnTimeFormatISO8601),
				endDatePoint:   time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format(keptnTimeFormatISO8601),
				timeframe:      "",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "start and end date provided (different time format) - return those",
			args: args{
				startDatePoint: time.Now().Round(time.Minute).UTC().Format("2006-01-02T15:04:05"),
				endDatePoint:   time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format("2006-01-02T15:04:05"),
				timeframe:      "",
				timeFormat:     "2006-01-02T15:04:05",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "start and timeframe - return startdate and startdate + timeframe",
			args: args{
				startDatePoint: time.Now().Round(time.Minute).UTC().Format(keptnTimeFormatISO8601),
				endDatePoint:   "",
				timeframe:      "10m",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(10 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "only timeframe provided - return time.Now - timeframe and time.Now",
			args: args{
				startDatePoint: "",
				endDatePoint:   "",
				timeframe:      "10m",
			},
			wantStart: time.Now().UTC().Add(-10 * time.Minute).Round(time.Minute).UTC(),
			wantEnd:   time.Now().UTC().Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "startDate > endDate provided - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				endDatePoint:   time.Now().UTC().Format(keptnTimeFormatISO8601),
				timeframe:      "",
			},
			wantErr: true,
		},
		{
			name: "startDate, endDate and timeframe provided - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				endDatePoint:   time.Now().UTC().Format(keptnTimeFormatISO8601),
				timeframe:      "5m",
			},
			wantErr: true,
		},
		{
			name: "startDate provided, but no endDate or timeframe - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				endDatePoint:   "",
				timeframe:      "",
			},
			wantErr: true,
		},
		{
			name: "endDate provided, but no startDate or timeframe - return error",
			args: args{
				startDatePoint: "",
				endDatePoint:   time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				timeframe:      "",
			},
			wantErr: true,
		},
		{
			name: "invalid timeframe string - return error",
			args: args{
				startDatePoint: "",
				endDatePoint:   "",
				timeframe:      "xyz",
			},
			wantErr: true,
		},
		{
			name: "invalid timeframe string - return error",
			args: args{
				startDatePoint: "",
				endDatePoint:   "",
				timeframe:      "xym",
			},
			wantErr: true,
		},
		{
			name: "invalid startDate string - return error",
			args: args{
				startDatePoint: "abc",
				endDatePoint:   "",
				timeframe:      "5m",
			},
			wantErr: true,
		},
		{
			name: "invalid endDate string - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(keptnTimeFormatISO8601),
				endDatePoint:   "abc",
				timeframe:      "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := GetStartEndTime(tt.args.startDatePoint, tt.args.endDatePoint, tt.args.timeframe, tt.args.timeFormat)
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
