package api

import (
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func getTestHTTPServer(handlerFunc func(writer http.ResponseWriter, request *http.Request)) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(handlerFunc))

	return server
}

func TestLogHandler_DeleteLogs(t *testing.T) {
	type args struct {
		params models.LogFilter
	}
	tests := []struct {
		name             string
		args             args
		httpResponseFunc func(writer http.ResponseWriter, request *http.Request)
		want             error
	}{
		{
			name: "deletion successful",
			args: args{
				params: models.LogFilter{},
			},
			httpResponseFunc: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(""))
			},
			want: nil,
		},
		{
			name: "deletion failed",
			args: args{
				params: models.LogFilter{},
			},
			httpResponseFunc: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(`{"code":0, "message":"oops"}`))
			},
			want: errors.New("oops"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := getTestHTTPServer(tt.httpResponseFunc)
			defer ts.Close()

			lh := NewLogHandler(ts.URL)

			got := lh.DeleteLogs(tt.args.params)
			require.Equal(t, tt.want, got)
		})
	}
}

func stringp(s string) *string {
	return &s
}

func TestLogHandler_Flush(t *testing.T) {
	tests := []struct {
		name             string
		httpResponseFunc func(writer http.ResponseWriter, request *http.Request)
		wantErr          bool
	}{
		{
			name: "writing logs successful",
			httpResponseFunc: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(""))
			},
			wantErr: false,
		},
		{
			name: "writing logs failed",
			httpResponseFunc: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(`{"code":0, "message":"oops"}`))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := getTestHTTPServer(tt.httpResponseFunc)
			defer ts.Close()

			lh := NewLogHandler(ts.URL)

			lh.LogCache = []models.LogEntry{
				{
					IntegrationID: "id",
				},
			}
			got := lh.Flush()
			if tt.wantErr {
				require.NotNil(t, got)
			} else {
				require.Nil(t, got)
			}
		})
	}
}

func TestLogHandler_GetLogs(t *testing.T) {
	tests := []struct {
		name             string
		httpResponseFunc func(writer http.ResponseWriter, request *http.Request)
		want             *models.GetLogsResponse
		wantErr          error
	}{
		{
			name: "retrieve logs",
			httpResponseFunc: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(`{"logs": [{"integrationid": "my-id", "message":"my-message"}]}`))
			},
			want: &models.GetLogsResponse{
				Logs: []models.LogEntry{
					{
						IntegrationID: "my-id",
						Message:       "my-message",
						Time:          time.Time{},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "retrieving logs failed",
			httpResponseFunc: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(`{"code":0, "message":"oops"}`))
			},
			wantErr: errors.New("oops"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := getTestHTTPServer(tt.httpResponseFunc)
			defer ts.Close()

			lh := NewLogHandler(ts.URL)

			got, err := lh.GetLogs(models.GetLogsParams{})
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestLogHandler_Log(t *testing.T) {
	lh := NewLogHandler("")

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i = i + 1 {
		go func() {
			lh.Log([]models.LogEntry{
				{
					IntegrationID: "my-id",
					Message:       "message",
				},
			})
			wg.Done()
		}()
	}
	wg.Wait()

	require.Len(t, lh.LogCache, 100)
}

func TestLogHandler_Start(t *testing.T) {
	endpointCalled := false
	ts := getTestHTTPServer(func(writer http.ResponseWriter, request *http.Request) {
		endpointCalled = true
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(""))
	})

	defer ts.Close()

	mockClock := clock.NewMock()
	lh := NewLogHandler(ts.URL)
	lh.TheClock = mockClock

	lh.Start(context.Background())

	mockClock.Add(60 * time.Second)

	require.Eventually(t, func() bool {
		return endpointCalled
	}, 5*time.Second, 1*time.Second)
}
