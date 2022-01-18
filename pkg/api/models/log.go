package models

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	IntegrationID string    `json:"integrationid" bson:"integrationid"`
	Message       string    `json:"message" bson:"message"`
	Time          time.Time `json:"time" bson:"time"`
	KeptnContext  string    `json:"shkeptncontext" bson:"shkeptncontext"`
	Task          string    `json:"task" bson:"task"`
	TriggeredID   string    `json:"triggeredid" bson:"triggeredid"`
	GitCommitID   string    `json:"gitcommitid" bson:"gitcommitid"`
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

type CreateLogsRequest struct {
	// logs
	Logs []LogEntry `form:"logs" json:"logs"`
}

// ToJSON converts object to JSON string
func (l *LogEntry) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

// FromJSON converts JSON string to object
func (l *LogEntry) FromJSON(b []byte) error {
	var res LogEntry
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*l = res
	return nil
}

// ToJSON converts object to JSON string
func (l *GetLogsResponse) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

// FromJSON converts JSON string to object
func (l *GetLogsResponse) FromJSON(b []byte) error {
	var res GetLogsResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*l = res
	return nil
}

// ToJSON converts object to JSON string
func (l *CreateLogsRequest) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

// FromJSON converts JSON string to object
func (l *CreateLogsRequest) FromJSON(b []byte) error {
	var res CreateLogsRequest
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*l = res
	return nil
}
