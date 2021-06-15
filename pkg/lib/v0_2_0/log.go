package v0_2_0

const ErrorLogEventName = "sh.keptn.log.error"

type ErrorLogEvent struct {
	Message       string `json:"message"`
	IntegrationID string `json:"integrationId"`
	Task          string `json:"task"`
}
