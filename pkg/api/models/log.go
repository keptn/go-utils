package models

import "time"

type LogEntry struct {
	IntegrationID string    `json:"integrationid" bson:"integrationid"`
	Message       string    `json:"message" bson:"message"`
	Time          time.Time `json:"time" bson:"time"`
	KeptnContext  string    `json:"shkeptncontext" bson:"shkeptncontext"`
	Task          string    `json:"task" bson:"task"`
	TriggeredID   string    `json:"triggeredid" bson:"triggeredid"`
}

type GetLogsParams struct {
	LogFilter
	PageSize    int
	NextPageKey int
}

type GetLogsResponse struct {
	NextPageKey int64      `json:"nextPageKey,omitempty"`
	PageSize    int64      `json:"pageSize,omitempty"`
	TotalCount  int64      `json:"totalCount,omitempty"`
	Logs        []LogEntry `json:"logs"`
}

type LogFilter struct {
	IntegrationID string
	FromTime      string
	BeforeTime    string
}
