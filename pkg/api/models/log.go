package models

import "time"

type LogEntry struct {
	IntegrationID string    `json:"integrationid" bson:"integrationid"`
	Message       string    `json:"message" bson:"message"`
	Time          time.Time `json:"time" bson:"time"`
}
