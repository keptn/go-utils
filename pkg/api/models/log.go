package models

import "time"

type LogEntry struct {
	IntegrationID string    `json:"integrationid" bson:"integrationid"`
	Message       string    `json:"message" bson:"message"`
	Date          time.Time `json:"time" bson:"time"`
}
