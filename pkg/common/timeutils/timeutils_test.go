package timeutils

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetStartEndTime(t *testing.T) {
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
				StartDate: time.Now().Round(time.Minute).UTC().Format(KeptnTimeFormatISO8601),
				EndDate:   time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format(KeptnTimeFormatISO8601),
				Timeframe: "",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "start and end date provided (different time format) - return those",
			args: GetStartEndTimeParams{
				StartDate:  time.Now().Round(time.Minute).UTC().Format("2006-01-02T15:04:05"),
				EndDate:    time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format("2006-01-02T15:04:05"),
				Timeframe:  "",
				TimeFormat: "2006-01-02T15:04:05",
			},
			wantStart: time.Now().Round(time.Minute).UTC(),
			wantEnd:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr:   false,
		},
		{
			name: "start and timeframe - return startdate and startdate + timeframe",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Round(time.Minute).UTC().Format(KeptnTimeFormatISO8601),
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
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(KeptnTimeFormatISO8601),
				EndDate:   time.Now().UTC().Format(KeptnTimeFormatISO8601),
				Timeframe: "",
			},
			wantErr: true,
		},
		{
			name: "startDate, endDate and timeframe provided - return error",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(KeptnTimeFormatISO8601),
				EndDate:   time.Now().UTC().Format(KeptnTimeFormatISO8601),
				Timeframe: "5m",
			},
			wantErr: true,
		},
		{
			name: "startDate provided, but no endDate or timeframe - return error",
			args: GetStartEndTimeParams{
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(KeptnTimeFormatISO8601),
				EndDate:   "",
				Timeframe: "",
			},
			wantErr: true,
		},
		{
			name: "endDate provided, but no startDate or timeframe - return error",
			args: GetStartEndTimeParams{
				StartDate: "",
				EndDate:   time.Now().Add(1 * time.Minute).UTC().Format(KeptnTimeFormatISO8601),
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
				StartDate: time.Now().Add(1 * time.Minute).UTC().Format(KeptnTimeFormatISO8601),
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

func TestParseTimestamp(t *testing.T) {

	correctISO8601Timestamp := "2020-01-02T15:04:05.000Z"
	correctFallbackTimestamp := "2020-01-02T15:04:05.000000000Z"

	timeObj, _ := time.Parse(KeptnTimeFormatISO8601, correctISO8601Timestamp)

	type args struct {
		timestamp string
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Time
		wantErr bool
	}{
		{
			name: "correct timestamp provided",
			args: args{
				timestamp: correctISO8601Timestamp,
			},
			want:    &timeObj,
			wantErr: false,
		},
		{
			name: "correct fallback timestamp provided",
			args: args{
				timestamp: correctFallbackTimestamp,
			},
			want:    &timeObj,
			wantErr: false,
		},
		{
			name: "incorrect timestamp provided",
			args: args{
				timestamp: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimestamp(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTimestamp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {

	correctISO8601Timestamp := "2020-01-02T15:04:05.000Z"
	correctFallbackTimestamp := "2020-01-02T15:04:05.000000000Z"

	timeObj, _ := time.Parse(KeptnTimeFormatISO8601, correctISO8601Timestamp)

	type args struct {
		timestamp string
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Time
		wantErr bool
	}{
		{
			name: "correct timestamp provided",
			args: args{
				timestamp: correctISO8601Timestamp,
			},
			want:    &timeObj,
			wantErr: false,
		},
		{
			name: "correct fallback timestamp provided",
			args: args{
				timestamp: correctFallbackTimestamp,
			},
			want:    &timeObj,
			wantErr: false,
		},
		{
			name: "incorrect timestamp provided",
			args: args{
				timestamp: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimestamp(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTimestamp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVariousTimeFormat(t *testing.T) {
	times := []string{
		"2020-01-02T15:00:00",
		"2020-01-02T15:00:00Z",
		"2020-01-02T15:00:00+10:00",
		"2020-01-02T15:00:00.000Z",
		"2020-01-02T15:00:00.000000000Z",
	}

	for _, time := range times {

		_, err := ParseTimestamp(time)
		require.Nil(t, err)
	}

}
