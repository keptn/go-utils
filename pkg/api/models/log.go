package models

type LogEntry struct {
	IntegrationID string `json:"integrationid" bson:"integrationid"`
	Message       string `json:"message" bson:"message"`
	Time          string `json:"time" bson:"time"`
}
